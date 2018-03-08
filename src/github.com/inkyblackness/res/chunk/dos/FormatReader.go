package dos

import (
	"fmt"
	"io"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/serial"
)

type formatReader struct {
	resourceIDs []res.ResourceID

	blockHolder map[res.ResourceID]chunk.BlockHolder
}

var errFormatMismatch = fmt.Errorf("Format mismatch")

// NewChunkProvider returns a chunk provider reading from a random access reader
// over a DOS format resource file.
func NewChunkProvider(source io.ReadSeeker) (provider chunk.Provider, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	if source == nil {
		panic(fmt.Errorf("source is nil"))
	}
	coder := serial.NewPositioningDecoder(source)

	skipAndVerifyHeaderString(coder)
	skipAndVerifyComment(coder)
	ids, addresses := readAndVerifyDirectory(coder)
	if coder.FirstError() != nil {
		return nil, coder.FirstError()
	}

	formatReader := &formatReader{
		resourceIDs: ids,
		blockHolder: make(map[res.ResourceID]chunk.BlockHolder)}

	for id, address := range addresses {
		blockHolder := &blockReader{coder: coder, address: address}
		formatReader.blockHolder[id] = blockHolder
	}

	provider = formatReader
	return
}

func (reader *formatReader) IDs() []res.ResourceID {
	return reader.resourceIDs
}

// Provide implements the chunk.Provider interface
func (reader *formatReader) Provide(id res.ResourceID) chunk.BlockHolder {
	return reader.blockHolder[id]
}

func skipAndVerifyHeaderString(coder serial.Coder) {
	headerStringBuffer := make([]byte, len(HeaderString))
	coder.Code(headerStringBuffer)
	if string(headerStringBuffer) != HeaderString {
		panic(errFormatMismatch)
	}
}

func skipAndVerifyComment(coder serial.PositioningCoder) {
	terminatorFound := false

	for remaining := ChunkDirectoryFileOffsetPos - coder.CurPos(); remaining > 0; remaining-- {
		temp := byte(0x00)
		coder.Code(&temp)
		if temp == CommentTerminator {
			terminatorFound = true
		}
	}
	if !terminatorFound {
		panic(errFormatMismatch)
	}
}

func readAndVerifyDirectory(coder serial.PositioningCoder) ([]res.ResourceID, map[res.ResourceID]*chunkAddress) {
	directoryFileOffset := uint32(0)
	directoryEntries := uint16(0)
	chunkFileOffset := uint32(0)

	coder.Code(&directoryFileOffset)
	coder.SetCurPos(directoryFileOffset)

	coder.Code(&directoryEntries)
	coder.Code(&chunkFileOffset)
	ids := make([]res.ResourceID, int(directoryEntries))
	addresses := make(map[res.ResourceID]*chunkAddress)

	for i := uint16(0); i < directoryEntries; i++ {
		resourceID := uint16(0xFFFF)
		address := &chunkAddress{}

		coder.Code(&resourceID)
		address.code(coder)

		address.startOffset = chunkFileOffset
		chunkFileOffset += address.chunkLength
		if chunkFileOffset%BoundarySize != 0 {
			chunkFileOffset += BoundarySize - chunkFileOffset%BoundarySize
		}

		ids[i] = res.ResourceID(resourceID)
		addresses[ids[i]] = address
	}

	return ids, addresses
}
