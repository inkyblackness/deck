package model

import (
	"fmt"

	"github.com/inkyblackness/shocked-model"
)

// Adapter is the central model adapter.
type Adapter struct {
	store model.DataStore

	message *observable

	activeProjectID     *observable
	availableArchiveIDs *observable
	activeArchiveID     *observable
	activeLevel         *LevelAdapter

	availableLevels   map[string]model.LevelProperties
	availableLevelIDs *observable

	palette        *observable
	textureAdapter *TextureAdapter
	objectsAdapter *ObjectsAdapter
}

// NewAdapter returns a new model adapter.
func NewAdapter(store model.DataStore) *Adapter {
	adapter := &Adapter{
		store:   store,
		message: newObservable(),

		activeProjectID:     newObservable(),
		availableArchiveIDs: newObservable(),
		activeArchiveID:     newObservable(),

		availableLevels:   make(map[string]model.LevelProperties),
		availableLevelIDs: newObservable(),

		palette: newObservable()}

	adapter.message.set("")
	adapter.activeLevel = newLevelAdapter(adapter, store)
	adapter.textureAdapter = newTextureAdapter(adapter, store)
	adapter.objectsAdapter = newObjectsAdapter(adapter, store)
	adapter.palette.set(&[256]model.Color{})

	return adapter
}

func (adapter *Adapter) simpleStoreFailure(info string) model.FailureFunc {
	return func() {
		adapter.SetMessage(fmt.Sprintf("Failed to process store query <%s>", info))
	}
}

// SetMessage sets the current global message.
func (adapter *Adapter) SetMessage(message string) {
	adapter.message.set(message)
}

// Message returns the current global message.
func (adapter *Adapter) Message() string {
	return adapter.message.orDefault("").(string)
}

// OnMessageChanged registers a callback for the global message.
func (adapter *Adapter) OnMessageChanged(callback func()) {
	adapter.message.addObserver(callback)
}

// ActiveProjectID returns the identifier of the current project.
func (adapter *Adapter) ActiveProjectID() string {
	return adapter.activeProjectID.orDefault("").(string)
}

// RequestProject sets the project to work on.
func (adapter *Adapter) RequestProject(projectID string) {
	adapter.textureAdapter.clear()
	adapter.objectsAdapter.clear()
	adapter.requestArchive("")
	adapter.availableArchiveIDs.set("")

	adapter.activeProjectID.set(projectID)
	if projectID != "" {
		adapter.availableArchiveIDs.set([]string{"archive"})
		adapter.requestArchive("archive")
		adapter.store.Palette(adapter.ActiveProjectID(), "game",
			adapter.onGamePalette, adapter.simpleStoreFailure("Palette"))
		adapter.objectsAdapter.refresh()
	}
}

func (adapter *Adapter) onGamePalette(colors [256]model.Color) {
	adapter.palette.set(&colors)
}

// GamePalette returns the main palette.
func (adapter *Adapter) GamePalette() *[256]model.Color {
	return adapter.palette.get().(*[256]model.Color)
}

// OnGamePaletteChanged registers a callback for updates.
func (adapter *Adapter) OnGamePaletteChanged(callback func()) {
	adapter.palette.addObserver(callback)
}

// TextureAdapter returns the adapter for textures.
func (adapter *Adapter) TextureAdapter() *TextureAdapter {
	return adapter.textureAdapter
}

// ObjectsAdapter returns the adapter for game objects.
func (adapter *Adapter) ObjectsAdapter() *ObjectsAdapter {
	return adapter.objectsAdapter
}

// ActiveArchiveID returns the identifier of the current archive.
func (adapter *Adapter) ActiveArchiveID() string {
	return adapter.activeArchiveID.orDefault("").(string)
}

func (adapter *Adapter) requestArchive(archiveID string) {
	adapter.RequestActiveLevel("")
	adapter.availableLevels = make(map[string]model.LevelProperties)
	adapter.availableLevelIDs.set(nil)

	adapter.activeArchiveID.set(archiveID)
	if archiveID != "" {
		adapter.store.Levels(adapter.ActiveProjectID(), adapter.ActiveArchiveID(),
			adapter.onLevels,
			adapter.simpleStoreFailure("Levels"))
	}
}

func (adapter *Adapter) onLevels(levels []model.Level) {
	availableLevelIDs := make([]string, len(levels))

	adapter.availableLevels = make(map[string]model.LevelProperties)
	for index, entry := range levels {
		availableLevelIDs[index] = entry.ID
		adapter.availableLevels[entry.ID] = entry.Properties
	}
	adapter.availableLevelIDs.set(availableLevelIDs)
}

// ActiveLevel returns the adapter for the currently active level.
func (adapter *Adapter) ActiveLevel() *LevelAdapter {
	return adapter.activeLevel
}

// RequestActiveLevel requests to set the specified level as the active one.
func (adapter *Adapter) RequestActiveLevel(levelID string) {
	levelProp, existing := adapter.availableLevels[levelID]
	adapter.activeLevel.isCyberspace = existing && levelProp.CyberspaceFlag
	adapter.activeLevel.requestByID(levelID)
}

// AvailableLevelIDs returns the list of identifier of available levels.
func (adapter *Adapter) AvailableLevelIDs() []string {
	return adapter.availableLevelIDs.get().([]string)
}

// OnAvailableLevelsChanged registers a callback for changes of available levels.
func (adapter *Adapter) OnAvailableLevelsChanged(callback func()) {
	adapter.availableLevelIDs.addObserver(callback)
}
