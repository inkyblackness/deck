package io

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
)

type dynamicBlockRetriever func(func(chunk.BlockStore))

type dynamicBlockStore struct {
	retrieve dynamicBlockRetriever
}

func newDynamicBlockStore(retriever dynamicBlockRetriever) *dynamicBlockStore {
	return &dynamicBlockStore{retrieve: retriever}
}

// Type returns the type of the chunk.
func (store *dynamicBlockStore) ChunkType() (result chunk.TypeID) {
	store.retrieve(func(wrapped chunk.BlockStore) {
		result = wrapped.ChunkType()
	})

	return
}

// ContentType returns the type of the data.
func (store *dynamicBlockStore) ContentType() (result res.DataTypeID) {
	store.retrieve(func(wrapped chunk.BlockStore) {
		result = wrapped.ContentType()
	})

	return
}

// BlockCount returns the number of blocks available in the chunk.
// Flat chunks must contain exactly one block.
func (store *dynamicBlockStore) BlockCount() (result uint16) {
	store.retrieve(func(wrapped chunk.BlockStore) {
		result = wrapped.BlockCount()
	})

	return
}

// BlockData returns the data for the requested block index.
func (store *dynamicBlockStore) BlockData(block uint16) (result []byte) {
	store.retrieve(func(wrapped chunk.BlockStore) {
		result = wrapped.BlockData(block)
	})

	return
}

// SetBlockData sets the data for the requested block index.
func (store *dynamicBlockStore) SetBlockData(block uint16, data []byte) {
	store.retrieve(func(wrapped chunk.BlockStore) {
		wrapped.SetBlockData(block, data)
	})
}
