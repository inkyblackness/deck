package chunk

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkIDValueReturnsOwnValue(t *testing.T) {
	assert.Equal(t, uint16(0), ID(0).Value())
	assert.Equal(t, uint16(123), ID(123).Value())
	assert.Equal(t, uint16(math.MaxUint16), ID(math.MaxUint16).Value(), "maximum of uint16 should be supported")
}

func TestChunkIDImplementsStringer(t *testing.T) {
	assert.Equal(t, "0FA0", fmt.Sprintf("%v", ID(0x0FA0)))
}
