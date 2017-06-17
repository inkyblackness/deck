package env

// Clipboard describes access to the clipboard content.
type Clipboard interface {
	// Text returns, if available, the current text content of the clipboard.
	Text() string
	// SetText changes the current text content of the clipboard.
	SetText(value string)
}
