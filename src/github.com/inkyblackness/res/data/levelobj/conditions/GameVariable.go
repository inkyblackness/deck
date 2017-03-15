package conditions

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var gameVariable = interpreters.New().
	With("VariableKey", 0, 2).
	With("Value", 2, 1).
	With("MessageIndex", 3, 1)

// GameVariable returns the description for a game variable condition
func GameVariable() *interpreters.Description {
	return gameVariable
}
