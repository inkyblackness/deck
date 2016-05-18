package store

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
)

type ModifiedCallback func()
type chunkRetriever func() chunk.BlockStore

func nullRetriever() chunk.BlockStore {
	return nil
}

// ProviderBacked is a chunk store allowing in-memory modification of chunks and their data.
type ProviderBacked struct {
	provider   chunk.Provider
	onModified ModifiedCallback

	ids    []res.ResourceID
	chunks map[res.ResourceID]chunkRetriever
}

// NewProviderBacked returns a new instance of a ProviderBacked chunk store.
// The onModified callback is called for any modification in the store.
func NewProviderBacked(provider chunk.Provider, onModified ModifiedCallback) *ProviderBacked {
	existingIDs := provider.IDs()
	backed := &ProviderBacked{
		provider:   provider,
		onModified: onModified,
		ids:        make([]res.ResourceID, len(existingIDs)),
		chunks:     make(map[res.ResourceID]chunkRetriever)}

	makeRetriever := func(id res.ResourceID) func() chunk.BlockStore {
		var cachedStore chunk.BlockStore

		return func() chunk.BlockStore {
			if cachedStore == nil {
				cachedStore = newBackedBlockStore(provider.Provide(id), backed.onModified)
			}

			return cachedStore
		}
	}

	for index, id := range existingIDs {
		backed.ids[index] = id
		backed.chunks[id] = makeRetriever(id)
	}

	return backed
}

// IDs returns a list of available IDs this store currently contains.
func (backed *ProviderBacked) IDs() []res.ResourceID {
	return backed.ids
}

// Del ensures the chunk with given identifier is deleted.
func (backed *ProviderBacked) Del(id res.ResourceID) {
	modified := false
	newIDs := make([]res.ResourceID, 0, len(backed.ids))

	for _, existingID := range backed.ids {
		if existingID != id {
			newIDs = append(newIDs, existingID)
		} else {
			delete(backed.chunks, id)
			modified = true
		}
	}
	backed.ids = newIDs
	if modified {
		backed.onModified()
	}
}

// Put (re-)assigns an identifier with data. If no chunk with given ID exists,
// then it is created. Existing chunks are overwritten with the provided data.
func (backed *ProviderBacked) Put(id res.ResourceID, holder chunk.BlockHolder) {
	blockStore := newBackedBlockStore(holder, func() {})

	if _, existing := backed.chunks[id]; !existing {
		backed.ids = append(backed.ids, id)
	}
	backed.chunks[id] = func() chunk.BlockStore { return blockStore }
	backed.onModified()
}

// Get returns a chunk for the given identifier.
func (backed *ProviderBacked) Get(id res.ResourceID) chunk.BlockStore {
	retriever, existing := backed.chunks[id]
	if !existing {
		retriever = nullRetriever
	}

	return retriever()
}
