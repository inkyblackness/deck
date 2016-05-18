package console

import (
	"fmt"
	"io"

	"github.com/inkyblackness/shocked-client/viewmodel"
)

// ValueSelectionNodeTexter is a texter for value selection.
type ValueSelectionNodeTexter struct {
	listener ViewModelListener
	node     *viewmodel.ValueSelectionNode

	onDetailChanged DetailDataChangeCallback
}

// NewValueSelectionNodeTexter returns a new instance of ValueSelectionNodeTexter.
func NewValueSelectionNodeTexter(node *viewmodel.ValueSelectionNode, listener ViewModelListener) *ValueSelectionNodeTexter {
	texter := &ValueSelectionNodeTexter{
		listener:        listener,
		node:            node,
		onDetailChanged: NullDetailChangeCallback}

	node.Selected().Subscribe(func(string) {
		texter.onSelectedChanged()
	})
	node.Subscribe(func([]string) {
		texter.onDetailChanged()
	})

	return texter
}

// Act implements the ViewModelNodeTexter interface.
func (texter *ValueSelectionNodeTexter) Act(viewFactory NodeDetailViewFactory) {
	values := texter.node.Values()
	selected := texter.node.Selected().Get()
	selectedIndex := 0

	for index, value := range values {
		if value == selected {
			selectedIndex = index
		}
	}
	texter.onDetailChanged = viewFactory.ForList(texter, selectedIndex)
}

// TextMain implements the ViewModelNodeTexter interface.
func (texter *ValueSelectionNodeTexter) TextMain(addLine ViewModelLiner) {
	addLine(texter.node.Label(), fmt.Sprintf("v [%s]", texter.node.Selected().Get()), texter)
}

func (texter *ValueSelectionNodeTexter) onSelectedChanged() {
	texter.listener.OnMainDataChanged()
	texter.onDetailChanged()
}

// Cancel implements the DetailController interface.
func (texter *ValueSelectionNodeTexter) Cancel() {
	texter.onDetailChanged = NullDetailChangeCallback
}

// WriteDetails implements the DetailController interface.
func (texter *ValueSelectionNodeTexter) WriteDetails(w io.Writer) {
	for _, value := range texter.node.Values() {
		fmt.Fprintf(w, "%v\n", value)
	}
}

// Confirm implements the ListDetailController interface.
func (texter *ValueSelectionNodeTexter) Confirm(index int) {
	values := texter.node.Values()
	selectedValue := ""

	texter.onDetailChanged = NullDetailChangeCallback
	if (index >= 0) && (index < len(values)) {
		selectedValue = values[index]
	}
	texter.node.Selected().Set(selectedValue)
}
