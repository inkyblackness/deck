package cmd

type staticSource struct {
	commands []string
	index    int
}

// NewStaticSource returns a source that returns the provided commands in sequence.
func NewStaticSource(commands ...string) Source {
	source := &staticSource{
		commands: commands,
		index:    0}

	return source
}

func (source *staticSource) Next() (cmd string, finished bool) {
	finished = source.index >= len(source.commands)
	if !finished {
		cmd = source.commands[source.index]
		source.index++
	}
	return
}
