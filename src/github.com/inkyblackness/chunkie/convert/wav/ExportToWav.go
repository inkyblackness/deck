package wav

import (
	"os"

	export "github.com/youpy/go-wav"

	"github.com/inkyblackness/res/audio"
)

// ExportToWav writes a file in the RIFF WAVE format, based on the provided sound data.
func ExportToWav(fileName string, soundData audio.SoundData) {
	file, _ := os.Create(fileName)
	defer file.Close()
	numSamples := soundData.SampleCount()
	numChannels := uint16(1)
	bitsPerSample := uint16(8)
	sampleRate := uint32(soundData.SampleRate())
	writer := export.NewWriter(file, uint32(numSamples), numChannels, sampleRate, bitsPerSample)

	inSamples := soundData.Samples(0, numSamples)
	outSamples := make([]export.Sample, numSamples)
	for index, sample := range inSamples {
		outSamples[index].Values[0] = int(sample)
	}
	writer.WriteSamples(outSamples)
}
