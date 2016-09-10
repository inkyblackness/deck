package model

import (
	"fmt"
)

// GameTextureKey is a key for game textures.
type GameTextureKey struct {
	id int
}

// GameTextureKeyFor returns a game texture key instance.
func GameTextureKeyFor(id int) GameTextureKey {
	return GameTextureKey{id}
}

// ID returns the actual identifier value.
func (key GameTextureKey) ID() int {
	return key.id
}

func (key GameTextureKey) String() string {
	return fmt.Sprintf("%d", key.ID())
}
