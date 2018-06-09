package cmd

import "github.com/inkyblackness/res/audio"

// SetAudioCommand changes an audio clip.
type SetAudioCommand struct {
	Setter   func(data audio.SoundData) error
	OldValue audio.SoundData
	NewValue audio.SoundData
}

// Do sets the new value.
func (cmd SetAudioCommand) Do() error {
	return cmd.Setter(cmd.NewValue)
}

// Undo sets the old value.
func (cmd SetAudioCommand) Undo() error {
	return cmd.Setter(cmd.OldValue)
}
