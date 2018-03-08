package command

import (
	"fmt"

	"github.com/inkyblackness/res/serial"
)

// Fixed is a value type for serialized coordinates.
type Fixed uint32

const fixedFactor = float32(0x10000)

// ToFixed creates a Fixed out of a floating point value.
func ToFixed(value float32) Fixed {
	return Fixed(uint32(value * fixedFactor))
}

// CodeFixed serializes a Fixed value with a coder.
func CodeFixed(coder serial.Coder, fixed *Fixed) {
	raw := uint32(*fixed)
	coder.Code(&raw)
	*fixed = Fixed(raw)
}

// Float returns the closest floating point value to the fixed value.
func (fixed Fixed) Float() float32 {
	return float32(int32(fixed)) / fixedFactor
}

// String returns the string presentation of the result of Float().
func (fixed Fixed) String() string {
	return fmt.Sprintf("%v", fixed.Float())
}
