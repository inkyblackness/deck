package store

type TestingProvider struct {
	data map[uint32][]byte
}

func NewTestingProvider() *TestingProvider {
	provider := &TestingProvider{data: make(map[uint32][]byte)}

	return provider
}

func (provider *TestingProvider) EntryCount() uint32 {
	return uint32(len(provider.data))
}

func (provider *TestingProvider) Provide(id uint32) []byte {
	return provider.data[id]
}

func (provider *TestingProvider) Consume(id uint32, data []byte) {
	provider.data[id] = data
}
