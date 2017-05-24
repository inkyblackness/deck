package model

import (
	"strings"

	"github.com/inkyblackness/res/data/gameobj"
	"github.com/inkyblackness/res/objprop"
	"github.com/inkyblackness/shocked-model"
)

// GameObject describes one object available in the game.
type GameObject struct {
	id ObjectID

	longName  [model.LanguageCount]string
	shortName [model.LanguageCount]string

	data objprop.ObjectData
}

// ID returns the identification of the object.
func (object *GameObject) ID() ObjectID {
	return object.id
}

// DisplayName returns the name of the object for editor purposes
func (object *GameObject) DisplayName() string {
	return strings.Replace(object.longName[0], "\n", " ", -1)
}

// CommonData returns the common data for this object
func (object *GameObject) CommonData() []byte {
	return object.data.Common
}

// GenericData returns the generic data for this object
func (object *GameObject) GenericData() []byte {
	return object.data.Generic
}

// SpecificData returns the specific data for this object
func (object *GameObject) SpecificData() []byte {
	return object.data.Specific
}

// CommonHitpoints returns the common value for hitpoints
func (object *GameObject) CommonHitpoints() int {
	return int(gameobj.CommonProperties(object.CommonData()).Get("DefaultHitpoints"))
}
