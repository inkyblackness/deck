package objprop

import "github.com/inkyblackness/res"

// Store represents a dynamically accessible container of object properties.
type Store interface {
	// Get returns the data for the requested ObjectID.
	Get(id res.ObjectID) ObjectData

	// Put takes the provided data and associates it with the given ID.
	Put(id res.ObjectID, data ObjectData)
}
