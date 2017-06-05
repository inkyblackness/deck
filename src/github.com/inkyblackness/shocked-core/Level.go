package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/res/logic"

	model "github.com/inkyblackness/shocked-model"
)

const surveillanceSources = 8

var tileTypes = map[data.TileType]model.TileType{
	data.Solid: model.Solid,
	data.Open:  model.Open,

	data.DiagonalOpenSouthEast: model.DiagonalOpenSouthEast,
	data.DiagonalOpenSouthWest: model.DiagonalOpenSouthWest,
	data.DiagonalOpenNorthWest: model.DiagonalOpenNorthWest,
	data.DiagonalOpenNorthEast: model.DiagonalOpenNorthEast,

	data.SlopeSouthToNorth: model.SlopeSouthToNorth,
	data.SlopeWestToEast:   model.SlopeWestToEast,
	data.SlopeNorthToSouth: model.SlopeNorthToSouth,
	data.SlopeEastToWest:   model.SlopeEastToWest,

	data.ValleySouthEastToNorthWest: model.ValleySouthEastToNorthWest,
	data.ValleySouthWestToNorthEast: model.ValleySouthWestToNorthEast,
	data.ValleyNorthWestToSouthEast: model.ValleyNorthWestToSouthEast,
	data.ValleyNorthEastToSouthWest: model.ValleyNorthEastToSouthWest,

	data.RidgeNorthWestToSouthEast: model.RidgeNorthWestToSouthEast,
	data.RidgeNorthEastToSouthWest: model.RidgeNorthEastToSouthWest,
	data.RidgeSouthEastToNorthWest: model.RidgeSouthEastToNorthWest,
	data.RidgeSouthWestToNorthEast: model.RidgeSouthWestToNorthEast}

var slopeControls = map[data.SlopeControl]model.SlopeControl{
	data.SlopeCeilingInverted: model.SlopeCeilingInverted,
	data.SlopeCeilingMirrored: model.SlopeCeilingMirrored,
	data.SlopeCeilingFlat:     model.SlopeCeilingFlat,
	data.SlopeFloorFlat:       model.SlopeFloorFlat}

func tileType(modelType model.TileType) (dataType data.TileType) {
	dataType = data.Solid

	for key, value := range tileTypes {
		if value == modelType {
			dataType = key
		}
	}

	return
}

func slopeControl(modelControl model.SlopeControl) (dataControl data.SlopeControl) {
	dataControl = data.SlopeCeilingInverted

	for key, value := range slopeControls {
		if value == modelControl {
			dataControl = key
		}
	}

	return
}

// Level is a structure holding level specific information.
type Level struct {
	id    int
	store chunk.Store

	mutex sync.Mutex

	tileMapStore chunk.BlockStore
	tileMap      *logic.TileMap

	objectListStore chunk.BlockStore
	objectList      []data.LevelObjectEntry
	objectChain     *logic.LevelObjectChain

	crossrefListStore chunk.BlockStore
	crossrefList      *logic.CrossReferenceList

	surveillanceSourceStore     chunk.BlockStore
	surveillanceDeathwatchStore chunk.BlockStore
}

// NewLevel returns a new instance of a Level structure.
func NewLevel(store chunk.Store, id int) *Level {
	baseStoreID := 4000 + id*100
	level := &Level{
		id:    id,
		store: store,

		tileMapStore: store.Get(res.ResourceID(baseStoreID + 5)),
		tileMap:      nil,

		objectListStore: store.Get(res.ResourceID(baseStoreID + 8)),

		crossrefListStore: store.Get(res.ResourceID(baseStoreID + 9)),

		surveillanceSourceStore:     store.Get(res.ResourceID(baseStoreID + 43)),
		surveillanceDeathwatchStore: store.Get(res.ResourceID(baseStoreID + 44))}

	level.tileMap = logic.DecodeTileMap(level.tileMapStore.BlockData(0), 64, 64)
	level.crossrefList = logic.DecodeCrossReferenceList(level.crossrefListStore.BlockData(0))

	{
		blockData := level.objectListStore.BlockData(0)
		level.objectList = make([]data.LevelObjectEntry, len(blockData)/data.LevelObjectEntrySize)
		reader := bytes.NewReader(blockData)
		binary.Read(reader, binary.LittleEndian, &level.objectList)

		level.objectChain = logic.NewLevelObjectChain(&level.objectList[0],
			func(index data.LevelObjectChainIndex) logic.LevelObjectChainLink {
				return &level.objectList[index]
			})
	}

	return level
}

func (level *Level) onTileDataChanged() {
	level.tileMapStore.SetBlockData(0, level.tileMap.Encode())
}

// ID returns the identifier of the level.
func (level *Level) ID() int {
	return level.id
}

func (level *Level) information() data.LevelInformation {
	blockData := level.store.Get(res.ResourceID(4000 + level.id*100 + 4)).BlockData(0)
	reader := bytes.NewReader(blockData)
	var info data.LevelInformation

	binary.Read(reader, binary.LittleEndian, &info)

	return info
}

func (level *Level) variables() data.LevelVariables {
	blockData := level.store.Get(res.ResourceID(4000 + level.id*100 + 45)).BlockData(0)
	reader := bytes.NewReader(blockData)
	var info data.LevelVariables

	binary.Read(reader, binary.LittleEndian, &info)

	return info
}

func (level *Level) isCyberspace() bool {
	info := level.information()
	return info.IsCyberspace()
}

