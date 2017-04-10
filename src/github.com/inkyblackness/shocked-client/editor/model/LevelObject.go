package model

import (
	"strconv"

	"github.com/inkyblackness/shocked-model"
)

// LevelObject describes one object within a level
type LevelObject struct {
	class            int
	properties       *model.LevelObjectProperties
	centerX, centerY float32
	index            int
}

func newLevelObject(data *model.LevelObject) *LevelObject {
	index, _ := strconv.ParseInt(data.ID, 10, 32)
	obj := &LevelObject{
		class: data.Class,
		index: int(index)}
	obj.onPropertiesChanged(&data.Properties)

	return obj
}

func (obj *LevelObject) onPropertiesChanged(properties *model.LevelObjectProperties) {
	obj.properties = properties
	obj.centerX = float32((*obj.properties.TileX << 8) + *obj.properties.FineX)
	obj.centerY = float32((*obj.properties.TileY << 8) + *obj.properties.FineY)
}

// Index returns the object's index within the level.
func (obj *LevelObject) Index() int {
	return obj.index
}

// ID returns the object ID of the object
func (obj *LevelObject) ID() ObjectID {
	return MakeObjectID(obj.class, *obj.properties.Subclass, *obj.properties.Type)
}

// ClassData returns the raw data for the level object.
func (obj *LevelObject) ClassData() []byte {
	return obj.properties.ClassData
}

// TileX returns the x-coordinate (tile)
func (obj *LevelObject) TileX() int {
	return *obj.properties.TileX
}

// FineX returns the x-coordinate (fine, within tile)
func (obj *LevelObject) FineX() int {
	return *obj.properties.FineX
}

// TileY returns the y-coordinate (tile)
func (obj *LevelObject) TileY() int {
	return *obj.properties.TileY
}

// FineY returns the y-coordinate (fine, within tile)
func (obj *LevelObject) FineY() int {
	return *obj.properties.FineY
}

// Z returns the z-coordinate (placement height) of the object
func (obj *LevelObject) Z() int {
	return *obj.properties.Z
}

// Center returns the location of the object within the map
func (obj *LevelObject) Center() (x, y float32) {
	return obj.centerX, obj.centerY
}

// RotationX returns the rotation around the X-axis
func (obj *LevelObject) RotationX() int {
	return *obj.properties.RotationX
}

// RotationY returns the rotation around the Y-axis
func (obj *LevelObject) RotationY() int {
	return *obj.properties.RotationY
}

// RotationZ returns the rotation around the Y-axis
func (obj *LevelObject) RotationZ() int {
	return *obj.properties.RotationZ
}

// Hitpoints returns the hitpoints of the object
func (obj *LevelObject) Hitpoints() int {
	return *obj.properties.Hitpoints
}
