package cmd

// SetIntPropertyCommand changes an integer property.
type SetIntPropertyCommand struct {
	Setter   func(value int) error
	OldValue int
	NewValue int
}

// Do sets the new value.
func (cmd SetIntPropertyCommand) Do() error {
	return cmd.Setter(cmd.NewValue)
}

// Undo sets the old value.
func (cmd SetIntPropertyCommand) Undo() error {
	return cmd.Setter(cmd.OldValue)
}
