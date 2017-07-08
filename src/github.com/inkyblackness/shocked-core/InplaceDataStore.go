package core

import (
	"encoding/base64"
	"fmt"
	"image/color"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/audio"
	"github.com/inkyblackness/res/image"
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

// Bitmap implements the model.DataStore interface
func (inplace *InplaceDataStore) Bitmap(projectID string, key model.ResourceKey,
	onSuccess func(model.ResourceKey, *model.RawBitmap), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			var imgBitmap image.Bitmap
			imgBitmap, err = project.Bitmaps().Image(key)

			if err == nil {
				rawBitmap := inplace.toRawBitmap(imgBitmap)
				inplace.out(func() { onSuccess(key, &rawBitmap) })
			}
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// SetBitmap implements the model.DataStore interface
func (inplace *InplaceDataStore) SetBitmap(projectID string, key model.ResourceKey, rawBitmap *model.RawBitmap,
	onSuccess func(model.ResourceKey, *model.RawBitmap), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			bitmaps := project.Bitmaps()
			imgBitmap := inplace.fromRawBitmap(rawBitmap)
			var resultKey model.ResourceKey

			resultKey, err = bitmaps.SetImage(key, imgBitmap)

			if err == nil {
				var imgResult image.Bitmap
				imgResult, err = project.Bitmaps().Image(resultKey)

				if err == nil {
					rawResult := inplace.toRawBitmap(imgResult)
					inplace.out(func() { onSuccess(resultKey, &rawResult) })
				}
			}
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// Text implements the model.DataStore interface
func (inplace *InplaceDataStore) Text(projectID string, key model.ResourceKey,
	onSuccess func(model.ResourceKey, string), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			var text string
			text, err = project.Texts().Text(key)

			if err == nil {
				inplace.out(func() { onSuccess(key, text) })
			}
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// SetText implements the model.DataStore interface
func (inplace *InplaceDataStore) SetText(projectID string, key model.ResourceKey, text string,
	onSuccess func(model.ResourceKey, string), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			texts := project.Texts()
			var resultKey model.ResourceKey

			resultKey, err = texts.SetText(key, text)

			if err == nil {
				inplace.out(func() { onSuccess(resultKey, text) })
			}
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// Audio implements the model.DataStore interface
func (inplace *InplaceDataStore) Audio(projectID string, key model.ResourceKey,
	onSuccess func(model.ResourceKey, audio.SoundData), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			var data audio.SoundData
			data, err = project.Sounds().Audio(key)

			if err == nil {
				inplace.out(func() { onSuccess(key, data) })
			}
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// SetAudio implements the model.DataStore interface
func (inplace *InplaceDataStore) SetAudio(projectID string, key model.ResourceKey, data audio.SoundData,
	onSuccess func(model.ResourceKey), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			sounds := project.Sounds()
			var resultKey model.ResourceKey

			resultKey, err = sounds.SetAudio(key, data)

			if err == nil {
				inplace.out(func() { onSuccess(resultKey) })
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

// SetGameObject implements the model.DataStore interface
func (inplace *InplaceDataStore) SetGameObject(projectID string, class, subclass, objType int, properties *model.GameObjectProperties,
	onSuccess func(properties *model.GameObjectProperties), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			gameObjects := project.GameObjects()
			objID := res.MakeObjectID(res.ObjectClass(class), res.ObjectSubclass(subclass), res.ObjectType(objType))
			var entity model.GameObjectProperties

			entity.Data = gameObjects.SetObjectData(objID, properties.Data)
			inplace.out(func() { onSuccess(&entity) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// ElectronicMessage implements the model.DataStore interface.
func (inplace *InplaceDataStore) ElectronicMessage(projectID string, messageType model.ElectronicMessageType, id int,
	onSuccess func(message model.ElectronicMessage), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			eMessages := project.ElectronicMessages()
			var message model.ElectronicMessage
			message, err = eMessages.Message(messageType, id)

			if err == nil {
				inplace.out(func() { onSuccess(message) })
			}
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// SetElectronicMessage implements the model.DataStore interface.
func (inplace *InplaceDataStore) SetElectronicMessage(projectID string, messageType model.ElectronicMessageType,
	id int, message model.ElectronicMessage,
	onSuccess func(message model.ElectronicMessage), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			eMessages := project.ElectronicMessages()
			eMessages.SetMessage(messageType, id, message)
			var result model.ElectronicMessage
			result, err = eMessages.Message(messageType, id)

			if err == nil {
				inplace.out(func() { onSuccess(result) })
			}
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// ElectronicMessageAudio implements the model.DataStore interface.
func (inplace *InplaceDataStore) ElectronicMessageAudio(projectID string,
	messageType model.ElectronicMessageType, id int, language model.ResourceLanguage,
	onSuccess func(data audio.SoundData), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			eMessages := project.ElectronicMessages()
			var data audio.SoundData
			data, err = eMessages.MessageAudio(messageType, id, language)

			if err == nil {
				inplace.out(func() { onSuccess(data) })
			}
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// SetElectronicMessageAudio implements the model.DataStore interface.
func (inplace *InplaceDataStore) SetElectronicMessageAudio(projectID string,
	messageType model.ElectronicMessageType, id int, language model.ResourceLanguage, data audio.SoundData,
	onSuccess func(), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			eMessages := project.ElectronicMessages()
			err = eMessages.SetMessageAudio(messageType, id, language, data)

			if err == nil {
				inplace.out(func() { onSuccess() })
			}
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
			archive := project.Archive()
			levelIDs := archive.LevelIDs()
			result := []model.Level{}

			for _, levelID := range levelIDs {
				var entry model.Level
				entry.ID = levelID
				result = append(result, entry)
			}

			inplace.out(func() { onSuccess(result) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// LevelProperties implements the model.DataStore interface.
func (inplace *InplaceDataStore) LevelProperties(projectID string, archiveID string, levelID int,
	onSuccess func(properties model.LevelProperties), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(levelID)
			properties := level.Properties()

			inplace.out(func() { onSuccess(properties) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// SetLevelProperties implements the model.DataStore interface.
func (inplace *InplaceDataStore) SetLevelProperties(projectID string, archiveID string, levelID int, properties model.LevelProperties,
	onSuccess func(properties model.LevelProperties), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(levelID)

			level.SetProperties(properties)
			result := level.Properties()

			inplace.out(func() { onSuccess(result) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
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

// LevelTextureAnimations implements the model.DataStore interface
func (inplace *InplaceDataStore) LevelTextureAnimations(projectID string, archiveID string, levelID int,
	onSuccess func(animations []model.TextureAnimation), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(levelID)
			animations := level.TextureAnimations()

			inplace.out(func() { onSuccess(animations) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// SetLevelTextureAnimation implements the model.DataStore interface
func (inplace *InplaceDataStore) SetLevelTextureAnimation(projectID string, archiveID string, levelID int,
	animationGroup int, properties model.TextureAnimation,
	onSuccess func(animations []model.TextureAnimation), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(levelID)
			level.SetTextureAnimation(animationGroup, properties)
			animations := level.TextureAnimations()

			inplace.out(func() { onSuccess(animations) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// Textures implements the model.DataStore interface
func (inplace *InplaceDataStore) Textures(projectID string,
	onSuccess func(textures []model.TextureProperties), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			textures := project.Textures()
			limit := textures.TextureCount()
			result := make([]model.TextureProperties, limit)

			for id := 0; id < limit; id++ {
				result[id] = textures.Properties(id)
			}

			inplace.out(func() { onSuccess(result) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// SetTextureProperties implements the model.DataStore interface.
func (inplace *InplaceDataStore) SetTextureProperties(projectID string, textureID int, newProperties *model.TextureProperties,
	onSuccess func(properties *model.TextureProperties), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			textures := project.Textures()

			textures.SetProperties(textureID, *newProperties)
			result := textures.Properties(textureID)

			inplace.out(func() { onSuccess(&result) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

func (inplace *InplaceDataStore) fromRawBitmap(rawBitmap *model.RawBitmap) image.Bitmap {
	var header image.BitmapHeader
	data, _ := base64.StdEncoding.DecodeString(rawBitmap.Pixels)

	header.Height = uint16(rawBitmap.Height)
	header.Width = uint16(rawBitmap.Width)
	header.Stride = header.Width
	header.Type = image.CompressedBitmap

	return image.NewMemoryBitmap(&header, data, nil)
}

func (inplace *InplaceDataStore) toRawBitmap(imgBitmap image.Bitmap) model.RawBitmap {
	var rawBitmap model.RawBitmap

	rawBitmap.Width = int(imgBitmap.ImageWidth())
	rawBitmap.Height = int(imgBitmap.ImageHeight())
	var pixel []byte

	for row := 0; row < rawBitmap.Height; row++ {
		pixel = append(pixel, imgBitmap.Row(row)...)
	}
	rawBitmap.Pixels = base64.StdEncoding.EncodeToString(pixel)

	return rawBitmap
}

// TextureBitmap implements the model.DataStore interface
func (inplace *InplaceDataStore) TextureBitmap(projectID string, textureID int, size string,
	onSuccess func(bmp *model.RawBitmap), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			bmp := project.Textures().Image(textureID, model.TextureSize(size))
			entity := inplace.toRawBitmap(bmp)

			inplace.out(func() { onSuccess(&entity) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// SetTextureBitmap implements the model.DataStore interface
func (inplace *InplaceDataStore) SetTextureBitmap(projectID string, textureID int, size string, rawBitmap *model.RawBitmap,
	onSuccess func(*model.RawBitmap), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			textures := project.Textures()
			imgBitmap := inplace.fromRawBitmap(rawBitmap)
			textures.SetImage(textureID, model.TextureSize(size), imgBitmap)
			imgResult := textures.Image(textureID, model.TextureSize(size))
			rawResult := inplace.toRawBitmap(imgResult)

			inplace.out(func() { onSuccess(&rawResult) })
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

			entity.Table = make([][]model.TileProperties, 64)
			for y := 0; y < 64; y++ {
				entity.Table[y] = make([]model.TileProperties, 64)
				for x := 0; x < 64; x++ {
					entity.Table[y][x] = level.TileProperties(x, y)
				}
			}

			inplace.out(func() { onSuccess(entity) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// Tile implements the model.DataStore interface
func (inplace *InplaceDataStore) Tile(projectID string, archiveID string, levelID int, x, y int,
	onSuccess func(properties model.TileProperties), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(int(levelID))
			properties := level.TileProperties(x, y)

			inplace.out(func() { onSuccess(properties) })
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
			result := level.TileProperties(x, y)
			inplace.out(func() { onSuccess(result) })
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

// SetLevelObject implements the model.DataStore interface.
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

// LevelSurveillanceObjects implements the model.DataStore interface.
func (inplace *InplaceDataStore) LevelSurveillanceObjects(projectID string, archiveID string, levelID int,
	onSuccess func(objects []model.SurveillanceObject), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(levelID)
			objects := level.LevelSurveillanceObjects()

			inplace.out(func() { onSuccess(objects) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}

// SetLevelSurveillanceObject implements the model.DataStore interface.
func (inplace *InplaceDataStore) SetLevelSurveillanceObject(projectID string, archiveID string, levelID int,
	surveillanceIndex int, data model.SurveillanceObject,
	onSuccess func(objects []model.SurveillanceObject), onFailure model.FailureFunc) {
	inplace.in(func() {
		project, err := inplace.workspace.Project(projectID)

		if err == nil {
			level := project.Archive().Level(levelID)

			level.SetLevelSurveillanceObject(surveillanceIndex, data)
			objects := level.LevelSurveillanceObjects()

			inplace.out(func() { onSuccess(objects) })
		}
		if err != nil {
			inplace.out(onFailure)
		}
	})
}
