package cmd

// SetLevelTexturesCommand sets the textures of a level.
type SetLevelTexturesCommand struct {
	Setter        func(textureIDs []int) error
	OldTextureIDs []int
	NewTextureIDs []int
}

// Do sets the new value.
func (cmd SetLevelTexturesCommand) Do() error {
	return cmd.Setter(cmd.NewTextureIDs)
}

// Undo sets the old value.
func (cmd SetLevelTexturesCommand) Undo() error {
	return cmd.Setter(cmd.OldTextureIDs)
}
