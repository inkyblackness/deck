package model

import (
	"github.com/inkyblackness/shocked-model"
)

// TextAdapter is the entry point for a text.
type TextAdapter struct {
	context archiveContext
	store   model.DataStore

	resourceKey model.ResourceKey
	data        *observable
}

func newTextAdapter(context archiveContext, store model.DataStore) *TextAdapter {
	adapter := &TextAdapter{
		context: context,
		store:   store,

		data: newObservable()}

	adapter.clear()

	return adapter
}

func (adapter *TextAdapter) clear() {
	adapter.resourceKey = model.ResourceKeyFromInt(0)
	adapter.publishText("")
}

func (adapter *TextAdapter) publishText(text string) {
	adapter.data.set(&text)
}

// OnTextChanged registers a callback for text changes.
func (adapter *TextAdapter) OnTextChanged(callback func()) {
	adapter.data.addObserver(callback)
}

// ResourceKey returns the key of the current text.
func (adapter *TextAdapter) ResourceKey() model.ResourceKey {
	return adapter.resourceKey
}

// RequestText requests to load the text of specified key.
func (adapter *TextAdapter) RequestText(key model.ResourceKey) {
	adapter.resourceKey = key
	adapter.publishText("")
	adapter.store.Text(adapter.context.ActiveProjectID(), key, adapter.onText,
		adapter.context.simpleStoreFailure("Text"))
}

// RequestTextChange requests to change the properties of the current text.
func (adapter *TextAdapter) RequestTextChange(text string) {
	if adapter.resourceKey.ToInt() > 0 {
		adapter.store.SetText(adapter.context.ActiveProjectID(), adapter.resourceKey, text, adapter.onText,
			adapter.context.simpleStoreFailure("SetText"))
	}
}

func (adapter *TextAdapter) onText(resourceKey model.ResourceKey, text string) {
	adapter.resourceKey = resourceKey
	adapter.publishText(text)
}

// Text returns the current text.
func (adapter *TextAdapter) Text() string {
	return *adapter.data.get().(*string)
}
