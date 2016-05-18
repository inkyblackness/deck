package io

import (
	"sync"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/objprop"
)

// DynamicObjPropStore wraps a object properties store instance that can be replaced during runtime.
// Access to the underlying store is synchronized.
type DynamicObjPropStore struct {
	mutex   sync.Mutex
	wrapped objprop.Store
}

// NewDynamicObjPropStore returns a new instance of a DynamicObjPropStore wrapping the given store
func NewDynamicObjPropStore(wrapped objprop.Store) *DynamicObjPropStore {
	store := &DynamicObjPropStore{
		wrapped: wrapped}

	return store
}

// Swap replaces the currently wrapped store with a new one with the help of a factory function.
// The factory function receives the old store and must return the new instance to wrap.
func (store *DynamicObjPropStore) Swap(factory func(oldStore objprop.Store) objprop.Store) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.wrapped = factory(store.wrapped)
}

// Get returns the data for the requested ObjectID.
func (store *DynamicObjPropStore) Get(id res.ObjectID) objprop.ObjectData {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	return store.wrapped.Get(id)
}

// Put takes the provided data and associates it with the given ID, overwriting the previous data.
func (store *DynamicObjPropStore) Put(id res.ObjectID, data objprop.ObjectData) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.wrapped.Put(id, data)
}
