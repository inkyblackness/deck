package store

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
)

type resourceResource struct {
	chunkIDs   []res.ResourceID
	chunksByID map[res.ResourceID]chunk.BlockHolder
}

func NewTestingResource() *resourceResource {
	resource := &resourceResource{
		chunksByID: make(map[res.ResourceID]chunk.BlockHolder)}

	return resource
}

// IDs returns a list of available chunk IDs
func (resource *resourceResource) IDs() []res.ResourceID {
	return resource.chunkIDs
}

// Consume adds the given chunk to the consumer.
func (resource *resourceResource) Consume(id res.ResourceID, holder chunk.BlockHolder) {
	if _, existing := resource.chunksByID[id]; !existing {
		resource.chunkIDs = append(resource.chunkIDs, id)
	}
	resource.chunksByID[id] = holder
}

// Provide returns the chunk for given ID if known.
func (resource *resourceResource) Provide(id res.ResourceID) (holder chunk.BlockHolder) {
	temp, existing := resource.chunksByID[id]

	if existing {
		holder = temp
	}
	return
}
