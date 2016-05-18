package command

import (
	"github.com/inkyblackness/res/geometry"
)

func SaveVertices(writer *Writer, container geometry.VertexContainer) {
	positions := make([]geometry.Vector, container.VertexCount())
	for i := 0; i < len(positions); i++ {
		positions[i] = container.Vertex(i).Position()
	}

	lastDissimilar := len(positions) - 1
	for lastDissimilar > 0 && findMostSimilarVector(positions[:lastDissimilar], positions[lastDissimilar]) < lastDissimilar {
		lastDissimilar--
	}

	if lastDissimilar > 0 {
		writer.WriteDefineVertices(positions[:lastDissimilar+1])
	} else if len(positions) > 0 {
		writer.WriteDefineVertex(positions[0])
	}

	for index := lastDissimilar + 1; index < len(positions); index++ {
		vector := positions[index]
		referenceIndex := findMostSimilarVector(positions[:index], vector)
		reference := positions[referenceIndex]

		if fixedEqual(reference.X(), vector.X()) {
			if fixedEqual(reference.Y(), vector.Y()) {
				writer.WriteDefineOneOffsetVertex(CmdDefineOffsetVertexZ, index, referenceIndex, vector.Z()-reference.Z())
			} else if fixedEqual(reference.Z(), vector.Z()) {
				writer.WriteDefineOneOffsetVertex(CmdDefineOffsetVertexY, index, referenceIndex, vector.Y()-reference.Y())
			} else {
				writer.WriteDefineTwoOffsetVertex(CmdDefineOffsetVertexYZ, index, referenceIndex,
					vector.Y()-reference.Y(), vector.Z()-reference.Z())
			}
		} else if fixedEqual(reference.Y(), vector.Y()) {
			if fixedEqual(reference.Z(), vector.Z()) {
				writer.WriteDefineOneOffsetVertex(CmdDefineOffsetVertexX, index, referenceIndex, vector.X()-reference.X())
			} else {
				writer.WriteDefineTwoOffsetVertex(CmdDefineOffsetVertexXZ, index, referenceIndex,
					vector.X()-reference.X(), vector.Z()-reference.Z())
			}
		} else {
			writer.WriteDefineTwoOffsetVertex(CmdDefineOffsetVertexXY, index, referenceIndex,
				vector.X()-reference.X(), vector.Y()-reference.Y())
		}
	}
}

func fixedEqual(a, b float32) bool {
	return ToFixed(a) == ToFixed(b)
}

func findMostSimilarVector(list []geometry.Vector, vector geometry.Vector) int {
	foundIndex := len(list)
	sameCoordinates := 0

	for index, other := range list {
		temp := 0
		if fixedEqual(other.X(), vector.X()) {
			temp++
		}
		if fixedEqual(other.Y(), vector.Y()) {
			temp++
		}
		if fixedEqual(other.Z(), vector.Z()) {
			temp++
		}
		if temp > sameCoordinates {
			foundIndex = index
			sameCoordinates = temp
		}
	}

	return foundIndex
}
