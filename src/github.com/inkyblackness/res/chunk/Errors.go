package chunk

import "fmt"

// ErrChunkDoesNotExist returns an error specifying the given ID doesn't
// have an associated chunk.
func ErrChunkDoesNotExist(id Identifier) error {
	return fmt.Errorf("chunk with ID %v does not exist", id)
}