// Properties returns the properties of the level.
func (level *Level) Properties() (result model.LevelProperties) {
	info := level.information()
	vars := level.variables()

	result.CyberspaceFlag = boolAsPointer(info.IsCyberspace())
	result.HeightShift = intAsPointer(int(info.HeightShift))
	result.CeilingHasRadiation = boolAsPointer(vars.RadiationRegister > 1)
	result.CeilingEffectLevel = intAsPointer(int(vars.Radiation))
	result.FloorHasBiohazard = boolAsPointer(vars.BioRegister > 1)
	result.FloorHasGravity = boolAsPointer(vars.GravitySwitch != 0)
	result.FloorEffectLevel = intAsPointer(int(vars.BioOrGravity))

	return
}

// SetProperties updates the properties of a level.
func (level *Level) SetProperties(properties model.LevelProperties) {
	{
		infoStore := level.store.Get(res.ResourceID(4000 + level.id*100 + 4))
		infoData := infoStore.BlockData(0)
		infoReader := bytes.NewReader(infoData)
		infoWriter := bytes.NewBuffer(nil)
		var info data.LevelInformation

		binary.Read(infoReader, binary.LittleEndian, &info)
		if properties.CyberspaceFlag != nil {
			info.CyberspaceFlag = 0
			if *properties.CyberspaceFlag {
				info.CyberspaceFlag = 1
			}
		}
		if properties.HeightShift != nil {
			info.HeightShift = uint32(*properties.HeightShift)
		}
		binary.Write(infoWriter, binary.LittleEndian, &info)
		infoStore.SetBlockData(0, infoWriter.Bytes())
	}
	{
		varsStore := level.store.Get(res.ResourceID(4000 + level.id*100 + 45))
		varsData := varsStore.BlockData(0)
		varsReader := bytes.NewReader(varsData)
		varsWriter := bytes.NewBuffer(nil)
		var vars data.LevelVariables

		binary.Read(varsReader, binary.LittleEndian, &vars)
		if properties.CeilingHasRadiation != nil {
			vars.RadiationRegister = 0
			if *properties.CeilingHasRadiation {
				vars.RadiationRegister = 2
			}
		}
		if properties.CeilingEffectLevel != nil {
			vars.Radiation = byte(*properties.CeilingEffectLevel)
		}
		if properties.FloorHasBiohazard != nil {
			vars.BioRegister = 0
			if *properties.FloorHasBiohazard {
				vars.BioRegister = 2
			}
		}
		if properties.FloorHasGravity != nil {
			vars.GravitySwitch = 0
			if *properties.FloorHasGravity {
				vars.GravitySwitch = 1
			}
		}
		if properties.FloorEffectLevel != nil {
			vars.BioOrGravity = byte(*properties.FloorEffectLevel)
		}
		binary.Write(varsWriter, binary.LittleEndian, &vars)
		varsStore.SetBlockData(0, varsWriter.Bytes())
	}
}

// Textures returns the texture identifier used in this level.
func (level *Level) Textures() (result []int) {
	blockData := level.store.Get(res.ResourceID(4000 + level.id*100 + 7)).BlockData(0)
	reader := bytes.NewReader(blockData)
	var ids [54]uint16

	binary.Read(reader, binary.LittleEndian, &ids)
	for _, id := range ids {
		result = append(result, int(id))
	}

	return
}

// SetTextures sets the texture identifier for this level.
func (level *Level) SetTextures(newIds []int) {
	blockStore := level.store.Get(res.ResourceID(4000 + level.id*100 + 7))
	var ids [54]uint16
	toCopy := len(ids)

	if len(newIds) < toCopy {
		toCopy = len(newIds)
	}
	for index := 0; index < len(ids); index++ {
		ids[index] = uint16(newIds[index])
	}

	buffer := bytes.NewBuffer(nil)
	binary.Write(buffer, binary.LittleEndian, &ids)
	blockStore.SetBlockData(0, buffer.Bytes())
}

// TextureAnimations returns the properties of the animation groups.
func (level *Level) TextureAnimations() (result []model.TextureAnimation) {
	level.mutex.Lock()
	defer level.mutex.Unlock()
	var rawEntries [4]data.TextureAnimationEntry

	result = make([]model.TextureAnimation, len(rawEntries))
	level.readTable(42, &rawEntries)
	for index := 0; index < len(rawEntries); index++ {
		resultEntry := &result[index]
		rawEntry := &rawEntries[index]

		resultEntry.FrameCount = intAsPointer(int(rawEntry.FrameCount))
		resultEntry.FrameTime = intAsPointer(int(rawEntry.FrameTime))
		resultEntry.LoopType = intAsPointer(int(rawEntry.LoopType))
	}
	return
}

// SetTextureAnimation modifies the properties of identified animation group.
func (level *Level) SetTextureAnimation(animationGroup int, properties model.TextureAnimation) {
	level.mutex.Lock()
	defer level.mutex.Unlock()

	var rawEntries [4]data.TextureAnimationEntry
	rawEntry := &rawEntries[animationGroup]

	level.readTable(42, &rawEntries)
	if properties.FrameCount != nil {
		rawEntry.FrameCount = byte(*properties.FrameCount)
	}
	if properties.FrameTime != nil {
		rawEntry.FrameTime = uint16(*properties.FrameTime)
	}
	if properties.LoopType != nil {
		rawEntry.LoopType = data.TextureAnimationLoopType(*properties.LoopType)
	}
	level.writeTable(42, &rawEntries)
}

