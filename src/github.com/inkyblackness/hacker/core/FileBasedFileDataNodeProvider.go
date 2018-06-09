package core

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/chunk/resfile"
	"github.com/inkyblackness/res/objprop"
	objDos "github.com/inkyblackness/res/objprop/dos"
	"github.com/inkyblackness/res/textprop"
	textDos "github.com/inkyblackness/res/textprop/dos"
)

type fileBasedFileDataNodeProvider struct {
	access fileAccess
}

func newFileDataNodeProvider(access fileAccess) FileDataNodeProvider {
	provider := &fileBasedFileDataNodeProvider{
		access: access}

	return provider
}

// Provide tries to resolve and return a DataNode for the given file.
func (provider *fileBasedFileDataNodeProvider) Provide(parentNode DataNode, filePath, fileName string) (node DataNode) {
	filePathName := filepath.Join(filePath, fileName)
	rawData, err := provider.access.readFile(filePathName)

	if err == nil {
		lowercaseFileName := strings.ToLower(fileName)
		reader := bytes.NewReader(rawData)

		if lowercaseFileName == "objprop.dat" {
			classes := objprop.StandardProperties()
			objProvider, objErr := objDos.NewProvider(reader, classes)

			if objErr == nil {
				consumerFactory := func() objprop.Consumer {
					outFile, _ := provider.access.createFile(filePathName)
					return objDos.NewConsumer(outFile, classes)
				}
				node = NewObjectPropertiesDataNode(parentNode, fileName, objProvider, classes, consumerFactory)
			}
		} else if lowercaseFileName == "textprop.dat" {
			propProvider, propErr := textDos.NewProvider(reader)

			if propErr == nil {
				consumerFactory := func() textprop.Consumer {
					outFile, _ := provider.access.createFile(filePathName)
					return textDos.NewConsumer(outFile)
				}
				node = NewTexturePropertiesDataNode(parentNode, fileName, propProvider, consumerFactory)
			}
		} else {
			chunkReader, chunkErr := resfile.ReaderFrom(reader)

			if chunkErr == nil {
				saver := func(delegate func(chunk.Store)) {
					store := chunk.NewProviderBackedStore(chunk.NullProvider())
					delegate(store)
					file, _ := provider.access.createFile(filePathName)
					resfile.Write(file, store)
					file.Close()
				}
				node = NewResourceDataNode(parentNode, fileName, chunkReader, saver)
			}
		}
	}

	return
}
