package core

import (
	"fmt"
	"os"

	"github.com/inkyblackness/hacker/styling"
)

func fileNames(files []os.FileInfo) (names []string) {
	for _, file := range files {
		if !file.IsDir() {
			names = append(names, file.Name())
		}
	}
	return
}

type styledData struct {
	value     byte
	styleFunc styling.StyleFunc
}

func createDump(data []styledData) (result string) {
	rightPad := func(input string, missing int) string {
		return fmt.Sprintf(fmt.Sprintf("%%s%%%ds", missing), input, "")
	}
	hexDump := ""
	hexLen := 0
	asciiDump := ""
	asciiLen := 0

	addLine := func(offset int) {
		result = result + fmt.Sprintf("%04X %s  %s\n", offset, rightPad(hexDump, 49-hexLen), rightPad(asciiDump, 17-asciiLen))
		hexDump = ""
		hexLen = 0
		asciiDump = ""
		asciiLen = 0
	}

	for index, entry := range data {
		value := entry.value

		if index == 0 {
		} else if (index % 16) == 0 {
			addLine(((index / 16) - 1) * 16)
		} else if (index % 8) == 0 {
			hexDump += " "
			asciiDump += " "
			hexLen++
			asciiLen++
		}

		hexDump += entry.styleFunc(fmt.Sprintf(" %02X", value))
		hexLen += 3
		if value >= 0x20 && value < 0x80 {
			asciiDump += entry.styleFunc(string(value))
		} else {
			asciiDump += entry.styleFunc(".")
		}
		asciiLen += 1
	}
	if hexDump != "" {
		addLine(((len(data) - 1) / 16) * 16)
	}
	return
}
