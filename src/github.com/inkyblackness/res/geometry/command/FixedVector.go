package command

import "github.com/inkyblackness/res/geometry"

type fixedVector struct {
	value Vector
}

// NewFixedVector returns a vector based on fixed values.
func NewFixedVector(value Vector) geometry.Vector {
	return &fixedVector{value}
}

func (vec *fixedVector) X() float32 {
	return vec.value.X.Float()
}

func (vec *fixedVector) Y() float32 {
	return vec.value.Y.Float()
}

func (vec *fixedVector) Z() float32 {
	return vec.value.Z.Float()
}
