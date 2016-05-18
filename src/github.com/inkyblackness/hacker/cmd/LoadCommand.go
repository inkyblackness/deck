package cmd

import (
	"regexp"
)

var loadCommandExpression = regexp.MustCompile(`^load[ ]+\"(?P<path1>[^"]+)\"[ ]*(\"(?P<path2>([^"]+))\")?$`)

func loadCommand(input string) (cmd commandFunction) {
	match := namedMatch(loadCommandExpression, input)

	if len(match) > 0 {
		cmd = func(target Target) string {
			return target.Load(match["path1"], match["path2"])
		}
	}

	return
}
