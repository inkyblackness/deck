package graphics

type testingTexture struct {
	disposed bool
}

func aTexture() Texture {
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
