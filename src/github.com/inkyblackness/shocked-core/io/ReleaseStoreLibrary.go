package io

import (
	"log"
	"time"

	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/chunk/dos"
	"github.com/inkyblackness/res/chunk/store"
	"github.com/inkyblackness/res/serial"
	"github.com/inkyblackness/shocked-core/release"
)

type ReleaseStoreLibrary struct {
	source      release.Release
	sink        release.Release
	timeoutMSec int
	chunkStores map[string]chunk.Store
}

// NewReleaseStoreLibrary returns a StoreLibrary that covers two Release container.
// Stores are first tried from the sink, then from source; They are saved in the sink.
func NewReleaseStoreLibrary(source release.Release, sink release.Release, timeoutMSec int) StoreLibrary {
	library := &ReleaseStoreLibrary{
		source:      source,
		sink:        sink,
		timeoutMSec: timeoutMSec,
		chunkStores: make(map[string]chunk.Store)}

	return library
}

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

func (library *ReleaseStoreLibrary) openChunkStoreFrom(rel release.Release, name string) (chunkStore chunk.Store, err error) {
	resource, err := rel.GetResource(name)

	if err == nil {
		var reader serial.SeekingReadCloser
		reader, err = resource.AsSource()
		if err == nil {
			var provider chunk.Provider
			provider, err = dos.NewChunkProvider(reader)
			if err == nil {
				chunkStore = library.createSavingChunkStore(provider, resource.Path(), name, func() { reader.Close() })
			}
		}
	}

	return
}

func (library *ReleaseStoreLibrary) createSavingChunkStore(provider chunk.Provider, path string, name string, closer func()) chunk.Store {
	storeChanged := make(chan interface{})
	onStoreChanged := func() { storeChanged <- nil }
	chunkStore := NewDynamicChunkStore(store.NewProviderBacked(provider, onStoreChanged))

	closeLastReader := closer
	saveAndSwap := func() {
		chunkStore.Swap(func(oldStore chunk.Store) chunk.Store {
			log.Printf("Saving resource <%s>/<%s>\n", path, name)
			data := library.serializeStore(oldStore)
			log.Printf("Serialized previous data, closing old reader")
			closeLastReader()

			log.Printf("Recreating new reader for new data")
			newProvider, newReader := library.saveAndReload(data, path, name)
			closeLastReader = func() { newReader.Close() }

			return store.NewProviderBacked(newProvider, onStoreChanged)
		})
	}

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

	return chunkStore
}

func (library *ReleaseStoreLibrary) serializeStore(store chunk.Store) []byte {
	buffer := serial.NewByteStore()
	consumer := dos.NewChunkConsumer(buffer)
	ids := store.IDs()

	for _, id := range ids {
		blockStore := store.Get(id)
		consumer.Consume(id, blockStore)
	}
	consumer.Finish()

	return buffer.Data()
}

func (library *ReleaseStoreLibrary) saveAndReload(data []byte, path string, name string) (provider chunk.Provider, reader serial.SeekingReadCloser) {
	var newResource release.Resource
	var err error

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

	if err == nil {
		reader, err = newResource.AsSource()
	}
	if err == nil {
		provider, err = dos.NewChunkProvider(reader)
		if err != nil {
			reader.Close()
		}
	}
	if err != nil {
		log.Printf("Failed to store in sink, buffering: %v\n", err)
		reader = serial.NewByteStoreFromData(data, func([]byte) {})
		provider, _ = dos.NewChunkProvider(reader)
	}

	return
}
