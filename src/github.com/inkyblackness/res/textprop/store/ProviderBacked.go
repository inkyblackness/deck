package store

import (
	"github.com/inkyblackness/res/textprop"
)

type ModifiedCallback func()
type dataRetriever func() []byte

type ProviderBacked struct {
	provider   textprop.Provider
	onModified ModifiedCallback

	retriever map[uint32]dataRetriever
}

func NewProviderBacked(provider textprop.Provider, onModified ModifiedCallback) *ProviderBacked {
	backed := &ProviderBacked{
		provider:   provider,
		onModified: onModified,
		retriever:  make(map[uint32]dataRetriever)}

	return backed
}

// EntryCount returns the amount of entries available
func (backed *ProviderBacked) EntryCount() uint32 {
	return backed.provider.EntryCount()
}

// Get returns the data for the requested ID.
func (backed *ProviderBacked) Get(id uint32) []byte {
	retriever, existing := backed.retriever[id]
	if !existing {
		return backed.provider.Provide(id)
	} else {
		return retriever()
	}
}

// Put takes the provided data and associates it with the given ID.
func (backed *ProviderBacked) Put(id uint32, data []byte) {
	backed.retriever[id] = func() []byte { return data }
	backed.onModified()
}
