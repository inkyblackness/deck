package cmd

import "github.com/inkyblackness/shocked-model"

// SetBitmapCommand changes an audio clip.
type SetBitmapCommand struct {
	Setter   func(bmp *model.RawBitmap) error
	OldValue *model.RawBitmap
	NewValue *model.RawBitmap
}

// Do sets the new value.
func (cmd SetBitmapCommand) Do() error {
	return cmd.Setter(cmd.NewValue)
}

// Undo sets the old value.
func (cmd SetBitmapCommand) Undo() error {
	return cmd.Setter(cmd.OldValue)
}
