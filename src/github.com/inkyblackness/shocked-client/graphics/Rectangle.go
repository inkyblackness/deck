package graphics

// Rectangle describes a rectangular area with four points.
type Rectangle interface {
	// Left returns the left edge of the rectangle.
	Left() float32
	// Top returns the top edge of the rectangle.
	Top() float32
	// Right returns the coordinate beyond the right edge of the rectangle.
	Right() float32
	// Bottom returns the coordinate beyond the bottom edge of the rectangle.
	Bottom() float32
}

type simpleRectangle [4]float32

// RectByCoord returns a rectangle instance by coordinates.
func RectByCoord(left, top, right, bottom float32) Rectangle {
	rect := &simpleRectangle{left, top, right, bottom}
	return rect
}

// Left impelements the Rectangle interface
func (rect *simpleRectangle) Left() float32 {
	return rect[0]
}

// Top impelements the Rectangle interface
func (rect *simpleRectangle) Top() float32 {
	return rect[1]
}

// Right impelements the Rectangle interface
func (rect *simpleRectangle) Right() float32 {
	return rect[2]
}

// Bottom impelements the Rectangle interface
func (rect *simpleRectangle) Bottom() float32 {
	return rect[3]
}
