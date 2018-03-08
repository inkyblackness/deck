package chunks

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/res/logic"
	"github.com/inkyblackness/res/serial"
)

// AddLevel adds one level to the consumer
func AddLevel(consumer chunk.Consumer, levelID int, solid bool, isCyberspace bool) {
	levelBaseID := res.ResourceID(4000 + 100*levelID)

	AddStaticChunk(consumer, levelBaseID+2, []byte{0x0B, 0x00, 0x00, 0x00})
	AddStaticChunk(consumer, levelBaseID+3, []byte{0x1B, 0x00, 0x00, 0x00})

	AddBasicLevelInformation(consumer, levelBaseID, isCyberspace)
	AddMap(consumer, levelBaseID, solid, levelID == 1)
	AddLevelTimer(consumer, levelBaseID)
	AddLevelTextures(consumer, levelBaseID)
	AddMasterObjectTables(consumer, levelBaseID)
	AddLevelObjects(consumer, levelBaseID)

	AddStaticChunk(consumer, levelBaseID+40, []byte{0x0D, 0x00, 0x00, 0x00})
	AddStaticChunk(consumer, levelBaseID+41, []byte{0x00})
	AddTypedStaticChunk(consumer, levelBaseID+42, chunk.BasicChunkType.WithCompression(), make([]byte, 0x1C))

	AddSurveillanceChunk(consumer, levelBaseID)
	AddLevelVariables(consumer, levelBaseID)
	AddMapNotes(consumer, levelBaseID)

	AddStaticChunk(consumer, levelBaseID+48, make([]byte, 0x30))
	AddTypedStaticChunk(consumer, levelBaseID+49, chunk.BasicChunkType.WithCompression(), make([]byte, 0x01C0))
	AddStaticChunk(consumer, levelBaseID+50, make([]byte, 2))

	AddLoopConfiguration(consumer, levelBaseID)

	// CD-Release only content
	AddStaticChunk(consumer, levelBaseID+52, make([]byte, 2))
	AddTypedStaticChunk(consumer, levelBaseID+53, chunk.BasicChunkType.WithCompression(), make([]byte, 0x40))
}

func addData(consumer chunk.Consumer, chunkID res.ResourceID, data interface{}) {
	addTypedData(consumer, chunkID, chunk.BasicChunkType, data)
}

func addTypedData(consumer chunk.Consumer, chunkID res.ResourceID, typeID chunk.TypeID, data interface{}) {
	store := serial.NewByteStore()
	coder := serial.NewEncoder(store)

	coder.Code(data)
	blocks := [][]byte{store.Data()}
	consumer.Consume(chunkID, chunk.NewBlockHolder(typeID, res.Map, blocks))
}

// AddBasicLevelInformation adds the basic level info block
func AddBasicLevelInformation(consumer chunk.Consumer, levelBaseID res.ResourceID, isCyberspace bool) {
	info := data.DefaultLevelInformation()

	if isCyberspace {
		info.CyberspaceFlag = 1
	}
	addTypedData(consumer, levelBaseID+4, chunk.BasicChunkType.WithCompression(), info)
}

// AddLevelTimer adds the basic timer list structure
func AddLevelTimer(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	AddStaticChunk(consumer, levelBaseID+6, make([]byte, data.TimerEntrySize))
}

// AddMap adds a map
func AddMap(consumer chunk.Consumer, levelBaseID res.ResourceID, solid bool, exceptStartingPosition bool) {
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
	addTypedData(consumer, levelBaseID+5, chunk.BasicChunkType.WithCompression(), table)
}

// AddLevelTextures adds level texture information
func AddLevelTextures(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	data := make([]byte, 54*2)
	data[0] = 101 // energ-light
	AddStaticChunk(consumer, levelBaseID+7, data)
}

// AddMasterObjectTables adds main object tables
func AddMasterObjectTables(consumer chunk.Consumer, levelBaseID res.ResourceID) {
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
		addTypedData(consumer, levelBaseID+8, chunk.BasicChunkType.WithCompression(), masterTable)
	}
	{
		crossrefList := logic.NewCrossReferenceList()

		crossrefList.Clear()
		addTypedData(consumer, levelBaseID+9, chunk.BasicChunkType.WithCompression(), crossrefList.Encode())
	}
}

// AddLevelObjects adds level object tables
func AddLevelObjects(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	for classID := 0; classID < 15; classID++ {
		meta := data.LevelObjectClassMetaEntry(res.ObjectClass(classID))
		addLevelObjectTables(consumer, levelBaseID, classID, meta.EntrySize, meta.EntryCount)
	}
}

func addLevelObjectTables(consumer chunk.Consumer, levelBaseID res.ResourceID, classID int, entrySize int, entryCount int) {
	requiresCompressed := classID != 0 && classID != 1 && classID != 4 && classID != 5 && classID != 6
	chunkType := chunk.BasicChunkType
	table := logic.NewLevelObjectClassTable(entrySize, entryCount)
	table.AsChain().Initialize(entryCount - 1)

	if requiresCompressed {
		chunkType = chunkType.WithCompression()
	}

	addTypedData(consumer, levelBaseID+10+res.ResourceID(classID), chunkType, table.Encode())
	AddTypedStaticChunk(consumer, levelBaseID+25+res.ResourceID(classID), chunkType, make([]byte, entrySize))
}

// AddSurveillanceChunk adds a chunk for surveillance information
func AddSurveillanceChunk(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	AddStaticChunk(consumer, levelBaseID+43, make([]byte, 8*2))
	AddStaticChunk(consumer, levelBaseID+44, make([]byte, 8*2))
}

// AddLevelVariables adds a chunk for level variables.
func AddLevelVariables(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	info := data.NewLevelVariables()

	addData(consumer, levelBaseID+45, info)
}

// AddMapNotes prepares empty map notes chunks
func AddMapNotes(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	AddTypedStaticChunk(consumer, levelBaseID+46, chunk.BasicChunkType.WithCompression(), make([]byte, 0x0800))
	AddStaticChunk(consumer, levelBaseID+47, make([]byte, 4))
}

// AddLoopConfiguration adds an empty loop configuration chunk
func AddLoopConfiguration(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	AddTypedStaticChunk(consumer, levelBaseID+51, chunk.BasicChunkType.WithCompression(), make([]byte, 0x03C0))
}
