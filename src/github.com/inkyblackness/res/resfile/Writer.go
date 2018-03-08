package resfile

import (
	"errors"
	"io"
	"math"

	"github.com/inkyblackness/res/resfile/compression"
	"github.com/inkyblackness/res/serial"
)

// Writer provides methods to write a new resource file from scratch.
// Chunks have to be created sequentially. The writer does not support
// concurrent creation and modification of chunks.
type Writer struct {
	encoder *serial.PositioningEncoder

	firstChunkOffset        uint32
	currentChunkStartOffset uint32
	currentChunk            chunkWriter

	directory []*chunkDirectoryEntry
}

var errTargetNil = errors.New("target is nil")

// NewWriter returns a new Writer instance prepared to add chunks.
// To finalize the created file, call Finish().
//
// This function will write initial information to the target and will return
// an error if the writer did. In such a case, the returned writer instance
// will produce invalid results and the state of the target is undefined.
func NewWriter(target io.WriteSeeker) (*Writer, error) {
	if target == nil {
		return nil, errTargetNil
	}

	encoder := serial.NewPositioningEncoder(target)
	writer := &Writer{encoder: encoder}
	writer.writeHeader()
	writer.firstChunkOffset = writer.encoder.CurPos()

	return writer, writer.encoder.FirstError()
}

var errWriterFinished = errors.New("writer is finished")

// CreateChunk adds a new single-block chunk to the current resource file.
// This chunk is closed by creating another chunk, or by finishing the writer.
func (writer *Writer) CreateChunk(id Identifier, contentType ContentType, compressed bool) (*BlockWriter, error) {
	if writer.encoder == nil {
		return nil, errWriterFinished
	}

	writer.finishLastChunk()
	if writer.encoder.FirstError() != nil {
		return nil, writer.encoder.FirstError()
	}

	var targetWriter io.Writer = serial.NewEncoder(writer.encoder)
	targetFinisher := func() {}
	chunkType := byte(0x00)
	if compressed {
		compressor := compression.NewCompressor(targetWriter)
		chunkType |= chunkTypeFlagCompressed
		targetWriter = compressor
		targetFinisher = func() { compressor.Close() } // nolint: errcheck
	}
	blockWriter := &BlockWriter{target: targetWriter, finisher: targetFinisher}
	writer.addNewChunk(id, contentType, chunkType, blockWriter)

	return blockWriter, nil
}

// CreateFragmentedChunk adds a new fragmented chunk to the current resource file.
// This chunk is closed by creating another chunk, or by finishing the writer.
func (writer *Writer) CreateFragmentedChunk(id Identifier, contentType ContentType, compressed bool) (*FragmentedChunkWriter, error) {
	if writer.encoder == nil {
		return nil, errWriterFinished
	}

	writer.finishLastChunk()
	if writer.encoder.FirstError() != nil {
		return nil, writer.encoder.FirstError()
	}

	chunkType := chunkTypeFlagFragmented
	if compressed {
		chunkType |= chunkTypeFlagCompressed
	}
	chunkWriter := &FragmentedChunkWriter{
		target:          serial.NewPositioningEncoder(writer.encoder),
		compressed:      compressed,
		dataPaddingSize: writer.dataPaddingSizeForFragmentedChunk(id)}
	writer.addNewChunk(id, contentType, chunkType, chunkWriter)

	return chunkWriter, nil
}

// Finish finalizes the resource file. After calling this function, the
// writer becomes unusable.
func (writer *Writer) Finish() (err error) {
	if writer.encoder == nil {
		return errWriterFinished
	}

	writer.finishLastChunk()

	directoryOffset := writer.encoder.CurPos()
	writer.encoder.SetCurPos(chunkDirectoryFileOffsetPos)
	writer.encoder.Code(directoryOffset)
	writer.encoder.SetCurPos(directoryOffset)
	writer.encoder.Code(uint16(len(writer.directory)))
	writer.encoder.Code(writer.firstChunkOffset)
	for _, entry := range writer.directory {
		writer.encoder.Code(entry)
	}

	err = writer.encoder.FirstError()
	writer.encoder = nil

	return
}

func (writer *Writer) writeHeader() {
	header := make([]byte, chunkDirectoryFileOffsetPos)
	for index, r := range headerString {
		header[index] = byte(r)
	}
	header[len(headerString)] = commentTerminator
	writer.encoder.Code(header)
	writer.encoder.Code(uint32(math.MaxUint32))
}

func (writer *Writer) addNewChunk(id Identifier, contentType ContentType, chunkType byte, newChunk chunkWriter) {
	entry := &chunkDirectoryEntry{ID: id.Value()}
	entry.setContentType(byte(contentType))
	entry.setChunkType(chunkType)
	writer.directory = append(writer.directory, entry)
	writer.currentChunk = newChunk
	writer.currentChunkStartOffset = writer.encoder.CurPos()
}

func (writer *Writer) finishLastChunk() {
	if writer.currentChunk != nil {
		currentEntry := writer.directory[len(writer.directory)-1]
		currentEntry.setUnpackedLength(writer.currentChunk.finish())
		currentEntry.setPackedLength(writer.encoder.CurPos() - writer.currentChunkStartOffset)

		writer.currentChunkStartOffset = 0
		writer.currentChunk = nil
	}
	writer.alignToBoundary()
}

func (writer *Writer) alignToBoundary() {
	extraBytes := writer.encoder.CurPos() % boundarySize
	if extraBytes > 0 {
		padding := make([]byte, boundarySize-extraBytes)
		writer.encoder.Code(padding)
	}
}

func (writer *Writer) dataPaddingSizeForFragmentedChunk(id Identifier) (padding int) {
	// Some directories have a 2byte padding before the actual data
	idValue := id.Value()
	if (idValue >= 0x08FC) && (idValue <= 0x094B) { // all chunks in obj3d.res
		padding = 2
	}
	return
}
