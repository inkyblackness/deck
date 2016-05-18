package core

type rawDataNode struct {
	parentNode DataNode
	id         string

	data []byte
}

func (node *rawDataNode) Parent() DataNode {
	return node.parentNode
}

func (node *rawDataNode) Children() []DataNode {
	return nil
}

func (node *rawDataNode) ID() string {
	return node.id
}

func (node *rawDataNode) Resolve(path string) DataNode {
	return nil
}

func (node *rawDataNode) Data() []byte {
	return node.data
}

func (node *rawDataNode) UnknownData() []byte {
	return node.Data()
}
