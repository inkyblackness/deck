package cmd

// SetEditorModeCommand changes the current editor mode.
type SetEditorModeCommand struct {
	Activator func(name string)
	OldMode   string
	NewMode   string
}

// Do activates the new mode.
func (cmd SetEditorModeCommand) Do() error {
	cmd.Activator(cmd.NewMode)
	return nil
}

// Undo activates the old mode.
func (cmd SetEditorModeCommand) Undo() error {
	cmd.Activator(cmd.OldMode)
	return nil
}
