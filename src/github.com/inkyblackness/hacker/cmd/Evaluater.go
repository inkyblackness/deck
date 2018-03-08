package cmd

import (
	"github.com/inkyblackness/hacker/styling"
)

type commandFunction func(Target) string
type commandParser func(string) commandFunction

// Evaluater wraps the Evaluate function to process some input
type Evaluater struct {
	style  styling.Style
	target Target

	commands []commandParser
}

// NewEvaluater returns an evaluater processing input strings.
func NewEvaluater(style styling.Style, target Target) *Evaluater {
	eval := &Evaluater{style: style, commands: []commandParser{}, target: target}

	eval.commands = append(eval.commands, loadCommand, saveCommand, infoCommand, changeDirectoryCommand,
		dumpCommand, diffCommand, putCommand)

	return eval
}

// Evaluate takes the given input, processes it and returns an evaluation result.
func (eval *Evaluater) Evaluate(input string) string {
	var cmd commandFunction
	var result string

	for i := 0; i < len(eval.commands) && cmd == nil; i++ {
		cmd = eval.commands[i](input)
	}
	if cmd != nil {
		result = cmd(eval.target)
	} else {
		result = eval.style.Error()("Unknown command: [", input, "]")
	}

	return result
}
