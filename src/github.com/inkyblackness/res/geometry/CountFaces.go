package geometry

type faceCounter struct {
	sum int
}

// CountFaces returns the number of faces found in the given node.
func CountFaces(node Node) int {
	counter := new(faceCounter)

	node.WalkAnchors(counter)

	return counter.result()
}

func (counter *faceCounter) result() int {
	return counter.sum
}

func (counter *faceCounter) Nodes(anchor NodeAnchor) {
	counter.sum += CountFaces(anchor.Left())
	counter.sum += CountFaces(anchor.Right())
}

func (counter *faceCounter) Faces(anchor FaceAnchor) {
	anchor.WalkFaces(counter)
}

func (counter *faceCounter) FlatColored(face FlatColoredFace) {
	counter.sum++
}

func (counter *faceCounter) ShadeColored(face ShadeColoredFace) {
	counter.sum++
}

func (counter *faceCounter) TextureMapped(face TextureMappedFace) {
	counter.sum++
}
