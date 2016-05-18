package cmd

import (
	"regexp"
)

var cdCommandExpression = regexp.MustCompile(`^cd[ ]+(?P<path>.+)$`)

func changeDirectoryCommand(input string) (cmd commandFunction) {
	match := namedMatch(cdCommandExpression, input)

	if len(match) > 0 {
		cmd = func(target Target) string {
			return target.ChangeDirectory(match["path"])
		}
	}

	return
}
