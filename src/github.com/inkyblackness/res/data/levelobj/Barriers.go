package levelobj

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var baseBarrier = interpreters.New().
	With("LockVariableIndex", 0, 2).As(interpreters.RangedValue(0, 0x1FF)).
	With("LockMessageIndex", 2, 1).
	With("ForceDoorColor", 3, 1).
	With("RequiredAccessLevel", 4, 1).
	With("AutoCloseTime", 5, 1).
	With("OtherObjectIndex", 6, 2).As(interpreters.ObjectIndex())

func initBarriers() interpreterRetriever {
	return newInterpreterLeaf(baseBarrier)
}
