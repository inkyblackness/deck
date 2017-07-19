package gameobj

import (
	"fmt"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data/interpreters"
)

var animationGenerics = interpreters.New().
	With("FrameTime", 0, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int64) string {
		return fmt.Sprintf("%3.0f millisec - raw: %d", (float64(value)*900)/255.0, value)
	})).
	With("EmitsLight", 1, 1).As(interpreters.EnumValue(map[uint32]string{0: "No", 1: "Yes"}))

func initAnimations() {
	objClass := res.ObjectClass(11)

	genericDescriptions[objClass] = animationGenerics

}