// Objects returns an array of all used objects.
func (level *Level) Objects() []model.LevelObject {
	level.mutex.Lock()
	defer level.mutex.Unlock()

	var result []model.LevelObject

	for index, rawEntry := range level.objectList {
		if rawEntry.IsInUse() {
			result = append(result, level.objectFromRawEntry(index, &rawEntry))
		}
	}

	return result
}

func intAsPointer(value int) (ptr *int) {
	ptr = new(int)
	*ptr = value
	return
}

func boolAsPointer(value bool) (ptr *bool) {
	ptr = new(bool)
	*ptr = value
	return
}

func (level *Level) objectFromRawEntry(index int, rawEntry *data.LevelObjectEntry) (entry model.LevelObject) {
	entry.Identifiable = model.Identifiable{ID: fmt.Sprintf("%d", index)}
	entry.Class = int(rawEntry.Class)

	entry.Properties.Subclass = intAsPointer(int(rawEntry.Subclass))
	entry.Properties.Type = intAsPointer(int(rawEntry.Type))

	entry.Properties.TileX = intAsPointer(int(rawEntry.X >> 8))
	entry.Properties.FineX = intAsPointer(int(rawEntry.X & 0xFF))
	entry.Properties.TileY = intAsPointer(int(rawEntry.Y >> 8))
	entry.Properties.FineY = intAsPointer(int(rawEntry.Y & 0xFF))
	entry.Properties.Z = intAsPointer(int(rawEntry.Z))

	entry.Properties.RotationX = intAsPointer(int(rawEntry.Rot1))
	entry.Properties.RotationY = intAsPointer(int(rawEntry.Rot3))
	entry.Properties.RotationZ = intAsPointer(int(rawEntry.Rot2))
	entry.Properties.Hitpoints = intAsPointer(int(rawEntry.Hitpoints))

	meta := data.LevelObjectClassMetaEntry(rawEntry.Class)
	classStore := level.store.Get(res.ResourceID(4000 + level.id*100 + 10 + entry.Class))
	blockData := classStore.BlockData(0)
	startOffset := meta.EntrySize * int(rawEntry.ClassTableIndex)
	if (startOffset + meta.EntrySize) > len(blockData) {
		fmt.Printf("!!!!! class %d meta says %d bytes size, can't reach index %d in blockData %d",
			int(entry.Class), meta.EntrySize, rawEntry.ClassTableIndex, len(blockData))
	} else {
		entry.Properties.ClassData = blockData[startOffset+data.LevelObjectPrefixSize : startOffset+meta.EntrySize]
	}

	return
}

// RemoveObject removes an object from the level.
func (level *Level) RemoveObject(objectIndex int) (err error) {
	level.mutex.Lock()
	defer level.mutex.Unlock()

	if (objectIndex > 0) && (objectIndex < len(level.objectList)) {
		objectEntry := &level.objectList[objectIndex]

		if objectEntry.IsInUse() {
			classMeta := data.LevelObjectClassMetaEntry(objectEntry.Class)
			classStore := level.store.Get(res.ResourceID(4000 + level.id*100 + 10 + int(objectEntry.Class)))
			classTable := logic.DecodeLevelObjectClassTable(classStore.BlockData(0), classMeta.EntrySize)
			classChain := classTable.AsChain()

			level.crossrefList.RemoveEntriesFromMap(logic.CrossReferenceListIndex(objectEntry.CrossReferenceTableIndex), level.tileMap)
			classChain.ReleaseLink(data.LevelObjectChainIndex(objectEntry.ClassTableIndex))
			level.objectChain.ReleaseLink(data.LevelObjectChainIndex(objectIndex))
			objectEntry.InUse = 0

			level.onObjectListChanged(classStore, classTable)
		}
	}

	return
}

// AddObject adds a new object at given tile.
func (level *Level) AddObject(template *model.LevelObjectTemplate) (entry model.LevelObject, err error) {
	level.mutex.Lock()
	defer level.mutex.Unlock()

	objID := res.MakeObjectID(res.ObjectClass(template.Class),
		res.ObjectSubclass(template.Subclass), res.ObjectType(template.Type))

	classMeta := data.LevelObjectClassMetaEntry(objID.Class)
	classStore := level.store.Get(res.ResourceID(4000 + level.id*100 + 10 + int(objID.Class)))
	classTable := logic.DecodeLevelObjectClassTable(classStore.BlockData(0), classMeta.EntrySize)
	classChain := classTable.AsChain()
	classIndex, classErr := classChain.AcquireLink()

	if classErr != nil {
		err = classErr
		return
	}
	classEntry := classTable.Entry(classIndex)
	classEntryData := classEntry.Data()
	for index := 0; index < len(classEntryData); index++ {
		classEntryData[index] = 0x00
	}

	objectIndex, objectErr := level.objectChain.AcquireLink()
	if objectErr != nil {
		classChain.ReleaseLink(classIndex)
		err = objectErr
		return
	}

	locations := []logic.TileLocation{logic.AtTile(uint16(template.TileX), uint16(template.TileY))}

	crossrefIndex, crossrefErr := level.crossrefList.AddObjectToMap(uint16(objectIndex), level.tileMap, locations)
	if crossrefErr != nil {
		classChain.ReleaseLink(classIndex)
		level.objectChain.ReleaseLink(objectIndex)
		err = crossrefErr
		return
	}
	crossrefEntry := level.crossrefList.Entry(crossrefIndex)

	objectEntry := &level.objectList[objectIndex]
	objectEntry.InUse = 1
	objectEntry.Class = objID.Class
	objectEntry.Subclass = objID.Subclass
	objectEntry.Type = objID.Type
	objectEntry.X = data.MapCoordinateOf(byte(template.TileX), byte(template.FineX))
	objectEntry.Y = data.MapCoordinateOf(byte(template.TileY), byte(template.FineY))
	objectEntry.Z = byte(template.Z)
	objectEntry.Rot1 = 0
	objectEntry.Rot2 = 0
	objectEntry.Rot3 = 0
	objectEntry.Hitpoints = uint16(template.Hitpoints)
	if objectEntry.Hitpoints == 0 {
		objectEntry.Hitpoints = 1
	}

	objectEntry.CrossReferenceTableIndex = uint16(crossrefIndex)
	crossrefEntry.LevelObjectTableIndex = uint16(objectIndex)

	objectEntry.ClassTableIndex = uint16(classIndex)
	classEntry.LevelObjectTableIndex = uint16(objectIndex)

	level.onObjectListChanged(classStore, classTable)
	entry = level.objectFromRawEntry(int(objectIndex), objectEntry)

	return
}

