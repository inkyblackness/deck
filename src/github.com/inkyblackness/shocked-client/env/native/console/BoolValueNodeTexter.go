package console

import (
	"fmt"

	"github.com/inkyblackness/shocked-client/viewmodel"
)

// BoolValueNodeTexter is a texter for boolean values.
type BoolValueNodeTexter struct {
	node *viewmodel.BoolValueNode
}

// NewBoolValueNodeTexter returns a new instance of BoolValueNodeTexter.
func NewBoolValueNodeTexter(node *viewmodel.BoolValueNode, listener ViewModelListener) *BoolValueNodeTexter {
	texter := &BoolValueNodeTexter{node: node}

	node.Subscribe(func(bool) {
		listener.OnMainDataChanged()
	})

	return texter
}

// Act implements the ViewModelNodeTexter interface.
func (texter *BoolValueNodeTexter) Act(viewFactory NodeDetailViewFactory) {
	texter.node.Set(!texter.node.Get())
}

// TextMain implements the ViewModelNodeTexter interface.
func (texter *BoolValueNodeTexter) TextMain(addLine ViewModelLiner) {
	state := map[bool]string{false: "_", true: "X"}
	addLine(texter.node.Label(), fmt.Sprintf("%v", state[texter.node.Get()]), texter)
}
