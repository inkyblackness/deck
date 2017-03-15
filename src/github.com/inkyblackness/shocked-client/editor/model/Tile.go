package model

import (
	"github.com/inkyblackness/shocked-model"
)

// Tile keeps properties about one map tile.
type Tile struct {
	properties *observable
}

// NewTile returns a new tile instance.
func NewTile() *Tile {
	tile := &Tile{properties: newObservable()}

	return tile
}

// SetProperties sets the current tile properties.
func (tile *Tile) setProperties(properties *model.TileProperties) {
	tile.properties.set(properties)
}

// Properties returns the current tile properties.
func (tile *Tile) Properties() *model.TileProperties {
	return tile.properties.get().(*model.TileProperties)
}

// OnPropertiesChanged registers a callback for change updates.
func (tile *Tile) OnPropertiesChanged(callback func()) {
	tile.properties.addObserver(callback)
}
