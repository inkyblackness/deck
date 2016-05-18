package core

import (
	"strings"
)

type locationDataNode struct {
	parentDataNode

	dataLocation DataLocation
	filePath     string
	fileNames    []string
}

func newLocationDataNode(parentNode DataNode, dataLocation DataLocation,
	filePath string, fileNames []string, fileDataNodeProvider FileDataNodeProvider) *locationDataNode {
	node := &locationDataNode{
		parentDataNode: makeParentDataNode(parentNode, dataLocation.String(), len(fileNames)),
		dataLocation:   dataLocation,
		filePath:       filePath,
		fileNames:      fileNames}

	node.setChildResolver(func(path string) (resolved DataNode) {
		if resolvedName := node.resolveFileName(path); resolvedName != "" {
			resolved = fileDataNodeProvider.Provide(node, node.filePath, resolvedName)
		}
		return
	})

	return node
}

func (node *locationDataNode) Info() string {
	info := "Location: " + string(node.dataLocation) + "\n"
	info = info + "FilePath: [" + node.filePath + "]\n"
	info = info + "Files:"
	for _, fileName := range node.fileNames {
		info = info + " " + fileName
	}

	return info
}

func (node *locationDataNode) resolveFileName(path string) (result string) {
	lowerPath := strings.ToLower(path)

	for _, knownName := range node.fileNames {
		if strings.ToLower(knownName) == lowerPath {
			result = knownName
		}
	}
	return
}

func (node *locationDataNode) save() (result string) {
	for _, child := range node.Children() {
		childSaveable := child.(saveable)
		result += childSaveable.save()
	}
	return
}
