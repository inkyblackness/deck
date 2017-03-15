package model

import (
	"sort"
	"strconv"

	"github.com/inkyblackness/shocked-model"
)

// LevelAdapter is the entry point for a level.
type LevelAdapter struct {
	context archiveContext
	store   model.DataStore

	id           *observable
	isCyberspace bool
	tileMap      *TileMap

	levelTextures *observable

	levelObjects *observable
}

func newLevelAdapter(context archiveContext, store model.DataStore) *LevelAdapter {
	adapter := &LevelAdapter{
		context: context,
		store:   store,

		id:      newObservable(),
		tileMap: NewTileMap(64, 64),

		levelTextures: newObservable(),
		levelObjects:  newObservable()}

	adapter.id.set("")

	return adapter
}

// ID returns the ID of the level.
func (adapter *LevelAdapter) ID() string {
	return adapter.id.orDefault("").(string)
}

func (adapter *LevelAdapter) storeLevelID() int {
	idAsString := adapter.ID()
	id := -1

	if idAsString != "" {
		parsed, _ := strconv.ParseInt(idAsString, 10, 16)
		id = int(parsed)
	}

	return id
}

// OnIDChanged registers a callback for changed IDs.
func (adapter *LevelAdapter) OnIDChanged(callback func()) {
	adapter.id.addObserver(callback)
}

func (adapter *LevelAdapter) requestByID(levelID string) {
	adapter.id.set("")
	adapter.tileMap.clear()
	textures := []int{}
	adapter.levelTextures.set(&textures)
	objects := make(map[int]*LevelObject)
	adapter.levelObjects.set(&objects)

	adapter.id.set(levelID)
	if levelID != "" {
		storeLevelID := adapter.storeLevelID()
		adapter.store.Tiles(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), storeLevelID,
			adapter.onTiles, adapter.context.simpleStoreFailure("Tiles"))
		adapter.store.LevelTextures(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), storeLevelID,
			adapter.onLevelTextures, adapter.context.simpleStoreFailure("LevelTextures"))
		adapter.store.LevelObjects(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), storeLevelID,
			adapter.onLevelObjects, adapter.context.simpleStoreFailure("LevelObjects"))
	}
}

// IsCyberspace returns true for cyberspace levels.
func (adapter *LevelAdapter) IsCyberspace() bool {
	return adapter.isCyberspace
}

// TileMap returns the map of tiles of the level.
func (adapter *LevelAdapter) TileMap() *TileMap {
	return adapter.tileMap
}

func (adapter *LevelAdapter) onTiles(tiles model.Tiles) {
	for y := 0; y < len(tiles.Table); y++ {
		row := tiles.Table[y]
		for x := 0; x < len(row); x++ {
			adapter.onTilePropertiesUpdated(TileCoordinateOf(x, y), &row[x].Properties)
		}
	}
}

func (adapter *LevelAdapter) onTilePropertiesUpdated(coord TileCoordinate, properties *model.TileProperties) {
	tile := adapter.tileMap.Tile(coord)
	tile.setProperties(properties)
}

// LevelTextureIDs returns the IDs of the level textures.
func (adapter *LevelAdapter) LevelTextureIDs() []int {
	return *adapter.levelTextures.get().(*[]int)
}

// LevelTextureID returns the texture ID for given level index.
func (adapter *LevelAdapter) LevelTextureID(index int) (id int) {
	ids := adapter.LevelTextureIDs()
	if index < len(ids) {
		id = ids[index]
	} else {
		id = -1
	}

	return
}

func (adapter *LevelAdapter) onLevelTextures(textureIDs []int) {
	adapter.levelTextures.set(&textureIDs)
}

// OnLevelTexturesChanged registers for updates of the level textures.
func (adapter *LevelAdapter) OnLevelTexturesChanged(callback func()) {
	adapter.levelTextures.addObserver(callback)
}

// RequestLevelTexturesChange requests to change the level textures list
func (adapter *LevelAdapter) RequestLevelTexturesChange(textureIDs []int) {
	adapter.store.SetLevelTextures(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), adapter.storeLevelID(),
		textureIDs, adapter.onLevelTextures, adapter.context.simpleStoreFailure("SetLevelTextures"))
}

// LevelObjects returns a sorted set of objects that match the provided filter.
func (adapter *LevelAdapter) LevelObjects(filter func(*LevelObject) bool) []*LevelObject {
	objects := *adapter.levelObjects.get().(*map[int]*LevelObject)
	indexList := make([]int, 0, len(objects))

	for key, obj := range objects {
		if filter(obj) {
			indexList = append(indexList, key)
		}
	}
	sort.Ints(indexList)
	result := make([]*LevelObject, len(indexList))
	for index, key := range indexList {
		result[index] = objects[key]
	}

	return result
}

// OnLevelObjectsChanged registers a callback for updates on the list of level objects.
func (adapter *LevelAdapter) OnLevelObjectsChanged(callback func()) {
	adapter.levelObjects.addObserver(callback)
}

func (adapter *LevelAdapter) onLevelObjects(objects *model.LevelObjects) {
	newMap := make(map[int]*LevelObject)
	for tableIndex := 0; tableIndex < len(objects.Table); tableIndex++ {
		obj := newLevelObject(&objects.Table[tableIndex])
		newMap[obj.Index()] = obj
	}
	adapter.levelObjects.set(&newMap)
}

// RequestTilePropertyChange requests the tiles at given coordinates to set provided properties.
func (adapter *LevelAdapter) RequestTilePropertyChange(coordinates []TileCoordinate, properties *model.TileProperties) {
	additionalQueries := make(map[TileCoordinate]bool)
	storeLevelID := adapter.storeLevelID()
	tileUpdateHandler := func(coord TileCoordinate) func(model.TileProperties) {
		return func(newProperties model.TileProperties) {
			adapter.onTilePropertiesUpdated(coord, &newProperties)
		}
	}
	for _, coord := range coordinates {
		x, y := coord.XY()
		additionalQueries[TileCoordinateOf(x-1, y)] = true
		additionalQueries[TileCoordinateOf(x+1, y)] = true
		additionalQueries[TileCoordinateOf(x, y-1)] = true
		additionalQueries[TileCoordinateOf(x, y+1)] = true
		adapter.store.SetTile(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), storeLevelID,
			x, y, *properties,
			tileUpdateHandler(coord), adapter.context.simpleStoreFailure("SetTile"))
	}
	for _, coord := range coordinates {
		delete(additionalQueries, coord)
	}
	for coord := range additionalQueries {
		x, y := coord.XY()
		if (x >= 0) && (x < 64) && (y >= 0) && (y < 64) {
			adapter.store.Tile(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), storeLevelID,
				x, y, tileUpdateHandler(coord), adapter.context.simpleStoreFailure("Tile"))
		}
	}
}
