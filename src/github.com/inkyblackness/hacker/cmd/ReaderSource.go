package cmd

import (
	"bufio"
	"io"
)

type readerSource struct {
	in *bufio.Scanner
}

// NewReaderSource returns a source wrapping an IO reader.
func NewReaderSource(in io.Reader) Source {
	source := &readerSource{in: bufio.NewScanner(in)}

	return source
}

func (source *readerSource) Next() (cmd string, finished bool) {
	finished = !source.in.Scan()
	cmd = source.in.Text()
	return
}
