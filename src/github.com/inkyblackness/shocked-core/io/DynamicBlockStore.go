package io

import (
	"io/ioutil"

	"github.com/inkyblackness/res/chunk"
)

type dynamicBlockRetriever func(func(*chunk.Chunk))

// DynamicBlockStore provides access to the blocks of a dynamic chunk store.
type DynamicBlockStore struct {
	retrieve       dynamicBlockRetriever
	changeCallback func()
}

func newDynamicBlockStore(retriever dynamicBlockRetriever, changeCallback func()) *DynamicBlockStore {
	return &DynamicBlockStore{retrieve: retriever, changeCallback: changeCallback}
}

// ContentType returns the type of the data.
func (store *DynamicBlockStore) ContentType() (result chunk.ContentType) {
	store.retrieve(func(wrapped *chunk.Chunk) {
		result = wrapped.ContentType
	})

	return
}

// BlockCount returns the number of blocks available in the chunk.
// Flat chunks must contain exactly one block.
func (store *DynamicBlockStore) BlockCount() (result uint16) {
	store.retrieve(func(wrapped *chunk.Chunk) {
		result = uint16(wrapped.BlockCount())
	})

	return
}

// BlockData returns the data for the requested block index.
func (store *DynamicBlockStore) BlockData(block uint16) (result []byte) {
	store.retrieve(func(wrapped *chunk.Chunk) {
		blockReader, blockErr := wrapped.Block(int(block))
		if blockErr != nil {
			return
		}
		result, _ = ioutil.ReadAll(blockReader) // nolint: gas
	})

	return
}

// SetBlockData sets the data for the requested block index.
func (store *DynamicBlockStore) SetBlockData(block uint16, data []byte) {
	store.retrieve(func(wrapped *chunk.Chunk) {
		wrapped.SetBlock(int(block), data)
		store.changeCallback()
	})
}
