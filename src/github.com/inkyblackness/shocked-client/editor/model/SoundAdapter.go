package model

import (
	"github.com/inkyblackness/res/audio"
	"github.com/inkyblackness/shocked-model"
)

// SoundAdapter is the entry point for a sound.
type SoundAdapter struct {
	context archiveContext
	store   model.DataStore

	resourceKey model.ResourceKey
	data        *observable
}

func newSoundAdapter(context archiveContext, store model.DataStore) *SoundAdapter {
	adapter := &SoundAdapter{
		context: context,
		store:   store,

		data: newObservable()}

	adapter.clear()

	return adapter
}

func (adapter *SoundAdapter) clear() {
	adapter.resourceKey = model.ResourceKeyFromInt(0)
	adapter.publishAudio(nil)
}

func (adapter *SoundAdapter) publishAudio(data audio.SoundData) {
	adapter.data.set(data)
}

// OnAudioChanged registers a callback for sound changes.
func (adapter *SoundAdapter) OnAudioChanged(callback func()) {
	adapter.data.addObserver(callback)
}

// ResourceKey returns the key of the current sound.
func (adapter *SoundAdapter) ResourceKey() model.ResourceKey {
	return adapter.resourceKey
}

// RequestAudio requests to load the audio of specified key.
func (adapter *SoundAdapter) RequestAudio(key model.ResourceKey) {
	adapter.resourceKey = key
	adapter.publishAudio(nil)
	adapter.store.Audio(adapter.context.ActiveProjectID(), key, adapter.onAudio,
		adapter.context.simpleStoreFailure("Audio"))
}

// RequestAudioChange requests to change the audio of the current sound.
func (adapter *SoundAdapter) RequestAudioChange(data audio.SoundData) {
	if adapter.resourceKey.ToInt() > 0 {
		adapter.store.SetAudio(adapter.context.ActiveProjectID(), adapter.resourceKey, data,
			func(resourceKey model.ResourceKey) { adapter.onAudio(resourceKey, data) },
			adapter.context.simpleStoreFailure("SetAudio"))
	}
}

func (adapter *SoundAdapter) onAudio(resourceKey model.ResourceKey, data audio.SoundData) {
	adapter.resourceKey = resourceKey
	adapter.publishAudio(data)
}

// Audio returns the current audio.
func (adapter *SoundAdapter) Audio() (data audio.SoundData) {
	ptr := adapter.data.get()
	if ptr != nil {
		data = ptr.(audio.SoundData)
	}
	return
}
