package core

import (
	"fmt"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/objprop"
)

type objectPropertyDataNode struct {
	parentDataNode
	objectID res.ObjectID
}

func newObjectPropertyDataNode(parentNode DataNode, id res.ObjectID, provider objprop.Provider) *objectPropertyDataNode {
	idString := fmt.Sprintf("%d-%d-%d", id.Class, id.Subclass, id.Type)
	node := &objectPropertyDataNode{
		parentDataNode: makeParentDataNode(parentNode, idString, 3),
		objectID:       id}

	objData := provider.Provide(id)
	node.addChild(newGenericPropertyDataNode(node, objData.Generic))
	node.addChild(newSpecificPropertyDataNode(node, objData.Specific))
	node.addChild(newCommonPropertyDataNode(node, objData.Common))

	return node
}

func (node *objectPropertyDataNode) Info() string {
	info := ""

	return info
}