// SetObject modifies the properties of identified object.
func (level *Level) SetObject(objectIndex int, newProperties *model.LevelObjectProperties) (properties model.LevelObjectProperties, err error) {
	level.mutex.Lock()
	defer level.mutex.Unlock()

	if (objectIndex > 0) && (objectIndex < len(level.objectList)) {
		objectEntry := &level.objectList[objectIndex]

		if objectEntry.IsInUse() {
			classMeta := data.LevelObjectClassMetaEntry(objectEntry.Class)
			classStore := level.store.Get(res.ResourceID(4000 + level.id*100 + 10 + int(objectEntry.Class)))
			classTable := logic.DecodeLevelObjectClassTable(classStore.BlockData(0), classMeta.EntrySize)
			changedTile := false

			if newProperties.Subclass != nil {
				objectEntry.Subclass = res.ObjectSubclass(*newProperties.Subclass)
			}
			if newProperties.Type != nil {
				objectEntry.Type = res.ObjectType(*newProperties.Type)
			}
			if newProperties.Z != nil {
				objectEntry.Z = byte(*newProperties.Z)
			}
			newTileX, newFineX := objectEntry.X.Tile(), objectEntry.X.Offset()
			if (newProperties.TileX != nil) && (newTileX != byte(*newProperties.TileX)) {
				newTileX = byte(*newProperties.TileX)
				changedTile = true
			}
			if newProperties.FineX != nil {
				newFineX = byte(*newProperties.FineX)
			}
			objectEntry.X = data.MapCoordinateOf(newTileX, newFineX)
			newTileY, newFineY := objectEntry.Y.Tile(), objectEntry.Y.Offset()
			if (newProperties.TileY != nil) && (newTileY != byte(*newProperties.TileY)) {
				newTileY = byte(*newProperties.TileY)
				changedTile = true
			}
			if newProperties.FineY != nil {
				newFineY = byte(*newProperties.FineY)
			}
			objectEntry.Y = data.MapCoordinateOf(newTileY, newFineY)

			if newProperties.RotationX != nil {
				objectEntry.Rot1 = byte(*newProperties.RotationX)
			}
			if newProperties.RotationY != nil {
				objectEntry.Rot3 = byte(*newProperties.RotationY)
			}
			if newProperties.RotationZ != nil {
				objectEntry.Rot2 = byte(*newProperties.RotationZ)
			}
			if newProperties.Hitpoints != nil {
				objectEntry.Hitpoints = uint16(*newProperties.Hitpoints)
			}

			if len(newProperties.ClassData) > 0 {
				classEntry := classTable.Entry(data.LevelObjectChainIndex(objectEntry.ClassTableIndex))

				copy(classEntry.Data(), newProperties.ClassData)
			}
			if changedTile {
				locations := []logic.TileLocation{logic.AtTile(uint16(newTileX), uint16(newTileY))}

				if objectEntry.CrossReferenceTableIndex != 0 {
					level.crossrefList.RemoveEntriesFromMap(logic.CrossReferenceListIndex(objectEntry.CrossReferenceTableIndex), level.tileMap)
					objectEntry.CrossReferenceTableIndex = 0
				}

				crossrefIndex, crossrefErr := level.crossrefList.AddObjectToMap(uint16(objectIndex), level.tileMap, locations)
				if crossrefErr != nil {
					// This is a kind of bad (and weird) situation.
					// The object, which was already stored, can not be stored anymore (?) and is furthermore left
					// in an incorrect state.
					err = crossrefErr
					return
				}
				crossrefEntry := level.crossrefList.Entry(crossrefIndex)

				objectEntry.CrossReferenceTableIndex = uint16(crossrefIndex)
				crossrefEntry.LevelObjectTableIndex = uint16(objectIndex)
			}

			level.onObjectListChanged(classStore, classTable)
			properties = level.objectFromRawEntry(int(objectIndex), objectEntry).Properties
		} else {
			err = fmt.Errorf("Object is not in use")
		}
	} else {
		err = fmt.Errorf("Invalid object index")
	}

	return
}

