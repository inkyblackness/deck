package data

// CommonObjectProperties describes properties every game object has.
type CommonObjectProperties struct {
	Unknown0000 [4]byte

	DefaultHitpoints uint16
	Armor            uint8
	RenderType       ObjectRenderType

	Unknown0008 byte
	Unknown0009 byte
	Unused000A  byte
	Unknown000B [3]byte

	Vulnerabilities        byte
	SpecialVulnerabilities byte

	Unused0010  [2]byte
	Unknown0011 byte
	Unknown0012 byte

	Flags      uint16
	ModelIndex uint16

	Unknown0018 byte

	Extra byte

	Unknown001A byte
}
