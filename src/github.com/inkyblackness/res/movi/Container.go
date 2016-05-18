package movi

import (
	"image/color"
)

// Container wraps the information and data of a MOVI container.
type Container interface {
	// MediaDuration returns the duration of the media in seconds.
	MediaDuration() float32

	// VideoWidth returns the width of a video in pixel.
	VideoWidth() uint16
	// VideoHeight returns the height of a video in pixel.
	VideoHeight() uint16
	// StartPalette returns the initial pallete of a video.
	StartPalette() color.Palette

	// AudioSampleRate returns the sample frequency used for audio entries.
	AudioSampleRate() uint16

	// EntryCount returns the number of available entries.
	EntryCount() int
	// Entry returns the entry for given index.
	Entry(index int) Entry
}
