package movi

import (
	"bytes"

	"github.com/inkyblackness/res/audio"
)

const audioEntrySize = 0x2000

// ContainSoundData packs a sound data into a container and encodes it.
func ContainSoundData(soundData audio.SoundData) []byte {
	builder := NewContainerBuilder()
	startOffset := 0
	entryStartTime := float32(0)
	timePerEntry := timeFromRaw(timeToRaw(float32(audioEntrySize) / soundData.SampleRate()))

	for (startOffset + audioEntrySize) <= soundData.SampleCount() {
		endOffset := startOffset + audioEntrySize
		builder.AddEntry(NewMemoryEntry(entryStartTime, Audio, soundData.Samples(startOffset, endOffset)))
		entryStartTime += timePerEntry
		startOffset = endOffset
	}
	if startOffset < soundData.SampleCount() {
		builder.AddEntry(NewMemoryEntry(entryStartTime, Audio, soundData.Samples(startOffset, soundData.SampleCount())))
		entryStartTime += timeFromRaw(timeToRaw(float32(soundData.SampleCount()-startOffset) / soundData.SampleRate()))
	}

	builder.MediaDuration(entryStartTime)
	builder.AudioSampleRate(uint16(soundData.SampleRate()))

	container := builder.Build()
	buffer := bytes.NewBuffer(nil)
	Write(buffer, container)
	return buffer.Bytes()
}
