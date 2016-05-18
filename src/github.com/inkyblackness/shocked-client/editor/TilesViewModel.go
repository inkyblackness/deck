package editor

import (
	"fmt"

	"github.com/inkyblackness/shocked-client/viewmodel"
	"github.com/inkyblackness/shocked-model"
)

// TilesViewModel contains the view model entries for the map tiles.
type TilesViewModel struct {
	root *viewmodel.SectionNode

	tileType      *viewmodel.ValueSelectionNode
	floorHeight   *viewmodel.ValueSelectionNode
	ceilingHeight *viewmodel.ValueSelectionNode
	slopeHeight   *viewmodel.ValueSelectionNode
	slopeControl  *viewmodel.ValueSelectionNode

	floorTexture   *viewmodel.ValueSelectionNode
	ceilingTexture *viewmodel.ValueSelectionNode
	wallTexture    *viewmodel.ValueSelectionNode

	floorTextureRotations   *viewmodel.ValueSelectionNode
	ceilingTextureRotations *viewmodel.ValueSelectionNode

	useAdjacentWallTexture *viewmodel.ValueSelectionNode
	wallTextureOffset      *viewmodel.ValueSelectionNode
}

func intStringList(start, stop int) (list []string) {
	for value := start; value <= stop; value++ {
		list = append(list, fmt.Sprintf("%d", value))
	}

	return append(list, "")
}

// NewTilesViewModel returns a new instance of a TilesViewModel.
func NewTilesViewModel(levelIsRealWorld *viewmodel.BoolValueNode) *TilesViewModel {
	vm := &TilesViewModel{}

	vm.tileType = viewmodel.NewValueSelectionNode("Tile Type", []string{string(model.Open), string(model.Solid),
		string(model.DiagonalOpenSouthEast), string(model.DiagonalOpenSouthWest), string(model.DiagonalOpenNorthWest), string(model.DiagonalOpenNorthEast),
		string(model.SlopeSouthToNorth), string(model.SlopeWestToEast), string(model.SlopeNorthToSouth), string(model.SlopeEastToWest),
		string(model.ValleySouthEastToNorthWest), string(model.ValleySouthWestToNorthEast), string(model.ValleyNorthWestToSouthEast), string(model.ValleyNorthEastToSouthWest),
		string(model.RidgeNorthWestToSouthEast), string(model.RidgeNorthEastToSouthWest), string(model.RidgeSouthEastToNorthWest), string(model.RidgeSouthWestToNorthEast),
		""},
		"")
	vm.floorHeight = viewmodel.NewValueSelectionNode("Floor Height Level", intStringList(0, 31), "")
	vm.ceilingHeight = viewmodel.NewValueSelectionNode("Ceiling Height Level", intStringList(1, 32), "")
	vm.slopeHeight = viewmodel.NewValueSelectionNode("Slope Height", intStringList(0, 31), "")
	vm.slopeControl = viewmodel.NewValueSelectionNode("Slope Control",
		[]string{model.SlopeCeilingInverted, model.SlopeCeilingMirrored, model.SlopeCeilingFlat, model.SlopeFloorFlat, ""},
		"")

	vm.floorTexture = viewmodel.NewValueSelectionNode("Floor Texture Index", []string{""}, "")
	vm.ceilingTexture = viewmodel.NewValueSelectionNode("Ceiling Texture Index", []string{""}, "")
	vm.wallTexture = viewmodel.NewValueSelectionNode("Wall Texture Index", []string{""}, "")

	vm.floorTextureRotations = viewmodel.NewValueSelectionNode("Floor Tex Rotations", intStringList(0, 3), "")
	vm.ceilingTextureRotations = viewmodel.NewValueSelectionNode("Ceiling Tex Rotations", intStringList(0, 3), "")
	vm.useAdjacentWallTexture = viewmodel.NewValueSelectionNode("Use Adj. Wall Tex", []string{"yes", "no", ""}, "")
	vm.wallTextureOffset = viewmodel.NewValueSelectionNode("Wall Texture Offset", intStringList(0, 31), "")

	realWorldSection := viewmodel.NewSectionNode("Real World",
		[]viewmodel.Node{vm.floorTexture, vm.ceilingTexture, vm.wallTexture,
			vm.floorTextureRotations, vm.ceilingTextureRotations,
			vm.useAdjacentWallTexture, vm.wallTextureOffset},
		levelIsRealWorld)

	vm.root = viewmodel.NewSectionNode("Tiles",
		[]viewmodel.Node{vm.tileType, vm.floorHeight, vm.ceilingHeight, vm.slopeHeight, vm.slopeControl,
			realWorldSection},
		viewmodel.NewBoolValueNode("", true))

	return vm
}

// SetLevelTextureCount registers the available amount of level textures.
func (vm *TilesViewModel) SetLevelTextureCount(count int) {
	values := []string{""}

	if count > 0 {
		values = intStringList(0, count-1)
	}
	vm.floorTexture.SetValues(values)
	vm.ceilingTexture.SetValues(values)
	vm.wallTexture.SetValues(values)
}

// TileType returns the tile type selection node.
func (vm *TilesViewModel) TileType() *viewmodel.ValueSelectionNode {
	return vm.tileType
}

// FloorHeight returns the floor height selection node.
func (vm *TilesViewModel) FloorHeight() *viewmodel.ValueSelectionNode {
	return vm.floorHeight
}

// CeilingHeight returns the ceiling height selection node.
func (vm *TilesViewModel) CeilingHeight() *viewmodel.ValueSelectionNode {
	return vm.ceilingHeight
}

// SlopeHeight returns the slope height selection node.
func (vm *TilesViewModel) SlopeHeight() *viewmodel.ValueSelectionNode {
	return vm.slopeHeight
}

// SlopeControl returns the slope control selection node.
func (vm *TilesViewModel) SlopeControl() *viewmodel.ValueSelectionNode {
	return vm.slopeControl
}

// FloorTexture returns the selection node for the floor texture.
func (vm *TilesViewModel) FloorTexture() *viewmodel.ValueSelectionNode {
	return vm.floorTexture
}

// CeilingTexture returns the selection node for the ceiling texture.
func (vm *TilesViewModel) CeilingTexture() *viewmodel.ValueSelectionNode {
	return vm.ceilingTexture
}

// WallTexture returns the selection node for the wall texture.
func (vm *TilesViewModel) WallTexture() *viewmodel.ValueSelectionNode {
	return vm.wallTexture
}

// FloorTextureRotations returns the selection node for floor texture rotations.
func (vm *TilesViewModel) FloorTextureRotations() *viewmodel.ValueSelectionNode {
	return vm.floorTextureRotations
}

// CeilingTextureRotations returns the selection node for ceiling texture rotations.
func (vm *TilesViewModel) CeilingTextureRotations() *viewmodel.ValueSelectionNode {
	return vm.ceilingTextureRotations
}

// UseAdjacentWallTexture returns the selection node for using the adjacent wall texture.
func (vm *TilesViewModel) UseAdjacentWallTexture() *viewmodel.ValueSelectionNode {
	return vm.useAdjacentWallTexture
}

// WallTextureOffset returns the selection node for the wall texture offset.
func (vm *TilesViewModel) WallTextureOffset() *viewmodel.ValueSelectionNode {
	return vm.wallTextureOffset
}
