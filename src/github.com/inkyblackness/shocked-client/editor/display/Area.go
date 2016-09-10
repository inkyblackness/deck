package display

// Area describes a rectangular area with a center.
type Area interface {
	Center() (x, y float32)
	Size() (width, height float32)
}
