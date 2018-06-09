package cmd

import (
	"github.com/inkyblackness/res/audio"
	"github.com/inkyblackness/shocked-model"
)

// ElectronicMessageStore is the simple interface for a message store.
type ElectronicMessageStore interface {
	RequestRemove()
	RequestMessageChange(model.ElectronicMessage)
	RequestAudioChange(model.ResourceLanguage, audio.SoundData)
}

// RemoveElectronicMessageCommand removes and restores an electronic message.
type RemoveElectronicMessageCommand struct {
	RestoreState func()
	Store        ElectronicMessageStore

	Properties model.ElectronicMessage
	Audio      [model.LanguageCount]audio.SoundData
}

// Do removes the message.
func (cmd RemoveElectronicMessageCommand) Do() error {
	cmd.RestoreState()
	cmd.Store.RequestRemove()
	return nil
}

// Undo restores the message.
func (cmd RemoveElectronicMessageCommand) Undo() error {
	cmd.RestoreState()
	cmd.Store.RequestMessageChange(cmd.Properties)
	for lang := 0; lang < model.LanguageCount; lang++ {
		if cmd.Audio[lang] != nil {
			cmd.Store.RequestAudioChange(model.LocalLanguages()[lang], cmd.Audio[lang])
		}
	}
	return nil
}
