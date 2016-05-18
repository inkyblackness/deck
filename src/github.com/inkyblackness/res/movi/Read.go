package movi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/inkyblackness/res/image"
	"github.com/inkyblackness/res/movi/format"
)

// Read tries to extract a MOVI container from the provided reader.
// On success the position of the reader is past the last data entry.
// On failure the position of the reader is undefined.
func Read(source io.ReadSeeker) (container Container, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	if source == nil {
		panic(fmt.Errorf("source is nil"))
	}

	var header format.Header
	startPos, _ := source.Seek(0, os.SEEK_CUR)

	binary.Read(source, binary.LittleEndian, &header)
	builder := NewContainerBuilder()
	verifyAndExtractHeader(source, builder, &header)
	readPalette(source, builder)
	readIndexAndEntries(source, startPos, builder, &header)

	container = builder.Build()
	return
}

func verifyAndExtractHeader(source io.Reader, builder *ContainerBuilder, header *format.Header) {
	if !bytes.Equal(header.Tag[:], bytes.NewBufferString(format.Tag).Bytes()) {
		panic(fmt.Errorf("Not a MOVI format"))
	}

	builder.MediaDuration(timeFromRaw(header.DurationSeconds, header.DurationFraction))
	builder.VideoHeight(header.VideoHeight)
	builder.VideoWidth(header.VideoWidth)
	builder.AudioSampleRate(header.SampleRate)
}

func readPalette(source io.Reader, builder *ContainerBuilder) {
	palette, err := image.LoadPalette(source)

	if err != nil {
		panic(err)
	}

	builder.StartPalette(palette)
}

func readIndexAndEntries(source io.ReadSeeker, startPos int64, builder *ContainerBuilder, header *format.Header) {
	indexEntries := make([]format.IndexTableEntry, header.IndexEntryCount)

	binary.Read(source, binary.LittleEndian, indexEntries)
	for index, indexEntry := range indexEntries {
		entryType := DataType(indexEntry.Type)

		if entryType != endOfMedia {
			timestamp := timeFromRaw(indexEntry.TimestampSecond, indexEntry.TimestampFraction)
			length := int(indexEntries[index+1].DataOffset - indexEntry.DataOffset)
			data := make([]byte, length)

			source.Seek(startPos+int64(indexEntry.DataOffset), 0)
			source.Read(data)

			builder.AddEntry(NewMemoryEntry(timestamp, entryType, data))
		}
	}
}
