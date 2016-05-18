package cmd

import (
	"regexp"
)

var queryCommandExpression = regexp.MustCompile(`^query[ ]+(?P<info>.+)$`)

func queryCommand(input string) (cmd commandFunction) {
	match := namedMatch(queryCommandExpression, input)

	if len(match) > 0 {
		cmd = func(target Target) string {
			return target.Query(match["info"])
		}
	}

	return
}
