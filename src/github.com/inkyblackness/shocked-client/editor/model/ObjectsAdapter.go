package model

import (
	"fmt"
	"sort"

	"github.com/inkyblackness/shocked-model"
)

// ObjectsAdapter is the adapter for all game objects
type ObjectsAdapter struct {
	context projectContext
	store   model.DataStore

	objects *observable
	icons   *Bitmaps

	bitmaps               *Bitmaps
	bitmapRequestsPending map[ObjectBitmapID]bool
}

func newObjectsAdapter(context projectContext, store model.DataStore) *ObjectsAdapter {
	adapter := &ObjectsAdapter{
		context: context,
		store:   store,

		objects: newObservable(),
		icons:   newBitmaps(),

		bitmaps:               newBitmaps(),
		bitmapRequestsPending: make(map[ObjectBitmapID]bool)}

	objectMap := make(map[ObjectID]*GameObject)
	adapter.objects.set(&objectMap)

	return adapter
}

func (adapter *ObjectsAdapter) clear() {
	objectMap := make(map[ObjectID]*GameObject)
	adapter.objects.set(&objectMap)
	adapter.icons.clear()
	adapter.bitmaps.clear()
	adapter.bitmapRequestsPending = make(map[ObjectBitmapID]bool)
}

func (adapter *ObjectsAdapter) refresh() {
	adapter.store.GameObjects(adapter.context.ActiveProjectID(),
		adapter.onNewGameObjects,
		adapter.context.simpleStoreFailure("GameObjects"))
}

func (adapter *ObjectsAdapter) objectMap() map[ObjectID]*GameObject {
	return *adapter.objects.get().(*map[ObjectID]*GameObject)
}

// Object returns the object with given ID.
func (adapter *ObjectsAdapter) Object(id ObjectID) *GameObject {
	return adapter.objectMap()[id]
}

// ObjectIDs returns a set of object identifier.
func (adapter *ObjectsAdapter) ObjectIDs() []ObjectID {
	objectMap := adapter.objectMap()
	result := make([]ObjectID, 0, len(objectMap))

	for key := range objectMap {
		result = append(result, key)
	}

	return result
}

// OnObjectsChanged registers a callback for updates.
func (adapter *ObjectsAdapter) OnObjectsChanged(callback func()) {
	adapter.objects.addObserver(callback)
}

func (adapter *ObjectsAdapter) onNewGameObjects(objects []model.GameObject) {
	objectMap := adapter.objectMap()

	for _, rawObject := range objects {
		objID := MakeObjectID(rawObject.Class, rawObject.Subclass, rawObject.Type)
		editorObject := &GameObject{
			id: objID}

		for i := 0; i < model.LanguageCount; i++ {
			editorObject.shortName[i] = *rawObject.Properties.ShortName[i]
			editorObject.longName[i] = *rawObject.Properties.LongName[i]
		}
		editorObject.data = rawObject.Properties.Data
		objectMap[objID] = editorObject
	}
	adapter.objects.notifyObservers()
}

// Icons returns a bitmap container for the objects. The key is based on the ObjectIDs.
func (adapter *ObjectsAdapter) Icons() *Bitmaps {
	return adapter.icons
}

// RequestIcon tries to load the icon for given object ID.
func (adapter *ObjectsAdapter) RequestIcon(id ObjectID) {
	adapter.store.GameObjectIcon(adapter.context.ActiveProjectID(),
		id.Class(), id.Subclass(), id.Type(),
		func(bmp *model.RawBitmap) {
			adapter.icons.setRawBitmap(id.ToInt(), bmp)
		},
		adapter.context.simpleStoreFailure(fmt.Sprintf("GameObjectIcon[%v]", id)))
}

// Bitmap returns the raw bitmap data for the given key (if loaded)
func (adapter *ObjectsAdapter) Bitmap(id ObjectBitmapID) *model.RawBitmap {
	return adapter.bitmaps.RawBitmap(id.ToInt())
}

// RequestBitmapChange will update the bitmap data for identified object.
func (adapter *ObjectsAdapter) RequestBitmapChange(id ObjectBitmapID, newBitmap *model.RawBitmap) {
	adapter.store.SetGameObjectBitmap(adapter.context.ActiveProjectID(),
		id.ObjectID.Class(), id.ObjectID.Subclass(), id.ObjectID.Type(), id.Index, newBitmap,
		func() {
			adapter.bitmaps.setRawBitmap(id.ToInt(), newBitmap)
		},
		func() {
			adapter.context.simpleStoreFailure(fmt.Sprintf("SetGameObjectBitmap[%v]", id))()
		})
}

// RequestBitmap will load the bitmap data for identified key.
func (adapter *ObjectsAdapter) RequestBitmap(id ObjectBitmapID) {
	if !adapter.bitmapRequestsPending[id] {
		adapter.bitmapRequestsPending[id] = true
		adapter.store.GameObjectBitmap(adapter.context.ActiveProjectID(),
			id.ObjectID.Class(), id.ObjectID.Subclass(), id.ObjectID.Type(), id.Index,
			func(bmp *model.RawBitmap) {
				adapter.bitmapRequestsPending[id] = false
				adapter.bitmaps.setRawBitmap(id.ToInt(), bmp)
			},
			func() {
				adapter.bitmapRequestsPending[id] = false
				adapter.context.simpleStoreFailure(fmt.Sprintf("GameObjectBitmap[%v]", id))()
			})
	}
}

// Bitmaps returns the container of bitmaps.
func (adapter *ObjectsAdapter) Bitmaps() *Bitmaps {
	return adapter.bitmaps
}

// ObjectsOfClass returns all objects matching the requested class.
func (adapter *ObjectsAdapter) ObjectsOfClass(class int) (objects []*GameObject) {
	return adapter.objectsByFilter(func(object *GameObject) bool { return object.id.Class() == class })
}

// Objects returns all objects
func (adapter *ObjectsAdapter) Objects() (objects []*GameObject) {
	return adapter.objectsByFilter(func(object *GameObject) bool { return true })
}

func (adapter *ObjectsAdapter) objectsByFilter(predicate func(*GameObject) bool) (objects []*GameObject) {
	for _, object := range adapter.objectMap() {
		if predicate(object) {
			objects = append(objects, object)
		}
	}
	sort.Slice(objects, func(a int, b int) bool { return objects[a].id.ToInt() < objects[b].id.ToInt() })

	return
}

// RequestObjectPropertiesChange requests to modify the properties of identifed object.
func (adapter *ObjectsAdapter) RequestObjectPropertiesChange(objectID ObjectID, properties *model.GameObjectProperties) {
	objectMap := adapter.objectMap()

	adapter.store.SetGameObject(adapter.context.ActiveProjectID(),
		objectID.Class(), objectID.Subclass(), objectID.Type(), properties,
		func(newProperties *model.GameObjectProperties) {
			objectMap[objectID].data = newProperties.Data
			adapter.objects.notifyObservers()
		},
		adapter.context.simpleStoreFailure(fmt.Sprintf("SetGameObject %v", objectID)))
}
