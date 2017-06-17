package data

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/text"

	check "gopkg.in/check.v1"
)

type ElectronicMessageSuite struct {
	cp text.Codepage
}

var _ = check.Suite(&ElectronicMessageSuite{})

func (suite *ElectronicMessageSuite) SetUpTest(c *check.C) {
	suite.cp = text.DefaultCodepage()
}

func (suite *ElectronicMessageSuite) TestEncodeBasicMessage(c *check.C) {
	message := NewElectronicMessage()

	message.SetTitle("1")
	message.SetSender("2")
	message.SetSubject("3")
	message.SetVerboseText("4")
	message.SetTerseText("5")

	holder := message.Encode(suite.cp)

	c.Assert(holder, check.NotNil)
	c.Check(holder.ChunkType(), check.Equals, chunk.BasicChunkType.WithDirectory())
	c.Check(holder.ContentType(), check.Equals, res.Text)
	c.Assert(holder.BlockCount(), check.Equals, uint16(8))
	c.Check(holder.BlockData(0), check.DeepEquals, []byte{0x00})
	c.Check(holder.BlockData(1), check.DeepEquals, []byte{0x31, 0x00})
	c.Check(holder.BlockData(2), check.DeepEquals, []byte{0x32, 0x00})
	c.Check(holder.BlockData(3), check.DeepEquals, []byte{0x33, 0x00})
	c.Check(holder.BlockData(4), check.DeepEquals, []byte{0x34, 0x00})
	c.Check(holder.BlockData(5), check.DeepEquals, []byte{0x00})
	c.Check(holder.BlockData(6), check.DeepEquals, []byte{0x35, 0x00})
	c.Check(holder.BlockData(7), check.DeepEquals, []byte{0x00})
}

func (suite *ElectronicMessageSuite) TestEncodeMeta_A(c *check.C) {
	message := NewElectronicMessage()

	message.SetNextMessage(0x20)
	message.SetColorIndex(0x13)
	message.SetLeftDisplay(30)
	message.SetRightDisplay(40)

	holder := message.Encode(suite.cp)

	c.Assert(holder, check.NotNil)
	c.Assert(holder.BlockCount() > 0, check.Equals, true)
	c.Check(holder.BlockData(0), check.DeepEquals, suite.cp.Encode("i20 c13 30,40"))
}

func (suite *ElectronicMessageSuite) TestEncodeMeta_B(c *check.C) {
	message := NewElectronicMessage()

	message.SetInterrupt(true)
	message.SetLeftDisplay(31)

	holder := message.Encode(suite.cp)

	c.Assert(holder, check.NotNil)
	c.Assert(holder.BlockCount() > 0, check.Equals, true)
	c.Check(holder.BlockData(0), check.DeepEquals, suite.cp.Encode("t 31"))
}

func (suite *ElectronicMessageSuite) TestEncodeMeta_C(c *check.C) {
	message := NewElectronicMessage()

	message.SetInterrupt(true)

	holder := message.Encode(suite.cp)

	c.Assert(holder, check.NotNil)
	c.Assert(holder.BlockCount() > 0, check.Equals, true)
	c.Check(holder.BlockData(0), check.DeepEquals, suite.cp.Encode("t"))
}

func (suite *ElectronicMessageSuite) TestEncodeCreatesNewBlocksPerNewLine(c *check.C) {
	message := NewElectronicMessage()

	message.SetVerboseText("line1\n\n\nline2")
	message.SetTerseText("terse1\n\n\nterse2")

	holder := message.Encode(suite.cp)

	c.Assert(holder, check.NotNil)
	c.Assert(holder.BlockCount() > 0, check.Equals, true)
	c.Check(holder.BlockData(4), check.DeepEquals, suite.cp.Encode("line1\n"))
	c.Check(holder.BlockData(5), check.DeepEquals, suite.cp.Encode("\n"))
	c.Check(holder.BlockData(6), check.DeepEquals, suite.cp.Encode("\n"))
	c.Check(holder.BlockData(7), check.DeepEquals, suite.cp.Encode("line2"))
	c.Check(holder.BlockData(9), check.DeepEquals, suite.cp.Encode("terse1\n"))
	c.Check(holder.BlockData(10), check.DeepEquals, suite.cp.Encode("\n"))
	c.Check(holder.BlockData(11), check.DeepEquals, suite.cp.Encode("\n"))
	c.Check(holder.BlockData(12), check.DeepEquals, suite.cp.Encode("terse2"))
}

func (suite *ElectronicMessageSuite) TestEncodeBreaksUpLinesAfterLimitCharacters(c *check.C) {
	message := NewElectronicMessage()

	message.SetVerboseText("aaaaaaaaa bbbbbbbbb ccccccccc ddddddddd eeeeeeeee fffffffff ggggggggg hhhhhhhhh iiiiiiiii jjjjjjjjj kkkkk")

	holder := message.Encode(suite.cp)

	c.Assert(holder, check.NotNil)
	c.Assert(holder.BlockCount() > 0, check.Equals, true)
	c.Check(holder.BlockData(4), check.DeepEquals,
		suite.cp.Encode("aaaaaaaaa bbbbbbbbb ccccccccc ddddddddd eeeeeeeee fffffffff ggggggggg hhhhhhhhh iiiiiiiii jjjjjjjjj "))
	c.Check(holder.BlockData(5), check.DeepEquals,
		suite.cp.Encode("kkkkk"))
}