func (level *Level) onObjectListChanged(classStore chunk.BlockStore, classTable *logic.LevelObjectClassTable) {
	classStore.SetBlockData(0, classTable.Encode())

	objWriter := bytes.NewBuffer(nil)
	binary.Write(objWriter, binary.LittleEndian, level.objectList)
	level.objectListStore.SetBlockData(0, objWriter.Bytes())

	level.crossrefListStore.SetBlockData(0, level.crossrefList.Encode())
	level.onTileDataChanged()
}

func (level *Level) readTable(levelBlockID int, value interface{}) {
	blockData := level.store.Get(res.ResourceID(4000 + level.id*100 + levelBlockID)).BlockData(0)
	reader := bytes.NewReader(blockData)

	binary.Read(reader, binary.LittleEndian, value)
}

func (level *Level) writeTable(levelBlockID int, value interface{}) {
	writer := bytes.NewBuffer(nil)

	binary.Write(writer, binary.LittleEndian, value)
	level.store.Get(res.ResourceID(4000+level.id*100+levelBlockID)).SetBlockData(0, writer.Bytes())
}

func (level *Level) isTileTypeValley(tileType data.TileType) bool {
	return tileType == data.ValleyNorthEastToSouthWest || tileType == data.ValleyNorthWestToSouthEast ||
		tileType == data.ValleySouthEastToNorthWest || tileType == data.ValleySouthWestToNorthEast
}

// Direction describes a flag field.
type Direction int

//
const (
	DirNone  = Direction(0)
	DirNorth = Direction(1)
	DirEast  = Direction(2)
	DirSouth = Direction(4)
	DirWest  = Direction(8)
)

var solidDirections = map[data.TileType]Direction{
	data.Solid: DirNorth | DirEast | DirSouth | DirWest,
	data.Open:  DirNone,

	data.DiagonalOpenSouthEast: DirNorth | DirWest,
	data.DiagonalOpenSouthWest: DirNorth | DirEast,
	data.DiagonalOpenNorthWest: DirEast | DirSouth,
	data.DiagonalOpenNorthEast: DirSouth | DirWest,

	data.SlopeSouthToNorth: DirNone,
	data.SlopeWestToEast:   DirNone,
	data.SlopeNorthToSouth: DirNone,
	data.SlopeEastToWest:   DirNone,

	data.ValleySouthEastToNorthWest: DirNone,
	data.ValleySouthWestToNorthEast: DirNone,
	data.ValleyNorthWestToSouthEast: DirNone,
	data.ValleyNorthEastToSouthWest: DirNone,

	data.RidgeNorthWestToSouthEast: DirNone,
	data.RidgeNorthEastToSouthWest: DirNone,
	data.RidgeSouthEastToNorthWest: DirNone,
	data.RidgeSouthWestToNorthEast: DirNone}

var slopedCeilingHeightDirections = map[data.TileType]Direction{
	data.Solid: DirNone,
	data.Open:  DirNone,

	data.DiagonalOpenSouthEast: DirNone,
	data.DiagonalOpenSouthWest: DirNone,
	data.DiagonalOpenNorthWest: DirNone,
	data.DiagonalOpenNorthEast: DirNone,

	data.SlopeSouthToNorth: DirSouth,
	data.SlopeWestToEast:   DirWest,
	data.SlopeNorthToSouth: DirNorth,
	data.SlopeEastToWest:   DirEast,

	data.ValleySouthEastToNorthWest: DirNone,
	data.ValleySouthWestToNorthEast: DirNone,
	data.ValleyNorthWestToSouthEast: DirNone,
	data.ValleyNorthEastToSouthWest: DirNone,

	data.RidgeNorthWestToSouthEast: DirNorth | DirWest,
	data.RidgeNorthEastToSouthWest: DirNorth | DirEast,
	data.RidgeSouthEastToNorthWest: DirEast | DirSouth,
	data.RidgeSouthWestToNorthEast: DirSouth | DirWest}

var slopedFloorHeightDirections = map[data.TileType]Direction{
	data.Solid: DirNone,
	data.Open:  DirNone,

	data.DiagonalOpenSouthEast: DirNone,
	data.DiagonalOpenSouthWest: DirNone,
	data.DiagonalOpenNorthWest: DirNone,
	data.DiagonalOpenNorthEast: DirNone,

	data.SlopeSouthToNorth: DirNorth | DirEast | DirWest,
	data.SlopeWestToEast:   DirNorth | DirEast | DirSouth,
	data.SlopeNorthToSouth: DirEast | DirSouth | DirWest,
	data.SlopeEastToWest:   DirNorth | DirSouth | DirWest,

	data.ValleySouthEastToNorthWest: DirNorth | DirEast | DirSouth | DirWest,
	data.ValleySouthWestToNorthEast: DirNorth | DirEast | DirSouth | DirWest,
	data.ValleyNorthWestToSouthEast: DirNorth | DirEast | DirSouth | DirWest,
	data.ValleyNorthEastToSouthWest: DirNorth | DirEast | DirSouth | DirWest,

	data.RidgeNorthWestToSouthEast: DirEast | DirSouth,
	data.RidgeNorthEastToSouthWest: DirSouth | DirWest,
	data.RidgeSouthEastToNorthWest: DirWest | DirNorth,
	data.RidgeSouthWestToNorthEast: DirNorth | DirEast}

