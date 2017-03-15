package model

import (
	"fmt"

	"github.com/inkyblackness/shocked-model"
)

// ObjectsAdapter is the adapter for all game objects
type ObjectsAdapter struct {
	context projectContext
	store   model.DataStore

	objects *observable
	icons   *Bitmaps
}

func newObjectsAdapter(context projectContext, store model.DataStore) *ObjectsAdapter {
	adapter := &ObjectsAdapter{
		context: context,
		store:   store,

		objects: newObservable(),
		icons:   newBitmaps()}

	objectMap := make(map[ObjectID]*GameObject)
	adapter.objects.set(&objectMap)

	return adapter
}

func (adapter *ObjectsAdapter) clear() {
	objectMap := make(map[ObjectID]*GameObject)
	adapter.objects.set(&objectMap)
	adapter.icons.clear()
}

func (adapter *ObjectsAdapter) refresh() {
	// TODO: query all object properties (names)
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
