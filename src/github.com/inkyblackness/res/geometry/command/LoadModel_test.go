package command

import (
	"bytes"

	"github.com/inkyblackness/res/geometry"

	check "gopkg.in/check.v1"
)

type LoadModelSuite struct {
	nodeAnchors []geometry.NodeAnchor
	nodes       []geometry.Node

	faceAnchors []geometry.FaceAnchor
	faces       []geometry.Face
}

var _ = check.Suite(&LoadModelSuite{})

func (suite *LoadModelSuite) SetUpTest(c *check.C) {
	suite.nodeAnchors = nil
	suite.nodes = nil

	suite.faceAnchors = nil
	suite.faces = nil
}

func (suite *LoadModelSuite) TestLoadModelReturnsErrorOnNil(c *check.C) {
	_, err := LoadModel(nil)

	c.Check(err, check.ErrorMatches, "source is nil")
}

func (suite *LoadModelSuite) TestLoadModelReturnsModelInstanceOnValidData(c *check.C) {
	source := bytes.NewReader(suite.aSimpleList())
	model, err := LoadModel(source)

	c.Assert(err, check.IsNil)
	c.Check(model, check.NotNil)
}

func (suite *LoadModelSuite) TestModelContainsSingleVertex(c *check.C) {
	source := bytes.NewReader(suite.anEmptyModelWith(func(writer *Writer) {
		writer.WriteDefineVertex(NewFixedVector(Vector{0, 0, 0}))
	}))
	model, err := LoadModel(source)

	c.Assert(err, check.IsNil)
	c.Check(model.VertexCount(), check.Equals, 1)
}

func (suite *LoadModelSuite) TestModelContainsMultipleVertices(c *check.C) {
	source := bytes.NewReader(suite.anEmptyModelWith(func(writer *Writer) {
		writer.WriteDefineVertices([]geometry.Vector{
			NewFixedVector(Vector{1, 0, 0}), NewFixedVector(Vector{0, 2, 0}), NewFixedVector(Vector{0, 0, 3})})
	}))
	model, err := LoadModel(source)

	c.Assert(err, check.IsNil)
	c.Check(model.VertexCount(), check.Equals, 3)
}

func (suite *LoadModelSuite) TestModelContainsSingleOffsetVertices(c *check.C) {
	source := bytes.NewReader(suite.anEmptyModelWith(func(writer *Writer) {
		writer.WriteDefineVertex(NewFixedVector(Vector{0, 0, 0}))
		writer.WriteDefineOneOffsetVertex(CmdDefineOffsetVertexX, 1, 0, 1.0)
		writer.WriteDefineOneOffsetVertex(CmdDefineOffsetVertexY, 2, 0, 1.0)
		writer.WriteDefineOneOffsetVertex(CmdDefineOffsetVertexZ, 3, 0, 1.0)
	}))
	model, err := LoadModel(source)

	c.Assert(err, check.IsNil)
	c.Check(model.VertexCount(), check.Equals, 4)
}

func (suite *LoadModelSuite) TestLoadModelReturnsErrorForSingleOffsetVertexWhenNewIndexIsNotEqualCurrentCount(c *check.C) {
	source := bytes.NewReader(suite.anEmptyModelWith(func(writer *Writer) {
		writer.WriteDefineVertex(NewFixedVector(Vector{0, 0, 0}))
		writer.WriteDefineOneOffsetVertex(CmdDefineOffsetVertexX, 2, 0, 1.0)
	}))
	_, err := LoadModel(source)

	c.Check(err, check.ErrorMatches, "Offset vertex uses invalid newIndex.*")
}

func (suite *LoadModelSuite) TestModelContainsDoubleOffsetVertices(c *check.C) {
	source := bytes.NewReader(suite.anEmptyModelWith(func(writer *Writer) {
		writer.WriteDefineVertex(NewFixedVector(Vector{0, 0, 0}))
		writer.WriteDefineTwoOffsetVertex(CmdDefineOffsetVertexXY, 1, 0, 1.0, 2.0)
		writer.WriteDefineTwoOffsetVertex(CmdDefineOffsetVertexXZ, 2, 0, 1.0, 2.0)
		writer.WriteDefineTwoOffsetVertex(CmdDefineOffsetVertexYZ, 3, 0, 1.0, 2.0)
	}))
	model, err := LoadModel(source)

	c.Assert(err, check.IsNil)
	c.Check(model.VertexCount(), check.Equals, 4)
}

