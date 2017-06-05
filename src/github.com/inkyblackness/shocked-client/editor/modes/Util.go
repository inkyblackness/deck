package modes

func cloneBytes(original []byte) []byte {
	count := len(original)
	clone := make([]byte, count)
	copy(clone, original)
	return clone
}

func intAsPointer(value int) (ptr *int) {
	ptr = new(int)
	*ptr = value
	return
}

func boolAsPointer(value bool) (ptr *bool) {
	ptr = new(bool)
	*ptr = value
	return
}
