package core

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

func stringAsPointer(value string) (ptr *string) {
	ptr = new(string)
	*ptr = value
	return
}
