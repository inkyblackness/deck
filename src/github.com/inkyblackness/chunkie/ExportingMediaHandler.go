package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"

	"github.com/inkyblackness/res/audio/mem"
	"github.com/inkyblackness/res/movi"

	"github.com/inkyblackness/chunkie/convert/wav"
)

var subtitleLanguages = map[movi.SubtitleControl]string{
	movi.SubtitleTextStd: "en",
	movi.SubtitleTextFrn: "fr",
	movi.SubtitleTextGer: "de"}

type subtitleEntry struct {
	file    io.WriteCloser
	counter int

	timestamp float32
	text      string
}

type exportingMediaHandler struct {
	mediaDuration float32
	fileBaseName  string

	sampleRate float32
	audio      []byte

	subtitles map[movi.SubtitleControl]*subtitleEntry

	frameCounter       int
	lastFrameTimestamp float32
	lastFrame          *image.Paletted

	framesPerSecond float32
}

func newExportingMediaHandler(fileBaseName string, mediaDuration float32, framesPerSecond float32, sampleRate float32) *exportingMediaHandler {
	return &exportingMediaHandler{
		mediaDuration:   mediaDuration,
		fileBaseName:    fileBaseName,
		subtitles:       make(map[movi.SubtitleControl]*subtitleEntry),
		framesPerSecond: framesPerSecond,
		sampleRate:      sampleRate}
}

func (handler *exportingMediaHandler) finish() {
	handler.writeLastFramesUntil(handler.mediaDuration)
	for _, entry := range handler.subtitles {
		handler.finishSubtitle(entry, handler.mediaDuration)
		entry.file.Close()
	}
	if len(handler.audio) > 0 {
		soundData := mem.NewL8SoundData(handler.sampleRate, handler.audio)
		wav.ExportToWav(handler.fileBaseName+".wav", soundData)
	}
}

func (handler *exportingMediaHandler) OnAudio(timestamp float32, samples []byte) {
	handler.audio = append(handler.audio, samples...)
}

func (handler *exportingMediaHandler) OnSubtitle(timestamp float32, control movi.SubtitleControl, text string) {
	if control != movi.SubtitleArea {
		entry := handler.subtitles[control]

		if entry == nil {
			file, _ := os.Create(handler.fileBaseName + "_" + subtitleLanguages[control] + ".srt")
			entry = &subtitleEntry{file: file}
			handler.subtitles[control] = entry
		}
		handler.finishSubtitle(entry, timestamp)
		entry.timestamp = timestamp
		entry.text = text
		entry.counter++
	}
}

func (handler *exportingMediaHandler) OnVideo(timestamp float32, frame *image.Paletted) {
	handler.writeLastFramesUntil(timestamp)

	handler.lastFrameTimestamp = timestamp
	handler.lastFrame = image.NewPaletted(frame.Bounds(), frame.Palette)
	copy(handler.lastFrame.Pix, frame.Pix)
}

func (handler *exportingMediaHandler) finishSubtitle(entry *subtitleEntry, endTime float32) {
	if entry.counter > 0 {
		fmt.Fprintf(entry.file, "%d\n", entry.counter)
		fmt.Fprintf(entry.file, "%s --> %s\n", handler.formatTimestamp(entry.timestamp), handler.formatTimestamp(endTime))
		fmt.Fprintf(entry.file, "%s\n", entry.text)
		fmt.Fprintf(entry.file, "\n")
	}
}

func (handler *exportingMediaHandler) formatTimestamp(timestamp float32) string {
	inMillis := uint64(timestamp * 1000)
	inSeconds := inMillis / 1000
	inMinutes := inSeconds / 60
	inHours := inMinutes / 60

	return fmt.Sprintf("%02d:%02d:%02d,%03d", inHours, inMinutes%60, inSeconds%60, inMillis%1000)
}

func (handler *exportingMediaHandler) writeLastFramesUntil(timestamp float32) {
	if handler.lastFrame != nil {
		if handler.framesPerSecond > 0 {
			limitFrameId := int((timestamp * handler.framesPerSecond) + 0.5)
			lastFrameId := int((handler.lastFrameTimestamp * handler.framesPerSecond) + 0.5)

			for lastFrameId < limitFrameId {
				name := fmt.Sprintf("%s_%04d.png", handler.fileBaseName, handler.frameCounter)
				handler.frameCounter++

				handler.writeFrame(handler.lastFrame, name)
				lastFrameId++
			}
		} else {
			name := handler.timedFileName(handler.lastFrameTimestamp)
			handler.writeFrame(handler.lastFrame, name)
		}
	}
}

func (handler *exportingMediaHandler) writeFrame(frame *image.Paletted, name string) {
	file, _ := os.Create(name)

	png.Encode(file, frame)
	file.Close()
}

func (handler *exportingMediaHandler) timedFileName(timestamp float32) string {
	inMillis := uint64(timestamp * 1000)
	inSeconds := inMillis / 1000

	return fmt.Sprintf("%s_%03d.%03d.png", handler.fileBaseName, inSeconds, inMillis%1000)
}
