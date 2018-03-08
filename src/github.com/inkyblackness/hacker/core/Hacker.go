package core

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/inkyblackness/hacker/styling"
)

// Hacker is the main entry point for the hacker logic.
type Hacker struct {
	style                styling.Style
	fileAccess           fileAccess
	fileDataNodeProvider FileDataNodeProvider

	root    *rootDataNode
	curNode DataNode
}

// NewHacker returns a hacker instance to work with.
func NewHacker(style styling.Style) *Hacker {
	access := realFileAccess
	hacker := &Hacker{
		style:                style,
		fileAccess:           access,
		fileDataNodeProvider: newFileDataNodeProvider(access)}

	return hacker
}

// Load tries to load the data files from the two given directories. The second directory
// is optional.
func (hacker *Hacker) Load(path1, path2 string) string {
	files1, err1 := hacker.fileAccess.readDir(path1)
	var release *ReleaseDesc
	var root *rootDataNode
	result := ""

	if err1 != nil {
		result = hacker.style.Error()("Can't access directories")
	} else if len(path2) == 0 {
		fileNames1 := fileNames(files1)
		release = FindRelease(fileNames1, nil)
		root = newRootDataNode(release)
		root.addChild(newLocationDataNode(root, HD, path1, fileNames1, hacker.fileDataNodeProvider))
	} else {
		files2, err2 := hacker.fileAccess.readDir(path2)

		if err2 == nil {
			fileNames1 := fileNames(files1)
			fileNames2 := fileNames(files2)

			release = FindRelease(fileNames1, fileNames2)
			if release == nil {
				release = FindRelease(fileNames2, fileNames1)
				root = newRootDataNode(release)
				root.addChild(newLocationDataNode(root, HD, path2, fileNames2, hacker.fileDataNodeProvider))
				root.addChild(newLocationDataNode(root, CD, path1, fileNames1, hacker.fileDataNodeProvider))
			} else {
				root = newRootDataNode(release)
				root.addChild(newLocationDataNode(root, HD, path1, fileNames1, hacker.fileDataNodeProvider))
				root.addChild(newLocationDataNode(root, CD, path2, fileNames2, hacker.fileDataNodeProvider))
			}
		} else {
			result = hacker.style.Error()("Can't access directories")
		}
	}
	if release != nil {
		hacker.root = root
		hacker.curNode = root
		result = hacker.style.Status()("Loaded release [", release.name, "]")
	} else if len(result) == 0 {
		result = hacker.style.Error()("Could not resolve release")
	}

	return result
}

// Save re-encodes all loaded data and overwrites the corresponding files.
func (hacker *Hacker) Save() (result string) {
	if hacker.root != nil {
		result = hacker.root.save()
	} else {
		result = hacker.style.Error()(`No data loaded. Use the [load "path1" "path2"] command.`)
	}
	return
}

// Info returns the status of the current node
func (hacker *Hacker) Info() string {
	var result string

	if hacker.curNode != nil {
		result = hacker.curNode.Info()
	} else {
		result = hacker.style.Error()(`No data loaded. Use the [load "path1" "path2"] command.`)
	}

	return result
}

// CurrentDirectory returns the absolute path to the current directory in string form
func (hacker *Hacker) CurrentDirectory() string {
	tempNode := hacker.curNode
	path := ""

	for tempNode != nil && tempNode != hacker.root {
		path = "/" + tempNode.ID() + path
		tempNode = tempNode.Parent()
	}

	return path
}

// ChangeDirectory changes the currently active node
func (hacker *Hacker) ChangeDirectory(path string) (result string) {
	resolved := hacker.resolve(path)

	if resolved != nil {
		hacker.curNode = resolved
		result = ""
	} else {
		result = hacker.style.Error()(`Directory not found: "`, path, `"`)
	}
	return
}

func (hacker *Hacker) resolve(path string) DataNode {
	return hacker.resolveFrom(hacker.curNode, path)
}

func (hacker *Hacker) resolveFrom(baseNode DataNode, path string) (resolved DataNode) {
	parts := strings.Split(path, "/")

	resolved = baseNode
	if parts[0] == "" {
		resolved = hacker.root
	}
	for _, part := range parts {
		if resolved != nil && part != "" {
			if part == ".." {
				resolved = resolved.Parent()
			} else {
				resolved = resolved.Resolve(part)
			}
		}
	}
	return
}

