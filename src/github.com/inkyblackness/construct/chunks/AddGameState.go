package chunks

import (
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/res/serial"
)

// AddGameState adds the chunk for the game state.
func AddGameState(chunkStore chunk.Store) {
	store := serial.NewByteStore()
	coder := serial.NewPositioningEncoder(store)
	state := data.DefaultGameState()

	coder.Code(make([]byte, data.GameStateSize))
	coder.SetCurPos(0)
	coder.Code(&state)

	chunkStore.Put(chunk.ID(0x0FA1), mapChunk(false, store.Data()))
}
