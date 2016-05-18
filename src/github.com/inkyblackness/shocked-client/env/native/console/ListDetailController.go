package console

// ListDetailController is a detail controller for simple lists.
type ListDetailController interface {
	DetailController

	// Confirm is called when the view has been confirmed to select the given index.
	Confirm(index int)
}
