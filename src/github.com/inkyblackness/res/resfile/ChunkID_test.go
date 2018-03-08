package resfile

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkIDValueReturnsOwnValue(t *testing.T) {
	var id ChunkID

	assert.Equal(t, uint16(0), id.Value(), "initial value should be zero")
	assert.Equal(t, uint16(123), ChunkID(123).Value())
	assert.Equal(t, uint16(math.MaxUint16), ChunkID(math.MaxUint16).Value(), "maximum of uint16 should be supported")
}

func TestChunkIDImplementsIdentifier(t *testing.T) {
	var id interface{} = ChunkID(123)
	_, ok := id.(Identifier)

	assert.True(t, ok)
}
