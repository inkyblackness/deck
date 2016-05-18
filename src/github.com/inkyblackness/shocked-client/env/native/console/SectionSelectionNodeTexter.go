package console

import (
	"github.com/inkyblackness/shocked-client/viewmodel"
)

// SectionSelectionNodeTexter is a texter for sections.
type SectionSelectionNodeTexter struct {
	node *viewmodel.SectionSelectionNode

	selectionTexter ViewModelNodeTexter
	subTexter       map[string]ViewModelNodeTexter
}

// NewSectionSelectionNodeTexter returns a new instance of SectionSelectionNodeTexter.
func NewSectionSelectionNodeTexter(node *viewmodel.SectionSelectionNode, listener ViewModelListener) *SectionSelectionNodeTexter {
	texter := &SectionSelectionNodeTexter{
		selectionTexter: NewValueSelectionNodeTexter(node.Selection(), listener),
		subTexter:       make(map[string]ViewModelNodeTexter),
		node:            node}

	sections := node.Sections()
	for key, section := range sections {
		visitor := NewViewModelTexterVisitor(listener)
		section.Specialize(visitor)
		texter.subTexter[key] = visitor.instance
	}

	return texter
}

// Act implements the ViewModelNodeTexter interface.
func (texter *SectionSelectionNodeTexter) Act(viewFactory NodeDetailViewFactory) {
}

// TextMain implements the ViewModelNodeTexter interface.
func (texter *SectionSelectionNodeTexter) TextMain(addLine ViewModelLiner) {
	texter.selectionTexter.TextMain(addLine)
	selectedTexter, _ := texter.subTexter[texter.node.Selection().Selected().Get()]

	if selectedTexter != nil {
		selectedTexter.TextMain(addLine)
	}
}