var mirroredSlopes = map[data.TileType]data.TileType{
	data.Solid: data.Solid,
	data.Open:  data.Open,

	data.DiagonalOpenSouthEast: data.DiagonalOpenSouthEast,
	data.DiagonalOpenSouthWest: data.DiagonalOpenSouthWest,
	data.DiagonalOpenNorthWest: data.DiagonalOpenNorthWest,
	data.DiagonalOpenNorthEast: data.DiagonalOpenNorthEast,

	data.SlopeSouthToNorth: data.SlopeNorthToSouth,
	data.SlopeWestToEast:   data.SlopeEastToWest,
	data.SlopeNorthToSouth: data.SlopeSouthToNorth,
	data.SlopeEastToWest:   data.SlopeWestToEast,

	data.ValleySouthEastToNorthWest: data.RidgeNorthWestToSouthEast,
	data.ValleySouthWestToNorthEast: data.RidgeNorthEastToSouthWest,
	data.ValleyNorthWestToSouthEast: data.RidgeSouthEastToNorthWest,
	data.ValleyNorthEastToSouthWest: data.RidgeSouthWestToNorthEast,

	data.RidgeNorthWestToSouthEast: data.ValleySouthEastToNorthWest,
	data.RidgeNorthEastToSouthWest: data.ValleySouthWestToNorthEast,
	data.RidgeSouthEastToNorthWest: data.ValleyNorthWestToSouthEast,
	data.RidgeSouthWestToNorthEast: data.ValleyNorthEastToSouthWest}

func (level *Level) calculatedFloorHeight(tile *data.TileMapEntry, dir Direction) (height int) {
	if (solidDirections[tile.Type] & dir) == 0 {
		slopeControl := (tile.Flags >> 10) & 0x3

		height = int(tile.Floor & 0x1F)
		if (slopeControl != 3) && ((slopedFloorHeightDirections[tile.Type] & dir) != 0) {
			height += int(tile.SlopeHeight)
		}
	} else {
		height = 0x20
	}

	return
}

func (level *Level) calculatedCeilingHeight(tile *data.TileMapEntry, dir Direction) (height int) {
	if (solidDirections[tile.Type] & dir) == 0 {
		slopeControl := (tile.Flags >> 10) & 0x3

		height = 0x20 - int(tile.Ceiling&0x1F)
		if ((slopeControl == 0) || (slopeControl == 3)) && ((slopedCeilingHeightDirections[tile.Type] & dir) != 0) {
			height -= int(tile.SlopeHeight)
		} else if (slopeControl == 1) && ((slopedCeilingHeightDirections[mirroredSlopes[tile.Type]] & dir) != 0) {
			height -= int(tile.SlopeHeight)
		}
	} else {
		height = 0x20
	}

	return
}

func (level *Level) getTileType(x, y int) data.TileType {
	tileType := data.Solid

	if (x >= 0) && (x < 64) && (y >= 0) && (y < 64) {
		entry := level.tileMap.Entry(logic.AtTile(uint16(x), uint16(y)))
		tileType = entry.Type
	}

	return tileType
}

func (level *Level) calculateWallHeight(thisTile *data.TileMapEntry, thisDir Direction, otherX, otherY int, otherDir Direction) float32 {
	calculatedHeight := float32(0.0)

	if (solidDirections[thisTile.Type] & thisDir) == 0 {
		otherType := level.getTileType(otherX, otherY)

		if (solidDirections[otherType] & otherDir) == 0 {
			otherTile := level.tileMap.Entry(logic.AtTile(uint16(otherX), uint16(otherY)))
			thisCeilingHeight := level.calculatedCeilingHeight(thisTile, thisDir)
			otherCeilingHeight := level.calculatedCeilingHeight(otherTile, otherDir)
			thisFloorHeight := level.calculatedFloorHeight(thisTile, thisDir)
			otherFloorHeight := level.calculatedFloorHeight(otherTile, otherDir)

			if (thisCeilingHeight < otherCeilingHeight) ||
				((thisCeilingHeight == otherCeilingHeight) && thisFloorHeight < otherFloorHeight) {

				if (thisFloorHeight >= otherCeilingHeight) || (otherFloorHeight >= thisCeilingHeight) {
					calculatedHeight = 1.0
				} else {
					minFloorHeight := thisFloorHeight
					maxFloorHeight := otherFloorHeight
					minCeilingHeight := thisCeilingHeight
					if minFloorHeight > otherFloorHeight {
						minFloorHeight = otherFloorHeight
						maxFloorHeight = thisFloorHeight
					}
					if minCeilingHeight > otherCeilingHeight {
						minCeilingHeight = otherCeilingHeight
					}
					calculatedHeight = (float32(maxFloorHeight) - float32(minFloorHeight)) / (float32(minCeilingHeight) - float32(minFloorHeight))
				}

			}
		} else {
			calculatedHeight = 1.0
		}
	}

	return calculatedHeight
}

