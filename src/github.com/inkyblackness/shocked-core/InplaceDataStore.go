package core

import (
	"encoding/base64"
	"fmt"
	"image/color"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/shocked-core/release"
	"github.com/inkyblackness/shocked-model"
)

// InplaceDataStore implements the model.DataStore interface with an
// inplace project. This store has only one project, called "(inplace)",
// which overwrites the source data files.
type InplaceDataStore struct {
	workspace *Workspace
	inQueue   chan func()
	outQueue  chan<- func()
}

// NewInplaceDataStore returns a new instance of an inplace data store.
func NewInplaceDataStore(source release.Release, outQueue chan<- func()) *InplaceDataStore {
	projects := release.NewStaticReleaseContainer(map[string]release.Release{"(inplace)": source})
	inplace := &InplaceDataStore{
		workspace: NewWorkspace(source, projects),
		inQueue:   make(chan func(), 100),
		outQueue:  outQueue}
	go inplace.processor(inplace.inQueue)

	return inplace
}

func (inplace *InplaceDataStore) processor(queue chan func()) {
	active := true
	for active {
		task := <-queue
		if task != nil {
			task()
		} else {
			active = false
		}
	}
}

func (inplace *InplaceDataStore) in(task func()) {
	//inplace.inQueue <- func() {
	task()
	//}
}

func (inplace *InplaceDataStore) out(task func()) {
	//inplace.outQueue <- func() {
	task()
	//}
}

// Projects implements the model.DataStore interface
func (inplace *InplaceDataStore) Projects(onSuccess func(projects []string), onFailure model.FailureFunc) {
	inplace.in(func() {
		projects := inplace.workspace.ProjectNames()
		inplace.out(func() {
			onSuccess(projects)
		})
	})
}

// NewProject implements the model.DataStore interface
func (inplace *InplaceDataStore) NewProject(projectID string,
	onSuccess func(), onFailure model.FailureFunc) {
	inplace.in(func() {
		inplace.out(onFailure)
	})
}

