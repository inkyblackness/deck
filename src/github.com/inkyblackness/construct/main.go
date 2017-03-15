package main

import (
	"fmt"
	"os"

	docopt "github.com/docopt/docopt-go"

	"github.com/inkyblackness/res/chunk/dos"

	"github.com/inkyblackness/construct/chunks"
)

const (
	// Version contains the current version number
	Version = "0.2.0"
	// Name is the name of the application
	Name = "InkyBlackness Construct"
	// Title contains a combined string of name and version
	Title = Name + " v." + Version
)

func main() {
	arguments, _ := docopt.Parse(usage(), nil, true, Title, false)
	outFileName := arguments["--file"].(string)
	levels := arguments["--levels"].(string)
	solid := arguments["--solid"].(bool)
	writer, errOut := os.Create(outFileName)
	if errOut != nil {
		fmt.Printf("Error creating destination: %v\n", errOut)
	}

	chunkConsumer := dos.NewChunkConsumer(writer)

	chunks.AddArchiveName(chunkConsumer, "Starting Game")
	chunks.AddGameState(chunkConsumer)
	for levelID, levelType := range levels {
		isRealWorld := levelType == 'r'
		isCyberspace := levelType == 'c'

		if isRealWorld || isCyberspace {
			chunks.AddLevel(chunkConsumer, levelID, solid, isCyberspace)
		}
	}

	chunkConsumer.Finish()
}

func usage() string {
	return Title + `

Usage:
  construct [--file=<file-name>] [--solid] [--levels=<levels>]
  construct -h | --help
  construct --version

Options:
  --file=<file-name>  Specifies the target file name. [default: archive.dat]
  --solid             Creates an entirely solid map; Exception: Starting tile on level 1.
  --levels=<levels>   Specifies which levels to create. String of 'n', 'r' and 'c'. [default: nr]
  -h --help           Show this screen.
  --version           Show version.
`
}
