package console

// ViewModelListener is the main interface for changes of the view model.
type ViewModelListener interface {
	// OnMainDataChanged is called when the main information changed.
	OnMainDataChanged()
}
