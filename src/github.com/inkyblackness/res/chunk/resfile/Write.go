package resfile

import (
	"fmt"
	"io"

	"github.com/inkyblackness/res/chunk"
)

// Write serializes the chunks from given provider into the target.
// It is a convenience function for using Writer.
func Write(target io.WriteSeeker, source chunk.Provider) error {
	writer, writerErr := NewWriter(target)
	if writerErr != nil {
		return fmt.Errorf("failed to create writer: %v", writerErr)
	}

	for _, id := range source.IDs() {
		entry, chunkErr := source.Chunk(id)
		if chunkErr != nil {
			return fmt.Errorf("failed to retrieve chunk %v: %v", id, chunkErr)
		}

		if entry.Fragmented {
			chunkWriter, chunkWriterErr := writer.CreateFragmentedChunk(id, entry.ContentType, entry.Compressed)
			if chunkWriterErr != nil {
				return fmt.Errorf("failed to create chunk %v: %v", id, chunkWriterErr)
			}
			copyErr := copyBlocks(entry, func() io.Writer { return chunkWriter.CreateBlock() })
			if copyErr != nil {
				return fmt.Errorf("failed to copy chunk %v: %v", id, copyErr)
			}
		} else if entry.BlockCount() == 1 {
			blockWriter, chunkWriterErr := writer.CreateChunk(id, entry.ContentType, entry.Compressed)
			if chunkWriterErr != nil {
				return fmt.Errorf("failed to create chunk %v, %v", id, chunkWriterErr)
			}
			copyErr := copyBlocks(entry, func() io.Writer { return blockWriter })
			if copyErr != nil {
				return fmt.Errorf("failed to copy chunk %v: %v", id, copyErr)
			}
		} else {
			return fmt.Errorf("unfragmented chunk %v has wrong number of blocks", id)
		}
	}

	finishErr := writer.Finish()
	if finishErr != nil {
		return fmt.Errorf("failed to finish writer: %v", finishErr)
	}
	return nil
}

func copyBlocks(source chunk.BlockProvider, nextWriter func() io.Writer) error {
	for blockIndex := 0; blockIndex < source.BlockCount(); blockIndex++ {
		blockReader, blockErr := source.Block(blockIndex)
		if blockErr != nil {
			return fmt.Errorf("failed to retrieve block %d: %v", blockIndex, blockErr)
		}
		_, copyErr := io.Copy(nextWriter(), blockReader)
		if copyErr != nil {
			return fmt.Errorf("failed to copy data %d: %v", blockIndex, copyErr)
		}
	}
	return nil
}
