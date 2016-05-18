package geometry

// NodeAnchor is an anchor for two further nodes.
type NodeAnchor interface {
	Anchor

	// Left returns the 'left' side of the binary tree.
	Left() Node
	// Right returns the 'right' side of the binary tree.
	Right() Node
}
