package data

// CommonObjectProperties describes properties every game object has.
type CommonObjectProperties struct {
	Mass             uint16
	Unused0002       [2]byte
	DefaultHitpoints uint16
	Armor            uint8

	RenderType  ObjectRenderType
	PhysicsType uint8
	Bounciness  int8
	Unused000A  byte

	VerticalFrameOffset byte
	Unknown000C         [2]byte

	Vulnerabilities        byte
	SpecialVulnerabilities byte

	Unused0010        [2]byte
	Defence           byte
	ReceiveDamageFlag byte

	Flags      uint16
	ModelIndex uint16

	Unknown0018 byte

	Extra byte

	DestructionEffect byte
}
