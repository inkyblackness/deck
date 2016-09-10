package io

import (
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/objprop"
)

// StoreLibrary wraps the methods to contain stores for various data
type StoreLibrary interface {
	// ChunkStore returns a chunk store for given name.
	ChunkStore(name string) (chunk.Store, error)

	// ObjpropStore returns a object properties store for given name.
	ObjpropStore(name string) (objprop.Store, error)
}
