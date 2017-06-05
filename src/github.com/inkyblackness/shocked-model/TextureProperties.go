package model

// TextureProperties contain all the behavioural settings of a texture.
type TextureProperties struct {
	// Name is the (language specific) name of the texture.
	Name [LanguageCount]*string
	// CantBeUsed is a (language specific) text for usage failures.
	CantBeUsed [LanguageCount]*string

	// Climbable determines whether the texture can be climbed, such as ladders.
	Climbable *bool
	// TransparencyControl determines how to interpret bitmap data.
	TransparencyControl *int
	// AnimationGroup relates textures for an animation.
	AnimationGroup *int
	// AnimationIndex identifies a texture within a group.
	AnimationIndex *int
}
