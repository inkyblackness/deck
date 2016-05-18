package command

import "github.com/inkyblackness/res/serial"

// Vector consists of three space coordinates.
type Vector struct {
	// X value
	X Fixed
	// Y value
	Y Fixed
	// Z value
	Z Fixed
}

// Code serializes the vector with given coder.
func (vec *Vector) Code(coder serial.Coder) {
	CodeFixed(coder, &vec.X)
	CodeFixed(coder, &vec.Y)
	CodeFixed(coder, &vec.Z)
}
