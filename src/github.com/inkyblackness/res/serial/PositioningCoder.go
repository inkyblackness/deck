package serial

// PositioningCoder is a coder that also knows about positioning
type PositioningCoder interface {
	Positioner
	Coder
}
