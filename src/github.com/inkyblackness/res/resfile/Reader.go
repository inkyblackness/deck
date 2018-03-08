package resfile

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/inkyblackness/res/resfile/compression"
	"github.com/inkyblackness/res/serial"
)

// Reader provides methods to extract resource data from a serialized form.
// Chunks may be accessed concurrently due to the nature of the underlying io.ReaderAt.
type Reader struct {
	source           io.ReaderAt
	firstChunkOffset uint32
	directory        []chunkDirectoryEntry
}

var errSourceNil = errors.New("source is nil")
var errFormatMismatch = errors.New("format mismatch")

// ReaderFrom accesses the provided source and creates a new Reader instance
// from it.
// Should the provided decoder not follow the resource file format, an error
// is returned.
func ReaderFrom(source io.ReaderAt) (reader *Reader, err error) {
	if source == nil {
		return nil, errSourceNil
	}

	var dirOffset uint32
	dirOffset, err = readAndVerifyHeader(io.NewSectionReader(source, 0, chunkDirectoryFileOffsetPos+4))
	if err != nil {
		return nil, err
	}
	firstChunkOffset, directory, err := readDirectoryAt(dirOffset, source)
	if err != nil {
		return nil, err
	}

	reader = &Reader{
		source:           source,
		firstChunkOffset: firstChunkOffset,
		directory:        directory}

	return
}

// IDs returns the chunk identifier available via this reader.
// The order in the slice is the same as in the underlying serialized form.
func (reader *Reader) IDs() []ChunkID {
	ids := make([]ChunkID, len(reader.directory))
	for index, entry := range reader.directory {
		ids[index] = ChunkID(entry.ID)
	}
	return ids
}

var errInvalidIdentifier = errors.New("chunk: invalid identifier")

// Chunk returns a reader for the specified chunk.
// An error is returned if either the ID is not known, or the chunk could not be prepared.
func (reader *Reader) Chunk(id Identifier) (*ChunkReader, error) {
	chunkStartOffset, entry := reader.findEntry(id.Value())
	if entry == nil {
		return nil, errInvalidIdentifier
	}
	chunkType := entry.chunkType()
	compressed := (chunkType & chunkTypeFlagCompressed) != 0
	fragmented := (chunkType & chunkTypeFlagFragmented) != 0
	contentType := ContentType(entry.contentType())

	if fragmented {
		return reader.newFragmentedChunkReader(entry, contentType, compressed, chunkStartOffset)
	}
	return reader.newSingleBlockChunkReader(entry, contentType, compressed, chunkStartOffset)
}

func readAndVerifyHeader(source io.ReadSeeker) (dirOffset uint32, err error) {
	coder := serial.NewPositioningDecoder(source)
	data := make([]byte, chunkDirectoryFileOffsetPos)
	coder.Code(data)
	coder.Code(&dirOffset)

	expected := make([]byte, len(headerString)+1)
	for index, r := range headerString {
		expected[index] = byte(r)
	}
	expected[len(headerString)] = commentTerminator
	if !bytes.Equal(data[:len(expected)], expected) {
		return 0, errFormatMismatch
	}

	return dirOffset, coder.FirstError()
}

func readDirectoryAt(dirOffset uint32, source io.ReaderAt) (firstChunkOffset uint32, directory []chunkDirectoryEntry, err error) {
	var header chunkDirectoryHeader
	headerSize := int64(binary.Size(&header))
	{
		headerCoder := serial.NewDecoder(io.NewSectionReader(source, int64(dirOffset), headerSize))
		headerCoder.Code(&header)
		if headerCoder.FirstError() != nil {
			return 0, nil, headerCoder.FirstError()
		}
	}

	firstChunkOffset = header.FirstChunkOffset
	directory = make([]chunkDirectoryEntry, header.ChunkCount)
	if header.ChunkCount > 0 {
		listCoder := serial.NewDecoder(io.NewSectionReader(source, int64(dirOffset)+headerSize, int64(binary.Size(directory))))
		listCoder.Code(directory)
		err = listCoder.FirstError()
	}
	return
}

