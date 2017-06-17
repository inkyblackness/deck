package model

import (
	"fmt"
	"sort"

	"github.com/inkyblackness/shocked-model"
)

// LevelAdapter is the entry point for a level.
type LevelAdapter struct {
	context        archiveContext
	store          model.DataStore
	objectsAdapter *ObjectsAdapter

	id      *observable
	tileMap *TileMap

	levelProperties *observable

	levelTextures               *observable
	levelTextureAnimationGroups *observable

	levelObjects *observable

	levelSurveillance *observable
}

func newLevelAdapter(context archiveContext, store model.DataStore, objectsAdapter *ObjectsAdapter) *LevelAdapter {
	adapter := &LevelAdapter{
		context:        context,
		store:          store,
		objectsAdapter: objectsAdapter,

		id:      newObservable(),
		tileMap: NewTileMap(64, 64),

		levelProperties:             newObservable(),
		levelTextures:               newObservable(),
		levelTextureAnimationGroups: newObservable(),
		levelObjects:                newObservable(),
		levelSurveillance:           newObservable()}

	adapter.id.set(-1)

	return adapter
}

// ID returns the ID of the level.
func (adapter *LevelAdapter) ID() int {
	return adapter.id.orDefault(-1).(int)
}

// OnIDChanged registers a callback for changed IDs.
func (adapter *LevelAdapter) OnIDChanged(callback func()) {
	adapter.id.addObserver(callback)
}

func (adapter *LevelAdapter) requestByID(levelID int) {
	adapter.id.set(-1)
	adapter.tileMap.clear()
	adapter.levelProperties.set(nil)
	textures := []int{}
	adapter.levelTextures.set(&textures)
	animationGroups := []*LevelTextureAnimationGroup{}
	adapter.levelTextureAnimationGroups.set(&animationGroups)
	objects := make(map[int]*LevelObject)
	adapter.levelObjects.set(&objects)
	objectIndices := []model.SurveillanceObject{}
	adapter.levelSurveillance.set(&objectIndices)

	adapter.id.set(levelID)
	if levelID >= 0 {
		storeLevelID := adapter.ID()
		adapter.store.LevelProperties(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), storeLevelID,
			adapter.onLevelProperties, adapter.context.simpleStoreFailure("LevelProperties"))
		adapter.store.Tiles(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), storeLevelID,
			adapter.onTiles, adapter.context.simpleStoreFailure("Tiles"))
		adapter.store.LevelTextures(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), storeLevelID,
			adapter.onLevelTextures, adapter.context.simpleStoreFailure("LevelTextures"))
		adapter.store.LevelTextureAnimations(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), storeLevelID,
			adapter.onLevelTextureAnimations, adapter.context.simpleStoreFailure("LevelTextureAnimations"))
		adapter.store.LevelObjects(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), storeLevelID,
			adapter.onLevelObjects, adapter.context.simpleStoreFailure("LevelObjects"))
		adapter.store.LevelSurveillanceObjects(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), storeLevelID,
			adapter.onLevelSurveillance, adapter.context.simpleStoreFailure("LevelSurveillance"))
	}
}

func (adapter *LevelAdapter) onLevelProperties(properties model.LevelProperties) {
	adapter.levelProperties.set(&properties)
}

func (adapter *LevelAdapter) properties() (properties *model.LevelProperties) {
	value := adapter.levelProperties.get()

	if value != nil {
		properties = value.(*model.LevelProperties)
	}
	return
}

// RequestLevelPropertiesChange requests to change the properties of the level.
func (adapter *LevelAdapter) RequestLevelPropertiesChange(modifier func(properties *model.LevelProperties)) {
	var properties model.LevelProperties

	modifier(&properties)
	adapter.store.SetLevelProperties(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), adapter.ID(),
		properties, adapter.onLevelProperties, adapter.context.simpleStoreFailure("SetLevelProperties"))
}

// OnLevelPropertiesChanged registers for updates of the level properties.
func (adapter *LevelAdapter) OnLevelPropertiesChanged(callback func()) {
	adapter.levelProperties.addObserver(callback)
}

// IsCyberspace returns true for cyberspace levels.
func (adapter *LevelAdapter) IsCyberspace() (result bool) {
	if properties := adapter.properties(); properties != nil {
		result = *properties.CyberspaceFlag
	}
	return
}

// HeightShift returns the height shift value, or -1 if not known.
func (adapter *LevelAdapter) HeightShift() int {
	result := -1
	if properties := adapter.properties(); properties != nil {
		result = *properties.HeightShift
	}
	return result
}

// CeilingEffect returns whether radiation is enabled and what its level is.
func (adapter *LevelAdapter) CeilingEffect() (radiation bool, level int) {
	level = -1
	if properties := adapter.properties(); properties != nil {
		radiation = *properties.CeilingHasRadiation
		level = *properties.CeilingEffectLevel
	}
	return
}

