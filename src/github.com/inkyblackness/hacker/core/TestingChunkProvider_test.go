package core

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
)

type TestingChunkProvider struct {
	chunkIDs   []res.ResourceID
	chunksByID map[res.ResourceID]chunk.BlockHolder
}

func NewTestingChunkProvider() *TestingChunkProvider {
	provider := &TestingChunkProvider{
		chunksByID: make(map[res.ResourceID]chunk.BlockHolder)}

	return provider
}

// IDs returns a list of available chunk IDs
func (provider *TestingChunkProvider) IDs() []res.ResourceID {
	return provider.chunkIDs
}

// Consume adds the given chunk to the consumer.
func (provider *TestingChunkProvider) Consume(id res.ResourceID, holder chunk.BlockHolder) {
	if _, existing := provider.chunksByID[id]; !existing {
		provider.chunkIDs = append(provider.chunkIDs, id)
	}
	provider.chunksByID[id] = holder
}

// Provide returns the chunk for given ID if known.
func (provider *TestingChunkProvider) Provide(id res.ResourceID) (holder chunk.BlockHolder) {
	temp, existing := provider.chunksByID[id]

	if existing {
		holder = temp
	}
	return
}
