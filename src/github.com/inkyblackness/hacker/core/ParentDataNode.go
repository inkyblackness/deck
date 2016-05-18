package core

type parentDataNode struct {
	parentNode DataNode
	id         string

	children     []DataNode
	childrenByID map[string]DataNode

	childResolver func(string) DataNode
}

func makeParentDataNode(parentNode DataNode, id string, childrenCount int) parentDataNode {
	node := parentDataNode{
		parentNode:    parentNode,
		id:            id,
		children:      make([]DataNode, 0, childrenCount),
		childrenByID:  make(map[string]DataNode),
		childResolver: func(string) DataNode { return nil }}

	return node
}

func (node *parentDataNode) addChild(childNode DataNode) {
	node.children = append(node.children, childNode)
	node.childrenByID[childNode.ID()] = childNode
}

func (node *parentDataNode) setChildResolver(resolver func(string) DataNode) {
	node.childResolver = resolver
}

func (node *parentDataNode) Parent() DataNode {
	return node.parentNode
}

func (node *parentDataNode) Children() []DataNode {
	return node.children
}

func (node *parentDataNode) ID() string {
	return node.id
}

func (node *parentDataNode) Resolve(path string) (resolved DataNode) {
	temp, existing := node.childrenByID[path]

	if existing {
		resolved = temp
	} else {
		temp = node.childResolver(path)
		if temp != nil {
			node.addChild(temp)
			resolved = temp
		}
	}

	return
}

func (node *parentDataNode) Data() []byte {
	return nil
}

func (node *parentDataNode) UnknownData() []byte {
	return nil
}
