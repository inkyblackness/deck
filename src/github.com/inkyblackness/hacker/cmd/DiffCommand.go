package cmd

import (
	"regexp"
)

var diffCommandExpression = regexp.MustCompile(`^diff[ ]+(?P<target>.+)$`)

func diffCommand(input string) (cmd commandFunction) {
	match := namedMatch(diffCommandExpression, input)

	if len(match) > 0 {
		cmd = func(target Target) string {
			return target.Diff(match["target"])
		}
	}

	return
}
