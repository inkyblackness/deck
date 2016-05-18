package main

import (
	"fmt"
	"os"

	docopt "github.com/docopt/docopt-go"

	"github.com/inkyblackness/hacker/cmd"
	"github.com/inkyblackness/hacker/core"
)

const (
	// Version contains the current version number
	Version = "0.1.0"
	// Name is the name of the application
	Name = "InkyBlackness Hacker"
	// Title contains a combined string of name and version
	Title = Name + " v." + Version
)

func main() {
	arguments, _ := docopt.Parse(usage(), nil, true, Title, false)

	sourceFiles := arguments["--run"].([]string)
	fileSources := make([]cmd.Source, len(sourceFiles))
	for index, sourceFile := range sourceFiles {
		file, _ := os.Open(sourceFile)
		fileSources[index] = cmd.NewReaderSource(file)
	}

	style := newStandardStyle()
	target := core.NewHacker(style)
	eval := cmd.NewEvaluater(style, target)

	prompter := func() {
		fmt.Printf(style.Prompt()(target.CurrentDirectory(), "> "))
	}

	source := cmd.NewCleaningSource(cmd.NewCombinedSource(
		cmd.NewCombinedSource(fileSources...),
		NewPromptSource(cmd.NewReaderSource(os.Stdin), prompter)))

	fmt.Println(style.Prompt()(Title))
	fmt.Println(style.Prompt()(`Type "quit" to exit`))
	fmt.Println(style.Prompt()(`Remember to keep backups! ...and to salt the fries!`))

	runCommands(source, eval)
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

func runCommands(source cmd.Source, eval *cmd.Evaluater) {
	quit := false

	for !quit {
		input, finished := source.Next()

		if finished || input == "quit" {
			quit = true
		} else if input != "" {
			result := eval.Evaluate(input)
			fmt.Println(result)
		}
	}
}
