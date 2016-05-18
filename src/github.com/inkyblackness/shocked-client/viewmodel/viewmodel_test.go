package viewmodel

import (
	"fmt"
	"math/rand"

	"testing"

	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

func someNode() Node {
	return NewBoolValueNode(fmt.Sprintf("node%v", rand.Intn(100)), false)
}

func someNodeList() []Node {
	nodes := make([]Node, rand.Intn(10))

	for index := range nodes {
		nodes[index] = someNode()
	}

	return nodes
}

func someNodeMap() map[string]Node {
	list := someNodeList()
	result := make(map[string]Node)

	for index, node := range list {
		result[fmt.Sprintf("key%v", index)] = node
	}

	return result
}
