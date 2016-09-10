package logic

import (
	"github.com/inkyblackness/res/data"
)

// LevelObjectChainLink is one entry of a chain of objects.
type LevelObjectChainLink interface {
	// NextIndex returns the index of the next object.
	NextIndex() data.LevelObjectChainIndex
	// SetNextIndex sets the index of the next object.
	SetNextIndex(index data.LevelObjectChainIndex)

	// PreviousIndex returns the index of the previous object.
	PreviousIndex() data.LevelObjectChainIndex
	// SetPreviousIndex sets the index of the previous object.
	SetPreviousIndex(index data.LevelObjectChainIndex)
}

// LevelObjectChainStart is the first entry of object chains.
// The previous entries are unused links.
type LevelObjectChainStart interface {
	LevelObjectChainLink

	// ReferenceIndex returns the index of a referenced object.
	ReferenceIndex() data.LevelObjectChainIndex
	// SetReferenceIndex sets the index of a referenced object.
	SetReferenceIndex(index data.LevelObjectChainIndex)
}