// TileProperties returns the properties of a specific tile.
func (level *Level) TileProperties(x, y int) (result model.TileProperties) {
	level.mutex.Lock()
	defer level.mutex.Unlock()

	entry := level.tileMap.Entry(logic.AtTile(uint16(x), uint16(y)))
	result.Type = new(model.TileType)
	*result.Type = tileTypes[entry.Type]
	result.SlopeHeight = new(model.HeightUnit)
	*result.SlopeHeight = model.HeightUnit(entry.SlopeHeight)
	result.FloorHeight = new(model.HeightUnit)
	*result.FloorHeight = model.HeightUnit(entry.Floor & 0x1F)
	result.CeilingHeight = new(model.HeightUnit)
	*result.CeilingHeight = model.HeightUnit(entry.Ceiling & 0x1F)
	result.SlopeControl = new(model.SlopeControl)
	*result.SlopeControl = slopeControls[data.SlopeControl((entry.Flags>>10)&0x3)]

	{
		result.CalculatedWallHeights = new(model.CalculatedWallHeights)
		result.CalculatedWallHeights.North = level.calculateWallHeight(entry, DirNorth, x, y+1, DirSouth)
		result.CalculatedWallHeights.East = level.calculateWallHeight(entry, DirEast, x+1, y, DirWest)
		result.CalculatedWallHeights.South = level.calculateWallHeight(entry, DirSouth, x, y-1, DirNorth)
		result.CalculatedWallHeights.West = level.calculateWallHeight(entry, DirWest, x-1, y, DirEast)
	}

	result.MusicIndex = intAsPointer(entry.Flags.MusicIndex())

	if !level.isCyberspace() {
		var properties model.RealWorldTileProperties
		var textureIDs = uint16(entry.Textures)

		properties.WallTexture = intAsPointer(int(textureIDs & 0x3F))
		properties.CeilingTexture = intAsPointer(int((textureIDs >> 6) & 0x1F))
		properties.CeilingTextureRotations = intAsPointer(int((entry.Ceiling >> 5) & 0x03))
		properties.FloorTexture = intAsPointer(int((textureIDs >> 11) & 0x1F))
		properties.FloorTextureRotations = intAsPointer(int((entry.Floor >> 5) & 0x03))

		properties.UseAdjacentWallTexture = boolAsPointer((entry.Flags & 0x00000100) != 0)
		properties.WallTextureOffset = new(model.HeightUnit)
		*properties.WallTextureOffset = model.HeightUnit(entry.Flags & 0x0000001F)
		properties.WallTexturePattern = intAsPointer(int((entry.Flags >> 5) & 0x00000003))

		properties.FloorHazard = boolAsPointer((entry.Floor & 0x80) != 0)
		properties.CeilingHazard = boolAsPointer((entry.Ceiling & 0x80) != 0)

		properties.FloorShadow = intAsPointer(entry.Flags.FloorShadow())
		properties.CeilingShadow = intAsPointer(entry.Flags.CeilingShadow())

		properties.SpookyMusic = boolAsPointer((entry.Flags & 0x00000200) != 0)

		result.RealWorld = &properties
	} else {
		var properties model.CyberspaceTileProperties
		var colors = uint16(entry.Textures)

		properties.FloorColorIndex = intAsPointer(int((colors >> 0) & 0x00FF))
		properties.CeilingColorIndex = intAsPointer(int((colors >> 8) & 0x00FF))

		properties.FlightPullType = intAsPointer(int((entry.Flags>>16)&0xF) + int((entry.Flags>>20)&0x10))
		properties.GameOfLifeSet = boolAsPointer((entry.Flags & 0x00000040) != 0)

		result.Cyberspace = &properties
	}

	return
}

