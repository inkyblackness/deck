package cmd

import (
	"regexp"
)

var infoCommandExpression = regexp.MustCompile(`^info$`)

func infoCommand(input string) (cmd commandFunction) {
	match := namedMatch(infoCommandExpression, input)

	if len(match) > 0 {
		cmd = func(target Target) string {
			return target.Info()
		}
	}

	return
}
