package controls

import (
	"github.com/inkyblackness/shocked-client/env"
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"
)

// ActionHandler is the callback for a firing button (actionable).
type ActionHandler func()

// TextButton is a button with a text label on it.
type TextButton struct {
	area         *ui.Area
	rectRenderer *graphics.RectangleRenderer

	label     *Label
	labelLeft ui.Anchor
	labelTop  ui.Anchor

	actionHandler ActionHandler

	idleColor     graphics.Color
	preparedColor graphics.Color

	prepared bool
	color    graphics.Color
}

// Dispose releases all resources.
func (button *TextButton) Dispose() {
	button.label.Dispose()
	button.area.Remove()
}

func (button *TextButton) onRender(area *ui.Area) {
	button.rectRenderer.Fill(area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(), button.color)
}

func (button *TextButton) onMouseDown(area *ui.Area, event events.Event) (consumed bool) {
	mouseEvent := event.(*events.MouseButtonEvent)

	if mouseEvent.Buttons() == env.MousePrimary {
		area.RequestFocus()
		button.prepare()
		consumed = true
	}

	return
}

func (button *TextButton) onMouseUp(area *ui.Area, event events.Event) (consumed bool) {
	mouseEvent := event.(*events.MouseButtonEvent)

	if button.area.HasFocus() && mouseEvent.AffectedButtons() == env.MousePrimary {
		area.ReleaseFocus()
		button.unprepare()
		if button.contains(mouseEvent) {
			button.callHandler()
		}
		consumed = true
	}

	return
}

func (button *TextButton) prepare() {
	if !button.prepared {
		button.color = button.preparedColor
		button.labelLeft.RequestValue(button.labelLeft.Value() + 5)
		button.labelTop.RequestValue(button.labelTop.Value() + 2)
		button.prepared = true
	}
}

func (button *TextButton) unprepare() {
	if button.prepared {
		button.color = button.idleColor
		button.labelLeft.RequestValue(button.labelLeft.Value() - 5)
		button.labelTop.RequestValue(button.labelTop.Value() - 2)
		button.prepared = false
	}
}

func (button *TextButton) contains(event events.PositionalEvent) bool {
	x, y := event.Position()

	return (x >= button.area.Left().Value()) && (x < button.area.Right().Value()) &&
		(y >= button.area.Top().Value()) && (y < button.area.Bottom().Value())
}

func (button *TextButton) callHandler() {
	button.actionHandler()
}
