package geometry

// DynamicNode is a Node implementation that allows addition of content
type DynamicNode struct {
	anchors []Anchor
}

// NewDynamicNode returns a new DynamicNode instance.
func NewDynamicNode() *DynamicNode {
	return new(DynamicNode)
}

// WalkAnchors is the Node implementation
func (node *DynamicNode) WalkAnchors(walker AnchorWalker) {
	for _, anchor := range node.anchors {
		anchor.Specialize(walker)
	}
}

// AddAnchor adds the given anchor to the node.
func (node *DynamicNode) AddAnchor(anchor Anchor) {
	node.anchors = append(node.anchors, anchor)
}
