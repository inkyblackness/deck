package objprop

import "github.com/inkyblackness/res"

// Provider wraps the Provide method.
type Provider interface {
	// Provide returns the data for the requested ObjectID.
	Provide(id res.ObjectID) ObjectData
}
