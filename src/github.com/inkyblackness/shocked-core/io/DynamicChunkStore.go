package io

import (
	"sync"

	"github.com/inkyblackness/res/chunk"
)

// DynamicChunkStore wraps a chunk store instance that can be replaced during runtime.
// Access to the underlying store is synchronized.
type DynamicChunkStore struct {
	mutex          sync.Mutex
	wrapped        chunk.Store
	changeCallback func()
}

// NewDynamicChunkStore returns a new instance of a dynamic chunk store.
func NewDynamicChunkStore(wrapped chunk.Store, changeCallback func()) *DynamicChunkStore {
	store := &DynamicChunkStore{
		wrapped:        wrapped,
		changeCallback: changeCallback}

	return store
}

// Swap replaces the currently wrapped store with a new one with the help of a factory function.
// The factory function receives the old store and must return the new instance to wrap.
func (store *DynamicChunkStore) Swap(factory func(oldStore chunk.Store) chunk.Store) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.wrapped = factory(store.wrapped)
}

// IDs returns a list of available IDs this store currently contains.
func (store *DynamicChunkStore) IDs() []chunk.Identifier {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	return store.wrapped.IDs()
}

// Get returns a chunk for the given identifier.
func (store *DynamicChunkStore) Get(id chunk.Identifier) (blockStore *DynamicBlockStore) {
	retriever := func(handler func(*chunk.Chunk)) {
		store.mutex.Lock()
		defer store.mutex.Unlock()
		retrieved, retrievedErr := store.wrapped.Chunk(id)
		if retrievedErr != nil {
			return
		}
		handler(retrieved)
	}
	exists := false

	retriever(func(existing *chunk.Chunk) {
		exists = existing != nil
	})
	if exists {
		blockStore = newDynamicBlockStore(retriever, store.changeCallback)
	}

	return blockStore
}

// Del ensures the chunk with given identifier is deleted.
func (store *DynamicChunkStore) Del(id chunk.Identifier) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.wrapped.Del(id)
	store.changeCallback()
}

// Put (re-)assigns an identifier with data. If no chunk with given ID exists,
// then it is created. Existing chunks are overwritten with the provided data.
func (store *DynamicChunkStore) Put(id chunk.Identifier, newChunk *chunk.Chunk) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.wrapped.Put(id, newChunk)
	store.changeCallback()
}
