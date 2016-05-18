package res

// ObjectClass is a first identification of an object
type ObjectClass byte

// ObjectSubclass is the second identification of an object
type ObjectSubclass byte

// ObjectType is the specific type of a class/subclass combination
type ObjectType byte

// ObjectID completely identifies a specific object
type ObjectID struct {
	Class    ObjectClass
	Subclass ObjectSubclass
	Type     ObjectType
}

// MakeObjectID returns an object ID with the provided values
func MakeObjectID(class ObjectClass, subclass ObjectSubclass, objType ObjectType) ObjectID {
	return ObjectID{Class: class, Subclass: subclass, Type: objType}
}
