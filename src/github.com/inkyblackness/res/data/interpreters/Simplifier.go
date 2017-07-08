package interpreters

import (
	"math"
)

// RawValueHandler is for a simple value range.
type RawValueHandler func(minValue, maxValue int64)

// EnumValueHandler is for enumerated (mapped) values.
type EnumValueHandler func(values map[uint32]string)

// BitfieldHandler is for bitfields.
type BitfieldHandler func(values map[uint32]string)

// ObjectIndexHandler is for object indices.
type ObjectIndexHandler func()

// SpecialHandler is for rare occasions.
type SpecialHandler func()

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

// Bitfield creates a field range describing bitfield values.
func Bitfield(values map[uint32]string) FieldRange {
	return func(simpl *Simplifier) bool {
		return simpl.bitfield(values)
	}
}

// ObjectIndex creates a field range describing object indices.
func ObjectIndex() FieldRange {
	return func(simpl *Simplifier) bool {
		return simpl.objectIndex()
	}
}

// SpecialValue creates a field range for special fields.
// Currently known special values:
// * BinaryCodedDecimal - for keypads storing their number as BCD
// * LevelTexture - index value into level texture list
// * VariableKey - for actions
// * VariableCondition - for action conditions
// * ObjectType - for 0x00CCSSTT selection
// * ObjectHeight - for level height value 0..255
// * MoveTileHeight - for change tile height action
// * Unknown - It is unclear whether this field would have any effect, none identified so far
// * Ignored - Although values have been found in this field, they don't appear to have any effect
// * Mistake - It is assumed that these values should have been placed somewhere else. Typical example: Container content
func SpecialValue(specialType string) FieldRange {
	return func(simpl *Simplifier) bool {
		return simpl.specialValue(specialType)
	}
}

// Simplifier forwards descriptions in a way the requester can use.
type Simplifier struct {
	rawValueHandler    RawValueHandler
	enumValueHandler   EnumValueHandler
	bitfieldHandler    BitfieldHandler
	objectIndexHandler ObjectIndexHandler
	specialHandler     map[string]SpecialHandler
}

// NewSimplifier returns a new instance of a simplifier, with the minimal
// handler set.
func NewSimplifier(rawValueHandler RawValueHandler) *Simplifier {
	return &Simplifier{
		rawValueHandler: rawValueHandler,
		specialHandler:  make(map[string]SpecialHandler)}
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

// SetBitfieldHandler registers the handler for bitfields.
func (simpl *Simplifier) SetBitfieldHandler(handler BitfieldHandler) {
	simpl.bitfieldHandler = handler
}

func (simpl *Simplifier) bitfield(values map[uint32]string) (result bool) {
	if simpl.bitfieldHandler != nil {
		simpl.bitfieldHandler(values)
		result = true
	}
	return
}

// SetObjectIndexHandler registers the handler for object indices.
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

// SetSpecialHandler registers the handler for special values.
func (simpl *Simplifier) SetSpecialHandler(specialType string, handler SpecialHandler) {
	simpl.specialHandler[specialType] = handler
}

func (simpl *Simplifier) specialValue(specialType string) (result bool) {
	handler, existing := simpl.specialHandler[specialType]
	if existing && (handler != nil) {
		handler()
		result = true
	}
	return
}
