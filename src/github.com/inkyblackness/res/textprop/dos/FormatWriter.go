package dos

import (
	"encoding/binary"

	"github.com/inkyblackness/res/serial"
	"github.com/inkyblackness/res/textprop"
)

type formatWriter struct {
	dest serial.SeekingWriteCloser
}

// NewConsumer wraps the provided Writer in a consumer for text properties.
func NewConsumer(dest serial.SeekingWriteCloser) textprop.Consumer {
	writer := &formatWriter{dest: dest}

	binary.Write(dest, binary.LittleEndian, MagicHeader)
	writer.dest.Write(make([]byte, 4000))

	return writer
}

// Consume takes the provided entry data and adds it to the stream for given ID.
func (writer *formatWriter) Consume(id uint32, data []byte) {
	writer.dest.Seek(int64(MagicHeaderSize+textprop.TexturePropertiesLength*id), 0)
	writer.dest.Write(data)
}

func (writer *formatWriter) Finish() {
	writer.dest.Close()
}
