package core

import (
	"strings"

	"github.com/inkyblackness/res/chunk"
)

type resourceDataNode struct {
	parentDataNode

	fileName string
	saver    func(func(chunk.Store))
}

func NewResourceDataNode(parentNode DataNode, name string,
	provider chunk.Provider, saver func(func(chunk.Store))) DataNode {
	ids := provider.IDs()
	node := &resourceDataNode{
		parentDataNode: makeParentDataNode(parentNode, strings.ToLower(name), len(ids)),
		fileName:       name,
		saver:          saver}

	addChunk := func(id chunk.Identifier) {
		containedChunk, chunkErr := provider.Chunk(id)
		if chunkErr != nil {
			return
		}
		node.addChild(newChunkDataNode(node, id, containedChunk))
	}
	for _, id := range ids {
		addChunk(id)
	}

	return node
}

func (node *resourceDataNode) Info() string {
	info := "ResourceFile: " + node.fileName + "\n"
	info += "IDs:"
	for _, node := range node.Children() {
		info += " " + node.ID()
	}

	return info
}

func (node *resourceDataNode) save() string {
	node.saver(func(target chunk.Store) {
		for _, child := range node.Children() {
			chunkNode := child.(*chunkDataNode)
			chunkNode.saveTo(target)
		}
	})

	return node.fileName + "\n"
}
