package store

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/objprop"
)

type ModifiedCallback func()
type dataRetriever func() objprop.ObjectData

type ProviderBacked struct {
	provider   objprop.Provider
	onModified ModifiedCallback

	retriever map[res.ObjectID]dataRetriever
}

func NewProviderBacked(provider objprop.Provider, onModified ModifiedCallback) *ProviderBacked {
	backed := &ProviderBacked{
		provider:   provider,
		onModified: onModified,
		retriever:  make(map[res.ObjectID]dataRetriever)}

	return backed
}

// Get returns the data for the requested ObjectID.
func (backed *ProviderBacked) Get(id res.ObjectID) objprop.ObjectData {
	retriever, existing := backed.retriever[id]
	if !existing {
		return backed.provider.Provide(id)
	} else {
		return retriever()
	}
}

// Put takes the provided data and associates it with the given ID.
func (backed *ProviderBacked) Put(id res.ObjectID, data objprop.ObjectData) {
	backed.retriever[id] = func() objprop.ObjectData { return data }
	backed.onModified()
}
