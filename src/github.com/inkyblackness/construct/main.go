package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"

	"github.com/inkyblackness/construct/chunks"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/chunk/resfile"
)

const (
	// Version contains the current version number
	Version = "1.1.0"
	// Name is the name of the application
	Name = "InkyBlackness Construct"
	// Title contains a combined string of name and version
	Title = Name + " v." + Version
)

func usage() string {
	return Title + `

Usage:
  construct archive [--file=<file-name>] [--solid] [--levels=<levels>]
  construct resource [--file=<file-name>]
  construct -h | --help
  construct --version

Options:
  --file=<file-name>  Specifies the target file name. Default: archive.dat or empty.res
  --solid             Creates an entirely solid map; Exception: Starting tile on level 1.
  --levels=<levels>   Specifies which levels to create. String of 'n', 'r' and 'c'. [default: nr]
  -h --help           Show this screen.
  --version           Show version.
`
}

func main() {
	arguments, argErr := docopt.Parse(usage(), nil, true, Title, false)
	if argErr != nil {
		fmt.Printf("Failed to parse arguments: %v\n", argErr)
		return
	}

	if arguments["archive"].(bool) {
		outFileName := orElse(arguments["--file"], "archive.dat").(string)
		levels := arguments["--levels"].(string)
		solid := arguments["--solid"].(bool)

		store := chunk.NewProviderBackedStore(chunk.NullProvider())

		chunks.AddArchiveName(store, "Starting Game")
		chunks.AddGameState(store)
		for levelID, levelType := range levels {
			isRealWorld := levelType == 'r'
			isCyberspace := levelType == 'c'

			if isRealWorld || isCyberspace {
				chunks.AddLevel(store, levelID, solid, isCyberspace)
			}
		}

		writeResourceFile(outFileName, store)
	} else if arguments["resource"].(bool) {
		outFileName := orElse(arguments["--file"], "empty.res").(string)

		writeResourceFile(outFileName, chunk.NullProvider())
	}
}

func orElse(optional, defaultValue interface{}) (result interface{}) {
	result = optional
	if result == nil {
		result = defaultValue
	}
	return
}

func writeResourceFile(fileName string, provider chunk.Provider) {
	writer, errOut := os.Create(fileName)
	if errOut != nil {
		fmt.Printf("Error creating file: %v\n", errOut)
		return
	}
	defer func() {
		errOut = writer.Close()
		if errOut != nil {
			fmt.Printf("Error closing file: %v\n", errOut)
		}
	}()

	err := resfile.Write(writer, provider)
	if err != nil {
		fmt.Printf("Error writing resource: %v\n", err)
	}
}
