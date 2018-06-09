package cmd

// SetBooleanPropertyCommand changes a boolean property.
type SetBooleanPropertyCommand struct {
	Setter   func(value bool) error
	OldValue bool
	NewValue bool
}

// Do sets the new value.
func (cmd SetBooleanPropertyCommand) Do() error {
	return cmd.Setter(cmd.NewValue)
}

// Undo sets the old value.
func (cmd SetBooleanPropertyCommand) Undo() error {
	return cmd.Setter(cmd.OldValue)
}
