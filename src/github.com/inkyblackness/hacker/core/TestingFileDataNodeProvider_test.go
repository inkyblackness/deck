package core

type TestingFileDataNodeProvider struct {
	nodesByFileName map[string]DataNode
}

func NewTestingFileDataNodeProvider() *TestingFileDataNodeProvider {
	provider := &TestingFileDataNodeProvider{
		nodesByFileName: make(map[string]DataNode)}

	return provider
}

// Provide tries to resolve and return a DataNode for the given file.
func (provider *TestingFileDataNodeProvider) Provide(parent DataNode, filePath, fileName string) (provided DataNode) {
	temp, existing := provider.nodesByFileName[fileName]

	if existing {
		provided = temp
	}

	return
}
