package io

import (
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/objprop"
	"github.com/inkyblackness/res/textprop"
)

// StoreLibrary wraps the methods to contain stores for various data
type StoreLibrary interface {
	// ChunkStore returns a chunk store for given name.
	ChunkStore(name string) (chunk.Store, error)

	// ObjpropStore returns an object properties store for given name.
	ObjpropStore(name string) (objprop.Store, error)

	// TextpropStore returns a texture properties store for given name.
	TextpropStore(name string) (textprop.Store, error)
}
