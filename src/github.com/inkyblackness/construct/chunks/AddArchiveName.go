package chunks

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/res/serial"
)

// AddArchiveName adds the chunk for the archive name with the provided information
func AddArchiveName(consumer chunk.Consumer, name string) {
	store := serial.NewByteStore()
	coder := serial.NewPositioningEncoder(store)

	coder.CodeBytes(make([]byte, 0x20))
	coder.SetCurPos(0)
	serial.MapData(&data.String{Value: name}, coder)

	blocks := [][]byte{store.Data()}
	consumer.Consume(res.ResourceID(0x0FA0), chunk.NewBlockHolder(chunk.BasicChunkType, res.Map, blocks))
}
