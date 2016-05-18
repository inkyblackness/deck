package query

import (
	"fmt"
)

// Local queries the data source about the position of the hacker and then prints
// all the objects that are in the current tile.
func Local(dataSource DataSource) (result string) {
	gameState := dataSource.GameState()
	tile := dataSource.Tile(int(gameState.HackerX.Tile()), int(gameState.HackerY.Tile()))

	result += fmt.Sprintf("Current Level: %v\n", gameState.CurrentLevel)
	result += fmt.Sprintf("Hacker Coord: X: %v, Y: %v\n", gameState.HackerX, gameState.HackerY)

	refIndex := tile.FirstObjectIndex
	for refIndex != 0 {
		ref := dataSource.LevelObjectCrossReference(refIndex)
		obj := dataSource.LevelObject(ref.LevelObjectTableIndex)

		result += "\n"
		result += fmt.Sprintf("Ref at: X: %v, Y: %v\n", ref.TileX, ref.TileY)
		result += fmt.Sprintf("Obj: Index %d\n%v", ref.LevelObjectTableIndex, obj)
		result += fmt.Sprintf("Path: %v\n", dataSource.ObjectEntryPath(int(obj.Class), int(obj.ClassTableIndex)))

		refIndex = ref.NextObjectIndex
	}

	return
}
