package model

import (
	"strconv"

	"github.com/inkyblackness/shocked-model"
)

// LevelObject describes one object within a level
type LevelObject struct {
	data             *model.LevelObject
	centerX, centerY float32
	index            int
}

func newLevelObject(data *model.LevelObject) *LevelObject {
	index, _ := strconv.ParseInt(data.ID, 10, 32)
	obj := &LevelObject{
		data:  data,
		index: int(index)}

	obj.centerX = float32((obj.data.BaseProperties.TileX << 8) + obj.data.BaseProperties.FineX)
	obj.centerY = float32((obj.data.BaseProperties.TileY << 8) + obj.data.BaseProperties.FineY)

	return obj
}

// Index returns the object's index within the level.
func (obj *LevelObject) Index() int {
	return obj.index
}

// ID returns the object ID of the object
func (obj *LevelObject) ID() ObjectID {
	return MakeObjectID(obj.data.Class, obj.data.Subclass, obj.data.Type)
}

// ClassData returns the raw data for the level object.
func (obj *LevelObject) ClassData() []byte {
	data := make([]byte, len(obj.data.ClassData))
	for index, value := range obj.data.ClassData {
		data[index] = byte(value)
	}
	return data
}

// Center returns the location of the object within the map
func (obj *LevelObject) Center() (x, y float32) {
	return obj.centerX, obj.centerY
}
