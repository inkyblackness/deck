package cmd

// SetActiveLevelCommand sets the currently active level.
type SetActiveLevelCommand struct {
	Setter   func(levelID int) error
	OldValue int
	NewValue int
}

// Do sets the new value.
func (cmd SetActiveLevelCommand) Do() error {
	return cmd.Setter(cmd.NewValue)
}

// Undo sets the old value.
func (cmd SetActiveLevelCommand) Undo() error {
	return cmd.Setter(cmd.OldValue)
}
