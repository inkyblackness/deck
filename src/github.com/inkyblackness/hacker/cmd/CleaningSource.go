package cmd

import (
	"strings"
)

type cleaningSource struct {
	source Source
}

// NewCleaningSource returns a source which cleans input by trimming whitespace
// and dropping comment lines
func NewCleaningSource(source Source) Source {
	return &cleaningSource{source: source}
}

func (source *cleaningSource) Next() (cmd string, finished bool) {
	cmd, finished = source.source.Next()
	cmd = strings.Trim(cmd, " ")
	if strings.IndexRune(cmd, '#') == 0 {
		cmd = ""
	}
	return
}
