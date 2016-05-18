package cmd

import (
	"regexp"
)

func namedMatch(exp *regexp.Regexp, input string) map[string]string {
	result := make(map[string]string)
	match := exp.FindStringSubmatch(input)

	if len(match) > 0 {
		for i, name := range exp.SubexpNames() {
			result[name] = match[i]
		}
	}

	return result
}
