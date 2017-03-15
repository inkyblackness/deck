package model

// FailureFunc is for failed queries.
type FailureFunc func()

// DataStore describes the necessary methods for querying and modifying model data.
type DataStore interface {
	// Projects queries the currently available projects.
	Projects(onSuccess func(projects []string), onFailure FailureFunc)
	// NewProject requests to create a new project with given ID.
	NewProject(projectID string, onSuccess func(), onFailure FailureFunc)

	// Font queries a specific font.
	Font(projectID string, fontID int, onSuccess func(font *Font), onFailure FailureFunc)

	// GameObjectIcon queries the icon bitmap of a game object.
	GameObjectIcon(projectID string, class, subclass, objType int, onSuccess func(bmp *RawBitmap), onFailure FailureFunc)

	// Palette queries a palette.
	Palette(projectID string, paletteID string, onSuccess func(colors [256]Color), onFailure FailureFunc)
	// Levels queries all levels of a project.
	Levels(projectID string, archiveID string, onSuccess func(levels []Level), onFailure FailureFunc)

	// LevelTextures queries the texture IDs for a level.
	LevelTextures(projectID string, archiveID string, levelID int, onSuccess func(textureIDs []int), onFailure FailureFunc)
	// SetLevelTextures requests to set the list of textures for a level.
	SetLevelTextures(projectID string, archiveID string, levelID int, textureIDs []int, onSuccess func(textureIDs []int), onFailure FailureFunc)

	// Textures queries all texture information of a project.
	Textures(projectID string, onSuccess func(textures []Texture), onFailure FailureFunc)
	// TextureBitmap queries the texture bitmap of a texture.
	TextureBitmap(projectID string, textureID int, size string, onSuccess func(bmp *RawBitmap), onFailure FailureFunc)

	// Tiles queries the complete tile map of a level.
	Tiles(projectID string, archiveID string, levelID int, onSuccess func(tiles Tiles), onFailure FailureFunc)

	// Tile requests the properties of a specific tile.
	Tile(projectID string, archiveID string, levelID int, x, y int,
		onSuccess func(properties TileProperties), onFailure FailureFunc)
	// SetTile requests to update properties of a specific tile.
	SetTile(projectID string, archiveID string, levelID int, x, y int, properties TileProperties,
		onSuccess func(properties TileProperties), onFailure FailureFunc)

	// LevelObjects requests all objects of the level.
	LevelObjects(projectID string, archiveID string, levelID int,
		onSuccess func(objects *LevelObjects), onFailure FailureFunc)
}
