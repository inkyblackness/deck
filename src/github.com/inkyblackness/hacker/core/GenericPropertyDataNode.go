package core

type genericPropertyDataNode struct {
	rawDataNode
}

func newGenericPropertyDataNode(parentNode DataNode, data []byte) *genericPropertyDataNode {
	node := &genericPropertyDataNode{rawDataNode{
		parentNode: parentNode,
		id:         "generic",
		data:       data}}

	return node
}

func (node *genericPropertyDataNode) Info() string {
	info := ""

	return info
}
