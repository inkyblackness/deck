package chunks

import (
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/serial"
	"github.com/inkyblackness/res/text"
)

// AddArchiveName adds the chunk for the archive name with the provided information.
func AddArchiveName(chunkStore chunk.Store, name string) {
	store := serial.NewByteStore()
	coder := serial.NewPositioningEncoder(store)
	cp := text.DefaultCodepage()
	encodedName := cp.Encode(name)
	nameData := make([]byte, 0x20)

	var copyLength int
	if len(nameData) > len(encodedName) {
		copyLength = len(encodedName) - 1
	} else {
		copyLength = len(nameData)
	}

	copy(nameData, encodedName[0:copyLength])
	coder.Code(nameData)

	chunkStore.Put(chunk.ID(0x0FA0), mapChunk(false, store.Data()))
}
