package browser

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
)

// ObjectMapper is a map for handles.
type ObjectMapper interface {
	put(value *js.Object) uint32
	get(key uint32) *js.Object
	del(key uint32) *js.Object
}

type objectMap struct {
	objects map[uint32]*js.Object
	counter uint32
}

// NewObjectMapper returns a new ObjectMapper instance
func NewObjectMapper() ObjectMapper {
	result := &objectMap{
		objects: make(map[uint32]*js.Object),
		counter: 0}

	return result
}

func (omap *objectMap) put(value *js.Object) uint32 {
	key := uint32(0)

	for key == 0 {
		_, exists := omap.objects[omap.counter]

		if (omap.counter == 0) || exists {
			omap.counter++
		} else {
			key = omap.counter
		}
	}
	omap.objects[key] = value

	return key
}

func (omap *objectMap) get(key uint32) *js.Object {
	value, ok := omap.objects[key]

	if !ok && (key != 0) {
		panic(fmt.Sprintf("Object with ID %u not known", key))
	}

	return value
}

func (omap *objectMap) del(key uint32) *js.Object {
	defer delete(omap.objects, key)

	return omap.get(key)
}
