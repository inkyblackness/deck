package cmd

import (
	"regexp"
	"strconv"
	"strings"
)

var putCommandExpression = regexp.MustCompile(`^put[ ]+(?P<offset>[0-9a-fA-F]+)[ ]+(?P<bytes>([0-9a-fA-F]{1,2}( [0-9a-fA-F]{1,2})*))$`)

func putCommand(input string) (cmd commandFunction) {
	match := namedMatch(putCommandExpression, input)

	if len(match) > 0 {
		bytes := []byte{}
		offset, offsetErr := strconv.ParseUint(match["offset"], 16, 32)
		bytesString := match["bytes"]

		for _, byteString := range strings.Split(bytesString, " ") {
			byteValue, _ := strconv.ParseUint(byteString, 16, 8)
			bytes = append(bytes, byte(byteValue))
		}

		if offsetErr == nil {
			cmd = func(target Target) string {
				return target.Put(uint32(offset), bytes)
			}
		}
	}

	return
}
