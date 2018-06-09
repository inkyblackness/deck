package model

import (
	"github.com/inkyblackness/shocked-model"
)

// GameTexture describes one texture available to the game.
type GameTexture struct {
	id int

	properties model.TextureProperties
}

func newGameTexture(id int) *GameTexture {
	return &GameTexture{id: id}
}

func nullGameTexture(id int) *GameTexture {
	texture := newGameTexture(id)
	valFalse := false
	valZero := 0
	valEmpty := ""

	for lang := 0; lang < model.LanguageCount; lang++ {
		texture.properties.Name[lang] = &valEmpty
		texture.properties.CantBeUsed[lang] = &valEmpty
	}
	texture.properties.Climbable = &valFalse
	texture.properties.AnimationGroup = &valZero
	texture.properties.AnimationIndex = &valZero
	texture.properties.TransparencyControl = &valZero
	return texture
}

// ID uniquely identifies the texture in the game.
func (texture *GameTexture) ID() int {
	return texture.id
}

// Climbable returns whether the texture can be climbed.
func (texture *GameTexture) Climbable() bool {
	return *texture.properties.Climbable
}

// TransparencyControl returns how the pixel data shall be interpreted.
func (texture *GameTexture) TransparencyControl() int {
	return *texture.properties.TransparencyControl
}

// AnimationGroup associates textures of an animation.
func (texture *GameTexture) AnimationGroup() int {
	return *texture.properties.AnimationGroup
}

// AnimationIndex places the texture within the group.
func (texture *GameTexture) AnimationIndex() int {
	return *texture.properties.AnimationIndex
}

// Name returns the name of the texture in given language.
func (texture *GameTexture) Name(language model.ResourceLanguage) string {
	return *texture.properties.Name[language.ToIndex()]
}

// UseText returns the usage text of the texture in given language.
func (texture *GameTexture) UseText(language model.ResourceLanguage) string {
	return *texture.properties.CantBeUsed[language.ToIndex()]
}
