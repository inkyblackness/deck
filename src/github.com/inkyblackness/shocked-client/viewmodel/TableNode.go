package viewmodel

// TableListener is the callback type for changes in an TableNode.
type TableListener func(newRows []*ContainerNode)

// TableNode is a node holding an Table of nodes.
type TableNode struct {
	label     string
	listeners []TableListener
	rows      []*ContainerNode
}

// NewTableNode returns a new instance of an TableNode.
func NewTableNode(label string, rows ...*ContainerNode) *TableNode {
	node := &TableNode{
		label: label,
		rows:  rows}

	return node
}

// Label is the Node interface implementation.
func (node *TableNode) Label() string {
	return node.label
}

// Specialize is the Node interface implementation.
func (node *TableNode) Specialize(visitor NodeVisitor) {
	visitor.Table(node)
}

// Subscribe registers the provided listener for table changes.
func (node *TableNode) Subscribe(listener TableListener) {
	node.listeners = append(node.listeners, listener)
}

// Get returns the current rows.
func (node *TableNode) Get() []*ContainerNode {
	return node.rows[:]
}

// Set changes the current rows
func (node *TableNode) Set(rows []*ContainerNode) {
	newRows := rows[:]

	node.rows = newRows
	for _, listener := range node.listeners {
		listener(newRows)
	}
}
