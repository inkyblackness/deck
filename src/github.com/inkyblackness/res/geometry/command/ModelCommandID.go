package command

// ModelCommandID identifies a definition command
type ModelCommandID uint16

const (
	CmdEndOfNode        = ModelCommandID(0x0000)
	CmdDefineNodeAnchor = ModelCommandID(0x0006)

	CmdDefineVertices = ModelCommandID(0x0003)
	CmdDefineVertex   = ModelCommandID(0x0015)

	CmdDefineOffsetVertexX  = ModelCommandID(0x000A)
	CmdDefineOffsetVertexY  = ModelCommandID(0x000B)
	CmdDefineOffsetVertexZ  = ModelCommandID(0x000C)
	CmdDefineOffsetVertexXY = ModelCommandID(0x000D)
	CmdDefineOffsetVertexXZ = ModelCommandID(0x000E)
	CmdDefineOffsetVertexYZ = ModelCommandID(0x000F)

	CmdDefineFaceAnchor = ModelCommandID(0x0001)

	CmdColoredFace      = ModelCommandID(0x0004)
	CmdSetColor         = ModelCommandID(0x0005)
	CmdSetColorAndShade = ModelCommandID(0x001C)

	CmdTextureMapping = ModelCommandID(0x0025)
	CmdTexturedFace   = ModelCommandID(0x0026)
)

const (
	cmdDefineNodeAnchorSize = 30
	cmdDefineFaceAnchorSize = 28
)
