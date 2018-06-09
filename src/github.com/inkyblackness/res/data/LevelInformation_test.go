package data

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevelInformationEncodedSize(t *testing.T) {
	info := DefaultLevelInformation()

	size := binary.Size(info)

	assert.Equal(t, 58, size)
}
