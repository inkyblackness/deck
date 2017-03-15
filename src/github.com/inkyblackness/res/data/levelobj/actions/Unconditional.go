package actions

import (
	"github.com/inkyblackness/res/data/interpreters"
)

// Unconditional returns the description of actions without a condition.
func Unconditional() *interpreters.Description {

	return unconditionalAction
}
