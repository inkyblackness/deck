package geometry

type simpleFlatColoredFace struct {
	abstractFace

	color ColorIndex
}

// NewSimpleFlatColoredFace returns a new instance of a FlatColoredFace.
func NewSimpleFlatColoredFace(vertices []int, color ColorIndex) FlatColoredFace {
	return &simpleFlatColoredFace{
		abstractFace: abstractFace{vertices: vertices},
		color:        color}
}

func (face *simpleFlatColoredFace) Specialize(walker FaceWalker) {
	walker.FlatColored(face)
}

func (face *simpleFlatColoredFace) Color() ColorIndex {
	return face.color
}
