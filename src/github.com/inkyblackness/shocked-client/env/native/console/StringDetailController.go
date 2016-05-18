package console

// StringDetailController is a detail controller for strings.
type StringDetailController interface {
	DetailController

	// Confirm is called when the view has been confirmed to set the provided value.
	Confirm(newValue string)
}
