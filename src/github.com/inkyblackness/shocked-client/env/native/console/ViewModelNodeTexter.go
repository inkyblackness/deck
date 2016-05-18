package console

// ViewModelLiner allows a texter to add an entry to the main view.
type ViewModelLiner func(label string, value string, texter ViewModelNodeTexter)

// ViewModelNodeTexter is a wrapper for a view model node. It acts as the
// interface to the text UI.
type ViewModelNodeTexter interface {
	// Act is called when the user selects the main entry.
	Act(viewFactory NodeDetailViewFactory)

	// TextMain is called to let the texter add its representation with the provided
	// liner.
	TextMain(addLiner ViewModelLiner)
}
