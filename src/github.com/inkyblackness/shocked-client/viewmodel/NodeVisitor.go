package viewmodel

// NodeVisitor will be called with appropriate node instances while
// walking through a tree in the view model.
type NodeVisitor interface {
	// Section will be called for any SectionNode.
	Section(node *SectionNode)
	// SectionSelection will be called for any SectionSelectionNode.
	SectionSelection(node *SectionSelectionNode)

	// ValueSelection will be called for any ValueSelectionNode.
	ValueSelection(node *ValueSelectionNode)

	// BoolValue will be called for any BoolValueNode.
	BoolValue(node *BoolValueNode)
	// StringValue will be called for any StringValueNode.
	StringValue(node *StringValueNode)
	// Container will be called for any ContainerNode.
	Container(node *ContainerNode)
	// Table will be called for any TableNode.
	Table(node *TableNode)

	// Action will be called for any ActionNode.
	Action(node *ActionNode)
}
