package geometry

type simpleNodeAnchor struct {
	abstractAnchor

	left  Node
	right Node
}

// NewSimpleNodeAnchor returns a NodeAnchor instance based on provided parameters.
func NewSimpleNodeAnchor(normal, reference Vector, left, right Node) NodeAnchor {
	return &simpleNodeAnchor{
		abstractAnchor: abstractAnchor{normal, reference},
		left:           left,
		right:          right}
}

func (anchor *simpleNodeAnchor) Specialize(walker AnchorWalker) {
	walker.Nodes(anchor)
}

func (anchor *simpleNodeAnchor) Left() Node {
	return anchor.left
}

func (anchor *simpleNodeAnchor) Right() Node {
	return anchor.right
}
