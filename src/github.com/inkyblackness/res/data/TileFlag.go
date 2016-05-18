package data

// TileFlag describes simple properties of a map tile.
type TileFlag uint32

const (
	// TileVisited specifies whether the tile has been seen.
	TileVisited TileFlag = 0x80000000
)

func (flag TileFlag) String() (result string) {
	texts := []string{}
	if (flag & TileVisited) != 0 {
		texts = append(texts, "TileVisited")
	}

	for _, text := range texts {
		if len(result) > 0 {
			result += "|"
		}
		result += text
	}

	return
}
