package geometry

// Node contains a list of anchors. These anchors can be iterated with a walk function.
type Node interface {
	// WalkAnchors iterates over the contained anchors and reports them (specialized) to the given walker.
	WalkAnchors(walker AnchorWalker)
}
