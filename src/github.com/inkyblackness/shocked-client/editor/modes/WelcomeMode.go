package modes

import (
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/graphics/controls"
	"github.com/inkyblackness/shocked-client/ui"
)

// WelcomeMode is a simple mode greeting the user and giving initial help.
type WelcomeMode struct {
	context Context

	infoArea  *ui.Area
	infoLabel *controls.Label
}

// NewWelcomeMode returns a new instance.
func NewWelcomeMode(context Context, parent *ui.Area) *WelcomeMode {
	mode := &WelcomeMode{context: context}
	welcomeText := `Hello and welcome to the editor!

At the top left you have a drop-down list to switch to different editor modes.
Drop-down lists are those with the three dots (...) on the right.

This is currently the "Welcome" mode.
Click on the drop-down list to open it and select another mode to get started. 
`
	scaled := func(value float32) float32 {
		return value * context.ControlFactory().Scale()
	}

	{
		horizontalCenter := ui.NewRelativeAnchor(parent.Left(), parent.Right(), 0.5)
		verticalCenter := ui.NewRelativeAnchor(parent.Top(), parent.Bottom(), 0.5)
		builder := ui.NewAreaBuilder()
		builder.SetParent(parent)
		builder.SetLeft(ui.NewOffsetAnchor(horizontalCenter, scaled(-300)))
		builder.SetTop(ui.NewOffsetAnchor(verticalCenter, scaled(-100)))
		builder.SetRight(ui.NewOffsetAnchor(horizontalCenter, scaled(300)))
		builder.SetBottom(ui.NewOffsetAnchor(verticalCenter, scaled(100)))
		builder.SetVisible(false)
		builder.OnRender(func(area *ui.Area) {
			context.ForGraphics().RectangleRenderer().Fill(
				area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
				graphics.RGBA(0.7, 0.0, 0.7, 0.1))
		})
		mode.infoArea = builder.Build()
	}
	{
		builder := context.ControlFactory().ForLabel()
		builder.SetParent(mode.infoArea)
		builder.SetLeft(ui.NewOffsetAnchor(mode.infoArea.Left(), scaled(20)))
		builder.SetTop(ui.NewOffsetAnchor(mode.infoArea.Top(), scaled(20)))
		builder.SetRight(ui.NewOffsetAnchor(mode.infoArea.Right(), scaled(-20)))
		builder.SetBottom(ui.NewOffsetAnchor(mode.infoArea.Bottom(), scaled(-20)))
		builder.AlignedHorizontallyBy(controls.LeftAligner)
		builder.AlignedVerticallyBy(controls.LeftAligner)
		mode.infoLabel = builder.Build()
		mode.infoLabel.SetText(welcomeText)
	}

	return mode
}

// SetActive implements the Mode interface.
func (mode *WelcomeMode) SetActive(active bool) {
	mode.infoArea.SetVisible(active)
}
