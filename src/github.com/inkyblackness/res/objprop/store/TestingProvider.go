package store

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/objprop"
)

type TestingProvider struct {
	data map[res.ObjectID]objprop.ObjectData
}

func NewTestingProvider() *TestingProvider {
	provider := &TestingProvider{data: make(map[res.ObjectID]objprop.ObjectData)}

	return provider
}

func (provider *TestingProvider) Provide(id res.ObjectID) objprop.ObjectData {
	return provider.data[id]
}

func (provider *TestingProvider) Consume(id res.ObjectID, data objprop.ObjectData) {
	provider.data[id] = data
}
