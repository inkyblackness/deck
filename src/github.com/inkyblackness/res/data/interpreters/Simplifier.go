package interpreters

import (
	"math"
)

// RawValueHandler is for a simple value range.
type RawValueHandler func(minValue, maxValue int64)

// EnumValueHandler is for enumerated (mapped) values.
type EnumValueHandler func(values map[uint32]string)

// ObjectIndexHandler is for object indices.
type ObjectIndexHandler func()

// RangedValue creates a field range for specific minimum and maximum values.
func RangedValue(minValue, maxValue int64) FieldRange {
	return func(simpl *Simplifier) bool {
		return simpl.rangedValue(minValue, maxValue)
	}
}

// EnumValue creates a field range describing enumerated values.
func EnumValue(values map[uint32]string) FieldRange {
	return func(simpl *Simplifier) bool {
		return simpl.enumValue(values)
	}
}

// ObjectIndex creates a field range describing object indices.
func ObjectIndex() FieldRange {
	return func(simpl *Simplifier) bool {
		return simpl.objectIndex()
	}
}

// Simplifier forwards descriptions in a way the requester can use.
type Simplifier struct {
	rawValueHandler    RawValueHandler
	enumValueHandler   EnumValueHandler
	objectIndexHandler ObjectIndexHandler
}

// NewSimplifier returns a new instance of a simplifier, with the minimal
// handler set.
func NewSimplifier(rawValueHandler RawValueHandler) *Simplifier {
	return &Simplifier{rawValueHandler: rawValueHandler}
}

func (simpl *Simplifier) rawValue(e *entry) {
	max := int64(math.Pow(2, float64(e.count*8)))
	if max == 256 {
		simpl.rawValueHandler(0, 255)
	} else {
		half := max / 2
		simpl.rawValueHandler(-1, half-1)
	}
}

func (simpl *Simplifier) rangedValue(minValue, maxValue int64) bool {
	simpl.rawValueHandler(minValue, maxValue)
	return true
}

// SetEnumValueHandler registers the handler for enumerations.
func (simpl *Simplifier) SetEnumValueHandler(handler EnumValueHandler) {
	simpl.enumValueHandler = handler
}

func (simpl *Simplifier) enumValue(values map[uint32]string) (result bool) {
	if simpl.enumValueHandler != nil {
		simpl.enumValueHandler(values)
		result = true
	}
	return
}

// SetObjectIndexHandler registers the handler for enumerations.
func (simpl *Simplifier) SetObjectIndexHandler(handler ObjectIndexHandler) {
	simpl.objectIndexHandler = handler
}

func (simpl *Simplifier) objectIndex() (result bool) {
	if simpl.objectIndexHandler != nil {
		simpl.objectIndexHandler()
		result = true
	}
	return
}
