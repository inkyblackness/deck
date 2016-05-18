package viewmodel

// SectionNode is a node holding other nodes.
type SectionNode struct {
	label     string
	available *BoolValueNode
	nodes     []Node
}

// NewSectionNode returns a new instance of a SectionNode.
func NewSectionNode(label string, nodes []Node, available *BoolValueNode) *SectionNode {
	node := &SectionNode{
		label:     label,
		available: available,
		nodes:     nodes}

	return node
}

// Label is the Node interface implementation.
func (node *SectionNode) Label() string {
	return node.label
}

// Specialize is the Node interface implementation.
func (node *SectionNode) Specialize(visitor NodeVisitor) {
	visitor.Section(node)
}

// Available returns the node about the section's availability.
func (node *SectionNode) Available() *BoolValueNode {
	return node.available
}

// Get returns the contained nodes.
func (node *SectionNode) Get() []Node {
	return node.nodes[:]
}
