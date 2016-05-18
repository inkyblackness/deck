package command

import (
	"fmt"

	"github.com/inkyblackness/res/geometry"
)

// AnchorWriter is a writer for anchors. It acts as an AnchorWalker.
type AnchorWriter struct {
	baseWriter        *Writer
	writer            *Writer
	pendingNodeAnchor geometry.NodeAnchor
}

// NewAnchorWriter returns a new instance.
func NewAnchorWriter() *AnchorWriter {
	baseWriter := NewWriter()
	return &AnchorWriter{baseWriter: baseWriter, writer: baseWriter}
}

// Finish writes any pending data to the internal buffer and returns the final byte array.
func (writer *AnchorWriter) Finish() []byte {
	writer.writer.WriteEndOfNode()
	if writer.writer != writer.baseWriter {
		extraData := writer.writer.Bytes()
		rightData := SaveNode(writer.pendingNodeAnchor.Right())
		leftData := SaveNode(writer.pendingNodeAnchor.Left())

		writer.baseWriter.WriteNodeAnchor(writer.pendingNodeAnchor.Normal(), writer.pendingNodeAnchor.Reference(),
			len(extraData)+len(rightData), len(extraData))
		writer.baseWriter.WriteBytes(extraData)
		writer.baseWriter.WriteBytes(rightData)
		writer.baseWriter.WriteBytes(leftData)
	}

	return writer.baseWriter.Bytes()
}

// Nodes writes a node anchor.
func (writer *AnchorWriter) Nodes(anchor geometry.NodeAnchor) {
	if writer.pendingNodeAnchor != nil {
		panic(fmt.Errorf("Commands don't support more than one NodeAnchor per node"))
	}

	writer.pendingNodeAnchor = anchor
	writer.writer = NewWriter()
}

// Faces writes a face anchor.
func (writer *AnchorWriter) Faces(anchor geometry.FaceAnchor) {
	faceWriter := NewFaceWriter()

	anchor.WalkFaces(faceWriter)
	data := faceWriter.Bytes()
	writer.writer.WriteFaceAnchor(anchor.Normal(), anchor.Reference(), len(data))
	writer.writer.WriteBytes(data)
}