// SetTileProperties sets the properties of a specific tile.
func (level *Level) SetTileProperties(x, y int, properties model.TileProperties) {
	level.mutex.Lock()
	defer level.mutex.Unlock()
	isCyberspace := level.isCyberspace()

	entry := level.tileMap.Entry(logic.AtTile(uint16(x), uint16(y)))
	flags := uint32(entry.Flags)
	if properties.Type != nil {
		entry.Type = tileType(*properties.Type)
	}
	if properties.FloorHeight != nil {
		entry.Floor = (entry.Floor & 0xE0) | (byte(*properties.FloorHeight) & 0x1F)
	}
	if properties.CeilingHeight != nil {
		entry.Ceiling = (entry.Ceiling & 0xE0) | (byte(*properties.CeilingHeight) & 0x1F)
	}
	if properties.SlopeHeight != nil {
		entry.SlopeHeight = byte(*properties.SlopeHeight)
	}
	if properties.SlopeControl != nil {
		flags = (flags & ^uint32(0x00000C00)) | (uint32(slopeControl(*properties.SlopeControl)) << 10)
	}
	if properties.MusicIndex != nil {
		flags = uint32(data.TileFlag(flags).WithMusicIndex(*properties.MusicIndex))
	}
	if !isCyberspace && (properties.RealWorld != nil) {
		var textureIDs = uint16(entry.Textures)

		if properties.RealWorld.FloorTexture != nil && (*properties.RealWorld.FloorTexture < 0x20) {
			textureIDs = (textureIDs & 0x07FF) | (uint16(*properties.RealWorld.FloorTexture) << 11)
		}
		if properties.RealWorld.FloorTextureRotations != nil {
			entry.Floor = (entry.Floor & 0x9F) | ((byte(*properties.RealWorld.FloorTextureRotations) & 0x3) << 5)
		}
		if properties.RealWorld.CeilingTexture != nil && (*properties.RealWorld.CeilingTexture < 0x20) {
			textureIDs = (textureIDs & 0xF83F) | (uint16(*properties.RealWorld.CeilingTexture) << 6)
		}
		if properties.RealWorld.CeilingTextureRotations != nil {
			entry.Ceiling = (entry.Ceiling & 0x9F) | ((byte(*properties.RealWorld.CeilingTextureRotations) & 0x3) << 5)
		}
		if properties.RealWorld.WallTexture != nil && (*properties.RealWorld.WallTexture < 0x40) {
			textureIDs = (textureIDs & 0xFFC0) | uint16(*properties.RealWorld.WallTexture)
		}
		if properties.RealWorld.UseAdjacentWallTexture != nil {
			flags = flags & ^uint32(0x00000100)
			if *properties.RealWorld.UseAdjacentWallTexture {
				flags |= 0x00000100
			}
		}
		if properties.RealWorld.WallTextureOffset != nil && *properties.RealWorld.WallTextureOffset < 0x20 {
			flags = (flags & ^uint32(0x0000001F)) | uint32(*properties.RealWorld.WallTextureOffset)
		}
		if properties.RealWorld.WallTexturePattern != nil {
			flags = (flags & ^uint32(0x00000060) | (uint32(*properties.RealWorld.WallTexturePattern) << 5))
		}
		if properties.RealWorld.FloorHazard != nil {
			entry.Floor &= 0x7F
			if *properties.RealWorld.FloorHazard {
				entry.Floor |= 0x80
			}
		}
		if properties.RealWorld.CeilingHazard != nil {
			entry.Ceiling &= 0x7F
			if *properties.RealWorld.CeilingHazard {
				entry.Ceiling |= 0x80
			}
		}
		if properties.RealWorld.FloorShadow != nil {
			flags = uint32(data.TileFlag(flags).WithFloorShadow(*properties.RealWorld.FloorShadow))
		}
		if properties.RealWorld.CeilingShadow != nil {
			flags = uint32(data.TileFlag(flags).WithCeilingShadow(*properties.RealWorld.CeilingShadow))
		}
		if properties.RealWorld.SpookyMusic != nil {
			flags = flags & ^uint32(0x00000200)
			if *properties.RealWorld.SpookyMusic {
				flags |= 0x00000200
			}
		}

		entry.Textures = data.TileTextureInfo(textureIDs)
	} else if isCyberspace && properties.Cyberspace != nil {
		var colors = uint16(entry.Textures)

		if properties.Cyberspace.FloorColorIndex != nil {
			colors = (colors & 0xFF00) | (uint16(*properties.Cyberspace.FloorColorIndex) << 0)
		}
		if properties.Cyberspace.CeilingColorIndex != nil {
			colors = (colors & 0x00FF) | (uint16(*properties.Cyberspace.CeilingColorIndex) << 8)
		}
		if properties.Cyberspace.FlightPullType != nil {
			flags = (flags & ^uint32(0x010F0000)) |
				((uint32(*properties.Cyberspace.FlightPullType) & 0xF) << 16) |
				((uint32(*properties.Cyberspace.FlightPullType) & 0x10) << 20)
		}
		if properties.Cyberspace.GameOfLifeSet != nil {
			flags = flags & ^uint32(0x00000040)
			if *properties.Cyberspace.GameOfLifeSet {
				flags |= 0x00000040
			}
		}

		entry.Textures = data.TileTextureInfo(colors)
	}
	entry.Flags = data.TileFlag(flags)

	level.onTileDataChanged()
}

// LevelSurveillanceObjects returns the surveillance objects of this level
func (level *Level) LevelSurveillanceObjects() []model.SurveillanceObject {
	level.mutex.Lock()
	defer level.mutex.Unlock()

	return level.makeSurveillanceObjects(level.readSurveillanceObjects())
}

// SetLevelSurveillanceObject updates one surveillance object
func (level *Level) SetLevelSurveillanceObject(index int, data model.SurveillanceObject) []model.SurveillanceObject {
	level.mutex.Lock()
	defer level.mutex.Unlock()

	sources, deathwatches := level.readSurveillanceObjects()

	if (index >= 0) && (index < len(sources)) {
		if data.SourceIndex != nil {
			sources[index] = int16(*data.SourceIndex)
		}
		if data.DeathwatchIndex != nil {
			deathwatches[index] = int16(*data.DeathwatchIndex)
		}

		storeArray := func(store chunk.BlockStore, data []int16) {
			writer := bytes.NewBuffer(nil)
			binary.Write(writer, binary.LittleEndian, data)
			store.SetBlockData(0, writer.Bytes())
		}
		storeArray(level.surveillanceSourceStore, sources)
		storeArray(level.surveillanceDeathwatchStore, deathwatches)
	}

	return level.makeSurveillanceObjects(sources, deathwatches)
}

func (level *Level) readSurveillanceObjects() (sources []int16, deathwatches []int16) {
	sources = make([]int16, surveillanceSources)
	deathwatches = make([]int16, surveillanceSources)
	sourceData := level.surveillanceSourceStore.BlockData(0)
	deathwatchData := level.surveillanceDeathwatchStore.BlockData(0)

	binary.Read(bytes.NewReader(sourceData), binary.LittleEndian, &sources)
	binary.Read(bytes.NewReader(deathwatchData), binary.LittleEndian, &deathwatches)

	return
}

func (level *Level) makeSurveillanceObjects(sources []int16, deathwatches []int16) []model.SurveillanceObject {
	objects := make([]model.SurveillanceObject, len(sources))

	for index := 0; index < len(sources); index++ {
		objects[index].SourceIndex = intAsPointer(int(sources[index]))
		objects[index].DeathwatchIndex = intAsPointer(int(deathwatches[index]))
	}

	return objects
}