func (hacker *Hacker) Dump() (result string) {
	if hacker.curNode != nil {
		data := hacker.curNode.Data()
		styled := make([]styledData, len(data))
		for index, value := range data {
			styled[index] = styledData{value: value, styleFunc: fmt.Sprint}
		}
		result = createDump(styled)
	}
	return
}

func (hacker *Hacker) Diff(source string) (result string) {
	sourceNode := hacker.resolve(source)
	targetNode := hacker.curNode

	if (targetNode != nil) && (sourceNode != nil) {
		sourceData := sourceNode.UnknownData()
		targetData := targetNode.UnknownData()

		if len(sourceData) > 0 || len(targetData) > 0 {
			result = hacker.diffData(sourceData, targetData)
		} else {
			result = hacker.diffNodes(source, sourceNode, hacker.CurrentDirectory(), targetNode)
		}
	} else {
		result = hacker.style.Error()("Failed to resolve node, check path.")
	}

	return result
}

func (hacker *Hacker) diffData(sourceData []byte, targetData []byte) string {
	sourceLen := len(sourceData)
	targetLen := len(targetData)
	styledSourceData := make([]styledData, sourceLen)
	styledTargetData := make([]styledData, targetLen)
	commonLimit := sourceLen
	if targetLen < commonLimit {
		commonLimit = targetLen
	}
	for index := 0; index < commonLimit; index++ {
		targetByte := targetData[index]
		sourceByte := sourceData[index]
		styledSourceData[index].value = sourceByte
		styledTargetData[index].value = targetByte
		if sourceByte == targetByte {
			styledSourceData[index].styleFunc = fmt.Sprint
			styledTargetData[index].styleFunc = fmt.Sprint
		} else {
			styledSourceData[index].styleFunc = hacker.style.Added()
			styledTargetData[index].styleFunc = hacker.style.Added()
		}
	}
	if sourceLen > targetLen {
		for index, sourceByte := range sourceData[targetLen:] {
			styledSourceData[targetLen+index].value = sourceByte
			styledSourceData[targetLen+index].styleFunc = hacker.style.Removed()
		}
	}
	if targetLen > sourceLen {
		for index, targetByte := range targetData[sourceLen:] {
			styledTargetData[sourceLen+index].value = targetByte
			styledTargetData[sourceLen+index].styleFunc = hacker.style.Removed()
		}
	}

	return createDump(styledSourceData) + "\n" + createDump(styledTargetData)
}

func (hacker *Hacker) diffNodes(sourcePath string, sourceNode DataNode, targetPath string, targetNode DataNode) (result string) {
	sourceChildren := sourceNode.Children()
	targetChildren := targetNode.Children()

	checkNodeExistence := func(basePath string, nodes []DataNode, refNodes []DataNode, diffSign string) {
		for _, node := range nodes {
			nodePath := basePath + "/" + node.ID()
			refNode := hacker.findNodeByID(refNodes, node.ID())

			if refNode == nil {
				result = result + diffSign + " " + nodePath + "\n"
			}
		}

	}
	checkNodeExistence(sourcePath, sourceChildren, targetChildren, "-")
	checkNodeExistence(targetPath, targetChildren, sourceChildren, "+")

	for _, targetChild := range targetChildren {
		sourceChild := hacker.findNodeByID(sourceChildren, targetChild.ID())
		if sourceChild != nil {
			targetChildPath := targetPath + "/" + targetChild.ID()
			result = result + hacker.diffNodes(sourcePath+"/"+sourceChild.ID(), sourceChild, targetChildPath, targetChild)
		}
	}

	if bytes.Compare(sourceNode.UnknownData(), targetNode.UnknownData()) != 0 {
		result = result + "M " + targetPath + "\n"
	}

	return
}

func (hacker *Hacker) findNodeByID(nodes []DataNode, id string) (found DataNode) {
	for _, node := range nodes {
		if node.ID() == id {
			found = node
		}
	}
	return
}

func (hacker *Hacker) Put(offset uint32, data []byte) (result string) {
	if hacker.curNode != nil {
		nodeData := hacker.curNode.Data()
		if int(offset)+len(data) <= len(nodeData) {
			oldData := make([]byte, len(nodeData))
			copy(oldData, nodeData)
			copy(nodeData[offset:], data)
			result = hacker.diffData(oldData, nodeData)
		} else {
			result = hacker.style.Error()(`Data length mismatch`)
		}
	} else {
		result = hacker.style.Error()(`No data loaded`)
	}
	return
}
