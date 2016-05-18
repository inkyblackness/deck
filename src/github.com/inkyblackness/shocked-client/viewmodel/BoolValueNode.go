package viewmodel

// BoolValueListener is the callback type for changes in a BoolValueNode.
type BoolValueListener func(newValue bool)

// BoolValueNode is a node holding a simple bool value.
type BoolValueNode struct {
	label     string
	listeners []BoolValueListener
	value     bool
}

// NewBoolValueNode returns a new instance of a BoolValueNode.
func NewBoolValueNode(label string, value bool) *BoolValueNode {
	node := &BoolValueNode{
		label: label,
		value: value}

	return node
}

// Label is the Node interface implementation.
func (node *BoolValueNode) Label() string {
	return node.label
}

// Specialize is the Node interface implementation.
func (node *BoolValueNode) Specialize(visitor NodeVisitor) {
	visitor.BoolValue(node)
}

// Subscribe registers the provided listener for value changes.
func (node *BoolValueNode) Subscribe(listener BoolValueListener) {
	node.listeners = append(node.listeners, listener)
}

// Get returns the current value.
func (node *BoolValueNode) Get() bool {
	return node.value
}

// Set requests to set a new value.
func (node *BoolValueNode) Set(value bool) {
	node.value = value
	for _, listener := range node.listeners {
		listener(value)
	}
}