func (suite *ElectronicMessageSuite) TestDecodeMeta(c *check.C) {
	message, err := DecodeElectronicMessage(suite.cp, suite.holderWithMeta("i20 c13 30,40"))

	c.Assert(err, check.IsNil)
	c.Assert(message, check.NotNil)
	c.Check(message.NextMessage(), check.Equals, 0x20)
	c.Check(message.ColorIndex(), check.Equals, 0x13)
	c.Check(message.LeftDisplay(), check.Equals, 30)
	c.Check(message.RightDisplay(), check.Equals, 40)
}

func (suite *ElectronicMessageSuite) TestDecodeMeta_Failure(c *check.C) {
	_, err := DecodeElectronicMessage(suite.cp, suite.holderWithMeta("i20 c 13 30,40"))

	c.Check(err, check.NotNil)
}

func (suite *ElectronicMessageSuite) TestDecodeMetaColorIs8BitUnsigned(c *check.C) {
	message, err := DecodeElectronicMessage(suite.cp, suite.holderWithMeta("cD1"))

	c.Assert(err, check.IsNil)
	c.Assert(message, check.NotNil)
	c.Check(message.ColorIndex(), check.Equals, 0xD1)
}

func (suite *ElectronicMessageSuite) TestDecodeMessage(c *check.C) {
	message, err := DecodeElectronicMessage(suite.cp, suite.holderWithMeta("10"))

	c.Assert(err, check.IsNil)
	c.Assert(message, check.NotNil)
	c.Check(message.Title(), check.Equals, "title")
	c.Check(message.Sender(), check.Equals, "sender")
	c.Check(message.Subject(), check.Equals, "subject")
	c.Check(message.VerboseText(), check.Equals, "verbose")
	c.Check(message.TerseText(), check.Equals, "terse")
}

func (suite *ElectronicMessageSuite) TestDecodeMessageIsPossibleForVanillaDummyMails(c *check.C) {
	message, err := DecodeElectronicMessage(suite.cp, suite.vanillaStubMail())

	c.Assert(err, check.IsNil)
	c.Assert(message, check.NotNil)
	c.Check(message.Title(), check.Equals, "")
	c.Check(message.Sender(), check.Equals, "")
	c.Check(message.Subject(), check.Equals, "")
	c.Check(message.VerboseText(), check.Equals, "stub emailstub email")
	c.Check(message.TerseText(), check.Equals, "")
}

func (suite *ElectronicMessageSuite) TestDecodeMessageIsPossibleForMissingTerminatingLine(c *check.C) {
	message, err := DecodeElectronicMessage(suite.cp, suite.holderWithMissingTerminatingLine())

	c.Assert(err, check.IsNil)
	c.Assert(message, check.NotNil)
	c.Check(message.Title(), check.Equals, "title")
	c.Check(message.Sender(), check.Equals, "sender")
	c.Check(message.Subject(), check.Equals, "subject")
	c.Check(message.VerboseText(), check.Equals, "verbose text")
	c.Check(message.TerseText(), check.Equals, "terse text")
}

func (suite *ElectronicMessageSuite) TestRecodeMessage(c *check.C) {
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

	c.Assert(err, check.IsNil)
	c.Assert(outMessage, check.NotNil)
	c.Check(outMessage, check.DeepEquals, inMessage)
}

func (suite *ElectronicMessageSuite) TestRecodeMessageWithMultipleNewLines(c *check.C) {
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

	c.Assert(err, check.IsNil)
	c.Assert(outMessage, check.NotNil)
	c.Check(outMessage, check.DeepEquals, inMessage)
}

func (suite *ElectronicMessageSuite) holderWithMeta(meta string) chunk.BlockHolder {
	blocks := [][]byte{
		suite.cp.Encode(meta),
		suite.cp.Encode("title"),
		suite.cp.Encode("sender"),
		suite.cp.Encode("subject"),
		suite.cp.Encode("verbose"),
		suite.cp.Encode(""),
		suite.cp.Encode("terse"),
		suite.cp.Encode("")}

	return chunk.NewBlockHolder(chunk.BasicChunkType.WithDirectory(), res.Text, blocks)
}

func (suite *ElectronicMessageSuite) vanillaStubMail() chunk.BlockHolder {
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

	return chunk.NewBlockHolder(chunk.BasicChunkType.WithDirectory(), res.Text, blocks)
}

func (suite *ElectronicMessageSuite) holderWithMissingTerminatingLine() chunk.BlockHolder {
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

	return chunk.NewBlockHolder(chunk.BasicChunkType.WithDirectory(), res.Text, blocks)
}
