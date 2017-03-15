package model

type observerFunc func()

type observable struct {
	value           interface{}
	observerCounter int
	observers       map[int]observerFunc
}

func newObservable() *observable {
	return &observable{
		observers: make(map[int]observerFunc)}
}

func (obj *observable) get() interface{} {
	return obj.value
}

func (obj *observable) orDefault(alt interface{}) interface{} {
	result := obj.value
	if result == nil {
		result = alt
	}
	return result
}

func (obj *observable) set(value interface{}) {
	if obj.value != value {
		obj.value = value
		obj.notifyObservers()
	}
}

func (obj *observable) addObserver(fn observerFunc) func() {
	key := obj.observerCounter
	obj.observerCounter++
	obj.observers[key] = fn

	return func() {
		delete(obj.observers, key)
	}
}

func (obj *observable) notifyObservers() {
	observers := make([]int, 0, len(obj.observers))
	for key := range obj.observers {
		observers = append(observers, key)
	}
	for _, key := range observers {
		if observer, existing := obj.observers[key]; existing {
			observer()
		}
	}
}
