package io

import (
	"bytes"
	"io/ioutil"
	"log"
	"time"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/chunk/resfile"
	"github.com/inkyblackness/res/objprop"
	dosObjprop "github.com/inkyblackness/res/objprop/dos"
	storeObjprop "github.com/inkyblackness/res/objprop/store"
	"github.com/inkyblackness/res/serial"
	"github.com/inkyblackness/res/textprop"
	dosTextprop "github.com/inkyblackness/res/textprop/dos"
	storeTextprop "github.com/inkyblackness/res/textprop/store"
	"github.com/inkyblackness/shocked-core/release"
)

// ReleaseStoreLibrary is a container with two releases: one source and one sink.
// Stores can be retrieved from this library, which access the source release for
// reading properties and the sink release for writing modified data.
type ReleaseStoreLibrary struct {
	source      release.Release
	sink        release.Release
	timeoutMSec int
	chunkStores map[string]*DynamicChunkStore

	descriptors   []objprop.ClassDescriptor
	objpropStores map[string]objprop.Store

	textpropStores map[string]textprop.Store

	saveChannels map[string]chan interface{}
}

// NewReleaseStoreLibrary returns a StoreLibrary that covers two Release container.
// Stores are first tried from the sink, then from source; They are saved in the sink.
func NewReleaseStoreLibrary(source release.Release, sink release.Release, timeoutMSec int) StoreLibrary {
	library := &ReleaseStoreLibrary{
		source:      source,
		sink:        sink,
		timeoutMSec: timeoutMSec,
		chunkStores: make(map[string]*DynamicChunkStore),

		descriptors:   objprop.StandardProperties(),
		objpropStores: make(map[string]objprop.Store),

		textpropStores: make(map[string]textprop.Store),

		saveChannels: make(map[string]chan interface{})}

	return library
}

// SaveAll requests all stores to save their current state to disk.
// This operation is performed asynchronously.
func (library *ReleaseStoreLibrary) SaveAll() {
	for _, save := range library.saveChannels {
		save <- true
	}
}

// ChunkStore implements the StoreLibrary interface.
func (library *ReleaseStoreLibrary) ChunkStore(name string) (chunkStore *DynamicChunkStore, err error) {
	chunkStore, exists := library.chunkStores[name]

	if !exists {
		if library.sink.HasResource(name) {
			chunkStore, err = library.openChunkStoreFrom(library.sink, name)
		} else if library.source.HasResource(name) {
			chunkStore, err = library.openChunkStoreFrom(library.source, name)
		} else {
			chunkStore = library.createSavingChunkStore(chunk.NullProvider(), "", name)
		}
		if err == nil {
			library.chunkStores[name] = chunkStore
		}
	}

	return
}

// ObjpropStore implements the StoreLibrary interface.
func (library *ReleaseStoreLibrary) ObjpropStore(name string) (objpropStore objprop.Store, err error) {
	objpropStore, exists := library.objpropStores[name]

	if !exists {
		if library.sink.HasResource(name) {
			objpropStore, err = library.openObjpropStoreFrom(library.sink, name)
		} else if library.source.HasResource(name) {
			objpropStore, err = library.openObjpropStoreFrom(library.source, name)
		} else {
			objpropStore = library.createSavingObjpropStore(objprop.NullProvider(library.descriptors), "", name, func() {})
		}
		if err == nil {
			library.objpropStores[name] = objpropStore
		}
	}

	return
}

// TextpropStore implements the StoreLibrary interface.
func (library *ReleaseStoreLibrary) TextpropStore(name string) (textpropStore textprop.Store, err error) {
	textpropStore, exists := library.textpropStores[name]

	if !exists {
		if library.sink.HasResource(name) {
			textpropStore, err = library.openTextpropStoreFrom(library.sink, name)
		} else if library.source.HasResource(name) {
			textpropStore, err = library.openTextpropStoreFrom(library.source, name)
		} else {
			textpropStore = library.createSavingTextpropStore(textprop.NullProvider(), "", name, func() {})
		}
		if err == nil {
			library.textpropStores[name] = textpropStore
		}
	}

	return
}

func (library *ReleaseStoreLibrary) openChunkStoreFrom(rel release.Release, name string) (chunkStore *DynamicChunkStore, err error) {
	resource, err := rel.GetResource(name)

	if err != nil {
		return
	}
	reader, err := resource.AsSource()
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(reader)
	_ = reader.Close()
	if err != nil {
		return
	}
	provider, err := resfile.ReaderFrom(bytes.NewReader(data))
	if err != nil {
		return
	}
	chunkStore = library.createSavingChunkStore(provider, resource.Path(), name)

	return
}

