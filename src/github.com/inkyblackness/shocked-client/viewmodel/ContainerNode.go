package viewmodel

// ContainerNode is a node holding other nodes by a string identifier.
type ContainerNode struct {
	label string
	nodes map[string]Node
}

// NewContainerNode returns a new instance of a ContainerNode.
func NewContainerNode(label string, nodes map[string]Node) *ContainerNode {
	node := &ContainerNode{
		label: label,
		nodes: nodes}

	return node
}

// Label is the Node interface implementation.
func (node *ContainerNode) Label() string {
	return node.label
}

// Specialize is the Node interface implementation.
func (node *ContainerNode) Specialize(visitor NodeVisitor) {
	visitor.Container(node)
}

// Get returns the contained nodes.
func (node *ContainerNode) Get() map[string]Node {
	return node.nodes
}
