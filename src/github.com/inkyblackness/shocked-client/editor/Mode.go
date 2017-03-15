package editor

// Mode is the general interface for an editor mode.
type Mode interface {
	// SetActive determines whether a mode is currently active.
	// Inactive modes should hide their UI and release any focus.
	SetActive(active bool)
}
