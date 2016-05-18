package cmd

import (
	"regexp"
)

var dumpCommandExpression = regexp.MustCompile(`^dump$`)

func dumpCommand(input string) (cmd commandFunction) {
	match := namedMatch(dumpCommandExpression, input)

	if len(match) > 0 {
		cmd = func(target Target) string {
			return target.Dump()
		}
	}

	return
}
