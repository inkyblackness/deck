package styling

// A StyleFunc wraps the provided data and produces a string in the given style.
type StyleFunc func(...interface{}) string

// A Style provides a collection of style functions.
type Style interface {
	// Prompt is for user queries.
	Prompt() StyleFunc
	// Error is for error indications.
	Error() StyleFunc
	// Status is for general status messages.
	Status() StyleFunc
	// Added is for added/modified data.
	Added() StyleFunc
	// Removed is for lost/removed data.
	Removed() StyleFunc
}
