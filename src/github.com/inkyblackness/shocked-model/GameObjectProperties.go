package model

// GameObjectProperties globally describe a game object.
type GameObjectProperties struct {
	// ShortName is the (language specific) short name of the object.
	ShortName [LanguageCount]*string `json:"shortName"`
	// LongName is the (language specific) long name of the object.
	LongName [LanguageCount]*string `json:"longName"`
}
