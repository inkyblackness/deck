package core

type rootDataNode struct {
	parentDataNode

	release *ReleaseDesc
}

func newRootDataNode(release *ReleaseDesc) *rootDataNode {
	node := &rootDataNode{
		parentDataNode: makeParentDataNode(nil, "", 2),
		release:        release}

	return node
}

func (node *rootDataNode) Info() string {
	info := "Release: [" + node.release.name + "]"
	info = info + "\nAvailable data locations:"
	for _, child := range node.Children() {
		info = info + " " + child.ID()
	}

	return info
}

func (node *rootDataNode) save() (result string) {
	for _, child := range node.Children() {
		locationNode := child.(saveable)
		result += locationNode.save()
	}
	return
}
