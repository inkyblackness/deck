package command

import (
	"github.com/inkyblackness/res/geometry"
)

// SaveNode uses an AnchorWriter to serialize a single node.
func SaveNode(node geometry.Node) []byte {
	writer := NewAnchorWriter()

	node.WalkAnchors(writer)

	return writer.Finish()
}
