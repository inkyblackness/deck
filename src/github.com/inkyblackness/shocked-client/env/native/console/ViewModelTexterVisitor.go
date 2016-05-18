package console

import (
	"github.com/inkyblackness/shocked-client/viewmodel"
)

// ViewModelTexterVisitor implements the viewmodel.NodeVisitor interface to create
// ViewModelNodeTexter instances
type ViewModelTexterVisitor struct {
	listener ViewModelListener
	instance ViewModelNodeTexter
}

// NewViewModelTexterVisitor returns a new instance of ViewModelTexterVisitor.
func NewViewModelTexterVisitor(listener ViewModelListener) *ViewModelTexterVisitor {
	return &ViewModelTexterVisitor{listener: listener}
}

// Section is the viewmodel.NodeVisitor implementation.
func (visitor *ViewModelTexterVisitor) Section(node *viewmodel.SectionNode) {
	visitor.instance = NewSectionNodeTexter(node, visitor.listener)
}

// SectionSelection is the viewmodel.NodeVisitor implementation.
func (visitor *ViewModelTexterVisitor) SectionSelection(node *viewmodel.SectionSelectionNode) {
	visitor.instance = NewSectionSelectionNodeTexter(node, visitor.listener)
}

// ValueSelection is the viewmodel.NodeVisitor implementation.
func (visitor *ViewModelTexterVisitor) ValueSelection(node *viewmodel.ValueSelectionNode) {
	visitor.instance = NewValueSelectionNodeTexter(node, visitor.listener)
}

// BoolValue is the viewmodel.NodeVisitor implementation.
func (visitor *ViewModelTexterVisitor) BoolValue(node *viewmodel.BoolValueNode) {
	visitor.instance = NewBoolValueNodeTexter(node, visitor.listener)
}

// StringValue is the viewmodel.NodeVisitor implementation.
func (visitor *ViewModelTexterVisitor) StringValue(node *viewmodel.StringValueNode) {
	visitor.instance = NewStringValueNodeTexter(node, visitor.listener)
}

// Container is the viewmodel.NodeVisitor implementation.
func (visitor *ViewModelTexterVisitor) Container(node *viewmodel.ContainerNode) {
}

// Table is the viewmodel.NodeVisitor implementation.
func (visitor *ViewModelTexterVisitor) Table(node *viewmodel.TableNode) {
}

// Action is the viewmodel.NodeVisitor implementation.
func (visitor *ViewModelTexterVisitor) Action(node *viewmodel.ActionNode) {
	visitor.instance = NewActionNodeTexter(node, visitor.listener)
}
