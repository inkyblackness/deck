package command

import (
	"fmt"
	"io"

	"github.com/inkyblackness/res/geometry"
	"github.com/inkyblackness/res/serial"
)

// LoadModel tries to decode a geometry model from a serialized list of model
// commands.
func LoadModel(source io.ReadSeeker) (model geometry.Model, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	if source == nil {
		panic(fmt.Errorf("source is nil"))
	}

	coder := serial.NewPositioningDecoder(source)

	unknownHeader := make([]byte, 6)
	expectedFaces := uint16(0)
	coder.CodeBytes(unknownHeader)
	coder.CodeUint16(&expectedFaces)

	dynamicModel := geometry.NewDynamicModel()

	loadNodeData(coder, dynamicModel, dynamicModel)

	model = dynamicModel

	return
}

func loadNodeData(coder serial.PositioningCoder, model *geometry.DynamicModel, node geometry.ExtensibleNode) {
	pendingNodes := make(map[uint32]geometry.ExtensibleNode)
	done := false

	for !done {
		startPos := coder.CurPos()
		rawCommand := uint16(0)

		coder.CodeUint16(&rawCommand)
		switch ModelCommandID(rawCommand) {
		case CmdEndOfNode:
			{
				done = true
			}
		case CmdDefineVertex:
			{
				unknown := uint16(0)
				vector := new(Vector)

				coder.CodeUint16(&unknown)
				vector.Code(coder)
				model.AddVertex(geometry.NewSimpleVertex(NewFixedVector(*vector)))
			}
		case CmdDefineVertices:
			{
				unknown := uint16(0)
				vertexCount := uint16(0)

				coder.CodeUint16(&vertexCount)
				coder.CodeUint16(&unknown)
				for i := uint16(0); i < vertexCount; i++ {
					vector := new(Vector)

					vector.Code(coder)
					model.AddVertex(geometry.NewSimpleVertex(NewFixedVector(*vector)))
				}
			}
		case CmdDefineOffsetVertexX:
			{
				loadDefineOffsetVertexOne(coder, model, AddingModifier, singleIdentityModifier, singleIdentityModifier)
			}
		case CmdDefineOffsetVertexY:
			{
				loadDefineOffsetVertexOne(coder, model, singleIdentityModifier, AddingModifier, singleIdentityModifier)
			}
		case CmdDefineOffsetVertexZ:
			{
				loadDefineOffsetVertexOne(coder, model, singleIdentityModifier, singleIdentityModifier, AddingModifier)
			}
		case CmdDefineOffsetVertexXY:
			{
				loadDefineOffsetVertexTwo(coder, model,
					func(offset1, offset2 float32) Modifier { return AddingModifier(offset1) },
					func(offset1, offset2 float32) Modifier { return AddingModifier(offset2) }, doubleIdentityModifier)
			}
		case CmdDefineOffsetVertexXZ:
			{
				loadDefineOffsetVertexTwo(coder, model,
					func(offset1, offset2 float32) Modifier { return AddingModifier(offset1) }, doubleIdentityModifier,
					func(offset1, offset2 float32) Modifier { return AddingModifier(offset2) })
			}
		case CmdDefineOffsetVertexYZ:
			{
				loadDefineOffsetVertexTwo(coder, model, doubleIdentityModifier,
					func(offset1, offset2 float32) Modifier { return AddingModifier(offset1) },
					func(offset1, offset2 float32) Modifier { return AddingModifier(offset2) })
			}
		case CmdDefineNodeAnchor:
			{
				normal := new(Vector)
				reference := new(Vector)
				leftOffset := uint16(0)
				rightOffset := uint16(0)

				normal.Code(coder)
				reference.Code(coder)
				coder.CodeUint16(&leftOffset)
				coder.CodeUint16(&rightOffset)

				left := geometry.NewDynamicNode()
				right := geometry.NewDynamicNode()
				anchor := geometry.NewSimpleNodeAnchor(NewFixedVector(*normal), NewFixedVector(*reference), left, right)
				node.AddAnchor(anchor)
				pendingNodes[startPos+uint32(leftOffset)] = left
				pendingNodes[startPos+uint32(rightOffset)] = right
			}
		case CmdDefineFaceAnchor:
			{
				normal := new(Vector)
				reference := new(Vector)
				size := uint16(0)

				coder.CodeUint16(&size)
				normal.Code(coder)
				reference.Code(coder)

				endPos := startPos + uint32(size)
				anchor := geometry.NewDynamicFaceAnchor(NewFixedVector(*normal), NewFixedVector(*reference))
				node.AddAnchor(anchor)
				loadFaces(coder, anchor, endPos)
			}
		default:
			{
				panic(fmt.Errorf("Unknown model command 0x%04X at offset 0x%X", rawCommand, startPos))
			}
		}
	}

	for curPos := coder.CurPos(); pendingNodes[curPos] != nil; curPos = coder.CurPos() {
		pendingNode := pendingNodes[curPos]
		delete(pendingNodes, curPos)
		loadNodeData(coder, model, pendingNode)
	}

	if len(pendingNodes) != 0 {
		panic(fmt.Errorf("Wrong offset values for node anchor"))
	}
}

