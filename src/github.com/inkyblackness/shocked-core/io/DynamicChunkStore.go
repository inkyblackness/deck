package io

import (
	"sync"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
)

// DynamicChunkStore wraps a chunk store instance that can be replaced during runtime.
// Access to the underlying store is synchronized.
type DynamicChunkStore struct {
	mutex   sync.Mutex
	wrapped chunk.Store
}

func NewDynamicChunkStore(wrapped chunk.Store) *DynamicChunkStore {
	store := &DynamicChunkStore{
		wrapped: wrapped}

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
func (store *DynamicChunkStore) IDs() []res.ResourceID {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	return store.wrapped.IDs()
}

// Get returns a chunk for the given identifier.
func (store *DynamicChunkStore) Get(id res.ResourceID) (blockStore chunk.BlockStore) {
	retriever := func(handler func(chunk.BlockStore)) {
		store.mutex.Lock()
		defer store.mutex.Unlock()

		handler(store.wrapped.Get(id))
	}
	exists := false

	retriever(func(existing chunk.BlockStore) {
		exists = existing != nil
	})
	if exists {
		blockStore = newDynamicBlockStore(retriever)
	}

	return blockStore
}

// Del ensures the chunk with given identifier is deleted.
func (store *DynamicChunkStore) Del(id res.ResourceID) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.wrapped.Del(id)
}

// Put (re-)assigns an identifier with data. If no chunk with given ID exists,
// then it is created. Existing chunks are overwritten with the provided data.
func (store *DynamicChunkStore) Put(id res.ResourceID, holder chunk.BlockHolder) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.wrapped.Put(id, holder)
}
