package gameobj

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data/interpreters"
)

var cyberItems = interpreters.New().
	Refining("ColorScheme", 0, 6, cyberColorScheme, interpreters.Always)

func initItems() {
	objClass := res.ObjectClass(8)

	setSpecific(objClass, 5, cyberItems)
}