// FloorEffect returns whether biohazard is enabled, floor has gravity, and what the level is.
func (adapter *LevelAdapter) FloorEffect() (biohazard bool, gravity bool, level int) {
	level = -1
	if properties := adapter.properties(); properties != nil {
		biohazard = *properties.FloorHasBiohazard
		gravity = *properties.FloorHasGravity
		level = *properties.FloorEffectLevel
	}
	return
}

// TileMap returns the map of tiles of the level.
func (adapter *LevelAdapter) TileMap() *TileMap {
	return adapter.tileMap
}

func (adapter *LevelAdapter) onTiles(tiles model.Tiles) {
	for y := 0; y < len(tiles.Table); y++ {
		row := tiles.Table[y]
		for x := 0; x < len(row); x++ {
			adapter.onTilePropertiesUpdated(TileCoordinateOf(x, y), &row[x])
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
	adapter.store.SetLevelTextures(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), adapter.ID(),
		textureIDs, adapter.onLevelTextures, adapter.context.simpleStoreFailure("SetLevelTextures"))
}

func (adapter *LevelAdapter) onLevelTextureAnimations(animations []model.TextureAnimation) {
	count := len(animations)
	groups := make([]*LevelTextureAnimationGroup, count)
	for index, animation := range animations {
		group := newLevelTextureAnimationGroup(index)
		groups[index] = group
		group.properties = animation
	}
	adapter.levelTextureAnimationGroups.set(&groups)
}

// OnLevelTextureAnimationsChanged registers for updates of the animation groups.
func (adapter *LevelAdapter) OnLevelTextureAnimationsChanged(callback func()) {
	adapter.levelTextureAnimationGroups.addObserver(callback)
}

func (adapter *LevelAdapter) levelTextureAnimationGroupList() []*LevelTextureAnimationGroup {
	return *adapter.levelTextureAnimationGroups.get().(*[]*LevelTextureAnimationGroup)
}

// TextureAnimationGroupCount returns how many groups are available.
func (adapter *LevelAdapter) TextureAnimationGroupCount() int {
	return len(adapter.levelTextureAnimationGroupList())
}

// TextureAnimationGroup returns the group instance for given group ID.
func (adapter *LevelAdapter) TextureAnimationGroup(id int) *LevelTextureAnimationGroup {
	return adapter.levelTextureAnimationGroupList()[id]
}

// RequestLevelTextureAnimationGroupChange requests the change properties of a texture animation group.
func (adapter *LevelAdapter) RequestLevelTextureAnimationGroupChange(id int, properties model.TextureAnimation) {
	adapter.store.SetLevelTextureAnimation(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), adapter.ID(),
		id, properties, adapter.onLevelTextureAnimations, adapter.context.simpleStoreFailure("SetLevelTextureAnimation"))
}

func (adapter *LevelAdapter) levelObjectsMap() map[int]*LevelObject {
	return *adapter.levelObjects.get().(*map[int]*LevelObject)
}

// LevelObjects returns a sorted set of objects that match the provided filter.
func (adapter *LevelAdapter) LevelObjects(filter func(*LevelObject) bool) []*LevelObject {
	objects := adapter.levelObjectsMap()
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

// RequestNewObject requests to add a new object at the given coordinate.
func (adapter *LevelAdapter) RequestNewObject(worldX, worldY float32, objectID ObjectID) {
	integerX, integerY := int(worldX), int(worldY)
	tileX, fineX := integerX>>8, integerX&0xFF
	tileY, fineY := integerY>>8, integerY&0xFF

	if (tileX >= 0) && (tileX < 64) && (tileY >= 0) && (tileY < 64) {
		tile := adapter.tileMap.Tile(TileCoordinateOf(tileX, tileY))
		z := int(*tile.Properties().FloorHeight) // TODO: take level.heightShift into account

		template := model.LevelObjectTemplate{
			Class:    objectID.Class(),
			Subclass: objectID.Subclass(),
			Type:     objectID.Type(),

			TileX: tileX,
			FineX: fineX,
			TileY: tileY,
			FineY: fineY,
			Z:     z,

			Hitpoints: adapter.objectsAdapter.Object(objectID).CommonHitpoints()}

		adapter.store.AddLevelObject(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), adapter.ID(),
			template, adapter.onLevelObjectAdded,
			adapter.context.simpleStoreFailure("AddLevelObject"))
	}
}

func (adapter *LevelAdapter) onLevelObjectAdded(object model.LevelObject) {
	objects := adapter.levelObjectsMap()
	obj := newLevelObject(&object)
	objects[obj.Index()] = obj
	adapter.levelObjects.notifyObservers()
}

