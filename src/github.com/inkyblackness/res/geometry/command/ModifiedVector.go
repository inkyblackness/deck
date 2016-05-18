package command

import "github.com/inkyblackness/res/geometry"

type modifiedVector struct {
	reference geometry.Vector

	xMod Modifier
	yMod Modifier
	zMod Modifier
}

// NewModifiedVector returns a vector that refers to another with modifications
func NewModifiedVector(reference geometry.Vector, xMod, yMod, zMod Modifier) geometry.Vector {
	return &modifiedVector{
		reference: reference,
		xMod:      xMod,
		yMod:      yMod,
		zMod:      zMod}
}

func (vec *modifiedVector) X() float32 {
	return vec.xMod(vec.reference.X())
}

func (vec *modifiedVector) Y() float32 {
	return vec.yMod(vec.reference.Y())
}

func (vec *modifiedVector) Z() float32 {
	return vec.zMod(vec.reference.Z())
}