func (suite *LoadModelSuite) TestLoadModelReturnsErrorForDoubleOffsetVertexWhenNewIndexIsNotEqualCurrentCount(c *check.C) {
	source := bytes.NewReader(suite.anEmptyModelWith(func(writer *Writer) {
		writer.WriteDefineVertex(NewFixedVector(Vector{0, 0, 0}))
		writer.WriteDefineTwoOffsetVertex(CmdDefineOffsetVertexXY, 2, 0, 1.0, 2.0)
	}))
	_, err := LoadModel(source)

	c.Check(err, check.ErrorMatches, "Offset vertex uses invalid newIndex.*")
}

func (suite *LoadModelSuite) TestModelCanHaveANodeAnchor(c *check.C) {
	source := bytes.NewReader(suite.anEmptyModelWith(func(writer *Writer) {
		writer.WriteDefineVertex(NewFixedVector(Vector{0, 0, 0}))
		writer.WriteNodeAnchor(NewFixedVector(Vector{1, 0, 0}), NewFixedVector(Vector{0, 0, 0}), 2, 4)
		writer.WriteEndOfNode() // end of root node
		writer.WriteEndOfNode() // end of left node
	}))
	model, err := LoadModel(source)

	c.Assert(err, check.IsNil)

	model.WalkAnchors(suite)
	c.Check(len(suite.nodeAnchors), check.Equals, 1)
}

func (suite *LoadModelSuite) TestLoadModelReturnsErrorForInvalidNodeAnchorOffsets(c *check.C) {
	source := bytes.NewReader(suite.anEmptyModelWith(func(writer *Writer) {
		writer.WriteDefineVertex(NewFixedVector(Vector{0, 0, 0}))
		writer.WriteNodeAnchor(NewFixedVector(Vector{1, 0, 0}), NewFixedVector(Vector{0, 0, 0}), 2, 5)
		writer.WriteEndOfNode() // end of root node
		writer.WriteEndOfNode() // end of left node
	}))
	_, err := LoadModel(source)

	c.Check(err, check.ErrorMatches, "Wrong offset values for node anchor")
}

func (suite *LoadModelSuite) TestModelCanHaveNestedNodeAnchors(c *check.C) {
	source := bytes.NewReader(suite.anEmptyModelWith(func(writer *Writer) {
		writer.WriteDefineVertex(NewFixedVector(Vector{0, 0, 0}))
		// root node
		writer.WriteNodeAnchor(NewFixedVector(Vector{1, 0, 0}), NewFixedVector(Vector{0, 0, 0}), 2, 2+cmdDefineNodeAnchorSize+2+2+2)
		writer.WriteEndOfNode() // end of root node
		// root.left node
		writer.WriteNodeAnchor(NewFixedVector(Vector{1, 0, 0}), NewFixedVector(Vector{0, 0, 0}), 4, 2)
		writer.WriteEndOfNode() // end of root.left node
		// root.left.right node
		writer.WriteEndOfNode() // end of root.left.right node
		// root.left.left node
		writer.WriteEndOfNode() // end of root.left.left node
		// root.right node (end provided by helper function)
	}))
	model, err := LoadModel(source)

	c.Assert(err, check.IsNil)

	model.WalkAnchors(suite)
	c.Check(len(suite.nodeAnchors), check.Equals, 2)
	c.Check(len(suite.nodes), check.Equals, 4)
}

func (suite *LoadModelSuite) TestModelCanHaveAFaceAnchor(c *check.C) {
	source := bytes.NewReader(suite.anEmptyModelWith(func(writer *Writer) {
		writer.WriteDefineVertex(NewFixedVector(Vector{0, 0, 0}))
		writer.WriteFaceAnchor(NewFixedVector(Vector{1, 0, 0}), NewFixedVector(Vector{0, 0, 0}), 0)
	}))
	model, err := LoadModel(source)

	c.Assert(err, check.IsNil)

	model.WalkAnchors(suite)
	c.Check(len(suite.faceAnchors), check.Equals, 1)
}

func (suite *LoadModelSuite) TestModelCanHaveAFlatColoredFace(c *check.C) {
	source := bytes.NewReader(suite.anEmptyModelWith(func(writer *Writer) {
		writer.WriteDefineVertices([]geometry.Vector{
			NewFixedVector(Vector{0, 0, 0}), NewFixedVector(Vector{1, 0, 0}), NewFixedVector(Vector{1, 1, 0})})
		writer.WriteFaceAnchor(NewFixedVector(Vector{1, 0, 0}), NewFixedVector(Vector{0, 0, 0}), 4+8)
		writer.WriteSetColor(0x00AB)
		writer.WriteColoredFace([]int{0, 1, 2})
	}))
	model, err := LoadModel(source)

	c.Assert(err, check.IsNil)

	model.WalkAnchors(suite)
	face := suite.faces[0].(geometry.FlatColoredFace)
	c.Check(face.Color(), check.Equals, geometry.ColorIndex(0x00AB))
	c.Check(face.Vertices(), check.DeepEquals, []int{0, 1, 2})
}

