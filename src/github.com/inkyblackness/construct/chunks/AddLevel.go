package chunks

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/res/logic"
	"github.com/inkyblackness/res/serial"
)

// AddLevel adds one level to the consumer
func AddLevel(chunkStore chunk.Store, levelID int, solid bool, isCyberspace bool) {
	levelBaseID := uint16(4000 + 100*levelID)

	chunkStore.Put(chunk.ID(levelBaseID+2), mapChunk(false, []byte{0x0B, 0x00, 0x00, 0x00}))
	chunkStore.Put(chunk.ID(levelBaseID+3), mapChunk(false, []byte{0x1B, 0x00, 0x00, 0x00}))

	AddBasicLevelInformation(chunkStore, levelBaseID, isCyberspace)
	AddMap(chunkStore, levelBaseID, solid, levelID == 1)
	AddLevelTimer(chunkStore, levelBaseID)
	AddLevelTextures(chunkStore, levelBaseID)
	AddMasterObjectTables(chunkStore, levelBaseID)
	AddLevelObjects(chunkStore, levelBaseID)

	chunkStore.Put(chunk.ID(levelBaseID+40), mapChunk(false, []byte{0x0D, 0x00, 0x00, 0x00}))
	chunkStore.Put(chunk.ID(levelBaseID+41), mapChunk(false, []byte{0x00}))
	chunkStore.Put(chunk.ID(levelBaseID+42), mapChunk(true, make([]byte, 0x1C)))

	AddSurveillanceChunk(chunkStore, levelBaseID)
	AddLevelVariables(chunkStore, levelBaseID)
	AddMapNotes(chunkStore, levelBaseID)

	chunkStore.Put(chunk.ID(levelBaseID+48), mapChunk(false, make([]byte, 0x30)))
	chunkStore.Put(chunk.ID(levelBaseID+49), mapChunk(true, make([]byte, 0x01C0)))
	chunkStore.Put(chunk.ID(levelBaseID+50), mapChunk(false, make([]byte, 2)))

	AddLoopConfiguration(chunkStore, levelBaseID)

	// CD-Release only content
	chunkStore.Put(chunk.ID(levelBaseID+52), mapChunk(false, make([]byte, 2)))
	chunkStore.Put(chunk.ID(levelBaseID+53), mapChunk(true, make([]byte, 0x40)))
}

func addTypedData(chunkStore chunk.Store, chunkID uint16, compressed bool, data interface{}) {
	store := serial.NewByteStore()
	coder := serial.NewEncoder(store)

	coder.Code(data)
	chunkStore.Put(chunk.ID(chunkID), mapChunk(compressed, store.Data()))
}

// AddBasicLevelInformation adds the basic level info block
func AddBasicLevelInformation(chunkStore chunk.Store, levelBaseID uint16, isCyberspace bool) {
	info := data.DefaultLevelInformation()

	if isCyberspace {
		info.CyberspaceFlag = 1
	}
	addTypedData(chunkStore, levelBaseID+4, true, info)
}

// AddLevelTimer adds the basic timer list structure
func AddLevelTimer(chunkStore chunk.Store, levelBaseID uint16) {
	chunkStore.Put(chunk.ID(levelBaseID+6), mapChunk(false, make([]byte, data.TimerEntrySize)))
}

// AddMap adds a map
func AddMap(chunkStore chunk.Store, levelBaseID uint16, solid bool, exceptStartingPosition bool) {
	tileFactory := func() interface{} {
		entry := data.DefaultTileMapEntry()

		if solid {
			entry.Type = data.Solid
		} else {
			entry.Type = data.Open
		}

		return entry
	}

	table := data.NewTable(64*64, tileFactory)
	for index, entry := range table.Entries {
		tile := entry.(*data.TileMapEntry)

		if solid && exceptStartingPosition && ((index % 64) == 30) && ((index / 64) == 22) {
			tile.Type = data.Open
		}
		// Block off outer border. Game locks up entering such tiles otherwise.
		if (index < 64) || ((index % 64) == 0) || ((index % 64) == 63) || (index > (64 * 63)) {
			tile.Type = data.Solid
		}
	}
	addTypedData(chunkStore, levelBaseID+5, true, table)
}

// AddLevelTextures adds level texture information
func AddLevelTextures(chunkStore chunk.Store, levelBaseID uint16) {
	data := make([]byte, 54*2)
	data[0] = 101 // energ-light
	chunkStore.Put(chunk.ID(levelBaseID+7), mapChunk(false, data))
}

// AddMasterObjectTables adds main object tables
func AddMasterObjectTables(chunkStore chunk.Store, levelBaseID uint16) {
	{
		masterCount := 872
		masters := make([]*data.LevelObjectEntry, masterCount)
		for index := range masters {
			master := data.DefaultLevelObjectEntry()

			masters[index] = master
			master.Previous = uint16((masterCount + index - 1) % masterCount)
		}

		masterTable := &data.Table{Entries: make([]interface{}, masterCount)}
		for i := range masters {
			masterTable.Entries[i] = masters[i]
		}
		addTypedData(chunkStore, levelBaseID+8, true, masterTable)
	}
	{
		crossrefList := logic.NewCrossReferenceList()

		crossrefList.Clear()
		addTypedData(chunkStore, levelBaseID+9, true, crossrefList.Encode())
	}
}

// AddLevelObjects adds level object tables
func AddLevelObjects(chunkStore chunk.Store, levelBaseID uint16) {
	for classID := uint16(0); classID < 15; classID++ {
		meta := data.LevelObjectClassMetaEntry(res.ObjectClass(classID))
		addLevelObjectTables(chunkStore, levelBaseID, classID, meta.EntrySize, meta.EntryCount)
	}
}

func addLevelObjectTables(chunkStore chunk.Store, levelBaseID uint16, classID uint16, entrySize int, entryCount int) {
	compressed := classID != 0 && classID != 1 && classID != 4 && classID != 5 && classID != 6
	table := logic.NewLevelObjectClassTable(entrySize, entryCount)
	table.AsChain().Initialize(entryCount - 1)

	addTypedData(chunkStore, levelBaseID+10+classID, compressed, table.Encode())
	chunkStore.Put(chunk.ID(levelBaseID+25+classID), mapChunk(compressed, make([]byte, entrySize)))
}

// AddSurveillanceChunk adds a chunk for surveillance information
func AddSurveillanceChunk(chunkStore chunk.Store, levelBaseID uint16) {
	chunkStore.Put(chunk.ID(levelBaseID+43), mapChunk(false, make([]byte, 8*2)))
	chunkStore.Put(chunk.ID(levelBaseID+44), mapChunk(false, make([]byte, 8*2)))
}

// AddLevelVariables adds a chunk for level variables.
func AddLevelVariables(chunkStore chunk.Store, levelBaseID uint16) {
	info := data.NewLevelVariables()

	addTypedData(chunkStore, levelBaseID+45, false, info)
}

// AddMapNotes prepares empty map notes chunks
func AddMapNotes(chunkStore chunk.Store, levelBaseID uint16) {
	chunkStore.Put(chunk.ID(levelBaseID+46), mapChunk(true, make([]byte, 0x0800)))
	chunkStore.Put(chunk.ID(levelBaseID+47), mapChunk(false, make([]byte, 4)))
}

// AddLoopConfiguration adds an empty loop configuration chunk
func AddLoopConfiguration(chunkStore chunk.Store, levelBaseID uint16) {
	chunkStore.Put(chunk.ID(levelBaseID+51), mapChunk(true, make([]byte, 0x03C0)))
}
