package objprop

import (
	"github.com/inkyblackness/res"
)

type nullProvider struct {
	descriptors []ClassDescriptor
}

// NullProvider returns a provider that has all properties reset to zero.
func NullProvider(descriptors []ClassDescriptor) (provider Provider) {
	return &nullProvider{descriptors: descriptors}
}

// Provide implements the Provider interface
func (provider *nullProvider) Provide(id res.ObjectID) ObjectData {
	classDesc := provider.descriptors[int(id.Class)]
	subclassDesc := classDesc.Subclasses[int(id.Subclass)]
	data := ObjectData{
		Generic:  make([]byte, classDesc.GenericDataLength),
		Specific: make([]byte, subclassDesc.SpecificDataLength),
		Common:   make([]byte, CommonPropertiesLength)}

	return data
}
