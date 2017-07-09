package modes

import (
	"fmt"
)

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

func stringAsPointer(value string) (ptr *string) {
	ptr = new(string)
	*ptr = value
	return
}

func heightToString(heightShift int, value int64, scale float64) (result string) {
	tileHeights := []float64{32.0, 16.0, 8.0, 4.0, 2.0, 1.0, 0.5, 0.25}
	if (heightShift >= 0) && (heightShift < len(tileHeights)) {
		result = fmt.Sprintf("%.3f tile(s)  - raw: %v", (float64(value)*tileHeights[heightShift])/scale, value)
	} else {
		result = fmt.Sprintf("??? - raw: %v", value)
	}
	return
}
