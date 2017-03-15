package display

type simpleArea struct {
	centerX, centerY float32
	width, height    float32
}

// NewSimpleArea returns an immutable area with given properties.
func NewSimpleArea(centerX, centerY float32, width, height float32) Area {
	return &simpleArea{
		centerX: centerX,
		centerY: centerY,
		width:   width,
		height:  height}
}

// Center implements the Area interface.
func (area *simpleArea) Center() (x, y float32) {
	return area.centerX, area.centerY
}

// Size implements the Area interface.
func (area *simpleArea) Size() (width, height float32) {
	return area.width, area.height
}
