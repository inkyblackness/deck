package geometry

// FaceWalker receives specialized faces from a FaceAnchor
type FaceWalker interface {
	// FlatColored is called for faces with a simple color.
	FlatColored(face FlatColoredFace)
	// ShadeColored is called for faces with a shaded color.
	ShadeColored(face ShadeColoredFace)
	// TextureMapped is called for faces with textures.
	TextureMapped(face TextureMappedFace)
}
