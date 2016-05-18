package geometry

type simpleShadeColoredFace struct {
	abstractFace

	color ColorIndex
	shade uint16
}

// NewSimpleShadeColoredFace returns a new instance of a ShadeColoredFace.
func NewSimpleShadeColoredFace(vertices []int, color ColorIndex, shade uint16) ShadeColoredFace {
	return &simpleShadeColoredFace{
		abstractFace: abstractFace{vertices: vertices},
		color:        color,
		shade:        shade}
}

func (face *simpleShadeColoredFace) Specialize(walker FaceWalker) {
	walker.ShadeColored(face)
}

func (face *simpleShadeColoredFace) Color() ColorIndex {
	return face.color
}

func (face *simpleShadeColoredFace) Shade() uint16 {
	return face.shade
}
