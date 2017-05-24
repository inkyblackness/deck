package modes

func cloneBytes(original []byte) []byte {
	count := len(original)
	clone := make([]byte, count)
	copy(clone, original)
	return clone
}
