package io

import (
	"log"
	"time"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	dosChunk "github.com/inkyblackness/res/chunk/dos"
	storeChunk "github.com/inkyblackness/res/chunk/store"
	"github.com/inkyblackness/res/objprop"
	dosObjprop "github.com/inkyblackness/res/objprop/dos"
	storeObjprop "github.com/inkyblackness/res/objprop/store"
	"github.com/inkyblackness/res/serial"
	"github.com/inkyblackness/shocked-core/release"
)

// ReleaseStoreLibrary is a container with two releases: one source and one sink.
// Stores can be retrieved from this library, which access the source release for
// reading properties and the sink release for writing modified data.
type ReleaseStoreLibrary struct {
	source      release.Release
	sink        release.Release
	timeoutMSec int
	chunkStores map[string]chunk.Store

	descriptors   []objprop.ClassDescriptor
	objpropStores map[string]objprop.Store
}

// NewReleaseStoreLibrary returns a StoreLibrary that covers two Release container.
// Stores are first tried from the sink, then from source; They are saved in the sink.
func NewReleaseStoreLibrary(source release.Release, sink release.Release, timeoutMSec int) StoreLibrary {
	library := &ReleaseStoreLibrary{
		source:      source,
		sink:        sink,
		timeoutMSec: timeoutMSec,
		chunkStores: make(map[string]chunk.Store),

		descriptors:   objprop.StandardProperties(),
		objpropStores: make(map[string]objprop.Store)}

	return library
}

// ChunkStore implements the StoreLibrary interface.
func (library *ReleaseStoreLibrary) ChunkStore(name string) (chunkStore chunk.Store, err error) {
	chunkStore, exists := library.chunkStores[name]

	if !exists {
		if library.sink.HasResource(name) {
			chunkStore, err = library.openChunkStoreFrom(library.sink, name)
		} else if library.source.HasResource(name) {
			chunkStore, err = library.openChunkStoreFrom(library.source, name)
		} else {
			chunkStore = library.createSavingChunkStore(chunk.NullProvider(), "", name, func() {})
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

func (library *ReleaseStoreLibrary) openChunkStoreFrom(rel release.Release, name string) (chunkStore chunk.Store, err error) {
	resource, err := rel.GetResource(name)

	if err == nil {
		var reader serial.SeekingReadCloser
		reader, err = resource.AsSource()
		if err == nil {
			var provider chunk.Provider
			provider, err = dosChunk.NewChunkProvider(reader)
			if err == nil {
				chunkStore = library.createSavingChunkStore(provider, resource.Path(), name, func() { reader.Close() })
			}
		}
	}

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
				objpropStore = library.createSavingObjpropStore(provider, resource.Path(), name, func() { reader.Close() })
			}
		}
	}

	return
}

func (library *ReleaseStoreLibrary) createSavingChunkStore(provider chunk.Provider, path string, name string, closer func()) chunk.Store {
	storeChanged := make(chan interface{})
	onStoreChanged := func() { storeChanged <- nil }
	chunkStore := NewDynamicChunkStore(storeChunk.NewProviderBacked(provider, onStoreChanged))

	closeLastReader := closer
	saveAndSwap := func() {
		chunkStore.Swap(func(oldStore chunk.Store) chunk.Store {
			log.Printf("Saving resource <%s>/<%s>\n", path, name)
			data := library.serializeChunkStore(oldStore)
			log.Printf("Serialized previous data, closing old reader")
			closeLastReader()

			log.Printf("Recreating new reader for new data")
			newProvider, newReader := library.saveAndReloadChunkData(data, path, name)
			closeLastReader = func() { newReader.Close() }

			return storeChunk.NewProviderBacked(newProvider, onStoreChanged)
		})
	}
	library.startSaverRoutine(storeChanged, saveAndSwap)

	return chunkStore
}

func (library *ReleaseStoreLibrary) serializeChunkStore(store chunk.Store) []byte {
	buffer := serial.NewByteStore()
	consumer := dosChunk.NewChunkConsumer(buffer)
	ids := store.IDs()

	for _, id := range ids {
		blockStore := store.Get(id)
		consumer.Consume(id, blockStore)
	}
	consumer.Finish()

	return buffer.Data()
}

func (library *ReleaseStoreLibrary) saveAndReloadChunkData(data []byte, path string, name string) (provider chunk.Provider, reader serial.SeekingReadCloser) {
	newResource, err := library.saveResource(data, path, name)

	if err == nil {
		reader, err = newResource.AsSource()
	}
	if err == nil {
		provider, err = dosChunk.NewChunkProvider(reader)
		if err != nil {
			reader.Close()
		}
	}
	if err != nil {
		log.Printf("Failed to store in sink, buffering: %v\n", err)
		reader = serial.NewByteStoreFromData(data, func([]byte) {})
		provider, _ = dosChunk.NewChunkProvider(reader)
	}

	return
}

func (library *ReleaseStoreLibrary) createSavingObjpropStore(provider objprop.Provider, path string, name string, closer func()) objprop.Store {
	storeChanged := make(chan interface{})
	onStoreChanged := func() { storeChanged <- nil }
	chunkStore := NewDynamicObjPropStore(storeObjprop.NewProviderBacked(provider, onStoreChanged))

	closeLastReader := closer
	saveAndSwap := func() {
		chunkStore.Swap(func(oldStore objprop.Store) objprop.Store {
			log.Printf("Saving resource <%s>/<%s>\n", path, name)
			data := library.serializeObjpropStore(oldStore)
			log.Printf("Serialized previous data, closing old reader")
			closeLastReader()

			log.Printf("Recreating new reader for new data")
			newProvider, newReader := library.saveAndReloadObjpropData(data, path, name)
			closeLastReader = func() { newReader.Close() }

			return storeObjprop.NewProviderBacked(newProvider, onStoreChanged)
		})
	}
	library.startSaverRoutine(storeChanged, saveAndSwap)

	return chunkStore
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
			reader.Close()
		}
	}
	if err != nil {
		log.Printf("Failed to store in sink, buffering: %v\n", err)
		reader = serial.NewByteStoreFromData(data, func([]byte) {})
		provider, _ = dosObjprop.NewProvider(reader, library.descriptors)
	}

	return
}

func (library *ReleaseStoreLibrary) startSaverRoutine(storeChanged <-chan interface{}, saveAndSwap func()) {
	go func() {
		for true {
			<-storeChanged
			for saved := false; !saved; {
				select {
				case <-storeChanged:
				case <-time.After(time.Duration(library.timeoutMSec) * time.Millisecond):
					saveAndSwap()
					saved = true
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
			newSink.Write(data)
			newSink.Close()
		}
	}

	return
}