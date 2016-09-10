package core

import (
	"bytes"
	"encoding/binary"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/res/image"
	"github.com/inkyblackness/res/objprop"
	"github.com/inkyblackness/res/text"
	"github.com/inkyblackness/shocked-core/io"
	model "github.com/inkyblackness/shocked-model"
)

// GameObjects provides access to game-global information about objects.
type GameObjects struct {
	cybstrng [model.LanguageCount]chunk.Store
	cp       text.Codepage

	desc          []objprop.ClassDescriptor
	objProperties objprop.Store
	objart        chunk.Store

	objIconOffsets map[res.ObjectID]int
	mapIconOffsets map[res.ObjectID]int
}

// NewGameObjects returns a new instance of GameObjects.
func NewGameObjects(library io.StoreLibrary) (gameObjects *GameObjects, err error) {
	var cybstrng [model.LanguageCount]chunk.Store
	var objart chunk.Store
	var objProperties objprop.Store

	if err == nil {
		objart, err = library.ChunkStore("objart.res")
	}
	if err == nil {
		objProperties, err = library.ObjpropStore("objprop.dat")
	}
	for i := 0; i < model.LanguageCount && err == nil; i++ {
		cybstrng[i], err = library.ChunkStore(localized[i].cybstrng)
	}
	if err == nil {
		gameObjects = &GameObjects{
			cybstrng:       cybstrng,
			cp:             text.DefaultCodepage(),
			desc:           objprop.StandardProperties(),
			objProperties:  objProperties,
			objart:         objart,
			objIconOffsets: make(map[res.ObjectID]int),
			mapIconOffsets: make(map[res.ObjectID]int)}

		offset := 1
		for classIndex, classDesc := range gameObjects.desc {
			for subclassIndex, subclassDesc := range classDesc.Subclasses {
				for typeIndex := uint32(0); typeIndex < subclassDesc.TypeCount; typeIndex++ {
					objID := res.MakeObjectID(res.ObjectClass(classIndex), res.ObjectSubclass(subclassIndex), res.ObjectType(typeIndex))
					commonProperties := gameObjects.commonProperties(objID)
					extraImages := int(commonProperties.Extra >> 4)

					gameObjects.objIconOffsets[objID] = offset
					offset = offset + 3 + extraImages
					gameObjects.mapIconOffsets[objID] = offset - 1
				}
			}
		}
	}

	return
}

func (gameObjects *GameObjects) commonProperties(id res.ObjectID) data.CommonObjectProperties {
	properties := gameObjects.objProperties.Get(id)
	commonData := bytes.NewReader(properties.Common)
	var commonProperties data.CommonObjectProperties
	binary.Read(commonData, binary.LittleEndian, &commonProperties)

	return commonProperties
}

// Icon returns the icon image of the specified game object.
// It first tries to return the bitmap for the map icon. If that is all transparent,
// the function reverts to the object icon.
func (gameObjects *GameObjects) Icon(id res.ObjectID) (bmp image.Bitmap) {
	mapIconBlockData := gameObjects.objart.Get(res.ResourceID(0x0546)).BlockData(uint16(gameObjects.mapIconOffsets[id]))

	bmp, _ = image.Read(bytes.NewReader(mapIconBlockData))

	allZero := true
	for row := 0; row < int(bmp.ImageHeight()); row++ {
		for _, b := range bmp.Row(row) {
			if b != 0x00 {
				allZero = false
			}
		}
	}
	if allZero {
		objIconBlockData := gameObjects.objart.Get(res.ResourceID(0x0546)).BlockData(uint16(gameObjects.objIconOffsets[id]))
		bmp, _ = image.Read(bytes.NewReader(objIconBlockData))
	}

	return
}

// Properties returns the current properties of a specific game object.
func (gameObjects *GameObjects) Properties(id res.ObjectID) model.GameObjectProperties {
	prop := model.GameObjectProperties{}
	index := objprop.ObjectIDToIndex(gameObjects.desc, id)

	for i := 0; i < model.LanguageCount; i++ {
		shortName := gameObjects.cybstrng[i].Get(res.ResourceID(0x086D))
		longName := gameObjects.cybstrng[i].Get(res.ResourceID(0x0024))

		prop.ShortName[i] = gameObjects.decodeString(shortName.BlockData(uint16(index)))
		prop.LongName[i] = gameObjects.decodeString(longName.BlockData(uint16(index)))
	}

	return prop
}

func (gameObjects *GameObjects) decodeString(data []byte) *string {
	value := gameObjects.cp.Decode(data)

	return &value
}

func (gameObjects *GameObjects) encodeString(value *string) []byte {
	data := gameObjects.cp.Encode(*value)

	return data
}
