package gameobj

import (
	"github.com/inkyblackness/res/data/interpreters"
)

var commonProperties = interpreters.New().
	With("Mass", 0x00, 2).As(interpreters.RangedValue(0, 10000)).
	With("DefaultHitpoints", 0x04, 2).As(interpreters.RangedValue(0, 10000)).
	With("Armor", 0x06, 1).
	With("RenderType", 0x07, 1).As(interpreters.EnumValue(map[uint32]string{
	0x01: "3D Object",
	0x02: "Sprite",
	0x03: "Screen",
	0x04: "Critter",
	0x06: "Fragments",
	0x07: "Invisible",
	0x08: "Oriented surface",
	0x0B: "Special",
	0x0C: "Force door"})).
	With("PhysicsType", 0x08, 1).As(interpreters.EnumValue(map[uint32]string{0x00: "Insubstantial", 0x01: "Regular", 0x02: "Special"})).
	With("Bounciness", 0x09, 1).As(interpreters.RangedValue(-128, 127)).
	With("VerticalFrameOffset", 0x0B, 1).
	With("Unknown000C", 0x0C, 1).As(interpreters.SpecialValue("Ignored")).
	With("Unknown000D", 0x0D, 1).As(interpreters.SpecialValue("Ignored")).
	With("Vulnerabilities", 0x0E, 1).As(damageType).
	With("SpecialVulnerabilities", 0x0F, 1).
	With("Defence", 0x12, 1).
	With("ReceiveDamageFlag", 0x13, 1).As(interpreters.EnumValue(map[uint32]string{0x00: "Yes", 0x03: "No", 0x04: "Unknown 0x04"})).
	With("Flags", 0x14, 2).As(interpreters.Bitfield(map[uint32]string{
	0x0001: "Useful",
	0x0002: "Solid",
	0x0004: "Triggerable",
	0x0008: "Unusable",
	0x0010: "Usable",
	0x0020: "Blocks3D",
	0x0040: "Unknown0040",
	0x0080: "IgnoreDarkness",
	0x0100: "SolidIfClosed",
	0x0200: "FlatSolid",
	0x0400: "LargeExplosion",
	0x0800: "DestroyOnContact",
	0x1000: "Unknown1000",
	0x2000: "Unknown2000",
	0x4000: "Unknown4000",
	0x8000: "Unknown8000"})).
	With("3DModelIndex", 0x16, 2).As(interpreters.RangedValue(0, 500)).
	With("Unknown0018", 0x18, 1).As(interpreters.EnumValue(map[uint32]string{0x00: "Unknown 0x00", 0x80: "Unknown 0x80"})).
	With("Extra", 0x19, 1).As(interpreters.Bitfield(map[uint32]string{
	0x0F: "Unknown",
	0xF0: "FrameCount"})).
	With("DestructionEffect", 0x1A, 1).As(interpreters.Bitfield(map[uint32]string{
	0x1F: "AnimationFrameIndex",
	0x20: "PlaySound",
	0x40: "PlaySound2",
	0x80: "Explosion"}))

// CommonProperties returns an interpreter about common object properties.
func CommonProperties(data []byte) *interpreters.Instance {
	return commonProperties.For(data)
}
