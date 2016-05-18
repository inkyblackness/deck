package movi

import (
	"bytes"
	"encoding/binary"
	goImage "image"
	"image/color"

	"github.com/inkyblackness/res/compress/rle"
	"github.com/inkyblackness/res/compress/video"
	"github.com/inkyblackness/res/image"
	"github.com/inkyblackness/res/text"
)

// MediaDispatcher iterates through the entries of a container and provides resulting media
// to a handler. The dispatcher handles intermediate container entries to create consumable media.
type MediaDispatcher struct {
	handler MediaHandler

	container Container
	nextIndex int

	codepage text.Codepage

	palette        color.Palette
	decoderBuilder *video.FrameDecoderBuilder
	frameBuffer    []byte
}

// NewMediaDispatcher returns a new instance of a dispatcher reading the provided container.
func NewMediaDispatcher(container Container, handler MediaHandler) *MediaDispatcher {
	width := int(container.VideoWidth())
	height := int(container.VideoHeight())
	dispatcher := &MediaDispatcher{
		handler:        handler,
		container:      container,
		codepage:       text.DefaultCodepage(),
		frameBuffer:    make([]byte, width*height),
		decoderBuilder: video.NewFrameDecoderBuilder(width, height)}

	dispatcher.setPalette(container.StartPalette())
	dispatcher.decoderBuilder.ForStandardFrame(dispatcher.frameBuffer, width)

	return dispatcher
}

// DispatchNext processes the next entries from the container to call the handler.
// Returns false if the dispatcher reached the end of the container.
func (dispatcher *MediaDispatcher) DispatchNext() (result bool, err error) {
	for !result && (dispatcher.nextIndex < dispatcher.container.EntryCount()) {
		entry := dispatcher.container.Entry(dispatcher.nextIndex)
		result, err = dispatcher.process(entry)
		dispatcher.nextIndex++
	}

	return
}

func (dispatcher *MediaDispatcher) process(entry Entry) (dispatched bool, err error) {
	switch entry.Type() {
	case Audio:
		{
			dispatcher.handler.OnAudio(entry.Timestamp(), entry.Data())
			dispatched = true
		}
	case Subtitle:
		{
			var subtitleHeader SubtitleHeader

			binary.Read(bytes.NewReader(entry.Data()), binary.LittleEndian, &subtitleHeader)
			subtitle := dispatcher.codepage.Decode(entry.Data()[SubtitleHeaderSize:])
			dispatcher.handler.OnSubtitle(entry.Timestamp(), subtitleHeader.Control, subtitle)
			dispatched = true
		}

	case Palette:
		{
			newPalette, palErr := image.LoadPalette(bytes.NewReader(entry.Data()))

			if palErr == nil {
				dispatcher.setPalette(newPalette)
				dispatcher.clearFrameBuffer()
			} else {
				err = palErr
			}
		}
	case ControlDictionary:
		{
			words, wordsErr := video.UnpackControlWords(entry.Data())

			if wordsErr == nil {
				dispatcher.decoderBuilder.WithControlWords(words)
			} else {
				err = wordsErr
			}
		}
	case PaletteLookupList:
		{
			dispatcher.decoderBuilder.WithPaletteLookupList(entry.Data())
		}

	case LowResVideo:
		{
			var videoHeader LowResVideoHeader
			reader := bytes.NewReader(entry.Data())

			binary.Read(reader, binary.LittleEndian, &videoHeader)
			frameErr := rle.Decompress(reader, dispatcher.frameBuffer)
			if frameErr == nil {
				dispatcher.notifyVideoFrame(entry.Timestamp())
				dispatched = true
			} else {
				err = frameErr
			}
		}
	case HighResVideo:
		{
			var videoHeader HighResVideoHeader
			reader := bytes.NewReader(entry.Data())

			binary.Read(reader, binary.LittleEndian, &videoHeader)
			bitstreamData := entry.Data()[HighResVideoHeaderSize:videoHeader.PixelDataOffset]
			maskstreamData := entry.Data()[videoHeader.PixelDataOffset:]
			decoder := dispatcher.decoderBuilder.Build()

			decoder.Decode(bitstreamData, maskstreamData)
			dispatcher.notifyVideoFrame(entry.Timestamp())
			dispatched = true
		}
	}

	return
}

func (dispatcher *MediaDispatcher) setPalette(newPalette color.Palette) {
	dispatcher.palette = make([]color.Color, len(newPalette))

	dispatcher.palette[0] = color.NRGBA{R: 0, G: 0, B: 0, A: 0xFF}
	copy(dispatcher.palette[1:], newPalette[1:])
}

func (dispatcher *MediaDispatcher) clearFrameBuffer() {
	for pixel := 0; pixel < len(dispatcher.frameBuffer); pixel++ {
		dispatcher.frameBuffer[pixel] = 0x00
	}
}

func (dispatcher *MediaDispatcher) notifyVideoFrame(timestamp float32) {
	rect := goImage.Rect(0, 0, int(dispatcher.container.VideoWidth()), int(dispatcher.container.VideoHeight()))
	paletted := goImage.NewPaletted(rect, dispatcher.palette)
	copy(paletted.Pix, dispatcher.frameBuffer)
	dispatcher.handler.OnVideo(timestamp, paletted)
}
