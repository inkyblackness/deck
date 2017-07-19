package core

import (
	"bytes"
	"encoding/binary"
	"fmt"

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

// Bitmap returns an image of the specified game object.
func (gameObjects *GameObjects) Bitmap(id res.ObjectID, index int) (bmp image.Bitmap, err error) {
	commonProperties := gameObjects.commonProperties(id)
	extraImages := int(commonProperties.Extra >> 4)
	available := 3 + extraImages

	if (index >= 0) && (available > index) {
		blockData := gameObjects.objart.Get(res.ResourceID(0x0546)).BlockData(uint16(gameObjects.objIconOffsets[id] + index))
		bmp, err = image.Read(bytes.NewReader(blockData))
	} else {
		err = fmt.Errorf("Object index out of range: %v[%d]", id, index)
	}
	return
}

// SetBitmap stores an image for the specified game object.
func (gameObjects *GameObjects) SetBitmap(id res.ObjectID, index int, bmp image.Bitmap) (err error) {
	commonProperties := gameObjects.commonProperties(id)
	extraImages := int(commonProperties.Extra >> 4)
	available := 3 + extraImages

	if (index >= 0) && (available > index) {
		blockIndex := uint16(gameObjects.objIconOffsets[id] + index)
		holder := gameObjects.objart.Get(res.ResourceID(0x0546)) //.BlockData(uint16(gameObjects.objIconOffsets[id] + index))
		buf := bytes.NewBuffer(nil)
		image.Write(buf, bmp, image.CompressedBitmap, true, 0)
		holder.SetBlockData(blockIndex, buf.Bytes())
	} else {
		err = fmt.Errorf("Object index out of range: %v[%d]", id, index)
	}
	return
}

// Objects returns an array of all objects
func (gameObjects *GameObjects) Objects() []model.GameObject {
	result := []model.GameObject{}
	linearIndex := uint16(0)

	for classIndex, classDesc := range gameObjects.desc {
		for subclassIndex, subclassDesc := range classDesc.Subclasses {
			for typeIndex := uint32(0); typeIndex < subclassDesc.TypeCount; typeIndex++ {
				var modelData model.GameObject

				modelData.Class = classIndex
				modelData.Subclass = subclassIndex
				modelData.Type = int(typeIndex)
				for i := 0; i < model.LanguageCount; i++ {
					shortName := gameObjects.cybstrng[i].Get(res.ResourceID(0x086D))
					longName := gameObjects.cybstrng[i].Get(res.ResourceID(0x0024))

					modelData.Properties.ShortName[i] = gameObjects.decodeString(shortName.BlockData(linearIndex))
					modelData.Properties.LongName[i] = gameObjects.decodeString(longName.BlockData(linearIndex))
				}
				modelData.Properties.Data = gameObjects.objProperties.Get(res.MakeObjectID(
					res.ObjectClass(classIndex), res.ObjectSubclass(subclassIndex), res.ObjectType(typeIndex)))

				result = append(result, modelData)
				linearIndex++
			}
		}
	}
	return result
}

// SetObjectData stores new object data
func (gameObjects *GameObjects) SetObjectData(id res.ObjectID, newData objprop.ObjectData) objprop.ObjectData {
	oldData := gameObjects.objProperties.Get(id)
	fusedData := oldData

	if newData.Common != nil {
		fusedData.Common = newData.Common
	}
	if newData.Generic != nil {
		fusedData.Generic = newData.Generic
	}
	if newData.Specific != nil {
		fusedData.Specific = newData.Specific
	}
	gameObjects.objProperties.Put(id, fusedData)

	return fusedData
}

func (gameObjects *GameObjects) decodeString(data []byte) *string {
	value := gameObjects.cp.Decode(data)

	return &value
}

func (gameObjects *GameObjects) encodeString(value *string) []byte {
	data := gameObjects.cp.Encode(*value)

	return data
}
