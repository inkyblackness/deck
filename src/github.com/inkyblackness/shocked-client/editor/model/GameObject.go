package model

import (
	"strings"

	"github.com/inkyblackness/shocked-model"
)

// GameObject describes one object available in the game.
type GameObject struct {
	id ObjectID

	longName  [model.LanguageCount]string
	shortName [model.LanguageCount]string
}

// ID returns the identification of the object.
func (object *GameObject) ID() ObjectID {
	return object.id
}

// DisplayName returns the name of the object for editor purposes
func (object *GameObject) DisplayName() string {
	return strings.Replace(object.longName[0], "\n", " ", -1)
}
