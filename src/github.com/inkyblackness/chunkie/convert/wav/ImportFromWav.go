package wav

import (
	"io"
	"os"

	wavLib "github.com/youpy/go-wav"

	"github.com/inkyblackness/res/audio"
	"github.com/inkyblackness/res/audio/mem"
)

func l16ToL8(sample int) byte {
	return byte(sample>>8 + 0x80)
}

func l8ToL8(sample int) byte {
	return byte(sample)
}

// ImportFromWav reads the file identified by given name and returns a SoundData instance
// that wraps the contained samples.
func ImportFromWav(fileName string) audio.SoundData {
	var samples []byte
	file, _ := os.Open(fileName)
	defer file.Close()
	reader := wavLib.NewReader(file)
	format, _ := reader.Format()
	converter := l16ToL8
	eof := false

	if format.BitsPerSample == 8 {
		converter = l8ToL8
	}
	for !eof {
		frame, err := reader.ReadSamples(1)
		if err == io.EOF {
			eof = true
		} else {
			samples = append(samples, converter(frame[0].Values[0]))
		}
	}

	return mem.NewL8SoundData(float32(format.SampleRate), samples)
}
