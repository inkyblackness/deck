package viewmodel

// ValueSelectionValuesListener is a callback to be called for changes of the value list.
type ValueSelectionValuesListener func(newValues []string)

// ValueSelectionNode is a node selecting a string out of a list.
type ValueSelectionNode struct {
	label    string
	selected *StringValueNode

	valuesListeners []ValueSelectionValuesListener
	values          []string
}

// NewValueSelectionNode returns a new instance of a ValueSelectionNode.
func NewValueSelectionNode(label string, values []string, selected string) *ValueSelectionNode {
	node := &ValueSelectionNode{
		label:    label,
		selected: NewStringValueNode("Selected", selected),
		values:   values}

	return node
}

// Label is the Node interface implementation.
func (node *ValueSelectionNode) Label() string {
	return node.label
}

// Specialize is the Node interface implementation.
func (node *ValueSelectionNode) Specialize(visitor NodeVisitor) {
	visitor.ValueSelection(node)
}

// Subscribe registers the provided callback on changes of the values list.
func (node *ValueSelectionNode) Subscribe(listener ValueSelectionValuesListener) {
	node.valuesListeners = append(node.valuesListeners, listener)
}

// Selected returns the node for the current selection.
func (node *ValueSelectionNode) Selected() *StringValueNode {
	return node.selected
}

// Values returns the list of available values.
func (node *ValueSelectionNode) Values() []string {
	return node.values[:]
}

// SetValues sets the new array of possible values.
func (node *ValueSelectionNode) SetValues(values []string) {
	newValues := values[:]
	node.values = newValues
	selected := node.selected.Get()
	found := false
	for _, value := range newValues {
		if value == selected {
			found = true
		}
	}
	if !found {
		node.selected.Set("")
	}
	for _, listener := range node.valuesListeners {
		listener(newValues)
	}
}
