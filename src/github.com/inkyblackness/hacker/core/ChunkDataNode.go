package core

import (
	"bytes"
	"fmt"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/res/image"
	moviFormat "github.com/inkyblackness/res/movi/format"
)

type chunkDataNode struct {
	parentDataNode

	chunkID res.ResourceID
	holder  chunk.BlockHolder
}

func isLevelChunk(chunkID res.ResourceID, relativeID int) bool {
	result := false

	if (chunkID >= res.ResourceID(4000)) && (int(chunkID)%100) == relativeID {
		result = true
	}

	return result
}

func newChunkDataNode(parentNode DataNode, chunkID res.ResourceID, holder chunk.BlockHolder) *chunkDataNode {
	node := &chunkDataNode{
		parentDataNode: makeParentDataNode(parentNode, fmt.Sprintf("%04X", chunkID), int(holder.BlockCount())),
		chunkID:        chunkID,
		holder:         holder}

	for i := uint16(0); i < holder.BlockCount(); i++ {
		blockData := holder.BlockData(i)
		table := getTableForBlock(chunkID, blockData)

		if table != nil {
			node.addChild(newTableDataNode(node, fmt.Sprintf("%d", i), blockData, table))
		} else {
			dataStruct := getDataStructForBlock(holder.ContentType(), chunkID, blockData)
			node.addChild(newBlockDataNode(node, i, blockData, dataStruct))
		}

	}

	return node
}

func getTableForBlock(chunkID res.ResourceID, blockData []byte) (table Table) {
	if isLevelChunk(chunkID, 5) {
		entryCount := 64 * 64
		table = data.NewTable(entryCount, func() interface{} { return data.DefaultTileMapEntry() })
	} else if isLevelChunk(chunkID, 8) {
		entryCount := len(blockData) / data.LevelObjectEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.DefaultLevelObjectEntry() })
	} else if isLevelChunk(chunkID, 9) {
		entryCount := len(blockData) / data.LevelObjectCrossReferenceSize
		table = data.NewTable(entryCount, func() interface{} { return data.DefaultLevelObjectCrossReference() })
	} else if isLevelChunk(chunkID, 10) {
		entryCount := len(blockData) / data.LevelWeaponEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelWeaponEntry() })
	} else if isLevelChunk(chunkID, 11) {
		entryCount := len(blockData) / data.LevelAmmoEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelAmmoEntry() })
	} else if isLevelChunk(chunkID, 12) {
		entryCount := len(blockData) / data.LevelProjectileEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelProjectileEntry() })
	} else if isLevelChunk(chunkID, 13) {
		entryCount := len(blockData) / data.LevelExplosiveEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelExplosiveEntry() })
	} else if isLevelChunk(chunkID, 14) {
		entryCount := len(blockData) / data.LevelPatchEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelPatchEntry() })
	} else if isLevelChunk(chunkID, 15) {
		entryCount := len(blockData) / data.LevelHardwareEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelHardwareEntry() })
	} else if isLevelChunk(chunkID, 16) {
		entryCount := len(blockData) / data.LevelSoftwareEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelSoftwareEntry() })
	} else if isLevelChunk(chunkID, 17) {
		entryCount := len(blockData) / data.LevelSceneryEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelSceneryEntry() })
	} else if isLevelChunk(chunkID, 18) {
		entryCount := len(blockData) / data.LevelItemEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelItemEntry() })
	} else if isLevelChunk(chunkID, 19) {
		entryCount := len(blockData) / data.LevelPanelEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelPanelEntry() })
	} else if isLevelChunk(chunkID, 20) {
		entryCount := len(blockData) / data.LevelBarrierEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelBarrierEntry() })
	} else if isLevelChunk(chunkID, 21) {
		entryCount := len(blockData) / data.LevelAnimationEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelAnimationEntry() })
	} else if isLevelChunk(chunkID, 22) {
		entryCount := len(blockData) / data.LevelMarkerEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelMarkerEntry() })
	} else if isLevelChunk(chunkID, 23) {
		entryCount := len(blockData) / data.LevelContainerEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelContainerEntry() })
	} else if isLevelChunk(chunkID, 24) {
		entryCount := len(blockData) / data.LevelCritterEntrySize
		table = data.NewTable(entryCount, func() interface{} { return data.NewLevelCritterEntry() })
	}

	return
}

func getDataStructForBlock(dataType res.DataTypeID, chunkID res.ResourceID, blockData []byte) (dataStruct interface{}) {
	if chunkID == res.ResourceID(0x0FA0) {
		dataStruct = data.NewString(bytes.IndexByte(blockData, 0x00) + 1)
	} else if chunkID == res.ResourceID(0x0FA1) {
		dataStruct = data.DefaultGameState()
	} else if isLevelChunk(chunkID, 4) {
		dataStruct = data.DefaultLevelInformation()
	} else if isLevelChunk(chunkID, 45) {
		dataStruct = data.NewLevelVariables()
	} else if dataType == res.Media {
		dataStruct = &moviFormat.Header{}
	} else if dataType == res.VideoClip {
		dataStruct = data.DefaultVideoClipSequence((len(blockData) - data.VideoClipSequenceBaseSize) / data.VideoClipSequenceEntrySize)
	} else if dataType == res.Bitmap {
		dataStruct = &image.BitmapHeader{}
	}

	return
}

func (node *chunkDataNode) Info() (info string) {
	info += fmt.Sprintf("Content type: 0x%02X\n", node.holder.ContentType())
	info += fmt.Sprintf("Available blocks: %d\n", node.holder.BlockCount())
	info += fmt.Sprintf("Chunk TypeID: 0x%02X (%v)\n", int(node.holder.ChunkType()), node.holder.ChunkType())

	return info
}

func (node *chunkDataNode) saveTo(consumer chunk.Consumer) {
	blockNodes := node.Children()
	blocks := make([][]byte, len(blockNodes))
	for index, blockNode := range blockNodes {
		blocks[index] = blockNode.Data()
	}
	newHolder := chunk.NewBlockHolder(node.holder.ChunkType(), node.holder.ContentType(), blocks)
	consumer.Consume(node.chunkID, newHolder)
}
