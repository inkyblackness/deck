package command

import (
	"github.com/inkyblackness/res/geometry"
)

// SaveModel saves the given model into a command list and returns the resulting bytes.
func SaveModel(model geometry.Model) []byte {
	commandWriter := NewWriter()

	commandWriter.WriteHeader(geometry.CountFaces(model))
	SaveVertices(commandWriter, model)
	commandWriter.WriteBytes(SaveNode(model))

	return commandWriter.Bytes()
}
