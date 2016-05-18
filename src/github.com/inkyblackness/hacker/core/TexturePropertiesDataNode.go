package core

import (
	"fmt"
	"strings"

	"github.com/inkyblackness/res/textprop"
)

type textpropConsumerFactory func() textprop.Consumer

type texturePropertiesDataNode struct {
	parentDataNode
	consumerFactory textpropConsumerFactory
}

func NewTexturePropertiesDataNode(parentNode DataNode, name string,
	provider textprop.Provider, consumerFactory textpropConsumerFactory) DataNode {
	node := &texturePropertiesDataNode{
		parentDataNode:  makeParentDataNode(parentNode, strings.ToLower(name), int(provider.EntryCount())),
		consumerFactory: consumerFactory}

	for i := uint32(0); i < provider.EntryCount(); i++ {
		blockData := provider.Provide(i)
		dataStruct := &textprop.Entry{}
		node.addChild(newBlockDataNode(node, uint16(i), blockData, dataStruct))
	}

	return node
}

func (node *texturePropertiesDataNode) Info() string {
	info := fmt.Sprintf("Textures available: %d", len(node.Children()))

	return info
}

func (node *texturePropertiesDataNode) save() string {
	consumer := node.consumerFactory()
	defer consumer.Finish()

	for index, child := range node.Children() {
		consumer.Consume(uint32(index), child.Data())
	}

	return node.ID() + "\n"
}
