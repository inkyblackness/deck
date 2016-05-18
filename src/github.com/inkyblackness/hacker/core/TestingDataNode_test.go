package core

type TestingDataNode struct {
	parentDataNode
	data []byte
}

func NewTestingDataNode(id string) *TestingDataNode {
	node := &TestingDataNode{parentDataNode: makeParentDataNode(nil, id, 0)}

	return node
}

// Info returns human readable information about this node.
func (node *TestingDataNode) Info() string {
	return "testing"
}

func (node *TestingDataNode) Data() []byte {
	return node.data
}

func (node *TestingDataNode) UnknownData() []byte {
	return node.data
}
