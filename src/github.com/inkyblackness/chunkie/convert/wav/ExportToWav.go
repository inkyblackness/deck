package wav

import (
	"os"

	"github.com/inkyblackness/res/audio"
	"github.com/inkyblackness/res/audio/wav"
)

// ExportToWav writes a file in the RIFF WAVE format, based on the provided sound data.
func ExportToWav(fileName string, soundData audio.SoundData) {
	file, _ := os.Create(fileName)
	defer file.Close()
	wav.Save(file, soundData.SampleRate(), soundData.Samples(0, soundData.SampleCount()))
}
