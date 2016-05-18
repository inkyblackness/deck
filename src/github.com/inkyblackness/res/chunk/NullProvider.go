package chunk

import "github.com/inkyblackness/res"

type nullProvider struct{}

// NullProvider returns a Provider instance that is empty.
// It contains no IDs and will not provide any holder.
func NullProvider() Provider {
	return &nullProvider{}
}

func (*nullProvider) IDs() []res.ResourceID {
	return nil
}

func (*nullProvider) Provide(id res.ResourceID) BlockHolder {
	return nil
}
