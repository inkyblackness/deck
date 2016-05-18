package browser

import (
	"github.com/Archs/js/gopherjs-ko"
	"github.com/gopherjs/jquery"

	"github.com/inkyblackness/shocked-client/env"
)

// Run initializes the browser environment and hooks up the provided application.
func Run(app env.Application) {
	canvas := jquery.NewJQuery("canvas#output")
	window, _ := NewWebGlWindow(canvas.Get(0))

	app.Init(window)

	root := newViewModelFiller()
	app.ViewModel().Specialize(root)
	vm := ko.ViewModelFromJS(root.object)

	ko.ApplyBindings(vm)
}
