package logic

import (
	"fmt"

	"github.com/inkyblackness/res/data"
)

// LevelObjectChainLinkGetter is a function to return links from a chain.
type LevelObjectChainLinkGetter func(index data.LevelObjectChainIndex) LevelObjectChainLink

// LevelObjectChain handles the logic for a chain of level objects.
type LevelObjectChain struct {
	start LevelObjectChainStart
	link  LevelObjectChainLinkGetter
}

// NewLevelObjectChain returns a new chain based on the given accessors.
func NewLevelObjectChain(start LevelObjectChainStart, linkGetter LevelObjectChainLinkGetter) *LevelObjectChain {
	return &LevelObjectChain{start: start, link: linkGetter}
}

// Initialize resets all entries to an clean state.
// All links will be added to the pool of available entries.
// The provided size is the number of possible links - excluding the start entry.
func (chain *LevelObjectChain) Initialize(size int) {
	chain.start.SetReferenceIndex(0)
	chain.start.SetNextIndex(0)
	chain.start.SetPreviousIndex(0)

	for counter := size; counter > 0; counter-- {
		index := data.LevelObjectChainIndex(counter)
		chain.addLinkToAvailablePool(index)
	}
}

// AcquireLink tries to reserve a new chain link from the chain.
// If the chain is exhausted, an error is returned.
func (chain *LevelObjectChain) AcquireLink() (index data.LevelObjectChainIndex, err error) {
	index = chain.start.PreviousIndex()

	if !index.IsStart() {
		link := chain.link(index)
		previousIndex := chain.start.ReferenceIndex()
		previous := chain.link(previousIndex)

		chain.start.SetPreviousIndex(link.PreviousIndex())
		chain.start.SetReferenceIndex(index)
		link.SetNextIndex(previous.NextIndex())
		previous.SetNextIndex(index)
		link.SetPreviousIndex(previousIndex)
	} else {
		err = fmt.Errorf("Object chain exhausted - Can not add more entries.")
	}

	return
}

// ReleaseLink releases a link from the chain.
func (chain *LevelObjectChain) ReleaseLink(index data.LevelObjectChainIndex) {
	link := chain.link(index)

	chain.link(link.PreviousIndex()).SetNextIndex(link.NextIndex())
	if link.NextIndex().IsStart() {
		chain.start.SetReferenceIndex(link.PreviousIndex())
	} else {
		chain.link(link.NextIndex()).SetPreviousIndex(link.PreviousIndex())
	}
	chain.addLinkToAvailablePool(index)
}

func (chain *LevelObjectChain) addLinkToAvailablePool(index data.LevelObjectChainIndex) {
	link := chain.link(index)

	link.SetPreviousIndex(chain.start.PreviousIndex())
	chain.start.SetPreviousIndex(index)
}
