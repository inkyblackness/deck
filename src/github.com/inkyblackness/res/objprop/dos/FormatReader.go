package dos

import (
	"fmt"
	"io"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/objprop"
	"github.com/inkyblackness/res/serial"
)

type formatReader struct {
	coder   serial.PositioningCoder
	entries map[res.ObjectID]*typeEntry
}

var errFormatMismatch = fmt.Errorf("Format mismatch")

// NewProvider wraps the provided ReadSeeker in a provider for object properties.
func NewProvider(source io.ReadSeeker, descriptors []objprop.ClassDescriptor) (provider objprop.Provider, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	if source == nil {
		panic(fmt.Errorf("source is nil"))
	}
	verifySourceLength(source, descriptors)

	provider = &formatReader{
		coder:   serial.NewPositioningDecoder(source),
		entries: calculateEntryValues(descriptors)}

	return
}

func (provider *formatReader) Provide(id res.ObjectID) objprop.ObjectData {
	entry := provider.entries[id]
	data := objprop.ObjectData{
		Generic:  make([]byte, entry.genericLength),
		Specific: make([]byte, entry.specificLength),
		Common:   make([]byte, objprop.CommonPropertiesLength)}

	codeObjectData(provider.coder, entry, &data)

	return data
}

func verifySourceLength(source io.Seeker, descriptors []objprop.ClassDescriptor) {
	sourceLength := getSeekerSize(source)
	expectedLength := expectedDataLength(descriptors)

	if expectedLength != sourceLength {
		panic(errFormatMismatch)
	}
}

func getSeekerSize(seeker io.Seeker) uint32 {
	length, err := seeker.Seek(0, 2)

	if err != nil {
		panic(err)
	}

	return uint32(length)
}
