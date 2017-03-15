package model

import (
	"github.com/inkyblackness/shocked-model"
)

// Bitmaps is a container of related bitmaps.
type Bitmaps struct {
	data map[int]*observable
}

func newBitmaps() *Bitmaps {
	bitmaps := &Bitmaps{data: make(map[int]*observable)}

	return bitmaps
}

func (bitmaps *Bitmaps) clear() {
	for _, data := range bitmaps.data {
		data.set(nil)
	}
}

func (bitmaps *Bitmaps) setRawBitmap(key int, bmp *model.RawBitmap) {
	bitmaps.ensureData(key).set(bmp)
}

// RawBitmap returns the raw bitmap information of the identified texture.
func (bitmaps *Bitmaps) RawBitmap(key int) (bmp *model.RawBitmap) {
	if data, existing := bitmaps.data[key]; existing {
		bmp = data.get().(*model.RawBitmap)
	}

	return
}

// OnBitmapChanged registers a callback for updates.
func (bitmaps *Bitmaps) OnBitmapChanged(key int, callback func()) {
	bitmaps.ensureData(key).addObserver(callback)
}

func (bitmaps *Bitmaps) ensureData(key int) *observable {
	data, existing := bitmaps.data[key]
	if !existing {
		data = newObservable()
		bitmaps.data[key] = data
	}
	return data
}