func (library *ReleaseStoreLibrary) openObjpropStoreFrom(rel release.Release, name string) (objpropStore objprop.Store, err error) {
	resource, err := rel.GetResource(name)

	if err == nil {
		var reader serial.SeekingReadCloser
		reader, err = resource.AsSource()
		if err == nil {
			var provider objprop.Provider
			provider, err = dosObjprop.NewProvider(reader, library.descriptors)
			if err == nil {
				objpropStore = library.createSavingObjpropStore(provider, resource.Path(), name, func() { _ = reader.Close() })
			}
		}
	}

	return
}

func (library *ReleaseStoreLibrary) openTextpropStoreFrom(rel release.Release, name string) (textpropStore textprop.Store, err error) {
	resource, err := rel.GetResource(name)

	if err == nil {
		var reader serial.SeekingReadCloser
		reader, err = resource.AsSource()
		if err == nil {
			var provider textprop.Provider
			provider, err = dosTextprop.NewProvider(reader)
			if err == nil {
				textpropStore = library.createSavingTextpropStore(provider, resource.Path(), name, func() { _ = reader.Close() })
			}
		}
	}

	return
}

func (library *ReleaseStoreLibrary) createSavingChunkStore(provider chunk.Provider, path string, name string) *DynamicChunkStore {
	storeChanged := make(chan interface{})
	onStoreChanged := func() { storeChanged <- nil }
	chunkStore := NewDynamicChunkStore(chunk.NewProviderBackedStore(provider), onStoreChanged)

	saveAndSwap := func() {
		chunkStore.Swap(func(oldStore chunk.Store) chunk.Store {
			log.Printf("Saving resource <%s>/<%s>\n", path, name)
			data := library.serializeChunkStore(oldStore)
			log.Printf("Serialized previous data, recreating new reader for new data")
			newProvider := library.saveAndReloadChunkData(data, path, name)

			return chunk.NewProviderBackedStore(newProvider)
		})
	}
	library.startSaverRoutine(name, storeChanged, saveAndSwap)

	return chunkStore
}

func (library *ReleaseStoreLibrary) serializeChunkStore(store chunk.Store) []byte {
	buffer := serial.NewByteStore()
	_ = resfile.Write(buffer, store)

	return buffer.Data()
}

func (library *ReleaseStoreLibrary) saveAndReloadChunkData(data []byte, path string, name string) (provider chunk.Provider) {
	newResource, err := library.saveResource(data, path, name)
	var reader serial.SeekingReadCloser

	if err == nil {
		reader, err = newResource.AsSource()
	}
	if err == nil {
		var newData []byte
		newData, err = ioutil.ReadAll(reader)
		_ = reader.Close()
		if err == nil {
			provider, err = resfile.ReaderFrom(bytes.NewReader(newData))
		}
	}
	if err != nil {
		log.Printf("Failed to store in sink, buffering: %v\n", err)
		provider, _ = resfile.ReaderFrom(bytes.NewReader(data))
	}

	return
}

func (library *ReleaseStoreLibrary) createSavingObjpropStore(provider objprop.Provider, path string, name string, closer func()) objprop.Store {
	storeChanged := make(chan interface{})
	onStoreChanged := func() { storeChanged <- nil }
	propStore := NewDynamicObjPropStore(storeObjprop.NewProviderBacked(provider, onStoreChanged))

	closeLastReader := closer
	saveAndSwap := func() {
		propStore.Swap(func(oldStore objprop.Store) objprop.Store {
			log.Printf("Saving resource <%s>/<%s>\n", path, name)
			data := library.serializeObjpropStore(oldStore)
			log.Printf("Serialized previous data, closing old reader")
			closeLastReader()

			log.Printf("Recreating new reader for new data")
			newProvider, newReader := library.saveAndReloadObjpropData(data, path, name)
			closeLastReader = func() { _ = newReader.Close() }

			return storeObjprop.NewProviderBacked(newProvider, onStoreChanged)
		})
	}
	library.startSaverRoutine(name, storeChanged, saveAndSwap)

	return propStore
}

func (library *ReleaseStoreLibrary) serializeObjpropStore(store objprop.Store) []byte {
	buffer := serial.NewByteStore()
	consumer := dosObjprop.NewConsumer(buffer, library.descriptors)

	for classIndex, classDesc := range library.descriptors {
		for subclassIndex, subclassDesc := range classDesc.Subclasses {
			for typeIndex := uint32(0); typeIndex < subclassDesc.TypeCount; typeIndex++ {
				objID := res.MakeObjectID(res.ObjectClass(classIndex), res.ObjectSubclass(subclassIndex), res.ObjectType(typeIndex))
				data := store.Get(objID)
				consumer.Consume(objID, data)
			}
		}
	}
	consumer.Finish()

	return buffer.Data()
}

