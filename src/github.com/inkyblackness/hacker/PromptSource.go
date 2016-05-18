package main

import (
	"github.com/inkyblackness/hacker/cmd"
)

type promptSource struct {
	source   cmd.Source
	prompter func()
}

// NewPromptSource returns a source that calls the given prompter whenever a next
// command is requested.
func NewPromptSource(source cmd.Source, prompter func()) cmd.Source {
	return &promptSource{source: source, prompter: prompter}
}

func (source *promptSource) Next() (cmd string, finished bool) {
	source.prompter()

	return source.source.Next()
}
