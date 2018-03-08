package main

import (
	"io"

	"github.com/chzyer/readline"
)

type prompterFunc func() string

// ReadLineSource is a command source based on a readline-like interface.
type ReadLineSource struct {
	rl       *readline.Instance
	prompter prompterFunc
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

// NewReadLineSource returns a new instance of a ReadLineSource.
func NewReadLineSource(prompter prompterFunc) *ReadLineSource {
	var completer = readline.NewPrefixCompleter(
		readline.PcItem("cd"),
		readline.PcItem("diff"),
		readline.PcItem("dump"),
		readline.PcItem("info"),
		readline.PcItem("load"),
		readline.PcItem("put"),
		readline.PcItem("save"),
		readline.PcItem("quit"),
	)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 prompter(),
		HistoryFile:            ".hacker-history",
		DisableAutoSaveHistory: false,
		HistoryLimit:           100,
		AutoComplete:           completer,
		InterruptPrompt:        "^C",
		EOFPrompt:              "quit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}

	return &ReadLineSource{rl: rl, prompter: prompter}
}

// Close releases the resources of this source.
func (source *ReadLineSource) Close() error {
	return source.rl.Close()
}

// Next queries the next command from the user.
func (source *ReadLineSource) Next() (cmd string, finished bool) {
	source.rl.SetPrompt(source.prompter())
	line, err := source.rl.Readline()

	if err == readline.ErrInterrupt || err == io.EOF {
		return "", true
	}
	return line, false
}
