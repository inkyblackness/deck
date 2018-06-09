package cmd

// SetStringPropertyCommand changes a textual property.
type SetStringPropertyCommand struct {
	Setter   func(value string) error
	OldValue string
	NewValue string
}

// Do sets the new value.
func (cmd SetStringPropertyCommand) Do() error {
	return cmd.Setter(cmd.NewValue)
}

// Undo sets the old value.
func (cmd SetStringPropertyCommand) Undo() error {
	return cmd.Setter(cmd.OldValue)
}
