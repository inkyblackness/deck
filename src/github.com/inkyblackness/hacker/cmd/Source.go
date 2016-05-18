package cmd

// Source desribes a command source
type Source interface {
	// Next returns the next command in sequence and reports whether more are available.
	// If the second return value is false, the returned string is not usable and the
	// source is exhausted.
	Next() (string, bool)
}
