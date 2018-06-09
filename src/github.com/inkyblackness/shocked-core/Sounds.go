package core

import (
	"bytes"
	"fmt"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/audio"
	memAudio "github.com/inkyblackness/res/audio/mem"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/movi"
	"github.com/inkyblackness/shocked-core/io"
	model "github.com/inkyblackness/shocked-model"
)

type soundInfo struct {
	contentType chunk.ContentType
	limit       uint16
}

var knownSounds = map[model.ResourceType]soundInfo{
	model.ResourceTypeTrapAudio: {chunk.Media, model.MaxTrapMessages}}

// Sounds is the adapter for general sounds.
type Sounds struct {
	citbark [model.LanguageCount]*io.DynamicChunkStore
}

// NewSounds returns a new Sounds instance, if possible.
func NewSounds(library io.StoreLibrary) (sounds *Sounds, err error) {
	var citbark [model.LanguageCount]*io.DynamicChunkStore

	for i := 0; i < model.LanguageCount && err == nil; i++ {
		citbark[i], err = library.ChunkStore(localized[i].citbark)
	}
	if err == nil {
		sounds = &Sounds{
			citbark: citbark}
	}

	return
}

func (sounds *Sounds) store(key model.ResourceKey) (store *io.DynamicChunkStore) {
	if key.Type == model.ResourceTypeTrapAudio {
		store = sounds.citbark[key.Language.ToIndex()]
	}
	return
}

// Audio returns the audio of identified sound.
func (sounds *Sounds) Audio(key model.ResourceKey) (data audio.SoundData, err error) {
	info, known := knownSounds[key.Type]
	if known && (key.Index < info.limit) && key.HasValidLanguage() {
		store := sounds.store(key)
		holder := store.Get(res.ResourceID(int(key.Type) + int(key.Index)))

		if (holder != nil) && (info.contentType == chunk.Media) {
			blockData := holder.BlockData(0)
			var container movi.Container
			container, err = movi.Read(bytes.NewReader(blockData))

			if err == nil {
				samples := []byte{}
				for index := 0; index < container.EntryCount(); index++ {
					entry := container.Entry(index)
					if entry.Type() == movi.Audio {
						samples = append(samples, entry.Data()...)
					}
				}
				data = memAudio.NewL8SoundData(float32(container.AudioSampleRate()), samples)
			}
		}
	} else {
		err = fmt.Errorf("Unsupported resource key: %v", key)
	}

	return
}

// SetAudio requests to set the audio of a sound resource.
func (sounds *Sounds) SetAudio(key model.ResourceKey, soundData audio.SoundData) (resultKey model.ResourceKey, err error) {
	info, known := knownSounds[key.Type]

	if known && (key.Index < info.limit) && key.HasValidLanguage() {
		store := sounds.store(key)

		if info.contentType == chunk.Media {
			resourceID := res.ResourceID(int(key.Type) + int(key.Index))

			if soundData != nil {
				encodedData := movi.ContainSoundData(soundData)
				store.Put(resourceID,
					&chunk.Chunk{
						ContentType:   info.contentType,
						BlockProvider: chunk.MemoryBlockProvider([][]byte{encodedData})})
			} else {
				store.Del(resourceID)
			}
		}
		resultKey = key
	} else {
		err = fmt.Errorf("Unsupported resource key: %v", key)
	}

	return
}
