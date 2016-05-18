package core

import (
	"fmt"
	"strings"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/objprop"
)

type objectPropertiesDataNode struct {
	parentDataNode
	consumerFactory func() objprop.Consumer
}

func NewObjectPropertiesDataNode(parentNode DataNode, name string,
	provider objprop.Provider, classes []objprop.ClassDescriptor, consumerFactory func() objprop.Consumer) DataNode {
	node := &objectPropertiesDataNode{
		parentDataNode:  makeParentDataNode(parentNode, strings.ToLower(name), 0),
		consumerFactory: consumerFactory}

	for classIndex, classDesc := range classes {
		for subclassIndex, subclassDesc := range classDesc.Subclasses {
			for typeIndex := uint32(0); typeIndex < subclassDesc.TypeCount; typeIndex++ {
				id := res.MakeObjectID(res.ObjectClass(classIndex), res.ObjectSubclass(subclassIndex), res.ObjectType(typeIndex))
				subnode := newObjectPropertyDataNode(node, id, provider)
				node.addChild(subnode)
			}
		}
	}

	return node
}

func (node *objectPropertiesDataNode) Info() string {
	info := fmt.Sprintf("Objects available: %d", len(node.Children()))

	return info
}

func (node *objectPropertiesDataNode) save() string {
	consumer := node.consumerFactory()
	defer consumer.Finish()

	for _, child := range node.Children() {
		objNode := child.(*objectPropertyDataNode)
		objData := objprop.ObjectData{
			Common:   child.Resolve("common").Data(),
			Generic:  child.Resolve("generic").Data(),
			Specific: child.Resolve("specific").Data()}

		consumer.Consume(objNode.objectID, objData)
	}

	return node.ID() + "\n"
}
