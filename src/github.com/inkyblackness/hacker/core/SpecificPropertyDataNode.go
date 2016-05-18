package core

type specificPropertyDataNode struct {
	rawDataNode
}

func newSpecificPropertyDataNode(parentNode DataNode, data []byte) *specificPropertyDataNode {
	node := &specificPropertyDataNode{rawDataNode{
		parentNode: parentNode,
		id:         "specific",
		data:       data}}

	return node
}

func (node *specificPropertyDataNode) Info() string {
	info := ""

	return info
}
