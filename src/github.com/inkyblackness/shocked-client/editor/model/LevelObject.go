package model

import (
	"strconv"

	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-model"
)

// LevelObject describes one object within a level
type LevelObject struct {
	data          *model.LevelObject
	index         int
	iconRetriever graphics.BitmapRetriever
}

// NewLevelObject returns a new instance of a LevelObject.
func NewLevelObject(data *model.LevelObject, iconRetriever graphics.BitmapRetriever) *LevelObject {
	index, _ := strconv.ParseInt(data.ID, 10, 32)
	obj := &LevelObject{
		data:          data,
		index:         int(index),
		iconRetriever: iconRetriever}

	return obj
}

// Index returns the object's index within the level.
func (obj *LevelObject) Index() int {
	return obj.index
}

// ID returns the object ID of the object
func (obj *LevelObject) ID() (class, subclass, objType int) {
	return obj.data.Class, obj.data.Subclass, obj.data.Type
}

// Center returns the location of the object within the map
func (obj *LevelObject) Center() (x, y float32) {
	x = (float32(obj.data.BaseProperties.TileX) + (float32(obj.data.BaseProperties.FineX) / float32(0xFF))) * 32.0
	y = (float32(63-obj.data.BaseProperties.TileY) + (float32(0xFF-obj.data.BaseProperties.FineY) / float32(0xFF))) * 32.0

	return
}

// Icon returns the icon bitmap for this level object.
func (obj *LevelObject) Icon() *graphics.BitmapTexture {
	return obj.iconRetriever()
}

// Size returns the dimension on the map.
func (obj *LevelObject) Size() (width, height float32) {
	return 16.0, 16.0
}
