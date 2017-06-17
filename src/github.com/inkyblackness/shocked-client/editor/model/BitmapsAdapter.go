package model

import (
	"fmt"

	"github.com/inkyblackness/shocked-model"
)

// BitmapsAdapter is the entry point for bitmaps.
type BitmapsAdapter struct {
	context projectContext
	store   model.DataStore

	bitmaps               *Bitmaps
	bitmapRequestsPending map[model.ResourceKey]bool
}

func newBitmapsAdapter(context projectContext, store model.DataStore) *BitmapsAdapter {
	adapter := &BitmapsAdapter{
		context: context,
		store:   store,

		bitmaps:               newBitmaps(),
		bitmapRequestsPending: make(map[model.ResourceKey]bool)}

	return adapter
}

func (adapter *BitmapsAdapter) clear() {
	adapter.bitmaps.clear()
}

func (adapter *BitmapsAdapter) refresh() {
}

// RequestBitmapChange will update the bitmap data for identified key.
func (adapter *BitmapsAdapter) RequestBitmapChange(key model.ResourceKey, newBitmap *model.RawBitmap) {
	adapter.store.SetBitmap(adapter.context.ActiveProjectID(), key, newBitmap,
		func(resultKey model.ResourceKey, bmp *model.RawBitmap) {
			adapter.bitmaps.setRawBitmap(key.ToInt(), bmp)
		},
		func() {
			adapter.context.simpleStoreFailure(fmt.Sprintf("SetBitmap[%v]", key))()
		})
}

// RequestBitmap will load the bitmap data for identified key.
func (adapter *BitmapsAdapter) RequestBitmap(key model.ResourceKey) {
	if !adapter.bitmapRequestsPending[key] {
		adapter.bitmapRequestsPending[key] = true
		adapter.store.Bitmap(adapter.context.ActiveProjectID(), key,
			func(resultKey model.ResourceKey, bmp *model.RawBitmap) {
				adapter.bitmapRequestsPending[key] = false
				adapter.bitmaps.setRawBitmap(key.ToInt(), bmp)
			},
			func() {
				adapter.bitmapRequestsPending[key] = false
				adapter.context.simpleStoreFailure(fmt.Sprintf("Bitmap[%v]", key))()
			})
	}
}

// Bitmaps returns the container of bitmaps.
func (adapter *BitmapsAdapter) Bitmaps() *Bitmaps {
	return adapter.bitmaps
}
