package editor

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/inkyblackness/shocked-model"
)

// RestDataStore is a REST based implementation of the DataStore interface.
type RestDataStore struct {
	transport RestTransport
}

// NewRestDataStore returns a new instance of a data store backed by a REST transport.
func NewRestDataStore(transport RestTransport) *RestDataStore {
	return &RestDataStore{transport: transport}
}

func (store *RestDataStore) get(url string, responseData interface{}, onSuccess func(), onFailure FailureFunc) {
	store.transport.Get(url, func(jsonString string) {
		json.Unmarshal(bytes.NewBufferString(jsonString).Bytes(), responseData)
		onSuccess()
	}, func() {
		onFailure()
	})
}

func (store *RestDataStore) put(url string, requestData interface{}, responseData interface{}, onSuccess func(), onFailure FailureFunc) {
	data, _ := json.Marshal(requestData)
	store.transport.Put(url, data, func(jsonString string) {
		json.Unmarshal(bytes.NewBufferString(jsonString).Bytes(), responseData)
		onSuccess()
	}, func() {
		onFailure()
	})
}

func (store *RestDataStore) post(url string, requestData interface{}, responseData interface{}, onSuccess func(), onFailure FailureFunc) {
	data, _ := json.Marshal(requestData)
	store.transport.Post(url, data, func(jsonString string) {
		json.Unmarshal(bytes.NewBufferString(jsonString).Bytes(), responseData)
		onSuccess()
	}, func() {
		onFailure()
	})
}

// NewProject implements the DataStore interface.
func (store *RestDataStore) NewProject(projectID string, onSuccess func(), onFailure FailureFunc) {
	url := fmt.Sprintf("/projects")
	var inData model.ProjectTemplate
	var outData model.Project

	inData.ID = projectID
	store.post(url, &inData, &outData, onSuccess, onFailure)
}

// Projects implements the DataStore interface.
func (store *RestDataStore) Projects(onSuccess func(projects []string), onFailure FailureFunc) {
	url := "/projects"
	var data model.Projects

	store.get(url, &data, func() {
		projectIDs := make([]string, len(data.Items))
		for index, item := range data.Items {
			projectIDs[index] = item.ID
		}
		onSuccess(projectIDs)
	}, onFailure)
}

// Font implements the DataStore interface.
func (store *RestDataStore) Font(projectID string, fontID int, onSuccess func(font *model.Font), onFailure FailureFunc) {
	url := fmt.Sprintf("/projects/%s/fonts/%v", projectID, fontID)
	var data model.Font

	store.get(url, &data, func() {
		onSuccess(&data)
	}, onFailure)
}

// GameObjectIcon implements the DataStore interface.
func (store *RestDataStore) GameObjectIcon(projectID string, class, subclass, objType int,
	onSuccess func(bmp *model.RawBitmap), onFailure FailureFunc) {
	url := fmt.Sprintf("/projects/%s/objects/%d/%d/%d/icon/raw", projectID, class, subclass, objType)
	var data model.RawBitmap

	store.get(url, &data, func() {
		onSuccess(&data)
	}, onFailure)
}

// Palette implements the DataStore interface.
func (store *RestDataStore) Palette(projectID string, paletteID string,
	onSuccess func(colors [256]model.Color), onFailure FailureFunc) {
	url := fmt.Sprintf("/projects/%s/palettes/%s", projectID, paletteID)
	var data model.Palette

	store.get(url, &data, func() {
		onSuccess(data.Colors)
	}, onFailure)
}

// Levels implements the DataStore interface.
func (store *RestDataStore) Levels(projectID string, archiveID string, onSuccess func(levels []model.Level), onFailure FailureFunc) {
	url := fmt.Sprintf("/projects/%s/%s/levels", projectID, archiveID)
	var data model.Levels

	store.get(url, &data, func() {
		onSuccess(data.List)
	}, onFailure)
}

// LevelTextures implements the DataStore interface.
func (store *RestDataStore) LevelTextures(projectID string, archiveID string, levelID int,
	onSuccess func(textureIDs []int), onFailure FailureFunc) {
	url := fmt.Sprintf("/projects/%s/%s/levels/%d/textures", projectID, archiveID, levelID)
	var data model.LevelTextures

	store.get(url, &data, func() {
		onSuccess(data.IDs)
	}, onFailure)
}

// SetLevelTextures implements the DataStore interface.
func (store *RestDataStore) SetLevelTextures(projectID string, archiveID string, levelID int, textureIDs []int,
	onSuccess func(textureIDs []int), onFailure FailureFunc) {
	url := fmt.Sprintf("/projects/%s/%s/levels/%d/textures", projectID, archiveID, levelID)
	var data model.LevelTextures

	store.put(url, textureIDs, &data, func() {
		onSuccess(data.IDs)
	}, onFailure)
}

// Textures implements the DataStore interface.
func (store *RestDataStore) Textures(projectID string, onSuccess func(textures []model.Texture), onFailure FailureFunc) {
	url := fmt.Sprintf("/projects/%s/textures", projectID)
	var data model.Textures

	store.get(url, &data, func() {
		onSuccess(data.List)
	}, onFailure)
}

// TextureBitmap implements the DataStore interface.
func (store *RestDataStore) TextureBitmap(projectID string, textureID int, size string,
	onSuccess func(bmp *model.RawBitmap), onFailure FailureFunc) {
	url := fmt.Sprintf("/projects/%s/textures/%d/%s/raw", projectID, textureID, size)
	var data model.RawBitmap

	store.get(url, &data, func() {
		onSuccess(&data)
	}, onFailure)
}

// Tiles implements the DataStore interface.
func (store *RestDataStore) Tiles(projectID string, archiveID string, levelID int,
	onSuccess func(tiles model.Tiles), onFailure FailureFunc) {
	url := fmt.Sprintf("/projects/%s/%s/levels/%d/tiles", projectID, archiveID, levelID)
	var data model.Tiles

	store.get(url, &data, func() {
		onSuccess(data)
	}, onFailure)
}

// Tile implements the DataStore interface.
func (store *RestDataStore) Tile(projectID string, archiveID string, levelID int, x, y int,
	onSuccess func(properties model.TileProperties), onFailure FailureFunc) {
	url := fmt.Sprintf("/projects/%s/%s/levels/%d/tiles/%d/%d", projectID, archiveID, levelID, y, x)
	var data model.Tile

	store.get(url, &data, func() {
		onSuccess(data.Properties)
	}, onFailure)
}

// SetTile implements the DataStore interface.
func (store *RestDataStore) SetTile(projectID string, archiveID string, levelID int, x, y int, properties model.TileProperties,
	onSuccess func(properties model.TileProperties), onFailure FailureFunc) {
	url := fmt.Sprintf("/projects/%s/%s/levels/%d/tiles/%d/%d", projectID, archiveID, levelID, y, x)
	var data model.Tile

	store.put(url, &properties, &data, func() {
		onSuccess(data.Properties)
	}, onFailure)
}

// LevelObjects implements the DataStore interface.
func (store *RestDataStore) LevelObjects(projectID string, archiveID string, levelID int,
	onSuccess func(objects *model.LevelObjects), onFailure FailureFunc) {
	url := fmt.Sprintf("/projects/%s/%s/levels/%d/objects", projectID, archiveID, levelID)
	var data model.LevelObjects

	store.get(url, &data, func() {
		onSuccess(&data)
	}, onFailure)
}
