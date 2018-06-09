package data

import (
	"io/ioutil"
	"testing"

	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ElectronicMessageSuite struct {
	suite.Suite
	cp text.Codepage
}

func TestElectronicMessageSuite(t *testing.T) {
	suite.Run(t, new(ElectronicMessageSuite))
}

func (suite *ElectronicMessageSuite) SetupTest() {
	suite.cp = text.DefaultCodepage()
}

func (suite *ElectronicMessageSuite) TestEncodeBasicMessage() {
	message := NewElectronicMessage()

	message.SetTitle("1")
	message.SetSender("2")
	message.SetSubject("3")
	message.SetVerboseText("4")
	message.SetTerseText("5")

	encoded := message.Encode(suite.cp)

	require.NotNil(suite.T(), encoded)
	assert.Equal(suite.T(), true, encoded.Fragmented)
	assert.Equal(suite.T(), chunk.Text, encoded.ContentType)
	assert.Equal(suite.T(), 8, encoded.BlockCount())

	suite.verifyBlock(0, encoded, []byte{0x00})
	suite.verifyBlock(1, encoded, []byte{0x31, 0x00})
	suite.verifyBlock(2, encoded, []byte{0x32, 0x00})
	suite.verifyBlock(3, encoded, []byte{0x33, 0x00})
	suite.verifyBlock(4, encoded, []byte{0x34, 0x00})
	suite.verifyBlock(5, encoded, []byte{0x00})
	suite.verifyBlock(6, encoded, []byte{0x35, 0x00})
	suite.verifyBlock(7, encoded, []byte{0x00})
}

func (suite *ElectronicMessageSuite) TestEncodeMeta_A() {
	message := NewElectronicMessage()

	message.SetNextMessage(0x20)
	message.SetColorIndex(0x13)
	message.SetLeftDisplay(30)
	message.SetRightDisplay(40)

	encoded := message.Encode(suite.cp)

	require.NotNil(suite.T(), encoded)
	require.Equal(suite.T(), true, encoded.BlockCount() > 0)
	suite.verifyBlock(0, encoded, suite.cp.Encode("i20 c13 30,40"))
}

func (suite *ElectronicMessageSuite) TestEncodeMeta_B() {
	message := NewElectronicMessage()

	message.SetInterrupt(true)
	message.SetLeftDisplay(31)

	encoded := message.Encode(suite.cp)

	require.NotNil(suite.T(), encoded)
	require.Equal(suite.T(), true, encoded.BlockCount() > 0)
	suite.verifyBlock(0, encoded, suite.cp.Encode("t 31"))
}

func (suite *ElectronicMessageSuite) TestEncodeMeta_C() {
	message := NewElectronicMessage()

	message.SetInterrupt(true)

	encoded := message.Encode(suite.cp)

	require.NotNil(suite.T(), encoded)
	require.Equal(suite.T(), true, encoded.BlockCount() > 0)
	suite.verifyBlock(0, encoded, suite.cp.Encode("t"))
}

func (suite *ElectronicMessageSuite) TestEncodeCreatesNewBlocksPerNewLine() {
	message := NewElectronicMessage()

	message.SetVerboseText("line1\n\n\nline2")
	message.SetTerseText("terse1\n\n\nterse2")

	encoded := message.Encode(suite.cp)

	require.NotNil(suite.T(), encoded)
	require.Equal(suite.T(), true, encoded.BlockCount() > 0)
	suite.verifyBlock(4, encoded, suite.cp.Encode("line1\n"))
	suite.verifyBlock(5, encoded, suite.cp.Encode("\n"))
	suite.verifyBlock(6, encoded, suite.cp.Encode("\n"))
	suite.verifyBlock(7, encoded, suite.cp.Encode("line2"))
	suite.verifyBlock(9, encoded, suite.cp.Encode("terse1\n"))
	suite.verifyBlock(10, encoded, suite.cp.Encode("\n"))
	suite.verifyBlock(11, encoded, suite.cp.Encode("\n"))
	suite.verifyBlock(12, encoded, suite.cp.Encode("terse2"))
}

func (suite *ElectronicMessageSuite) TestEncodeBreaksUpLinesAfterLimitCharacters() {
	message := NewElectronicMessage()

	message.SetVerboseText("aaaaaaaaa bbbbbbbbb ccccccccc ddddddddd eeeeeeeee fffffffff ggggggggg hhhhhhhhh iiiiiiiii jjjjjjjjj kkkkk")

	encoded := message.Encode(suite.cp)

	require.NotNil(suite.T(), encoded)
	require.Equal(suite.T(), true, encoded.BlockCount() > 0)
	suite.verifyBlock(4, encoded,
		suite.cp.Encode("aaaaaaaaa bbbbbbbbb ccccccccc ddddddddd eeeeeeeee fffffffff ggggggggg hhhhhhhhh iiiiiiiii jjjjjjjjj "))
	suite.verifyBlock(5, encoded,
		suite.cp.Encode("kkkkk"))
}

func (suite *ElectronicMessageSuite) TestDecodeMeta() {
	message, err := DecodeElectronicMessage(suite.cp, suite.holderWithMeta("i20 c13 30,40"))

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), message)
	assert.Equal(suite.T(), 0x20, message.NextMessage())
	assert.Equal(suite.T(), 0x13, message.ColorIndex())
	assert.Equal(suite.T(), 30, message.LeftDisplay())
	assert.Equal(suite.T(), 40, message.RightDisplay())
}

