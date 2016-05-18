package store

import (
	"crypto/rand"
	simple "math/rand"
	"testing"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"

	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

func emptyBlockHolder() chunk.BlockHolder {
	return chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{nil})
}

func randomBlockHolder(blockCount int) chunk.BlockHolder {
	blocks := make([][]byte, blockCount)

	for index := range blocks {
		block := make([]byte, simple.Intn(5))
		rand.Read(block)
		blocks[index] = block
	}

	return chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, blocks)
}
