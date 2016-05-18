package movi

import (
	"bytes"

	"github.com/inkyblackness/res/audio"
	"github.com/inkyblackness/res/audio/mem"
)

// ExtractAudio decodes the given data array as a MOVI container and
// extracts the audio track.
func ExtractAudio(data []byte) (soundData audio.SoundData, err error) {
	container, err := Read(bytes.NewReader(data))

	if container != nil {
		var samples []byte

		for i := 0; i < container.EntryCount(); i++ {
			entry := container.Entry(i)

			if entry.Type() == Audio {
				samples = append(samples, entry.Data()...)
			}
		}
		soundData = mem.NewL8SoundData(float32(container.AudioSampleRate()), samples)
	}
	return
}