func (suite *ElectronicMessageSuite) TestDecodeMeta_Failure() {
	_, err := DecodeElectronicMessage(suite.cp, suite.holderWithMeta("i20 c 13 30,40"))

	assert.NotNil(suite.T(), err)
}

func (suite *ElectronicMessageSuite) TestDecodeMetaColorIs8BitUnsigned() {
	message, err := DecodeElectronicMessage(suite.cp, suite.holderWithMeta("cD1"))

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), message)
	assert.Equal(suite.T(), 0xD1, message.ColorIndex())
}

func (suite *ElectronicMessageSuite) TestDecodeMessage() {
	message, err := DecodeElectronicMessage(suite.cp, suite.holderWithMeta("10"))

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), message)
	assert.Equal(suite.T(), "title", message.Title())
	assert.Equal(suite.T(), "sender", message.Sender())
	assert.Equal(suite.T(), "subject", message.Subject())
	assert.Equal(suite.T(), "verbose", message.VerboseText())
	assert.Equal(suite.T(), "terse", message.TerseText())
}

func (suite *ElectronicMessageSuite) TestDecodeMessageIsPossibleForVanillaDummyMails() {
	message, err := DecodeElectronicMessage(suite.cp, suite.vanillaStubMail())

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), message)
	assert.Equal(suite.T(), "", message.Title())
	assert.Equal(suite.T(), "", message.Sender())
	assert.Equal(suite.T(), "", message.Subject())
	assert.Equal(suite.T(), "stub emailstub email", message.VerboseText())
	assert.Equal(suite.T(), "", message.TerseText())
}

func (suite *ElectronicMessageSuite) TestDecodeMessageIsPossibleForMissingTerminatingLine() {
	message, err := DecodeElectronicMessage(suite.cp, suite.holderWithMissingTerminatingLine())

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), message)
	assert.Equal(suite.T(), "title", message.Title())
	assert.Equal(suite.T(), "sender", message.Sender())
	assert.Equal(suite.T(), "subject", message.Subject())
	assert.Equal(suite.T(), "verbose text", message.VerboseText())
	assert.Equal(suite.T(), "terse text", message.TerseText())
}

func (suite *ElectronicMessageSuite) TestRecodeMessage() {
	inMessage := NewElectronicMessage()
	inMessage.SetInterrupt(true)
	inMessage.SetNextMessage(0x10)
	inMessage.SetColorIndex(0x20)
	inMessage.SetLeftDisplay(40)
	inMessage.SetRightDisplay(50)
	inMessage.SetVerboseText("abcd\nefgh\nsome")
	inMessage.SetTerseText("\n")

	holder := inMessage.Encode(suite.cp)
	outMessage, err := DecodeElectronicMessage(suite.cp, holder)

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), outMessage)
	assert.Equal(suite.T(), inMessage, outMessage)
}

func (suite *ElectronicMessageSuite) TestRecodeMessageWithMultipleNewLines() {
	inMessage := NewElectronicMessage()
	inMessage.SetInterrupt(true)
	inMessage.SetNextMessage(0x10)
	inMessage.SetColorIndex(0x20)
	inMessage.SetLeftDisplay(40)
	inMessage.SetRightDisplay(50)
	inMessage.SetVerboseText("first\n\n\nsecond")
	inMessage.SetTerseText("terse\n")

	holder := inMessage.Encode(suite.cp)
	outMessage, err := DecodeElectronicMessage(suite.cp, holder)

	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), outMessage)
	assert.Equal(suite.T(), inMessage, outMessage)
}

func (suite *ElectronicMessageSuite) verifyBlock(index int, provider chunk.BlockProvider, expected []byte) {
	reader, readerErr := provider.Block(index)
	require.Nil(suite.T(), readerErr)
	data, dataErr := ioutil.ReadAll(reader)
	require.Nil(suite.T(), dataErr)
	assert.Equal(suite.T(), expected, data)
}

func (suite *ElectronicMessageSuite) holderWithMeta(meta string) chunk.BlockProvider {
	blocks := [][]byte{
		suite.cp.Encode(meta),
		suite.cp.Encode("title"),
		suite.cp.Encode("sender"),
		suite.cp.Encode("subject"),
		suite.cp.Encode("verbose"),
		suite.cp.Encode(""),
		suite.cp.Encode("terse"),
		suite.cp.Encode("")}

	return chunk.MemoryBlockProvider(blocks)
}

func (suite *ElectronicMessageSuite) vanillaStubMail() chunk.BlockProvider {
	// The string resources contain a few mails which aren't used.
	// They are missing the terminating line for the verbose text.
	blocks := [][]byte{
		suite.cp.Encode(""),
		suite.cp.Encode(""),
		suite.cp.Encode(""),
		suite.cp.Encode(""),
		suite.cp.Encode("stub email"),
		suite.cp.Encode("stub email"),
		suite.cp.Encode("")}

	return chunk.MemoryBlockProvider(blocks)
}

func (suite *ElectronicMessageSuite) holderWithMissingTerminatingLine() chunk.BlockProvider {
	// This case is encountered once in gerstrng.res
	blocks := [][]byte{
		suite.cp.Encode(""),
		suite.cp.Encode("title"),
		suite.cp.Encode("sender"),
		suite.cp.Encode("subject"),
		suite.cp.Encode("verbose "),
		suite.cp.Encode("text"),
		suite.cp.Encode(""),
		suite.cp.Encode("terse "),
		suite.cp.Encode("text")}

	return chunk.MemoryBlockProvider(blocks)
}
