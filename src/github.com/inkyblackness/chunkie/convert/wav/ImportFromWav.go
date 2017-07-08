package wav

import (
	"os"

	"github.com/inkyblackness/res/audio"
	"github.com/inkyblackness/res/audio/wav"
)

// ImportFromWav reads the file identified by given name and returns a SoundData instance
// that wraps the contained samples.
func ImportFromWav(fileName string) audio.SoundData {
	file, _ := os.Open(fileName)
	defer file.Close()

	data, _ := wav.Load(file)

	return data
}
