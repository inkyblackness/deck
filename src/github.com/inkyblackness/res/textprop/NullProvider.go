package textprop

type nullProvider struct{}

// NullProvider returns a Provider instance that is empty/reset.
func NullProvider() Provider {
	return &nullProvider{}
}

// EntryCount implements the Provider interface.
func (provider *nullProvider) EntryCount() uint32 {
	return 0
}

// Provide implements the Provider interface.
func (provider *nullProvider) Provide(id uint32) []byte {
	return make([]byte, TexturePropertiesLength)
}
