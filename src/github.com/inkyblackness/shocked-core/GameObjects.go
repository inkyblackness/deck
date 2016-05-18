package core

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/objprop"
	"github.com/inkyblackness/res/text"
	"github.com/inkyblackness/shocked-core/io"
	model "github.com/inkyblackness/shocked-model"
)

type GameObjects struct {
	desc     []objprop.ClassDescriptor
	cybstrng [model.LanguageCount]chunk.Store
	cp       text.Codepage
}

func NewGameObjects(library io.StoreLibrary) (gameObjects *GameObjects, err error) {
	var cybstrng [model.LanguageCount]chunk.Store

	for i := 0; i < model.LanguageCount && err == nil; i++ {
		cybstrng[i], err = library.ChunkStore(localized[i].cybstrng)
	}
	if err == nil {
		gameObjects = &GameObjects{cybstrng: cybstrng, cp: text.DefaultCodepage()}
		gameObjects.desc = objprop.StandardProperties()
	}

	return
}

func (gameObjects *GameObjects) Properties(id res.ObjectID) model.GameObjectProperties {
	prop := model.GameObjectProperties{}
	index := objprop.ObjectIDToIndex(gameObjects.desc, id)

	for i := 0; i < model.LanguageCount; i++ {
		shortName := gameObjects.cybstrng[i].Get(res.ResourceID(0x086D))
		longName := gameObjects.cybstrng[i].Get(res.ResourceID(0x0024))

		prop.ShortName[i] = gameObjects.DecodeString(shortName.BlockData(uint16(index)))
		prop.LongName[i] = gameObjects.DecodeString(longName.BlockData(uint16(index)))
	}

	return prop
}

func (gameObjects *GameObjects) DecodeString(data []byte) *string {
	value := gameObjects.cp.Decode(data)

	return &value
}

func (gameObjects *GameObjects) EncodeString(value *string) []byte {
	data := gameObjects.cp.Encode(*value)

	return data
}
