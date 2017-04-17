package levelobj

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var baseHardware = interpreters.New().
	With("Version", 0, 1).As(interpreters.RangedValue(0, 4))
