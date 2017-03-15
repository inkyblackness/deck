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

		crossrefListStore: store.Get(res.ResourceID(baseStoreID + 9))}

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

// Properties returns the properties of the level.
func (level *Level) Properties() (result model.LevelProperties) {
	blockData := level.store.Get(res.ResourceID(4000 + level.id*100 + 4)).BlockData(0)
	reader := bytes.NewReader(blockData)
	var info data.LevelInformation

	binary.Read(reader, binary.LittleEndian, &info)
	result.CyberspaceFlag = info.IsCyberspace()
	result.HeightShift = int(info.HeightShift)

	return
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

func bytesToIntArray(bs []byte) []int {
	result := make([]int, len(bs))
	for index, value := range bs {
		result[index] = int(value)
	}

	return result
}

// Objects returns an array of all used objects.
func (level *Level) Objects() []model.LevelObject {
	level.mutex.Lock()
	defer level.mutex.Unlock()

	var result []model.LevelObject

	for index, rawEntry := range level.objectList {
		if rawEntry.IsInUse() {
			entry := model.LevelObject{
				Identifiable: model.Identifiable{ID: fmt.Sprintf("%d", index)},
				Class:        int(rawEntry.Class),
				Subclass:     int(rawEntry.Subclass),
				Type:         int(rawEntry.Type)}

			entry.BaseProperties.TileX = int(rawEntry.X >> 8)
			entry.BaseProperties.FineX = int(rawEntry.X & 0xFF)
			entry.BaseProperties.TileY = int(rawEntry.Y >> 8)
			entry.BaseProperties.FineY = int(rawEntry.Y & 0xFF)
			entry.BaseProperties.Z = int(rawEntry.Z)

			meta := data.LevelObjectClassMetaEntry(rawEntry.Class)
			classStore := level.store.Get(res.ResourceID(4000 + level.id*100 + 10 + entry.Class))
			blockData := classStore.BlockData(0)
			startOffset := meta.EntrySize * int(rawEntry.ClassTableIndex)
			if (startOffset + meta.EntrySize) > len(blockData) {
				fmt.Printf("!!!!! class %d meta says %d bytes size, can't reach index %d in blockData %d",
					int(entry.Class), meta.EntrySize, rawEntry.ClassTableIndex, len(blockData))
			} else {
				entry.ClassData = bytesToIntArray(blockData[startOffset+data.LevelObjectPrefixSize : startOffset+meta.EntrySize])
			}

			result = append(result, entry)
		}
	}

	return result
}

// AddObject adds a new object at given tile.
func (level *Level) AddObject(template *model.LevelObjectTemplate) (createdIndex int, err error) {
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
	createdIndex = int(objectIndex)

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
	objectEntry.Hitpoints = 1

	objectEntry.CrossReferenceTableIndex = uint16(crossrefIndex)
	crossrefEntry.LevelObjectTableIndex = uint16(objectIndex)

	objectEntry.ClassTableIndex = uint16(classIndex)
	classEntry.LevelObjectTableIndex = uint16(objectIndex)

	classStore.SetBlockData(0, classTable.Encode())

	objWriter := bytes.NewBuffer(nil)
	binary.Write(objWriter, binary.LittleEndian, level.objectList)
	level.objectListStore.SetBlockData(0, objWriter.Bytes())

	level.crossrefListStore.SetBlockData(0, level.crossrefList.Encode())
	level.onTileDataChanged()

	return
}

func (level *Level) readTable(levelBlockID int, value interface{}) {
	blockData := level.store.Get(res.ResourceID(4000 + level.id*100 + levelBlockID)).BlockData(0)
	reader := bytes.NewReader(blockData)

	binary.Read(reader, binary.LittleEndian, value)
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

	{
		var properties model.RealWorldTileProperties
		var textureIDs = uint16(entry.Textures)

		properties.WallTexture = new(int)
		*properties.WallTexture = int(textureIDs & 0x3F)
		properties.CeilingTexture = new(int)
		*properties.CeilingTexture = int((textureIDs >> 6) & 0x1F)
		properties.CeilingTextureRotations = new(int)
		*properties.CeilingTextureRotations = int((entry.Ceiling >> 5) & 0x03)
		properties.FloorTexture = new(int)
		*properties.FloorTexture = int((textureIDs >> 11) & 0x1F)
		properties.FloorTextureRotations = new(int)
		*properties.FloorTextureRotations = int((entry.Floor >> 5) & 0x03)

		properties.UseAdjacentWallTexture = new(bool)
		*properties.UseAdjacentWallTexture = (entry.Flags & 0x00000100) != 0
		properties.WallTextureOffset = new(model.HeightUnit)
		*properties.WallTextureOffset = model.HeightUnit(entry.Flags & 0x0000001F)

		result.RealWorld = &properties
	}

	return
}

// SetTileProperties sets the properties of a specific tile.
func (level *Level) SetTileProperties(x, y int, properties model.TileProperties) {
	level.mutex.Lock()
	defer level.mutex.Unlock()

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
	if properties.RealWorld != nil {
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

		entry.Textures = data.TileTextureInfo(textureIDs)
	}
	entry.Flags = data.TileFlag(flags)

	level.onTileDataChanged()
}
