package cmd

import (
	"regexp"
)

var saveCommandExpression = regexp.MustCompile(`^save$`)

func saveCommand(input string) (cmd commandFunction) {
	match := namedMatch(saveCommandExpression, input)

	if len(match) > 0 {
		cmd = func(target Target) string {
			return target.Save()
		}
	}

	return
}
