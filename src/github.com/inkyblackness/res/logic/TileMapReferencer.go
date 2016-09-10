package logic

// TileMapReferencer is an interface for keeping cross-references in a tile map.
type TileMapReferencer interface {
	// ReferenceIndex returns the index of the first cross-reference for the given tile.
	ReferenceIndex(location TileLocation) CrossReferenceListIndex
	// SetReferenceIndex sets the index of the first cross-reference for the given tile.
	SetReferenceIndex(location TileLocation, index CrossReferenceListIndex)
}
