package model

func safeString(value *string) (result string) {
	if value != nil {
		result = *value
	}
	return
}

func safeInt(value *int, defaultValue int) (result int) {
	result = defaultValue
	if value != nil {
		result = *value
	}
	return
}
