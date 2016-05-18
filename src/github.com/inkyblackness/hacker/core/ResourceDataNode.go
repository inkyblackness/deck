package core

import (
	"strings"

	"github.com/inkyblackness/res/chunk"
)

type resourceDataNode struct {
	parentDataNode

	fileName        string
	consumerFactory func() chunk.Consumer
}

func NewResourceDataNode(parentNode DataNode, name string,
	provider chunk.Provider, consumerFactory func() chunk.Consumer) DataNode {
	ids := provider.IDs()
	node := &resourceDataNode{
		parentDataNode:  makeParentDataNode(parentNode, strings.ToLower(name), len(ids)),
		fileName:        name,
		consumerFactory: consumerFactory}

	for _, id := range ids {
		node.addChild(newChunkDataNode(node, id, provider.Provide(id)))
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
	consumer := node.consumerFactory()
	defer consumer.Finish()

	for _, child := range node.Children() {
		chunkNode := child.(*chunkDataNode)
		chunkNode.saveTo(consumer)
	}

	return node.fileName + "\n"
}