type singleModifierFactory func(float32) Modifier
type doubleModifierFactory func(float32, float32) Modifier

func singleIdentityModifier(offset float32) Modifier {
	return IdentityModifier
}

func doubleIdentityModifier(offset1 float32, offset2 float32) Modifier {
	return IdentityModifier
}

func loadDefineOffsetVertexOne(coder serial.PositioningCoder, model *geometry.DynamicModel,
	xModFactory singleModifierFactory, yModFactory singleModifierFactory, zModFactory singleModifierFactory) {
	newIndex := uint16(0)
	referenceIndex := uint16(0)
	fixedOffset := Fixed(0)

	coder.CodeUint16(&newIndex)
	coder.CodeUint16(&referenceIndex)
	CodeFixed(coder, &fixedOffset)

	if int(newIndex) != model.VertexCount() {
		panic(fmt.Errorf("Offset vertex uses invalid newIndex (%d)", newIndex))
	}

	reference := model.Vertex(int(referenceIndex))
	offset := fixedOffset.Float()
	model.AddVertex(geometry.NewSimpleVertex(NewModifiedVector(reference.Position(),
		xModFactory(offset), yModFactory(offset), zModFactory(offset))))
}

func loadDefineOffsetVertexTwo(coder serial.PositioningCoder, model *geometry.DynamicModel,
	xModFactory doubleModifierFactory, yModFactory doubleModifierFactory, zModFactory doubleModifierFactory) {
	newIndex := uint16(0)
	referenceIndex := uint16(0)
	fixedOffset1 := Fixed(0)
	fixedOffset2 := Fixed(0)

	coder.CodeUint16(&newIndex)
	coder.CodeUint16(&referenceIndex)
	CodeFixed(coder, &fixedOffset1)
	CodeFixed(coder, &fixedOffset2)

	if int(newIndex) != model.VertexCount() {
		panic(fmt.Errorf("Offset vertex uses invalid newIndex (%d)", newIndex))
	}

	reference := model.Vertex(int(referenceIndex))
	offset1 := fixedOffset1.Float()
	offset2 := fixedOffset2.Float()
	model.AddVertex(geometry.NewSimpleVertex(NewModifiedVector(reference.Position(),
		xModFactory(offset1, offset2), yModFactory(offset1, offset2), zModFactory(offset1, offset2))))
}

func loadFaces(coder serial.PositioningCoder, anchor *geometry.DynamicFaceAnchor, endPos uint32) {
	currentColor := uint16(0)
	currentShade := uint16(0xFFFF)
	var textureCoordinates []geometry.TextureCoordinate

	for startPos := coder.CurPos(); startPos < endPos; startPos = coder.CurPos() {
		rawCommand := uint16(0)

		coder.CodeUint16(&rawCommand)
		switch ModelCommandID(rawCommand) {
		case CmdSetColor:
			{
				coder.CodeUint16(&currentColor)
			}
		case CmdSetColorAndShade:
			{
				coder.CodeUint16(&currentColor)
				coder.CodeUint16(&currentShade)
			}
		case CmdColoredFace:
			{
				vertexCount := uint16(0)
				coder.CodeUint16(&vertexCount)
				vertices := make([]int, int(vertexCount))
				for i := uint16(0); i < vertexCount; i++ {
					index := uint16(0)
					coder.CodeUint16(&index)
					vertices[i] = int(index)
				}
				if currentShade == 0xFFFF {
					anchor.AddFace(geometry.NewSimpleFlatColoredFace(vertices, geometry.ColorIndex(currentColor)))
				} else {
					anchor.AddFace(geometry.NewSimpleShadeColoredFace(vertices, geometry.ColorIndex(currentColor), currentShade))
				}
			}
		case CmdTextureMapping:
			{
				entryCount := uint16(0)
				coder.CodeUint16(&entryCount)
				textureCoordinates = make([]geometry.TextureCoordinate, int(entryCount))
				for i := uint16(0); i < entryCount; i++ {
					vertex := uint16(0)
					u := Fixed(0)
					v := Fixed(0)

					coder.CodeUint16(&vertex)
					CodeFixed(coder, &u)
					CodeFixed(coder, &v)
					textureCoordinates[i] = geometry.NewSimpleTextureCoordinate(int(vertex), u.Float(), v.Float())
				}
			}
		case CmdTexturedFace:
			{
				textureId := uint16(0)
				vertexCount := uint16(0)
				coder.CodeUint16(&textureId)
				coder.CodeUint16(&vertexCount)
				vertices := make([]int, int(vertexCount))
				for i := uint16(0); i < vertexCount; i++ {
					index := uint16(0)
					coder.CodeUint16(&index)
					vertices[i] = int(index)
				}
				anchor.AddFace(geometry.NewSimpleTextureMappedFace(vertices, textureId, textureCoordinates))
			}
		default:
			{
				panic(fmt.Errorf("Unknown face command 0x%04X at offset 0x%X", rawCommand, startPos))
			}
		}
	}
}
