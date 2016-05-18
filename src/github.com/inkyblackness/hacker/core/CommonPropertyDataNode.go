package core

type commonPropertyDataNode struct {
	rawDataNode
}

func newCommonPropertyDataNode(parentNode DataNode, data []byte) *commonPropertyDataNode {
	node := &commonPropertyDataNode{rawDataNode{
		parentNode: parentNode,
		id:         "common",
		data:       data}}

	return node
}

func (node *commonPropertyDataNode) Info() string {
	info := ""

	return info
}
