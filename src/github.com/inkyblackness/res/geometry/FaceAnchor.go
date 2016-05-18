package geometry

// FaceAnchor is an anchor for rendered faces.
type FaceAnchor interface {
	Anchor

	// WalkFaces iterates over the contained faces and reports them specialized to the walker.
	WalkFaces(walker FaceWalker)
}
