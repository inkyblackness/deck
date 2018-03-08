package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"

	"github.com/inkyblackness/hacker/cmd"
	"github.com/inkyblackness/hacker/core"
)

const (
	// Version contains the current version number
	Version = "1.1.0"
	// Name is the name of the application
	Name = "InkyBlackness Hacker"
	// Title contains a combined string of name and version
	Title = Name + " v." + Version
)

func main() {
	arguments, optErr := docopt.Parse(usage(), nil, true, Title, false)
	if optErr != nil {
		fmt.Printf("Couldn't parse arguments: %v\n", optErr)
		return
	}

	sourceFiles := arguments["--run"].([]string)
	fileSources := make([]cmd.Source, len(sourceFiles))
	for index, sourceFile := range sourceFiles {
		file, sourceErr := os.Open(sourceFile)
		if sourceErr != nil {
			fmt.Printf("Couldn't Open source %v\n", sourceFile)
			return
		}
		fileSources[index] = cmd.NewReaderSource(file)
	}

	style := newStandardStyle()
	target := core.NewHacker(style)
	eval := cmd.NewEvaluater(style, target)

	prompter := func() string {
		return style.Prompt()(target.CurrentDirectory(), "> ")
	}
	readLineSource := NewReadLineSource(prompter)
	defer readLineSource.Close() // nolint:errcheck

	source := cmd.NewCleaningSource(cmd.NewCombinedSource(
		cmd.NewCombinedSource(fileSources...),
		readLineSource))

	style.Println(style.Prompt()(Title))
	style.Println(style.Prompt()(`Type "quit" to exit`))
	style.Println(style.Prompt()(`Remember to keep backups! ...and to salt the fries!`))

	runCommands(style, source, eval)
}

func usage() string {
	return Title + `

Usage:
  hacker [--run <file>...]
  hacker -h | --help
  hacker --version

Options:
  -h --help     Show this screen.
  --version     Show version.
  --run <file>  Run the specified file. Can be repeated to run several in sequence.`
}

func runCommands(style *standardStyle, source cmd.Source, eval *cmd.Evaluater) {
	quit := false

	for !quit {
		input, finished := source.Next()

		if finished || input == "quit" {
			quit = true
		} else if input != "" {
			result := eval.Evaluate(input)
			style.Println(result)
		}
	}
}
