package command

import (
	"bytes"

	"github.com/inkyblackness/res/geometry"

	check "gopkg.in/check.v1"
)

type ResaveSuite struct {
	model *geometry.DynamicModel
}

var _ = check.Suite(&ResaveSuite{})

func (suite *ResaveSuite) SetUpTest(c *check.C) {
	suite.model = geometry.NewDynamicModel()
}

func (suite *ResaveSuite) TestEmptyModel(c *check.C) {
	suite.verifyResave(c)
}

func (suite *ResaveSuite) TestSingleVertex(c *check.C) {
	suite.model.AddVertex(geometry.NewSimpleVertex(NewFixedVector(Vector{0, 0, 0})))
	suite.verifyResave(c)
}

func (suite *ResaveSuite) TestVertices(c *check.C) {
	suite.model.AddVertex(geometry.NewSimpleVertex(NewFixedVector(Vector{0, 0, 0})))
	suite.model.AddVertex(geometry.NewSimpleVertex(NewFixedVector(Vector{1, 2, 3})))
	suite.model.AddVertex(geometry.NewSimpleVertex(NewFixedVector(Vector{0, 0, 10})))
	suite.model.AddVertex(geometry.NewSimpleVertex(NewFixedVector(Vector{0, 10, 10})))
	suite.model.AddVertex(geometry.NewSimpleVertex(NewFixedVector(Vector{11, 0, 0})))
	suite.model.AddVertex(geometry.NewSimpleVertex(NewFixedVector(Vector{11, 11, 0})))
	suite.model.AddVertex(geometry.NewSimpleVertex(NewFixedVector(Vector{12, 0, 12})))
	suite.model.AddVertex(geometry.NewSimpleVertex(NewFixedVector(Vector{0, 13, 0})))
	suite.verifyResave(c)
}

func (suite *ResaveSuite) TestNodeAnchor(c *check.C) {
	node1 := geometry.NewDynamicNode()
	node11 := geometry.NewDynamicNode()
	node12 := geometry.NewDynamicNode()
	node1.AddAnchor(geometry.NewSimpleNodeAnchor(suite.aNormal(), suite.aReference(), node11, node12))
	node2 := geometry.NewDynamicNode()
	suite.model.AddAnchor(geometry.NewSimpleNodeAnchor(suite.aNormal(), suite.aReference(), node1, node2))
	suite.verifyResave(c)
}

func (suite *ResaveSuite) TestEmptyFaceAnchor(c *check.C) {
	suite.model.AddAnchor(geometry.NewDynamicFaceAnchor(suite.aNormal(), suite.aReference()))
	suite.verifyResave(c)
}

func (suite *ResaveSuite) TestFaces(c *check.C) {
	suite.model.AddVertex(geometry.NewSimpleVertex(NewFixedVector(Vector{0, 0, 0})))
	suite.model.AddVertex(geometry.NewSimpleVertex(NewFixedVector(Vector{1, 0, 0})))
	suite.model.AddVertex(geometry.NewSimpleVertex(NewFixedVector(Vector{1, 1, 0})))

	anchor := geometry.NewDynamicFaceAnchor(suite.aNormal(), suite.aReference())
	anchor.AddFace(geometry.NewSimpleFlatColoredFace([]int{0, 1, 2}, geometry.ColorIndex(0x23)))
	anchor.AddFace(geometry.NewSimpleShadeColoredFace([]int{0, 1, 2}, geometry.ColorIndex(0x23), 2))
	anchor.AddFace(geometry.NewSimpleTextureMappedFace([]int{0, 1, 2}, 0x1234, suite.someTextureCoordinates()))

	suite.model.AddAnchor(anchor)
	suite.verifyResave(c)
}

func (suite *ResaveSuite) verifyResave(c *check.C) {
	saved1 := SaveModel(suite.model)
	newModel, err := LoadModel(bytes.NewReader(saved1))
	c.Assert(err, check.IsNil)
	saved2 := SaveModel(newModel)
	c.Assert(saved2, check.DeepEquals, saved1)
}

func (suite *ResaveSuite) aNormal() geometry.Vector {
	return NewFixedVector(Vector{1, 0, 0})
}

func (suite *ResaveSuite) aReference() geometry.Vector {
	return NewFixedVector(Vector{0, 0, 0})
}

func (suite *ResaveSuite) someTextureCoordinates() []geometry.TextureCoordinate {
	return []geometry.TextureCoordinate{
		geometry.NewSimpleTextureCoordinate(0, 0, 0),
		geometry.NewSimpleTextureCoordinate(1, 1, 0),
		geometry.NewSimpleTextureCoordinate(2, 1, 1)}
}
