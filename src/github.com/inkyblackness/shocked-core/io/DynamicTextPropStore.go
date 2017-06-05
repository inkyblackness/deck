package io

import (
	"sync"

	"github.com/inkyblackness/res/textprop"
)

// DynamicTextPropStore wraps a object properties store instance that can be replaced during runtime.
// Access to the underlying store is synchronized.
type DynamicTextPropStore struct {
	mutex   sync.Mutex
	wrapped textprop.Store
}

// NewDynamicTextPropStore returns a new instance of a DynamicTextPropStore wrapping the given store
func NewDynamicTextPropStore(wrapped textprop.Store) *DynamicTextPropStore {
	store := &DynamicTextPropStore{
		wrapped: wrapped}

	return store
}

// Swap replaces the currently wrapped store with a new one with the help of a factory function.
// The factory function receives the old store and must return the new instance to wrap.
func (store *DynamicTextPropStore) Swap(factory func(oldStore textprop.Store) textprop.Store) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.wrapped = factory(store.wrapped)
}

// EntryCount returns the number of currently stored properties.
func (store *DynamicTextPropStore) EntryCount() uint32 {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	return store.wrapped.EntryCount()
}

// Get returns the data for the requested id.
func (store *DynamicTextPropStore) Get(id uint32) []byte {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	return store.wrapped.Get(id)
}

// Put takes the provided data and associates it with the given ID, overwriting the previous data.
func (store *DynamicTextPropStore) Put(id uint32, data []byte) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.wrapped.Put(id, data)
}
