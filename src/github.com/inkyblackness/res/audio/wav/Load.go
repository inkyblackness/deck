package wav

import (
	"fmt"
	"io"

	"github.com/inkyblackness/res/audio/mem"
)

var errNotASupportedWave = fmt.Errorf("Not a supported WAV")

// Load reads from the provided source and returns the data.
func Load(source io.Reader) (data *mem.L8SoundData, err error) {
	if source == nil {
		err = fmt.Errorf("source is nil")
	} else {
		var loader waveLoader

		loader.load(source)
		err = loader.err
		if err == nil {
			data = mem.NewL8SoundData(loader.sampleRate, loader.samples)
		}
	}

	return
}
