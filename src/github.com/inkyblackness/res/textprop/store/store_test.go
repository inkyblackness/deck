package store

import (
	"crypto/rand"

	"github.com/inkyblackness/res/textprop"

	"testing"

	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

func randomProperties() []byte {
	data := make([]byte, textprop.TexturePropertiesLength)

	rand.Read(data)

	return data
}
