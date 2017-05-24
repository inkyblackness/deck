package data

import (
	"fmt"
)

// TileFlag describes simple properties of a map tile.
type TileFlag uint32

const (
	// TileVisited specifies whether the tile has been seen.
	TileVisited TileFlag = 0x80000000

	// MusicIndexTileFlagMask is the mask for the music identifier.
	MusicIndexTileFlagMask = 0x0000F000

	// FloorShadowTileFlagMask is the mask for the floor shadow.
	FloorShadowTileFlagMask = 0x000F0000
	// CeilingShadowTileFlagMask is the mask for the ceiling shadow.
	CeilingShadowTileFlagMask = 0x0F000000
)

// CeilingShadow returns the shadow value for the ceiling.
func (flag TileFlag) CeilingShadow() int {
	return int((flag & CeilingShadowTileFlagMask) >> 24)
}

// WithCeilingShadow returns a new flag combination having given shadow value.
func (flag TileFlag) WithCeilingShadow(value int) TileFlag {
	cleared := uint32(flag) & ^uint32(CeilingShadowTileFlagMask)
	newValue := (uint32(value) << 24) & CeilingShadowTileFlagMask
	return TileFlag(cleared | newValue)
}

// FloorShadow returns the shadow value for the floor.
func (flag TileFlag) FloorShadow() int {
	return int((flag & FloorShadowTileFlagMask) >> 16)
}

// WithFloorShadow returns a new flag combination having given shadow value.
func (flag TileFlag) WithFloorShadow(value int) TileFlag {
	cleared := uint32(flag) & ^uint32(FloorShadowTileFlagMask)
	newValue := (uint32(value) << 16) & FloorShadowTileFlagMask
	return TileFlag(cleared | newValue)
}

// MusicIndex returns the music identifier for the tile.
func (flag TileFlag) MusicIndex() int {
	return int((flag & MusicIndexTileFlagMask) >> 12)
}

// WithMusicIndex returns a new flag combination having given music identifier.
func (flag TileFlag) WithMusicIndex(value int) TileFlag {
	cleared := uint32(flag) & ^uint32(MusicIndexTileFlagMask)
	newValue := (uint32(value) << 12) & MusicIndexTileFlagMask
	return TileFlag(cleared | newValue)
}

func (flag TileFlag) String() (result string) {
	texts := []string{}
	if (flag & TileVisited) != 0 {
		texts = append(texts, "TileVisited")
	}
	texts = append(texts, fmt.Sprintf("MusicIndex=%v", flag.MusicIndex()))

	for _, text := range texts {
		if len(result) > 0 {
			result += "|"
		}
		result += text
	}

	return
}
