package env

// Application represents the public interface between the environment and the
// actual application core.
type Application interface {
	// Init sets up the application for the given OpenGL window.
	Init(window OpenGlWindow)
}
