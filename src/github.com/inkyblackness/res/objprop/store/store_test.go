package store

import (
	"crypto/rand"

	"github.com/inkyblackness/res/objprop"

	"testing"

	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

func randomObjectData() objprop.ObjectData {
	data := objprop.ObjectData{
		Generic:  make([]byte, 2),
		Specific: make([]byte, 4),
		Common:   make([]byte, 3)}

	rand.Read(data.Generic)
	rand.Read(data.Specific)
	rand.Read(data.Common)

	return data
}
