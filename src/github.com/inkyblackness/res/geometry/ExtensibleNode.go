package geometry

// ExtensibleNode is a node which allows addition of extra anchors
type ExtensibleNode interface {
	Node

	// AddAnchor adds the given anchor to the node.
	AddAnchor(anchor Anchor)
}
