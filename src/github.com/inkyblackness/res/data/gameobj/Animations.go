package gameobj

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data/interpreters"
)

var animationGenerics = interpreters.New().
	With("FrameTime", 0, 1).
	With("EmitsLight", 1, 1).As(interpreters.EnumValue(map[uint32]string{0: "No", 1: "Yes"}))

func initAnimations() {
	objClass := res.ObjectClass(11)

	genericDescriptions[objClass] = animationGenerics

}