// RequestRemoveObjects requests to remove all identified objects.
func (adapter *LevelAdapter) RequestRemoveObjects(objectIndices []int) {
	levelID := adapter.ID()
	objects := adapter.levelObjectsMap()
	successHandler := func(objectIndex int) func() {
		return func() {
			delete(objects, objectIndex)
			adapter.levelObjects.notifyObservers()
		}
	}

	for _, objectIndex := range objectIndices {
		adapter.store.RemoveLevelObject(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), levelID,
			objectIndex, successHandler(objectIndex),
			adapter.context.simpleStoreFailure(fmt.Sprintf("RemoveLevelObject %v", objectIndex)))
	}
}

// RequestObjectPropertiesChange requests to modify identified objects.
func (adapter *LevelAdapter) RequestObjectPropertiesChange(objectIndices []int, properties *model.LevelObjectProperties) {
	levelID := adapter.ID()
	objects := adapter.levelObjectsMap()
	successHandler := func(objectIndex int) func(newProperties *model.LevelObjectProperties) {
		return func(newProperties *model.LevelObjectProperties) {
			objects[objectIndex].onPropertiesChanged(newProperties)
			adapter.levelObjects.notifyObservers()
		}
	}

	for _, objectIndex := range objectIndices {
		adapter.store.SetLevelObject(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), levelID,
			objectIndex, properties, successHandler(objectIndex),
			adapter.context.simpleStoreFailure(fmt.Sprintf("SetLevelObject %v", objectIndex)))
	}
}

// RequestTilePropertyChange requests the tiles at given coordinates to set provided properties.
func (adapter *LevelAdapter) RequestTilePropertyChange(coordinates []TileCoordinate, properties *model.TileProperties) {
	additionalQueries := make(map[TileCoordinate]bool)
	storeLevelID := adapter.ID()
	tileUpdateHandler := func(coord TileCoordinate) func(model.TileProperties) {
		return func(newProperties model.TileProperties) {
			adapter.onTilePropertiesUpdated(coord, &newProperties)
		}
	}
	for _, coord := range coordinates {
		x, y := coord.XY()
		additionalQueries[TileCoordinateOf(x, y)] = true
		additionalQueries[TileCoordinateOf(x-1, y)] = true
		additionalQueries[TileCoordinateOf(x+1, y)] = true
		additionalQueries[TileCoordinateOf(x, y-1)] = true
		additionalQueries[TileCoordinateOf(x, y+1)] = true
		adapter.store.SetTile(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), storeLevelID,
			x, y, *properties,
			func(model.TileProperties) {}, adapter.context.simpleStoreFailure("SetTile"))
	}
	for coord := range additionalQueries {
		x, y := coord.XY()
		if (x >= 0) && (x < 64) && (y >= 0) && (y < 64) {
			adapter.store.Tile(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), storeLevelID,
				x, y, tileUpdateHandler(coord), adapter.context.simpleStoreFailure("Tile"))
		}
	}
}

// OnLevelSurveillanceChanged registers for updates about the surveillance.
func (adapter *LevelAdapter) OnLevelSurveillanceChanged(callback func()) {
	adapter.levelSurveillance.addObserver(callback)
}

func (adapter *LevelAdapter) onLevelSurveillance(objects []model.SurveillanceObject) {
	adapter.levelSurveillance.set(&objects)
}

// ObjectSurveillanceCount returns how many surveillance objects there are.
func (adapter *LevelAdapter) ObjectSurveillanceCount() int {
	objects := *adapter.levelSurveillance.get().(*[]model.SurveillanceObject)

	return len(objects)
}

// ObjectSurveillanceInfo returns the current source and deathwatch indices for given surveillance index.
func (adapter *LevelAdapter) ObjectSurveillanceInfo(index int) (source int, deathwatch int) {
	objects := *adapter.levelSurveillance.get().(*[]model.SurveillanceObject)

	if (index >= 0) && (index < len(objects)) {
		source = *objects[index].SourceIndex
		deathwatch = *objects[index].DeathwatchIndex
	}

	return
}

// RequestObjectSurveillance requests to modify the surveillance information of identified object.
func (adapter *LevelAdapter) RequestObjectSurveillance(surveillanceIndex int, sourceObject *int, deathwatchObject *int) {
	var data model.SurveillanceObject

	data.SourceIndex = sourceObject
	data.DeathwatchIndex = deathwatchObject

	adapter.store.SetLevelSurveillanceObject(adapter.context.ActiveProjectID(), adapter.context.ActiveArchiveID(), adapter.ID(),
		surveillanceIndex, data,
		adapter.onLevelSurveillance, adapter.context.simpleStoreFailure("SetLevelSurveillanceObject"))
}
