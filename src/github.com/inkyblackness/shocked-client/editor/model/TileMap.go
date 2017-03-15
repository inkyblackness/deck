package model

// TileMap contains the tiles with their properties.
type TileMap struct {
	tiles map[TileCoordinate]*Tile
}

// NewTileMap returns a new tile map instance
func NewTileMap(width, height int) *TileMap {
	tileMap := &TileMap{
		tiles: make(map[TileCoordinate]*Tile)}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			coord := TileCoordinateOf(x, y)
			tileMap.tiles[coord] = NewTile()
		}
	}

	return tileMap
}

func (tileMap *TileMap) clear() {
	for _, tile := range tileMap.tiles {
		tile.setProperties(nil)
	}
}

// Tile returns the tile at the given coordinate
func (tileMap *TileMap) Tile(coord TileCoordinate) *Tile {
	return tileMap.tiles[coord]
}
