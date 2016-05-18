package env

import (
	"github.com/inkyblackness/shocked-client/viewmodel"
)

// Application represents the public interface between the environment and the
// actual application core.
type Application interface {
	// ViewModel returns the root node of the view model of the application.
	ViewModel() viewmodel.Node
	// Init sets up the application for the given OpenGL window.
	Init(window OpenGlWindow)
}
