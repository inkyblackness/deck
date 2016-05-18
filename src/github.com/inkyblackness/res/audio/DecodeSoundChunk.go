package audio

import (
	"bytes"

	"github.com/inkyblackness/res/audio/voc"
)

// DecodeSoundChunk decodes the data of a chunk of type SoundData.
func DecodeSoundChunk(data []byte) (SoundData, error) {
	return voc.Load(bytes.NewReader(data))
}