func (suite *LoadModelSuite) TestModelCanHaveAShadeColoredFace(c *check.C) {
	source := bytes.NewReader(suite.anEmptyModelWith(func(writer *Writer) {
		writer.WriteDefineVertices([]geometry.Vector{
			NewFixedVector(Vector{0, 0, 0}), NewFixedVector(Vector{1, 0, 0}), NewFixedVector(Vector{1, 1, 0})})
		writer.WriteFaceAnchor(NewFixedVector(Vector{1, 0, 0}), NewFixedVector(Vector{0, 0, 0}), 6+8)
		writer.WriteSetColorAndShade(0x0012, 0x0002)
		writer.WriteColoredFace([]int{0, 1, 2})
	}))
	model, err := LoadModel(source)

	c.Assert(err, check.IsNil)

	model.WalkAnchors(suite)
	face := suite.faces[0].(geometry.ShadeColoredFace)
	c.Check(face.Color(), check.Equals, geometry.ColorIndex(0x0012))
	c.Check(face.Shade(), check.Equals, uint16(2))
	c.Check(face.Vertices(), check.DeepEquals, []int{0, 1, 2})
}

func (suite *LoadModelSuite) TestModelCanHaveATextureMappedFace(c *check.C) {
	textureCoordinates := []geometry.TextureCoordinate{geometry.NewSimpleTextureCoordinate(0, 0.0, 0.0),
		geometry.NewSimpleTextureCoordinate(1, 1.0, 0.0), geometry.NewSimpleTextureCoordinate(2, 1.0, 1.0)}
	source := bytes.NewReader(suite.anEmptyModelWith(func(writer *Writer) {
		writer.WriteDefineVertices([]geometry.Vector{
			NewFixedVector(Vector{0, 0, 0}), NewFixedVector(Vector{1, 0, 0}), NewFixedVector(Vector{1, 1, 0})})
		writer.WriteFaceAnchor(NewFixedVector(Vector{1, 0, 0}), NewFixedVector(Vector{0, 0, 0}), 34+10)
		writer.WriteTextureMapping(textureCoordinates)
		writer.WriteTexturedFace([]int{0, 1, 2}, 0x1234)
	}))
	model, err := LoadModel(source)

	c.Assert(err, check.IsNil)

	model.WalkAnchors(suite)
	face := suite.faces[0].(geometry.TextureMappedFace)
	c.Check(face.Vertices(), check.DeepEquals, []int{0, 1, 2})
	c.Check(face.TextureID(), check.Equals, uint16(0x1234))
	c.Check(face.TextureCoordinates(), check.DeepEquals, textureCoordinates)
}

func (suite *LoadModelSuite) Nodes(anchor geometry.NodeAnchor) {
	suite.nodeAnchors = append(suite.nodeAnchors, anchor)
	suite.nodes = append(suite.nodes, anchor.Left(), anchor.Right())
	anchor.Left().WalkAnchors(suite)
	anchor.Right().WalkAnchors(suite)
}

func (suite *LoadModelSuite) Faces(anchor geometry.FaceAnchor) {
	suite.faceAnchors = append(suite.faceAnchors, anchor)
	anchor.WalkFaces(suite)
}

func (suite *LoadModelSuite) FlatColored(face geometry.FlatColoredFace) {
	suite.faces = append(suite.faces, face)
}

func (suite *LoadModelSuite) ShadeColored(face geometry.ShadeColoredFace) {
	suite.faces = append(suite.faces, face)
}

func (suite *LoadModelSuite) TextureMapped(face geometry.TextureMappedFace) {
	suite.faces = append(suite.faces, face)
}

func (suite *LoadModelSuite) aSimpleList() []byte {
	writer := NewWriter()

	writer.WriteHeader(0)
	writer.WriteEndOfNode()

	return writer.Bytes()
}

func (suite *LoadModelSuite) anEmptyModelWith(vertices func(writer *Writer)) []byte {
	writer := NewWriter()

	writer.WriteHeader(0)
	vertices(writer)
	writer.WriteEndOfNode()

	return writer.Bytes()
}
