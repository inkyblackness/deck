package console

import "io"

// DetailController is the base interface for an active detail controller.
type DetailController interface {
	// Cancel is called when the view has been closed without confirmation.
	Cancel()
	// WriteDetails requests the controller to write the current state.
	WriteDetails(w io.Writer)
}