func (library *ReleaseStoreLibrary) saveAndReloadObjpropData(data []byte, path string, name string) (provider objprop.Provider, reader serial.SeekingReadCloser) {
	newResource, err := library.saveResource(data, path, name)

	if err == nil {
		reader, err = newResource.AsSource()
	}
	if err == nil {
		provider, err = dosObjprop.NewProvider(reader, library.descriptors)
		if err != nil {
			_ = reader.Close()
		}
	}
	if err != nil {
		log.Printf("Failed to store in sink, buffering: %v\n", err)
		reader = serial.NewByteStoreFromData(data, func([]byte) {})
		provider, _ = dosObjprop.NewProvider(reader, library.descriptors)
	}

	return
}

func (library *ReleaseStoreLibrary) createSavingTextpropStore(provider textprop.Provider, path string, name string, closer func()) textprop.Store {
	storeChanged := make(chan interface{})
	onStoreChanged := func() { storeChanged <- nil }
	propStore := NewDynamicTextPropStore(storeTextprop.NewProviderBacked(provider, onStoreChanged))

	closeLastReader := closer
	saveAndSwap := func() {
		propStore.Swap(func(oldStore textprop.Store) textprop.Store {
			log.Printf("Saving resource <%s>/<%s>\n", path, name)
			data := library.serializeTextpropStore(oldStore)
			log.Printf("Serialized previous data, closing old reader")
			closeLastReader()

			log.Printf("Recreating new reader for new data")
			newProvider, newReader := library.saveAndReloadTextpropData(data, path, name)
			closeLastReader = func() { _ = newReader.Close() }

			return storeTextprop.NewProviderBacked(newProvider, onStoreChanged)
		})
	}
	library.startSaverRoutine(name, storeChanged, saveAndSwap)

	return propStore
}

func (library *ReleaseStoreLibrary) serializeTextpropStore(store textprop.Store) []byte {
	buffer := serial.NewByteStore()
	consumer := dosTextprop.NewConsumer(buffer)

	for textureIndex := uint32(0); textureIndex < store.EntryCount(); textureIndex++ {
		data := store.Get(textureIndex)
		consumer.Consume(textureIndex, data)
	}
	consumer.Finish()

	return buffer.Data()
}

func (library *ReleaseStoreLibrary) saveAndReloadTextpropData(data []byte, path string, name string) (provider textprop.Provider, reader serial.SeekingReadCloser) {
	newResource, err := library.saveResource(data, path, name)

	if err == nil {
		reader, err = newResource.AsSource()
	}
	if err == nil {
		provider, err = dosTextprop.NewProvider(reader)
		if err != nil {
			_ = reader.Close()
		}
	}
	if err != nil {
		log.Printf("Failed to store in sink, buffering: %v\n", err)
		reader = serial.NewByteStoreFromData(data, func([]byte) {})
		provider, _ = dosTextprop.NewProvider(reader)
	}

	return
}

func (library *ReleaseStoreLibrary) startSaverRoutine(name string, storeChanged <-chan interface{}, saveAndSwap func()) {
	saveNow := make(chan interface{})
	library.saveChannels[name] = saveNow

	go func() {
		for {
			select {
			case <-saveNow:
			case <-storeChanged:
				for saved := false; !saved; {
					doSave := func() {
						saveAndSwap()
						saved = true
					}
					select {
					case <-storeChanged:
					case <-saveNow:
						doSave()
					case <-time.After(time.Duration(library.timeoutMSec) * time.Millisecond):
						doSave()
					}
				}
			}
		}
	}()
}

func (library *ReleaseStoreLibrary) saveResource(data []byte, path string, name string) (newResource release.Resource, err error) {
	if library.sink.HasResource(name) {
		log.Printf("Sink has resource <%s>, acquiring\n", name)
		newResource, err = library.sink.GetResource(name)
	} else {
		log.Printf("Creating sink resource <%s>\n", name)
		newResource, err = library.sink.NewResource(name, path)
	}
	if err == nil {
		var newSink serial.SeekingWriteCloser

		newSink, err = newResource.AsSink()
		if err == nil {
			_, _ = newSink.Write(data)
			_ = newSink.Close()
		}
	}

	return
}
