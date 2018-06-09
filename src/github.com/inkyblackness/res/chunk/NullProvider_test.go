package chunk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullProvider(t *testing.T) {
	provider := NullProvider()
	verifyError := func(id uint16) {
		_, err := provider.Chunk(ID(id))
		assert.Error(t, err)
	}

	assert.Equal(t, 0, len(provider.IDs()))
	verifyError(0)
	verifyError(1)
}
