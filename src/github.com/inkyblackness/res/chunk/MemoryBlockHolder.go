package chunk

import (
	"github.com/inkyblackness/res"
)

type memoryBlockHolder struct {
	typeID     TypeID
	dataTypeID res.DataTypeID

	blocks [][]byte
}

// NewBlockHolder returns an in-memory block holder.
func NewBlockHolder(typeID TypeID, dataTypeID res.DataTypeID, blocks [][]byte) BlockHolder {
	holder := &memoryBlockHolder{
		typeID:     typeID,
		dataTypeID: dataTypeID,
		blocks:     blocks}

	return holder
}

// Type returns the type of the chunk.
func (holder *memoryBlockHolder) ChunkType() TypeID {
	return holder.typeID
}

// ContentType returns the type of the data.
func (holder *memoryBlockHolder) ContentType() res.DataTypeID {
	return holder.dataTypeID
}

// BlockCount returns the number of blocks available in the chunk.
// Flat chunks must contain exactly one block.
func (holder *memoryBlockHolder) BlockCount() uint16 {
	return uint16(len(holder.blocks))
}

// BlockData returns the data for the requested block index.
func (holder *memoryBlockHolder) BlockData(block uint16) []byte {
	return holder.blocks[block]
}
