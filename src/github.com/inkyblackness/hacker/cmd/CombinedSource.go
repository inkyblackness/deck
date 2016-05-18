package cmd

type combinedSource struct {
	sources []Source
	index   int
}

// NewCombinedSource wraps a list of sources and reads them sequentially.
func NewCombinedSource(sources ...Source) Source {
	combined := &combinedSource{
		sources: sources,
		index:   0}

	return combined
}

func (source *combinedSource) Next() (cmd string, finished bool) {
	finished = source.index >= len(source.sources)
	if !finished {
		cmd, finished = source.sources[source.index].Next()
		if finished {
			source.index++

			cmd, finished = source.Next()
		}
	}

	return
}
