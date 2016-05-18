package geometry

type DynamicFaceAnchor struct {
	abstractAnchor

	faces []Face
}

// NewSimpleNodeAnchor returns a NodeAnchor instance based on provided parameters.
func NewDynamicFaceAnchor(normal, reference Vector) *DynamicFaceAnchor {
	return &DynamicFaceAnchor{
		abstractAnchor: abstractAnchor{normal, reference},
		faces:          nil}
}

// Specialize is the Anchor implementation.
func (anchor *DynamicFaceAnchor) Specialize(walker AnchorWalker) {
	walker.Faces(anchor)
}

// WalkFaces is the Anchor implementation.
func (anchor *DynamicFaceAnchor) WalkFaces(walker FaceWalker) {
	for _, face := range anchor.faces {
		face.Specialize(walker)
	}
}

// AddFace adds the given face to the anchor
func (anchor *DynamicFaceAnchor) AddFace(face Face) {
	anchor.faces = append(anchor.faces, face)
}
