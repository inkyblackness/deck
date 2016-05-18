package model

// TextureProperties contain all the behavioural settings of a texture.
type TextureProperties struct {
	// Name is the (language specific) name of the texture.
	Name [LanguageCount]*string `json:"name"`
	// CantBeUsed is a (language specific) text for usage failures.
	CantBeUsed [LanguageCount]*string `json:"cantBeUsed"`

	//Climbable bool `json:"climbable"`
}
