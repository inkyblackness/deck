package chunks

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/res/serial"
)

// AddLevel adds one level to the consumer
func AddLevel(consumer chunk.Consumer, levelID int, solid bool) {
	levelBaseID := res.ResourceID(4000 + 100*levelID)

	AddStaticChunk(consumer, levelBaseID+2, []byte{0x0B, 0x00, 0x00, 0x00})
	AddStaticChunk(consumer, levelBaseID+3, []byte{0x1B, 0x00, 0x00, 0x00})

	AddBasicLevelInformation(consumer, levelBaseID)
	AddMap(consumer, levelBaseID, solid, levelID == 1)
	AddStaticChunk(consumer, levelBaseID+6, make([]byte, 8))
	AddLevelTextures(consumer, levelBaseID)
	AddMasterObjectTables(consumer, levelBaseID)
	AddLevelObjects(consumer, levelBaseID)

	AddStaticChunk(consumer, levelBaseID+40, []byte{0x0D, 0x00, 0x00, 0x00})
	AddStaticChunk(consumer, levelBaseID+41, []byte{0x00})
	AddStaticChunk(consumer, levelBaseID+42, make([]byte, 0x1C))

	AddSurveillanceChunk(consumer, levelBaseID)
	AddLevelVariables(consumer, levelBaseID)
	AddMapNotes(consumer, levelBaseID)

	AddStaticChunk(consumer, levelBaseID+48, make([]byte, 0x30))
	AddStaticChunk(consumer, levelBaseID+49, make([]byte, 0x01C0))
	AddStaticChunk(consumer, levelBaseID+50, make([]byte, 2))

	AddLoopConfiguration(consumer, levelBaseID)

	// CD-Release only content
	AddStaticChunk(consumer, levelBaseID+52, make([]byte, 2))
	AddStaticChunk(consumer, levelBaseID+53, make([]byte, 0x40))
}

func addData(consumer chunk.Consumer, chunkID res.ResourceID, data interface{}) {
	addTypedData(consumer, chunkID, chunk.BasicChunkType, data)
}

func addTypedData(consumer chunk.Consumer, chunkID res.ResourceID, typeID chunk.TypeID, data interface{}) {
	store := serial.NewByteStore()
	coder := serial.NewPositioningEncoder(store)

	serial.MapData(data, coder)
	blocks := [][]byte{store.Data()}
	consumer.Consume(chunkID, chunk.NewBlockHolder(typeID, res.Map, blocks))
}

// AddBasicLevelInformation adds the basic level info block
func AddBasicLevelInformation(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	info := data.DefaultLevelInformation()

	addTypedData(consumer, levelBaseID+4, chunk.BasicChunkType.WithCompression(), info)
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
			master.Next = uint16((index + 1) % masterCount)
			master.Previous = uint16((masterCount + index - 1) % masterCount)
		}

		masterTable := &data.Table{Entries: make([]interface{}, masterCount)}
		for i := range masters {
			masterTable.Entries[i] = masters[i]
		}
		addTypedData(consumer, levelBaseID+8, chunk.BasicChunkType.WithCompression(), masterTable)
	}
	{
		refCount := 1600
		references := make([]*data.LevelObjectCrossReference, refCount)
		for index := range references {
			ref := data.DefaultLevelObjectCrossReference()
			references[index] = ref
			ref.NextObjectIndex = uint16((index + 1) % refCount)
		}

		refTable := &data.Table{Entries: make([]interface{}, refCount)}
		for i := range references {
			refTable.Entries[i] = references[i]
		}
		addTypedData(consumer, levelBaseID+9, chunk.BasicChunkType.WithCompression(), refTable)
	}
}

// AddLevelObjects adds level object tables
func AddLevelObjects(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	for class := 0; class < 15; class++ {
		meta := data.LevelObjectClassMetaEntry(res.ObjectClass(class))
		addLevelObjectTables(consumer, levelBaseID, class, meta.EntrySize, meta.EntryCount)
	}
}

type tempStruct struct {
	data.LevelObjectPrefix
	Extra []byte
}

func addLevelObjectTables(consumer chunk.Consumer, levelBaseID res.ResourceID, classID int, entrySize int, entryCount int) {
	table := data.Table{Entries: make([]interface{}, entryCount)}

	for i := range table.Entries {
		table.Entries[i] = &tempStruct{
			LevelObjectPrefix: data.LevelObjectPrefix{
				Next:                  uint16((i + 1) % entryCount),
				Previous:              uint16((entryCount + i - 1) % entryCount),
				LevelObjectTableIndex: 0},
			Extra: make([]byte, entrySize-data.LevelObjectPrefixSize)}
	}
	addData(consumer, levelBaseID+10+res.ResourceID(classID), table)
	//AddStaticChunk(consumer, levelBaseID+10+res.ResourceID(classID), make([]byte, entrySize*entryCount))
	AddStaticChunk(consumer, levelBaseID+25+res.ResourceID(classID), make([]byte, entrySize))
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
	AddStaticChunk(consumer, levelBaseID+46, make([]byte, 0x0800))
	AddStaticChunk(consumer, levelBaseID+47, make([]byte, 4))
}

// AddLoopConfiguration adds an empty loop configuration chunk
func AddLoopConfiguration(consumer chunk.Consumer, levelBaseID res.ResourceID) {
	AddStaticChunk(consumer, levelBaseID+51, make([]byte, 0x03C0))
}
