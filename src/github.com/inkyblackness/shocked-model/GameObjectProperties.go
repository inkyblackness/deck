package model

import "github.com/inkyblackness/res/objprop"

// GameObjectProperties globally describe a game object.
type GameObjectProperties struct {
	// ShortName is the (language specific) short name of the object.
	ShortName [LanguageCount]*string
	// LongName is the (language specific) long name of the object.
	LongName [LanguageCount]*string

	Data objprop.ObjectData
}
