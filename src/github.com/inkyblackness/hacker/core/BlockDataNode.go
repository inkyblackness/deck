package core

import (
	"bytes"
	"fmt"

	"github.com/inkyblackness/res/serial"
)

type blockDataNode struct {
	rawDataNode

	dataStruct interface{}
}

func newBlockDataNode(parentNode DataNode, blockIndex uint16, data []byte, dataStruct interface{}) *blockDataNode {
	node := &blockDataNode{
		rawDataNode: rawDataNode{
			parentNode: parentNode,
			id:         fmt.Sprintf("%d", blockIndex),
			data:       data},
		dataStruct: dataStruct}

	return node
}

func (node *blockDataNode) Info() string {
	info := ""
	if node.dataStruct != nil {
		serial.NewDecoder(bytes.NewReader(node.Data())).Code(node.dataStruct)
		info = fmt.Sprintf("%v", node.dataStruct)
	}

	return info
}

func (node *blockDataNode) UnknownData() []byte {
	originalData := node.Data()
	maskedData := originalData

	/*
		if node.dataStruct != nil {
			maskedData = make([]byte, len(originalData))
			copy(maskedData, originalData)
		}
	*/

	return maskedData
}
