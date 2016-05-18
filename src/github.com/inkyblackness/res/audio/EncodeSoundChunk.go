package audio

import (
	"bytes"

	"github.com/inkyblackness/res/audio/voc"
)

// EncodeSoundChunk encodes the provided sound data into a byte array for a chunk.
func EncodeSoundChunk(soundData SoundData) []byte {
	writer := bytes.NewBuffer(nil)

	voc.Save(writer, soundData.SampleRate(), soundData.Samples(0, soundData.SampleCount()))

	return writer.Bytes()
}
