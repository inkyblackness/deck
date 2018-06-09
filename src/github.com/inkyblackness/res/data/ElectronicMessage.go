package data

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/text"
)

const (
	// textLineLimit was determined through a search of the maximum line length found in the resources.
	textLineLimit = 104

	metaExpressionInterrupt    = "(t)"
	metaExpressionNextMessage  = "(?:i([0-9a-fA-F]{2}))"
	metaExpressionColorIndex   = "(?:c([0-9a-fA-F]{2}))"
	metaExpressionLeftDisplay  = "([0-9]+)"
	metaExpressionRightDisplay = "(?:,[ ]*([0-9]+))"
)

var metaExpression = regexp.MustCompile("^[ ]*(?:(?:" +
	metaExpressionInterrupt + "|" +
	metaExpressionNextMessage + "|" +
	metaExpressionColorIndex + "|" +
	metaExpressionLeftDisplay + "|" +
	metaExpressionRightDisplay + ")[ ]*)*$")

// ElectronicMessage describes one message.
type ElectronicMessage struct {
	nextMessage  int
	isInterrupt  bool
	colorIndex   int
	leftDisplay  int
	rightDisplay int

	title       string
	sender      string
	subject     string
	verboseText string
	terseText   string
}

// NewElectronicMessage returns a new instance of an electronic message.
func NewElectronicMessage() *ElectronicMessage {
	message := &ElectronicMessage{
		nextMessage:  -1,
		colorIndex:   -1,
		leftDisplay:  -1,
		rightDisplay: -1}

	return message
}

// DecodeElectronicMessage tries to decode a message from given block holder.
func DecodeElectronicMessage(cp text.Codepage, provider chunk.BlockProvider) (message *ElectronicMessage, err error) {
	blockIndex := 0
	nextBlockString := func() (line string) {
		if err != nil {
			return
		}
		if blockIndex < provider.BlockCount() {
			blockReader, blockErr := provider.Block(blockIndex)
			if blockErr != nil {
				err = fmt.Errorf("failed to access block %v: %v", blockIndex, blockErr)
				return
			}
			data, dataErr := ioutil.ReadAll(blockReader)
			if dataErr != nil {
				err = fmt.Errorf("failed to read data from block %v: %v", blockIndex, dataErr)
			}
			line = cp.Decode(data)
			blockIndex++
		}
		return
	}
	nextText := func() (text string) {
		for line := nextBlockString(); len(line) > 0; line = nextBlockString() {
			text += line
		}
		return
	}

	metaString := nextBlockString()
	metaData := metaExpression.FindStringSubmatch(metaString)
	parseInt := func(metaIndex int, base, bits int) (result int) {
		var value uint64
		result = -1
		if (err == nil) && (len(metaData[metaIndex]) > 0) {
			value, err = strconv.ParseUint(metaData[metaIndex], base, bits)
			if err == nil {
				result = int(value)
			}
		}
		return
	}

	message = NewElectronicMessage()
	if (metaData == nil) || (len(metaData[0]) != len(metaString)) {
		err = fmt.Errorf("Failed to parse meta string: <%v>", metaString)
	}
	if (err == nil) && (len(metaData[1]) > 0) {
		message.isInterrupt = true
	}
	message.nextMessage = parseInt(2, 16, 16)
	message.colorIndex = parseInt(3, 16, 8)
	message.leftDisplay = parseInt(4, 10, 16)
	message.rightDisplay = parseInt(5, 10, 16)
	message.title = nextBlockString()
	message.sender = nextBlockString()
	message.subject = nextBlockString()
	message.verboseText = nextText()
	message.terseText = nextText()

	return
}

// Encode serializes the message into a block holder.
func (message *ElectronicMessage) Encode(cp text.Codepage) *chunk.Chunk {
	blocks := [][]byte{}

	blocks = append(blocks, cp.Encode(message.metaString()))
	blocks = append(blocks, cp.Encode(message.title))
	blocks = append(blocks, cp.Encode(message.sender))
	blocks = append(blocks, cp.Encode(message.subject))
	for _, line := range message.splitText(message.verboseText) {
		blocks = append(blocks, cp.Encode(line))
	}
	blocks = append(blocks, []byte{0x00})
	for _, line := range message.splitText(message.terseText) {
		blocks = append(blocks, cp.Encode(line))
	}
	blocks = append(blocks, []byte{0x00})

	return &chunk.Chunk{
		Compressed:    false,
		ContentType:   chunk.Text,
		Fragmented:    true,
		BlockProvider: chunk.MemoryBlockProvider(blocks)}
}

