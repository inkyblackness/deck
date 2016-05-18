package io

import (
	"github.com/inkyblackness/res/chunk"
)

// StoreLibrary wraps the methods to contain stores for various data
type StoreLibrary interface {
	// ChunkStore returns a chunk store for given name.
	ChunkStore(name string) (chunk.Store, error)
}
