package model

import (
	"github.com/inkyblackness/shocked-model"
)

// Tile keeps properties about one map tile.
type Tile struct {
	properties *model.TileProperties
}

// NewTile returns a new tile instance.
func NewTile() *Tile {
	return &Tile{}
}

// SetProperties sets the current tile properties.
func (tile *Tile) SetProperties(properties *model.TileProperties) {
	tile.properties = properties
}

// Properties returns the current tile properties.
func (tile *Tile) Properties() *model.TileProperties {
	return tile.properties
}
