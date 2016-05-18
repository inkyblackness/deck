package core

import (
	"sync"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"

	"github.com/inkyblackness/shocked-core/io"
)

// Archive wraps a map archive
type Archive struct {
	mutex sync.Mutex

	store chunk.Store

	levels [MaximumLevelsPerArchive]*Level
}

// NewArchive creates a new archive wrapper for given store name. This wrapper is
// applicable for the "starting game" archive, as well as savegames.
func NewArchive(library io.StoreLibrary, storeName string) (archive *Archive, err error) {
	var store chunk.Store

	store, err = library.ChunkStore(storeName)

	if err == nil {
		archive = &Archive{store: store}
	}

	return
}

// HasLevel returns true when given level ID (0..15) refers to a valid level.
func (archive *Archive) HasLevel(id int) bool {
	return archive.store.Get(res.ResourceID(4000+id*100+4)) != nil
}

func (archive *Archive) LevelIDs() (result []int) {
	archive.mutex.Lock()
	defer archive.mutex.Unlock()

	for i := 0; i < MaximumLevelsPerArchive; i++ {
		if archive.HasLevel(i) {
			result = append(result, i)
		}
	}

	return
}

// Level returns a level wrapper should the given ID refer to a valid level.
func (archive *Archive) Level(id int) (level *Level) {
	archive.mutex.Lock()
	defer archive.mutex.Unlock()

	if archive.HasLevel(id) {
		level = archive.levels[id]

		if level == nil {
			level = NewLevel(archive.store, id)
			archive.levels[id] = level
		}
	}

	return
}