// Font implements the model.DataStore interface
func (inplace *InplaceDataStore) Font(projectID string, fontID int,
	onSuccess func(font *model.Font), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			var font *model.Font
			font, err = project.Fonts().Font(res.ResourceID(fontID))
			if err == nil {
				inplace.out(func() { onSuccess(font) })
			}
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// GameObjects implements the model.DataStore interface
func (inplace *InplaceDataStore) GameObjects(projectID string,
	onSuccess func(objects []model.GameObject), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			gameObjects := project.GameObjects()
			objects := gameObjects.Objects()

			inplace.out(func() { onSuccess(objects) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// GameObjectIcon implements the model.DataStore interface
func (inplace *InplaceDataStore) GameObjectIcon(projectID string, class, subclass, objType int,
	onSuccess func(bmp *model.RawBitmap), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			objID := res.MakeObjectID(res.ObjectClass(class), res.ObjectSubclass(subclass), res.ObjectType(objType))
			bmp := project.GameObjects().Icon(objID)
			var entity model.RawBitmap

			entity.Width = int(bmp.ImageWidth())
			entity.Height = int(bmp.ImageHeight())
			var pixel []byte

			for row := 0; row < entity.Height; row++ {
				pixel = append(pixel, bmp.Row(row)...)
			}
			entity.Pixels = base64.StdEncoding.EncodeToString(pixel)

			inplace.out(func() { onSuccess(&entity) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// Palette implements the model.DataStore interface
func (inplace *InplaceDataStore) Palette(projectID string, paletteID string,
	onSuccess func(colors [256]model.Color), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			if paletteID == "game" {
				var palette color.Palette
				palette, err = project.Palettes().GamePalette()
				if err == nil {
					var entity model.Palette

					inplace.encodePalette(&entity.Colors, palette)
					inplace.out(func() { onSuccess(entity.Colors) })
				}
			} else {
				err = fmt.Errorf("Wrong palette")
			}
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

func (inplace *InplaceDataStore) encodePalette(out *[256]model.Color, palette color.Palette) {
	for index, inColor := range palette {
		outColor := &out[index]
		r, g, b, _ := inColor.RGBA()

		outColor.Red = int(r >> 8)
		outColor.Green = int(g >> 8)
		outColor.Blue = int(b >> 8)
	}
}

// Levels implements the model.DataStore interface
func (inplace *InplaceDataStore) Levels(projectID string, archiveID string,
	onSuccess func(levels []model.Level), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			var entity model.Levels
			archive := project.Archive()
			levelIDs := archive.LevelIDs()

			for _, id := range levelIDs {
				entry := inplace.getLevelEntity(project, archive, id)

				entity.List = append(entity.List, entry)
			}

			inplace.out(func() { onSuccess(entity.List) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

func (inplace *InplaceDataStore) getLevelEntity(project *Project, archive *Archive, levelID int) (entity model.Level) {
	entity.ID = fmt.Sprintf("%d", levelID)
	entity.Href = "/projects/" + project.Name() + "/archive/levels/" + entity.ID
	level := archive.Level(levelID)
	entity.Properties = level.Properties()

	entity.Links = []model.Link{}
	entity.Links = append(entity.Links, model.Link{Rel: "tiles", Href: entity.Href + "/tiles/{y}/{x}"})
	if !entity.Properties.CyberspaceFlag {
		entity.Links = append(entity.Links, model.Link{Rel: "textures", Href: entity.Href + "/textures"})
	}

	return
}

// LevelTextures implements the model.DataStore interface
func (inplace *InplaceDataStore) LevelTextures(projectID string, archiveID string, levelID int,
	onSuccess func(textureIDs []int), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(levelID)
			textureIDs := level.Textures()

			inplace.out(func() { onSuccess(textureIDs) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// SetLevelTextures implements the model.DataStore interface
func (inplace *InplaceDataStore) SetLevelTextures(projectID string, archiveID string, levelID int, textureIDs []int,
	onSuccess func(textureIDs []int), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(levelID)
			level.SetTextures(textureIDs)
			result := level.Textures()

			inplace.out(func() { onSuccess(result) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// Textures implements the model.DataStore interface
func (inplace *InplaceDataStore) Textures(projectID string,
	onSuccess func(textures []model.Texture), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			textures := project.Textures()
			limit := textures.TextureCount()
			var entity model.Textures

			entity.List = make([]model.Texture, limit)
			for id := 0; id < limit; id++ {
				entity.List[id] = inplace.textureEntity(project, id)
			}

			inplace.out(func() { onSuccess(entity.List) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

func (inplace *InplaceDataStore) textureEntity(project *Project, textureID int) (entity model.Texture) {
	entity.ID = fmt.Sprintf("%d", textureID)
	entity.Href = "/projects/" + project.Name() + "/textures/" + entity.ID
	entity.Properties = project.Textures().Properties(textureID)
	for _, size := range model.TextureSizes() {
		entity.Images = append(entity.Images, model.Link{Rel: string(size), Href: entity.Href + "/" + string(size)})
	}

	return
}

// TextureBitmap implements the model.DataStore interface
func (inplace *InplaceDataStore) TextureBitmap(projectID string, textureID int, size string,
	onSuccess func(bmp *model.RawBitmap), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			bmp := project.Textures().Image(textureID, model.TextureSize(size))
			var entity model.RawBitmap

			entity.Width = int(bmp.ImageWidth())
			entity.Height = int(bmp.ImageHeight())
			var pixel []byte

			for row := 0; row < entity.Height; row++ {
				pixel = append(pixel, bmp.Row(row)...)
			}
			entity.Pixels = base64.StdEncoding.EncodeToString(pixel)

			inplace.out(func() { onSuccess(&entity) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// Tiles implements the model.DataStore interface
func (inplace *InplaceDataStore) Tiles(projectID string, archiveID string, levelID int,
	onSuccess func(tiles model.Tiles), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(levelID)
			var entity model.Tiles

			entity.Table = make([][]model.Tile, 64)
			for y := 0; y < 64; y++ {
				entity.Table[y] = make([]model.Tile, 64)
				for x := 0; x < 64; x++ {
					entity.Table[y][x] = inplace.getLevelTileEntity(project, level, x, y)
				}
			}

			inplace.out(func() { onSuccess(entity) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

func (inplace *InplaceDataStore) getLevelTileEntity(project *Project, level *Level, x int, y int) (entity model.Tile) {
	entity.Href = "/projects/" + project.Name() + "/archive/levels/" + fmt.Sprintf("%d", level.ID()) +
		"/tiles/" + fmt.Sprintf("%d", y) + "/" + fmt.Sprintf("%d", x)
	entity.Properties = level.TileProperties(int(x), int(y))

	return
}

// Tile implements the model.DataStore interface
func (inplace *InplaceDataStore) Tile(projectID string, archiveID string, levelID int, x, y int,
	onSuccess func(properties model.TileProperties), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(int(levelID))
			entity := inplace.getLevelTileEntity(project, level, x, y)

			inplace.out(func() { onSuccess(entity.Properties) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// SetTile implements the model.DataStore interface
func (inplace *InplaceDataStore) SetTile(projectID string, archiveID string, levelID int, x, y int, properties model.TileProperties,
	onSuccess func(properties model.TileProperties), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(levelID)

			level.SetTileProperties(x, y, properties)
			result := inplace.getLevelTileEntity(project, level, x, y)
			inplace.out(func() { onSuccess(result.Properties) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// LevelObjects implements the model.DataStore interface
func (inplace *InplaceDataStore) LevelObjects(projectID string, archiveID string, levelID int,
	onSuccess func(objects *model.LevelObjects), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(levelID)
			var entity model.LevelObjects

			entity.Table = level.Objects()

			inplace.out(func() { onSuccess(&entity) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// AddLevelObject implements the model.DataStore interface
func (inplace *InplaceDataStore) AddLevelObject(projectID string, archiveID string, levelID int, template model.LevelObjectTemplate,
	onSuccess func(object model.LevelObject), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(levelID)
			var entity model.LevelObject

			entity, err = level.AddObject(&template)
			if err == nil {
				inplace.out(func() { onSuccess(entity) })
			}
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// RemoveLevelObject implements the model.DataStore interface
func (inplace *InplaceDataStore) RemoveLevelObject(projectID string, archiveID string, levelID int, objectID int,
	onSuccess func(), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(levelID)

			err = level.RemoveObject(objectID)
			if err == nil {
				inplace.out(func() { onSuccess() })
			}
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// SetLevelObject implements the model.DataStore interface
func (inplace *InplaceDataStore) SetLevelObject(projectID string, archiveID string, levelID int, objectID int,
	properties *model.LevelObjectProperties, onSuccess func(properties *model.LevelObjectProperties), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(levelID)
			var newProperties model.LevelObjectProperties

			newProperties, err = level.SetObject(objectID, properties)
			if err == nil {
				inplace.out(func() { onSuccess(&newProperties) })
			}
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}
