package command

import (
	"bytes"
	"encoding/binary"

	"github.com/inkyblackness/res/geometry"
)

// Writer implements helper functions to write the commands
type Writer struct {
	buf *bytes.Buffer
}

// NewWriter returns a writer instance.
func NewWriter() *Writer {
	return &Writer{buf: bytes.NewBuffer(nil)}
}

func (writer *Writer) write16(value uint16) {
	binary.Write(writer.buf, binary.LittleEndian, &value)
}

func (writer *Writer) write32(value uint32) {
	binary.Write(writer.buf, binary.LittleEndian, &value)
}

func (writer *Writer) writeVector(value geometry.Vector) {
	writer.write32(uint32(ToFixed(value.X())))
	writer.write32(uint32(ToFixed(value.Y())))
	writer.write32(uint32(ToFixed(value.Z())))
}

// Bytes returns the current byte buffer of the writer
func (writer *Writer) Bytes() []byte {
	return writer.buf.Bytes()
}

// WriteBytes adds an array of arbitrary bytes to the buffer.
func (writer *Writer) WriteBytes(bytes []byte) {
	writer.buf.Write(bytes)
}

// WriteHeader writes the command header.
func (writer *Writer) WriteHeader(faceCount int) {
	writer.write16(0x0027)
	writer.write16(0x0008)
	writer.write16(0x0002)
	writer.write16(uint16(faceCount))
}

// WriteDefineVertex writes the command.
func (writer *Writer) WriteDefineVertex(vector geometry.Vector) {
	writer.write16(uint16(CmdDefineVertex))
	writer.write16(uint16(0))
	writer.writeVector(vector)
}

// WriteDefineVertices writes the command.
func (writer *Writer) WriteDefineVertices(vectors []geometry.Vector) {
	writer.write16(uint16(CmdDefineVertices))
	writer.write16(uint16(len(vectors)))
	writer.write16(uint16(0))
	for _, vector := range vectors {
		writer.writeVector(vector)
	}
}

// WriteDefineOneOffsetVertex writes the command.
func (writer *Writer) WriteDefineOneOffsetVertex(cmd ModelCommandID, newIndex int, referenceIndex int, offset float32) {
	writer.write16(uint16(cmd))
	writer.write16(uint16(newIndex))
	writer.write16(uint16(referenceIndex))
	writer.write32(uint32(ToFixed(offset)))
}

// WriteDefineTwoOffsetVertex writes the command.
func (writer *Writer) WriteDefineTwoOffsetVertex(cmd ModelCommandID, newIndex int, referenceIndex int, offset1 float32, offset2 float32) {
	writer.write16(uint16(cmd))
	writer.write16(uint16(newIndex))
	writer.write16(uint16(referenceIndex))
	writer.write32(uint32(ToFixed(offset1)))
	writer.write32(uint32(ToFixed(offset2)))
}

// WriteEndOfNode writes the command.
func (writer *Writer) WriteEndOfNode() {
	writer.write16(uint16(CmdEndOfNode))
}

// WriteNodeAnchor writes the command. The left and right offset values are excluding the size of this command.
func (writer *Writer) WriteNodeAnchor(normal geometry.Vector, reference geometry.Vector, leftOffset int, rightOffset int) {
	writer.write16(uint16(CmdDefineNodeAnchor))
	writer.writeVector(normal)
	writer.writeVector(reference)
	writer.write16(uint16(cmdDefineNodeAnchorSize + leftOffset))
	writer.write16(uint16(cmdDefineNodeAnchorSize + rightOffset))
}

// WriteFaceAnchor writes the command. The length parameter is excluding the size of this command
func (writer *Writer) WriteFaceAnchor(normal geometry.Vector, reference geometry.Vector, size int) {
	writer.write16(uint16(CmdDefineFaceAnchor))
	writer.write16(uint16(cmdDefineFaceAnchorSize + size))
	writer.writeVector(normal)
	writer.writeVector(reference)
}

// WriteSetColor writes the command.
func (writer *Writer) WriteSetColor(color uint16) {
	writer.write16(uint16(CmdSetColor))
	writer.write16(color)
}

// WriteSetColorAndShade writes the command.
func (writer *Writer) WriteSetColorAndShade(color uint16, shade uint16) {
	writer.write16(uint16(CmdSetColorAndShade))
	writer.write16(color)
	writer.write16(shade)
}

// WriteColoredFace writes the command.
func (writer *Writer) WriteColoredFace(vertices []int) {
	writer.write16(uint16(CmdColoredFace))
	writer.write16(uint16(len(vertices)))
	for _, vertex := range vertices {
		writer.write16(uint16(vertex))
	}
}

// WriteColoredFace writes the command.
func (writer *Writer) WriteTextureMapping(textureCoordinates []geometry.TextureCoordinate) {
	writer.write16(uint16(CmdTextureMapping))
	writer.write16(uint16(len(textureCoordinates)))
	for _, coord := range textureCoordinates {
		writer.write16(uint16(coord.Vertex()))
		writer.write32(uint32(ToFixed(coord.U())))
		writer.write32(uint32(ToFixed(coord.V())))
	}
}

// WriteTexturedFace writes the command.
func (writer *Writer) WriteTexturedFace(vertices []int, textureId uint16) {
	writer.write16(uint16(CmdTexturedFace))
	writer.write16(textureId)
	writer.write16(uint16(len(vertices)))
	for _, vertex := range vertices {
		writer.write16(uint16(vertex))
	}
}
