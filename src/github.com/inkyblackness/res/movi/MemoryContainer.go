package movi

import (
	"image/color"
)

type memoryContainer struct {
	mediaDuration float32

	videoWidth   uint16
	videoHeight  uint16
	startPalette color.Palette

	audioSampleRate uint16

	entries []Entry
}

func (container *memoryContainer) MediaDuration() float32 {
	return container.mediaDuration
}

func (container *memoryContainer) VideoWidth() uint16 {
	return container.videoWidth
}

func (container *memoryContainer) VideoHeight() uint16 {
	return container.videoHeight
}

func (container *memoryContainer) StartPalette() color.Palette {
	return container.startPalette
}

func (container *memoryContainer) AudioSampleRate() uint16 {
	return container.audioSampleRate
}

func (container *memoryContainer) EntryCount() int {
	return len(container.entries)
}

func (container *memoryContainer) Entry(index int) Entry {
	return container.entries[index]
}
