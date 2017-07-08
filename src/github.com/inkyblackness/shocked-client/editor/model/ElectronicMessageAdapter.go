package model

import (
	"github.com/inkyblackness/res/audio"

	"github.com/inkyblackness/shocked-model"
)

// ElectronicMessageAdapter is the entry point for an electronic message.
type ElectronicMessageAdapter struct {
	context archiveContext
	store   model.DataStore

	messageType model.ElectronicMessageType
	id          int
	data        *observable

	audio [model.LanguageCount]*observable
}

func newElectronicMessageAdapter(context archiveContext, store model.DataStore) *ElectronicMessageAdapter {
	adapter := &ElectronicMessageAdapter{
		context: context,
		store:   store,

		data: newObservable()}

	for i := 0; i < model.LanguageCount; i++ {
		adapter.audio[i] = newObservable()
	}
	adapter.clear()

	return adapter
}

func (adapter *ElectronicMessageAdapter) clear() {
	adapter.id = -1
	var message model.ElectronicMessage
	adapter.data.set(&message)
	for i := 0; i < model.LanguageCount; i++ {
		adapter.audio[i].set(nil)
	}
}

func (adapter *ElectronicMessageAdapter) messageData() *model.ElectronicMessage {
	return adapter.data.get().(*model.ElectronicMessage)
}

// OnMessageDataChanged registers a callback for data changes.
func (adapter *ElectronicMessageAdapter) OnMessageDataChanged(callback func()) {
	adapter.data.addObserver(callback)
}

// OnMessageAudioChanged registers a callback for audio changes.
func (adapter *ElectronicMessageAdapter) OnMessageAudioChanged(callback func()) {
	for i := 0; i < model.LanguageCount; i++ {
		adapter.audio[i].addObserver(callback)
	}
}

// ID returns the ID of the electronic message.
func (adapter *ElectronicMessageAdapter) ID() int {
	return adapter.id
}

// RequestMessage requests to load the message data of specified ID.
func (adapter *ElectronicMessageAdapter) RequestMessage(messageType model.ElectronicMessageType, id int) {
	adapter.clear()
	adapter.id = id
	adapter.messageType = messageType
	adapter.store.ElectronicMessage(adapter.context.ActiveProjectID(), messageType, id,
		func(message model.ElectronicMessage) { adapter.onMessageData(messageType, id, message) },
		adapter.context.simpleStoreFailure("ElectronicMessage"))
	if (messageType == model.ElectronicMessageTypeLog) || (messageType == model.ElectronicMessageTypeMail) {
		for _, language := range model.LocalLanguages() {
			adapter.requestAudio(messageType, id, language)
		}
	}
}

func (adapter *ElectronicMessageAdapter) requestAudio(messageType model.ElectronicMessageType, id int, language model.ResourceLanguage) {
	adapter.store.ElectronicMessageAudio(adapter.context.ActiveProjectID(), messageType, id, language,
		func(data audio.SoundData) {
			adapter.audio[language.ToIndex()].set(data)
		}, adapter.context.simpleStoreFailure("ElectronicMessageAudio"))
}

// RequestMessageChange requests to change the properties of the current message.
func (adapter *ElectronicMessageAdapter) RequestMessageChange(properties model.ElectronicMessage) {
	if adapter.id >= 0 {
		adapter.store.SetElectronicMessage(adapter.context.ActiveProjectID(), adapter.messageType, adapter.id,
			properties,
			func(message model.ElectronicMessage) { adapter.onMessageData(adapter.messageType, adapter.id, message) },
			adapter.context.simpleStoreFailure("SetElectronicMessage"))
	}
}

// RequestAudioChange requests to change the audio of the current message.
func (adapter *ElectronicMessageAdapter) RequestAudioChange(language model.ResourceLanguage, data audio.SoundData) {
	if adapter.id >= 0 {
		adapter.store.SetElectronicMessageAudio(adapter.context.ActiveProjectID(), adapter.messageType, adapter.id, language, data,
			func() { adapter.audio[language.ToIndex()].set(data) },
			adapter.context.simpleStoreFailure("SetElectronicMessageAudio"))
	}
}

func (adapter *ElectronicMessageAdapter) onMessageData(messageType model.ElectronicMessageType, id int, message model.ElectronicMessage) {
	if (adapter.messageType == messageType) && (adapter.id == id) {
		adapter.data.set(&message)
	}
}

// Audio returns the audio of the message.
func (adapter *ElectronicMessageAdapter) Audio(language int) (data audio.SoundData) {
	ptr := adapter.audio[language].get()
	if ptr != nil {
		data = ptr.(audio.SoundData)
	}
	return
}

// Title returns the title of the message.
func (adapter *ElectronicMessageAdapter) Title(language int) string {
	return safeString(adapter.messageData().Title[language])
}

// Sender returns the sender of the message.
func (adapter *ElectronicMessageAdapter) Sender(language int) string {
	return safeString(adapter.messageData().Sender[language])
}

// Subject returns the subject of the message.
func (adapter *ElectronicMessageAdapter) Subject(language int) string {
	return safeString(adapter.messageData().Subject[language])
}

// VerboseText returns the text in long form of the message.
func (adapter *ElectronicMessageAdapter) VerboseText(language int) string {
	return safeString(adapter.messageData().VerboseText[language])
}

// TerseText returns the text in short form of the message.
func (adapter *ElectronicMessageAdapter) TerseText(language int) string {
	return safeString(adapter.messageData().TerseText[language])
}

// NextMessage returns the identifier of an interrupting message. Or -1 if no interrupt.
func (adapter *ElectronicMessageAdapter) NextMessage() int {
	return safeInt(adapter.messageData().NextMessage, -1)
}

// IsInterrupt returns true if this message is an interrupt of another.
func (adapter *ElectronicMessageAdapter) IsInterrupt() bool {
	var isInterrupt = adapter.messageData().IsInterrupt
	return (isInterrupt != nil) && *isInterrupt
}

// ColorIndex returns the color index for the header text. -1 for default color.
func (adapter *ElectronicMessageAdapter) ColorIndex() int {
	return safeInt(adapter.messageData().ColorIndex, -1)
}

// LeftDisplay returns the display index for the left side. -1 for no display.
func (adapter *ElectronicMessageAdapter) LeftDisplay() int {
	return safeInt(adapter.messageData().LeftDisplay, -1)
}

// RightDisplay returns the display index for the right side. -1 for no display.
func (adapter *ElectronicMessageAdapter) RightDisplay() int {
	return safeInt(adapter.messageData().RightDisplay, -1)
}
