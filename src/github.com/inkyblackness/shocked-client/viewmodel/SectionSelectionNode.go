package viewmodel

// SectionSelectionNode is a node selecting section nodes using a string identifier.
type SectionSelectionNode struct {
	label    string
	selected *ValueSelectionNode
	sections map[string]*SectionNode
}

// NewSectionSelectionNode returns a new instance of a SectionSelectionNode.
func NewSectionSelectionNode(label string, sections map[string]*SectionNode, selected string) *SectionSelectionNode {
	node := &SectionSelectionNode{
		label:    label,
		sections: sections}
	node.selected = NewValueSelectionNode(label, node.availableSectionKeys(), selected)
	for _, section := range sections {
		section.Available().Subscribe(func(bool) { node.selected.SetValues(node.availableSectionKeys()) })
	}

	return node
}

// Label is the Node interface implementation.
func (node *SectionSelectionNode) Label() string {
	return node.label
}

// Specialize is the Node interface implementation.
func (node *SectionSelectionNode) Specialize(visitor NodeVisitor) {
	visitor.SectionSelection(node)
}

// Sections returns the map of all contained sections.
func (node *SectionSelectionNode) Sections() map[string]*SectionNode {
	return node.sections
}

// Selection returns the node for the current selection.
func (node *SectionSelectionNode) Selection() *ValueSelectionNode {
	return node.selected
}

func (node *SectionSelectionNode) availableSectionKeys() []string {
	var values []string

	for key, section := range node.sections {
		if section.Available().Get() {
			values = append(values, key)
		}
	}

	return values
}