func (message *ElectronicMessage) metaString() string {
	result := ""
	append := func(sep, part string) {
		if len(result) > 0 {
			result += sep
		}
		result += part
	}

	if message.isInterrupt {
		append("", "t")
	}
	if message.nextMessage >= 0 {
		append(" ", fmt.Sprintf("i%02X", message.nextMessage))
	}
	if message.colorIndex >= 0 {
		append(" ", fmt.Sprintf("c%02X", message.colorIndex))
	}
	if message.leftDisplay >= 0 {
		append(" ", fmt.Sprintf("%d", message.leftDisplay))
	}
	if message.rightDisplay >= 0 {
		append("", fmt.Sprintf(",%d", message.rightDisplay))
	}

	return result
}

func (message *ElectronicMessage) splitText(input string) []string {
	result := []string{}
	textLines := strings.Split(input, "\n")
	resultLine := ""
	newLine := func() {
		if len(resultLine) > 0 {
			result = append(result, resultLine)
			resultLine = ""
		}
	}

	for textLineIndex, textLine := range textLines {
		words := strings.Split(textLine, " ")
		for wordIndex, word := range words {
			if wordIndex > 0 {
				resultLine += " "
			}
			if (len(resultLine) + len(word)) > textLineLimit {
				newLine()
			}
			resultLine += word
		}
		if textLineIndex < (len(textLines) - 1) {
			resultLine += "\n"
		}
		newLine()
	}

	return result
}

// NextMessage returns the identifier of an interrupting message. Or -1 if no interrupt.
func (message *ElectronicMessage) NextMessage() int {
	return message.nextMessage
}

// SetNextMessage sets the identifier of the interrupting message. -1 to have no interrupt.
func (message *ElectronicMessage) SetNextMessage(id int) {
	message.nextMessage = id
}

// IsInterrupt returns true if this message is an interrupt of another.
func (message *ElectronicMessage) IsInterrupt() bool {
	return message.isInterrupt
}

// SetInterrupt sets whether the message shall be an interrupting message.
func (message *ElectronicMessage) SetInterrupt(value bool) {
	message.isInterrupt = value
}

// ColorIndex returns the color index for the header text. -1 for default color.
func (message *ElectronicMessage) ColorIndex() int {
	return message.colorIndex
}

// SetColorIndex sets the color index for the header text. -1 for default color.
func (message *ElectronicMessage) SetColorIndex(value int) {
	message.colorIndex = value
}

// LeftDisplay returns the display index for the left side. -1 for no display.
func (message *ElectronicMessage) LeftDisplay() int {
	return message.leftDisplay
}

// SetLeftDisplay sets the display index for the left side. -1 for no display.
func (message *ElectronicMessage) SetLeftDisplay(value int) {
	message.leftDisplay = value
}

// RightDisplay returns the display index for the right side. -1 for no display.
func (message *ElectronicMessage) RightDisplay() int {
	return message.rightDisplay
}

// SetRightDisplay sets the display index for the right side. -1 for no display.
func (message *ElectronicMessage) SetRightDisplay(value int) {
	message.rightDisplay = value
}

// Title returns the title of the message.
func (message *ElectronicMessage) Title() string {
	return message.title
}

// SetTitle sets the title of the message.
func (message *ElectronicMessage) SetTitle(value string) {
	message.title = value
}

// Sender returns the sender of the message.
func (message *ElectronicMessage) Sender() string {
	return message.sender
}

// SetSender sets the sender of the message.
func (message *ElectronicMessage) SetSender(value string) {
	message.sender = value
}

// Subject returns the subject of the message.
func (message *ElectronicMessage) Subject() string {
	return message.subject
}

// SetSubject sets the subject of the message.
func (message *ElectronicMessage) SetSubject(value string) {
	message.subject = value
}

// VerboseText returns the verbose text of the message.
func (message *ElectronicMessage) VerboseText() string {
	return message.verboseText
}

// SetVerboseText sets the verbose text of the message.
func (message *ElectronicMessage) SetVerboseText(value string) {
	message.verboseText = value
}

// TerseText returns the terse text of the message.
func (message *ElectronicMessage) TerseText() string {
	return message.terseText
}

// SetTerseText sets the terse text of the message.
func (message *ElectronicMessage) SetTerseText(value string) {
	message.terseText = value
}
