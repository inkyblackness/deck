package dos

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/inkyblackness/res/textprop"
)

type formatReader struct {
	source     io.ReadSeeker
	entryCount uint32
}

var errFormatMismatch = fmt.Errorf("Format mismatch")

// NewProvider wraps the provided ReadSeeker in a provider for texture properties.
func NewProvider(source io.ReadSeeker) (provider textprop.Provider, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	if source == nil {
		panic(fmt.Errorf("source is nil"))
	}
	count := readAndVerifyEntryCount(source)

	provider = &formatReader{source: source, entryCount: count}

	return
}

func (provider *formatReader) EntryCount() uint32 {
	return provider.entryCount
}

func (provider *formatReader) Provide(id uint32) []byte {
	data := make([]byte, int(textprop.TexturePropertiesLength))

	provider.source.Seek(int64(MagicHeaderSize+textprop.TexturePropertiesLength*id), 0)
	binary.Read(provider.source, binary.LittleEndian, data)

	return data
}

func readAndVerifyEntryCount(source io.ReadSeeker) uint32 {
	sourceLength := getSeekerSize(source)
	count := uint32((sourceLength - MagicHeaderSize) / textprop.TexturePropertiesLength)
	header := uint32(0)

	binary.Read(source, binary.LittleEndian, &header)

	if header != MagicHeader {
		panic(errFormatMismatch)
	}

	return count
}

func getSeekerSize(seeker io.Seeker) uint32 {
	length, err := seeker.Seek(0, 2)

	if err != nil {
		panic(err)
	}
	seeker.Seek(0, 0)

	return uint32(length)
}
