package core

import (
	"github.com/inkyblackness/res/serial"
	"github.com/inkyblackness/res/textprop"

	check "gopkg.in/check.v1"
)

type FileBasedFileDataNodeProviderSuite struct {
	parentNode DataNode

	provider FileDataNodeProvider

	input  []byte
	output *serial.ByteStore
}

var _ = check.Suite(&FileBasedFileDataNodeProviderSuite{})

func (suite *FileBasedFileDataNodeProviderSuite) SetUpTest(c *check.C) {
	access := fileAccess{
		readDir:  nil,
		readFile: func(filename string) ([]byte, error) { return suite.input, nil },
		createFile: func(filename string) (serial.SeekingWriteCloser, error) {
			suite.output = serial.NewByteStore()
			return suite.output, nil
		}}

	suite.provider = newFileDataNodeProvider(access)
}

func (suite *FileBasedFileDataNodeProviderSuite) TestProviderCanOpenTextureProperties(c *check.C) {
	suite.input = []byte{0x09, 0x00, 0x00, 0x00}
	suite.input = append(suite.input, make([]byte, textprop.TexturePropertiesLength)...)

	node := suite.provider.Provide(suite.parentNode, ".", "textprop.dat")

	c.Check(node, check.Not(check.IsNil))
}

func (suite *FileBasedFileDataNodeProviderSuite) TestProviderForwardsWriterForSaving(c *check.C) {
	suite.input = []byte{0x09, 0x00, 0x00, 0x00}
	suite.input = append(suite.input, make([]byte, textprop.TexturePropertiesLength)...)

	node := suite.provider.Provide(suite.parentNode, ".", "textprop.dat")
	saver := node.(saveable)
	saver.save()

	c.Check(suite.output, check.Not(check.IsNil))
}