func (reader *Reader) findEntry(id uint16) (startOffset uint32, entry *chunkDirectoryEntry) {
	startOffset = reader.firstChunkOffset
	for index := 0; (index < len(reader.directory)) && (entry == nil); index++ {
		cur := &reader.directory[index]
		if cur.ID == id {
			entry = cur
		} else {
			startOffset += cur.packedLength()
			startOffset += (boundarySize - (startOffset % boundarySize)) % boundarySize
		}
	}
	return
}

type blockListEntry struct {
	start uint32
	size  uint32
}

func (reader *Reader) newFragmentedChunkReader(entry *chunkDirectoryEntry,
	contentType ContentType, compressed bool, chunkStartOffset uint32) (*ChunkReader, error) {
	chunkDataReader := io.NewSectionReader(reader.source, int64(chunkStartOffset), int64(entry.packedLength()))

	firstBlockOffset, blockList, err := reader.readBlockList(chunkDataReader)
	if err != nil {
		fmt.Printf("Fail with size: %v at start offset 0x%08X\n", chunkDataReader.Size(), chunkStartOffset)
		return nil, err
	}
	blockCount := len(blockList)

	rawBlockDataReader := io.NewSectionReader(chunkDataReader, int64(firstBlockOffset), chunkDataReader.Size()-int64(firstBlockOffset))
	var uncompressedReader io.ReaderAt
	if !compressed {
		uncompressedReader = rawBlockDataReader
	}

	blockReader := func(index int) (io.Reader, error) {
		if (index < 0) || (index >= blockCount) {
			return nil, fmt.Errorf("block index wrong: %v/%v", index, blockCount)
		}

		if compressed && (uncompressedReader == nil) {
			decompressor := compression.NewDecompressor(rawBlockDataReader)
			decompressedData, err := ioutil.ReadAll(decompressor)
			if err != nil {
				return nil, err
			}
			uncompressedReader = bytes.NewReader(decompressedData)
		}

		entry := blockList[index]
		reader := io.NewSectionReader(uncompressedReader, int64(entry.start)-int64(firstBlockOffset), int64(entry.size))
		return reader, nil
	}

	return &ChunkReader{
		fragmented:  true,
		blockCount:  len(blockList),
		contentType: contentType,
		compressed:  compressed,
		blockReader: blockReader}, nil
}

func (reader *Reader) readBlockList(source io.Reader) (uint32, []blockListEntry, error) {
	listDecoder := serial.NewDecoder(source)
	var blockCount uint16
	listDecoder.Code(&blockCount)
	var firstBlockOffset uint32
	listDecoder.Code(&firstBlockOffset)
	lastBlockEndOffset := firstBlockOffset
	blockList := make([]blockListEntry, blockCount)
	for blockIndex := uint16(0); blockIndex < blockCount; blockIndex++ {
		var endOffset uint32
		listDecoder.Code(&endOffset)
		blockList[blockIndex].start = lastBlockEndOffset
		blockList[blockIndex].size = endOffset - lastBlockEndOffset
		lastBlockEndOffset = endOffset
	}

	if listDecoder.FirstError() != nil {
		fmt.Printf("reading of block list failed. count: %v, first offset: %v\n",
			blockCount, firstBlockOffset)
	}

	return firstBlockOffset, blockList, listDecoder.FirstError()
}

func (reader *Reader) newSingleBlockChunkReader(entry *chunkDirectoryEntry,
	contentType ContentType, compressed bool, chunkStartOffset uint32) (*ChunkReader, error) {
	blockReader := func(index int) (io.Reader, error) {
		if index != 0 {
			return nil, fmt.Errorf("block index wrong: %v/%v", index, 1)
		}
		chunkSize := entry.packedLength()
		var chunkSource io.Reader = io.NewSectionReader(reader.source, int64(chunkStartOffset), int64(entry.packedLength()))
		if compressed {
			chunkSize = entry.unpackedLength()
			chunkSource = compression.NewDecompressor(chunkSource)
		}
		return io.LimitReader(chunkSource, int64(chunkSize)), nil
	}

	return &ChunkReader{
		fragmented:  false,
		blockCount:  1,
		contentType: contentType,
		compressed:  compressed,
		blockReader: blockReader}, nil
}
