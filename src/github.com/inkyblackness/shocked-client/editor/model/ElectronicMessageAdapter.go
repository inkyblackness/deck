package model

import (
	"github.com/inkyblackness/shocked-model"
)

// ElectronicMessageAdapter is the entry point for an electronic message.
type ElectronicMessageAdapter struct {
	context archiveContext
	store   model.DataStore

	messageType model.ElectronicMessageType
	id          int
	data        *observable
}

func newElectronicMessageAdapter(context archiveContext, store model.DataStore) *ElectronicMessageAdapter {
	adapter := &ElectronicMessageAdapter{
		context: context,
		store:   store,

		data: newObservable()}

	adapter.clear()

	return adapter
}

func (adapter *ElectronicMessageAdapter) clear() {
	adapter.id = -1
	var message model.ElectronicMessage
	adapter.data.set(&message)
}

func (adapter *ElectronicMessageAdapter) messageData() *model.ElectronicMessage {
	return adapter.data.get().(*model.ElectronicMessage)
}

// OnMessageDataChanged registers a callback for data changes.
func (adapter *ElectronicMessageAdapter) OnMessageDataChanged(callback func()) {
	adapter.data.addObserver(callback)
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

func (adapter *ElectronicMessageAdapter) onMessageData(messageType model.ElectronicMessageType, id int, message model.ElectronicMessage) {
	if (adapter.messageType == messageType) && (adapter.id == id) {
		adapter.data.set(&message)
	}
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
