package console

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/jroimartin/gocui"

	"github.com/inkyblackness/shocked-client/env"
	"github.com/inkyblackness/shocked-client/env/native"
)

type appRunner struct {
	gui *gocui.Gui
	app env.Application

	rootTexter         ViewModelNodeTexter
	highlightedTexter  ViewModelNodeTexter
	activeDetailTexter ViewModelNodeTexter
	mainControlLines   int

	activeListDetailController   ListDetailController
	activeStringDetailController StringDetailController
}

// Run initializes the environment to run the given application within.
func Run(app env.Application, deferrer <-chan func()) {
	runtime.LockOSThread()

	gui := gocui.NewGui()
	if err := gui.Init(); err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	runner := &appRunner{gui: gui, app: app}

	gui.Cursor = true
	gui.SelBgColor = gocui.ColorGreen
	gui.SelFgColor = gocui.ColorBlack
	gui.SetLayout(runner.layout)

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("mainControls", gocui.KeyArrowDown, gocui.ModNone, runner.cursorDown); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("mainControls", gocui.KeyArrowUp, gocui.ModNone, runner.cursorUp); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("mainControls", gocui.KeyEnter, gocui.ModNone, runner.actMainControl); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("listDetail", gocui.KeyEnter, gocui.ModNone, runner.confirmListDetail); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("listDetail", gocui.KeyBackspace, gocui.ModNone, runner.cancelListDetail); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("listDetail", gocui.KeyBackspace2, gocui.ModNone, runner.cancelListDetail); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("listDetail", gocui.KeyEsc, gocui.ModNone, runner.cancelListDetail); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("listDetail", gocui.KeyArrowDown, gocui.ModNone, runner.cursorDown); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("listDetail", gocui.KeyArrowUp, gocui.ModNone, runner.cursorUp); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("stringDetail", gocui.KeyEnter, gocui.ModNone, runner.confirmStringDetail); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("stringDetail", gocui.KeyEsc, gocui.ModNone, runner.cancelStringDetail); err != nil {
		log.Panicln(err)
	}

	var window *native.OpenGlWindow
	{
		var err error
		window, err = native.NewOpenGlWindow()
		if err != nil {
			log.Panicln(err)
		}
	}

	app.Init(window)

	{
		visitor := NewViewModelTexterVisitor(runner)
		app.ViewModel().Specialize(visitor)
		runner.rootTexter = visitor.instance
	}

	startDeferrerRoutine(gui, deferrer)

	gui.Execute(getWindowUpdater(window))
	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	window.Close()
}

func getWindowUpdater(window *native.OpenGlWindow) (updater func(*gocui.Gui) error) {
	updater = func(gui *gocui.Gui) error {
		window.Update()
		gui.Execute(updater)

		return nil
	}

	return
}

func startDeferrerRoutine(gui *gocui.Gui, deferrer <-chan func()) {
	go func() {
		for task := range deferrer {
			deferredTask := task
			gui.Execute(func(*gocui.Gui) error {
				deferredTask()

				return nil
			})
		}
	}()
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (runner *appRunner) OnMainDataChanged() {
	runner.gui.Execute(runner.layout)
}

func (runner *appRunner) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	view, _ := g.SetView("mainControls", -1, -1, maxX/2, maxY)
	_, originY := view.Origin()
	_, cursorY := view.Cursor()

	view.Clear()
	runner.mainControlLines = 0
	runner.highlightedTexter = nil

	if g.CurrentView() == nil {
		g.SetCurrentView("mainControls")
	}
	view.Highlight = g.CurrentView() == view

	runner.rootTexter.TextMain(func(label, line string, texter ViewModelNodeTexter) {
		paddedLabel := fmt.Sprintf("%20s", label)
		fmt.Fprintf(view, "%s %v\n", paddedLabel[len(paddedLabel)-20:len(paddedLabel)], line)
		if (originY + cursorY) == runner.mainControlLines {
			runner.highlightedTexter = texter
		}
		runner.mainControlLines++
	})

	return nil
}

func (runner *appRunner) actMainControl(g *gocui.Gui, v *gocui.View) error {
	if runner.highlightedTexter != nil {
		runner.highlightedTexter.Act(runner)
	}

	return nil
}

func (runner *appRunner) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
		runner.layout(g)
	}
	return nil
}

func (runner *appRunner) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
		runner.layout(g)
	}
	return nil
}

func (runner *appRunner) ForList(controller ListDetailController, index int) DetailDataChangeCallback {
	runner.activeListDetailController = controller

	redrawDetails := func(view *gocui.View) {
		view.Clear()
		controller.WriteDetails(view)
	}

	runner.gui.Execute(func(*gocui.Gui) error {
		maxX, maxY := runner.gui.Size()
		visibleLines := maxY - 2

		view, _ := runner.gui.SetView("listDetail", maxX/2, 0, maxX-1, maxY-1)
		view.Highlight = true
		redrawDetails(view)

		if index < visibleLines {
			view.SetCursor(0, index)
		} else {
			halfHeight := visibleLines / 2
			view.SetOrigin(0, index-halfHeight)
			view.SetCursor(0, halfHeight)
		}
		runner.gui.SetCurrentView(view.Name())
		return nil
	})

	return func() {
		runner.gui.Execute(func(*gocui.Gui) error {
			view, _ := runner.gui.View("listDetail")
			redrawDetails(view)
			return nil
		})
	}
}

func (runner *appRunner) ForString(controller StringDetailController) DetailDataChangeCallback {
	runner.activeStringDetailController = controller

	redrawDetails := func(view *gocui.View) {
		view.Clear()
		controller.WriteDetails(view)
	}

	runner.gui.Execute(func(*gocui.Gui) error {
		maxX, maxY := runner.gui.Size()

		view, _ := runner.gui.SetView("stringDetail", maxX/2, 0, maxX-1, maxY-1)
		view.Editable = true
		view.Frame = true
		redrawDetails(view)

		runner.gui.SetCurrentView(view.Name())
		return nil
	})

	return func() {
		runner.gui.Execute(func(*gocui.Gui) error {
			view, _ := runner.gui.View("stringDetail")
			redrawDetails(view)
			return nil
		})
	}
}

func (runner *appRunner) restoreMainView() {
	runner.activeListDetailController = nil
	runner.activeStringDetailController = nil
	runner.gui.SetCurrentView("mainControls")
}

func (runner *appRunner) confirmListDetail(gui *gocui.Gui, view *gocui.View) error {
	_, originY := view.Origin()
	_, cursorY := view.Cursor()
	runner.activeListDetailController.Confirm(originY + cursorY)
	runner.restoreMainView()
	runner.gui.DeleteView(view.Name())
	return nil
}

func (runner *appRunner) cancelListDetail(gui *gocui.Gui, view *gocui.View) error {
	runner.activeListDetailController.Cancel()
	runner.restoreMainView()
	runner.gui.DeleteView(view.Name())
	return nil
}

func (runner *appRunner) confirmStringDetail(gui *gocui.Gui, view *gocui.View) error {
	text := strings.TrimSpace(view.Buffer())

	runner.activeStringDetailController.Confirm(text)
	runner.restoreMainView()
	runner.gui.DeleteView(view.Name())
	return nil
}

func (runner *appRunner) cancelStringDetail(gui *gocui.Gui, view *gocui.View) error {
	runner.activeStringDetailController.Cancel()
	runner.restoreMainView()
	runner.gui.DeleteView(view.Name())
	return nil
}
