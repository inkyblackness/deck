package geometry

type abstractFace struct {
	vertices []int
}

func (face *abstractFace) Vertices() []int {
	result := make([]int, len(face.vertices))
	copy(result, face.vertices)
	return result
}
