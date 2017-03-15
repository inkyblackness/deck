package controls

// Aligner is a function to calculate the left position of an element within
// a container. The returned value is the relative offset from left.
type Aligner func(containerSize float32, elementSize float32) float32

// LeftAligner always positions the element at zero.
func LeftAligner(containerSize float32, elementSize float32) float32 {
	return 0
}

// CenterAligner positions the element centered within the container.
func CenterAligner(containerSize float32, elementSize float32) float32 {
	return (containerSize / 2.0) - (elementSize / 2.0)
}

// RightAligner positions the element right bound within the container.
func RightAligner(containerSize float32, elementSize float32) float32 {
	return containerSize - elementSize
}
