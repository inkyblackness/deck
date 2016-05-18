package command

import (
	"github.com/inkyblackness/res/geometry"
)

// FaceWriter writes face data into a byte array. It acts as a FaceWalker.
type FaceWriter struct {
	writer       *Writer
	currentColor uint16
	currentShade uint16
}

// NewFaceWriter creates a new instance.
func NewFaceWriter() *FaceWriter {
	return &FaceWriter{
		writer:       NewWriter(),
		currentColor: 0xFFFF,
		currentShade: 0xFFFF}
}

// Bytes returns the created byte buffer.
func (writer *FaceWriter) Bytes() []byte {
	return writer.writer.Bytes()
}

// FlatColored writes a flat colored face
func (writer *FaceWriter) FlatColored(face geometry.FlatColoredFace) {
	color := uint16(face.Color())
	if writer.currentColor != color {
		writer.writer.WriteSetColor(color)
		writer.currentColor = color
	}
	writer.writer.WriteColoredFace(face.Vertices())
}

// ShadeColored writes a shade colored face
func (writer *FaceWriter) ShadeColored(face geometry.ShadeColoredFace) {
	color := uint16(face.Color())
	shade := face.Shade()
	if writer.currentColor != color || writer.currentShade != shade {
		writer.writer.WriteSetColorAndShade(color, shade)
		writer.currentColor = color
		writer.currentShade = shade
	}
	writer.writer.WriteColoredFace(face.Vertices())
}

// TextureMapped writes a texture mapped face
func (writer *FaceWriter) TextureMapped(face geometry.TextureMappedFace) {
	writer.writer.WriteTextureMapping(face.TextureCoordinates())
	writer.writer.WriteTexturedFace(face.Vertices(), face.TextureID())
}
