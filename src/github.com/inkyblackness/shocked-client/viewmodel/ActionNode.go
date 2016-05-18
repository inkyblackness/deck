package viewmodel

// ActionListener is the callback type for acting on an ActionNode.
type ActionListener func()

// ActionNode is a node for executing an action.
type ActionNode struct {
	label     string
	listeners []ActionListener
}

// NewActionNode returns a new instance of a ActionNode.
func NewActionNode(label string) *ActionNode {
	node := &ActionNode{
		label: label}

	return node
}

// Label is the Node interface implementation.
func (node *ActionNode) Label() string {
	return node.label
}

// Specialize is the Node interface implementation.
func (node *ActionNode) Specialize(visitor NodeVisitor) {
	visitor.Action(node)
}

// Subscribe registers the provided listener for actions.
func (node *ActionNode) Subscribe(listener ActionListener) {
	node.listeners = append(node.listeners, listener)
}

// Act fires all registered listeners.
func (node *ActionNode) Act() {
	for _, listener := range node.listeners {
		listener()
	}
}
