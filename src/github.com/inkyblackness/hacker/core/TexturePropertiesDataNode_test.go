package core

import (
	"github.com/inkyblackness/res/textprop"

	check "gopkg.in/check.v1"
)

type TexturePropertiesDataNodeSuite struct {
	parentNode DataNode
	name       string

	output          *TestingTexturePropertiesProvidingConsumer
	consumerFactory func() textprop.Consumer
}

var _ = check.Suite(&TexturePropertiesDataNodeSuite{})

func (suite *TexturePropertiesDataNodeSuite) SetUpTest(c *check.C) {
	suite.name = "textprop.dat"
	suite.consumerFactory = func() textprop.Consumer {
		suite.output = &TestingTexturePropertiesProvidingConsumer{}
		return suite.output
	}
}

func (suite *TexturePropertiesDataNodeSuite) makeProperties(filler byte) []byte {
	data := make([]byte, textprop.TexturePropertiesLength)

	for index := range data {
		data[index] = filler
	}

	return data
}

func (suite *TexturePropertiesDataNodeSuite) TestSaveWritesTextureData(c *check.C) {
	prop1 := suite.makeProperties(1)
	prop2 := suite.makeProperties(2)
	provider := &TestingTexturePropertiesProvidingConsumer{textureData: [][]byte{prop1, prop2}}

	node := NewTexturePropertiesDataNode(suite.parentNode, suite.name, provider, suite.consumerFactory)
	saver := node.(saveable)
	saver.save()

	c.Check(suite.output.textureData, check.DeepEquals, provider.textureData)
}
