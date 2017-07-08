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
	dataType  res.DataTypeID
	chunkType chunk.TypeID
	limit     uint16
}

var knownSounds = map[model.ResourceType]soundInfo{
	model.ResourceTypeTrapAudio: {res.Media, chunk.BasicChunkType, model.MaxTrapMessages}}

// Sounds is the adapter for general sounds.
type Sounds struct {
	citbark [model.LanguageCount]chunk.Store
}

// NewSounds returns a new Sounds instance, if possible.
func NewSounds(library io.StoreLibrary) (sounds *Sounds, err error) {
	var citbark [model.LanguageCount]chunk.Store

	for i := 0; i < model.LanguageCount && err == nil; i++ {
		citbark[i], err = library.ChunkStore(localized[i].citbark)
	}
	if err == nil {
		sounds = &Sounds{
			citbark: citbark}
	}

	return
}

func (sounds *Sounds) store(key model.ResourceKey) (store chunk.Store) {
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

		if (holder != nil) && (info.dataType == res.Media) {
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
func (sounds *Sounds) SetAudio(key model.ResourceKey, data audio.SoundData) (resultKey model.ResourceKey, err error) {
	info, known := knownSounds[key.Type]

	if known && (key.Index < info.limit) && key.HasValidLanguage() {
		store := sounds.store(key)

		if info.dataType == res.Media {
			data := movi.ContainSoundData(data)
			store.Put(res.ResourceID(int(key.Type)+int(key.Index)), chunk.NewBlockHolder(info.chunkType, res.Media, [][]byte{data}))
		}
		resultKey = key
	} else {
		err = fmt.Errorf("Unsupported resource key: %v", key)
	}

	return
}
