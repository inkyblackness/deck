package console

import (
	"fmt"

	"github.com/inkyblackness/shocked-client/viewmodel"
)

// ActionNodeTexter is a texter for action nodes.
type ActionNodeTexter struct {
	node *viewmodel.ActionNode
}

// NewActionNodeTexter returns a new instance of ActionNodeTexter.
func NewActionNodeTexter(node *viewmodel.ActionNode, listener ViewModelListener) *ActionNodeTexter {
	texter := &ActionNodeTexter{node: node}

	return texter
}

// Act implements the ViewModelNodeTexter interface.
func (texter *ActionNodeTexter) Act(viewFactory NodeDetailViewFactory) {
	texter.node.Act()
}

// TextMain implements the ViewModelNodeTexter interface.
func (texter *ActionNodeTexter) TextMain(addLine ViewModelLiner) {
	addLine(texter.node.Label(), fmt.Sprintf("  <Execute>"), texter)
}
