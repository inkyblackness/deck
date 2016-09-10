package data

// ObjectRenderType describes how an object should be visually represented.
type ObjectRenderType byte

const (
	// RenderType3D is for all three-dimensional objects.
	RenderType3D = ObjectRenderType(0x01)
	// RenderTypeSprite is for two-dimensional sprite-based objects.
	RenderTypeSprite = ObjectRenderType(0x02)
	// RenderTypeScreen is for animated content
	RenderTypeScreen = ObjectRenderType(0x03)
	// RenderTypeCritter is for enemies.
	RenderTypeCritter = ObjectRenderType(0x04)
	// RenderTypeFragment is for programs in cyberspace.
	RenderTypeFragment = ObjectRenderType(0x06)
	// RenderTypeInvisible is for infrastructure items (triggers, ...)
	RenderTypeInvisible = ObjectRenderType(0x07)
	// RenderTypeOrientedSurface is for items with a fixed orientation.
	RenderTypeOrientedSurface = ObjectRenderType(0x08)
	// RenderTypeSpecial require special handling for display.
	RenderTypeSpecial = ObjectRenderType(0x0B)
	// RenderTypeForceDoor is for semi-transparent barriers.
	RenderTypeForceDoor = ObjectRenderType(0x0C)
)
