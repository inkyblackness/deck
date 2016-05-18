package model

import (
	"github.com/inkyblackness/shocked-client/graphics"
)

type testingTexture struct {
	disposed bool
}

func aTexture() graphics.Texture {
	return aTestingTexture()
}

func aTestingTexture() *testingTexture {
	return &testingTexture{}
}

func (tex *testingTexture) Dispose() {
	tex.disposed = true
}

func (tex *testingTexture) Handle() uint32 {
	return 0
}
