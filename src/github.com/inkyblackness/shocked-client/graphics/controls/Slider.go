package controls

import (
	"fmt"

	"github.com/inkyblackness/shocked-client/env"
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"
)

// SliderChangeHandler is a callback for notifying the current value.
type SliderChangeHandler func(value int64)

// Slider is a control for selecting a numerical value with a slider.
type Slider struct {
	area         *ui.Area
	rectRenderer *graphics.RectangleRenderer

	valueLabel *Label

	sliderChangeHandler SliderChangeHandler

	valueMin int64
	valueMax int64

	valueUndefined bool
	value          int64
}

// Dispose releases all resources and removes the area from the tree.
func (slider *Slider) Dispose() {
	slider.valueLabel.Dispose()
	slider.area.Remove()
}

// SetRange sets the minimum and maximum of valid values.
func (slider *Slider) SetRange(min, max int64) {
	slider.valueMin, slider.valueMax = min, max
}

// SetValueUndefined clears the current value.
func (slider *Slider) SetValueUndefined() {
	slider.valueUndefined = true
	slider.value = 0
	slider.valueLabel.SetText("")
}

// SetValue updates the current value.
func (slider *Slider) SetValue(value int64) {
	slider.valueUndefined = false
	slider.value = value
	slider.valueLabel.SetText(fmt.Sprintf("%v", value))
}

func (slider *Slider) onRender(area *ui.Area) {
	withinLimits := (slider.value >= slider.valueMin) && (slider.value <= slider.valueMax)
	areaLeft := area.Left().Value()
	areaTop := area.Top().Value()
	areaRight := area.Right().Value()
	areaBottom := area.Bottom().Value()

	if slider.valueUndefined || withinLimits {
		slider.rectRenderer.Fill(areaLeft, areaTop, areaRight, areaBottom, graphics.RGBA(0.31, 0.56, 0.34, 0.8))
	} else if !withinLimits {
		slider.rectRenderer.Fill(areaLeft, areaTop, areaRight, areaBottom, graphics.RGBA(0.56, 0.0, 0.34, 0.8))
	}

	if !slider.valueUndefined && withinLimits {
		sliderCenter := areaLeft + (float32(slider.value-slider.valueMin)/float32(slider.valueMax-slider.valueMin))*(areaRight-areaLeft)
		if (sliderCenter - 1) >= areaLeft {
			slider.rectRenderer.Fill(sliderCenter-1, areaTop, sliderCenter, areaBottom, graphics.RGBA(1.0, 0.0, 0.34, 0.5))
		}
		if (sliderCenter + 1) < areaRight {
			slider.rectRenderer.Fill(sliderCenter+1, areaTop, sliderCenter+2, areaBottom, graphics.RGBA(1.0, 0.0, 0.34, 0.5))
		}
		slider.rectRenderer.Fill(sliderCenter, areaTop, sliderCenter+1, areaBottom, graphics.RGBA(1.0, 0.0, 0.34, 1.0))
	}
}

func (slider *Slider) onMouseButtonDown(area *ui.Area, event events.Event) bool {
	mouseEvent := event.(*events.MouseButtonEvent)
	if mouseEvent.AffectedButtons() == env.MousePrimary {
		area.RequestFocus()
		slider.updateValueOnMouseEvent(mouseEvent)
	}
	return true
}

func (slider *Slider) onMouseButtonUp(area *ui.Area, event events.Event) bool {
	mouseEvent := event.(*events.MouseButtonEvent)
	if slider.area.HasFocus() && (mouseEvent.AffectedButtons() == env.MousePrimary) {
		area.ReleaseFocus()
		slider.updateValueOnMouseEvent(mouseEvent)
		slider.onValueChange(slider.value)
	}
	return true
}

func (slider *Slider) onMouseMove(area *ui.Area, event events.Event) bool {
	mouseEvent := event.(*events.MouseMoveEvent)
	if area.HasFocus() && (mouseEvent.Buttons() == env.MousePrimary) {
		slider.updateValueOnMouseEvent(mouseEvent)
	}
	return true
}

func (slider *Slider) onMouseScroll(area *ui.Area, event events.Event) bool {
	mouseEvent := event.(*events.MouseScrollEvent)

	if !slider.valueUndefined {
		_, dy := mouseEvent.Deltas()

		if (dy < 0) && (slider.value > slider.valueMin) {
			slider.onValueChange(slider.value - 1)
		} else if (dy > 0) && (slider.value < slider.valueMax) {
			slider.onValueChange(slider.value + 1)
		}
	}

	return true
}

func (slider *Slider) updateValueOnMouseEvent(mouseEvent events.PositionalEvent) {
	areaLeft := slider.area.Left().Value()
	areaRight := slider.area.Right().Value()
	mouseX, _ := mouseEvent.Position()
	newValue := slider.value

	if mouseX <= areaLeft {
		newValue = slider.valueMin
	} else if mouseX >= areaRight {
		newValue = slider.valueMax
	} else {
		newValue = slider.valueMin + int64(((mouseX-areaLeft)/(areaRight-areaLeft))*float32(slider.valueMax-slider.valueMin))
	}
	slider.SetValue(newValue)
}

func (slider *Slider) onValueChange(newValue int64) {
	slider.SetValue(newValue)
	slider.sliderChangeHandler(newValue)
}
