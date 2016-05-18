package console

import (
	"fmt"
	"io"
	"strings"

	"github.com/inkyblackness/shocked-client/viewmodel"
)

// StringValueNodeTexter is a texter for string values.
type StringValueNodeTexter struct {
	node *viewmodel.StringValueNode

	onDetailChanged DetailDataChangeCallback
}

// NewStringValueNodeTexter returns a new instance of StringValueNodeTexter.
func NewStringValueNodeTexter(node *viewmodel.StringValueNode, listener ViewModelListener) *StringValueNodeTexter {
	texter := &StringValueNodeTexter{
		node:            node,
		onDetailChanged: NullDetailChangeCallback}

	node.Subscribe(func(string) {
		listener.OnMainDataChanged()
		texter.onDetailChanged()
	})

	return texter
}

// Act implements the ViewModelNodeTexter interface.
func (texter *StringValueNodeTexter) Act(viewFactory NodeDetailViewFactory) {
	if texter.node.Editable() {
		texter.onDetailChanged = viewFactory.ForString(texter)
	}
}

// TextMain implements the ViewModelNodeTexter interface.
func (texter *StringValueNodeTexter) TextMain(addLine ViewModelLiner) {
	printText := texter.node.Get()
	prefix := "="

	if texter.node.Editable() {
		prefix = "?"
	}
	printText = strings.Replace(printText, "\n", "\\n", -1)
	printText = strings.Replace(printText, "\r", "\\r", -1)
	addLine(texter.node.Label(), fmt.Sprintf("%s [%s]", prefix, printText), texter)
}

// Cancel implements the DetailController interface.
func (texter *StringValueNodeTexter) Cancel() {
	texter.onDetailChanged = NullDetailChangeCallback
}

// WriteDetails implements the DetailController interface.
func (texter *StringValueNodeTexter) WriteDetails(w io.Writer) {
	fmt.Fprintf(w, "%s", texter.node.Get())
}

// Confirm implements the StringDetailController interface.
func (texter *StringValueNodeTexter) Confirm(newValue string) {
	texter.onDetailChanged = NullDetailChangeCallback
	texter.node.Set(newValue)
}
