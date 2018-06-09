package controls

import (
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"
)

// SliderBuilder is a builder for Slider instances.
type SliderBuilder struct {
	areaBuilder  *ui.AreaBuilder
	rectRenderer *graphics.RectangleRenderer
	labelBuilder *LabelBuilder

	sliderChangeHandler SliderChangeHandler

	valueMin       int64
	valueMax       int64
	invertedScroll bool
}

// NewSliderBuilder returns a new SliderBuilder instance.
func NewSliderBuilder(labelBuilder *LabelBuilder, rectRenderer *graphics.RectangleRenderer) *SliderBuilder {
	builder := &SliderBuilder{
		areaBuilder:         ui.NewAreaBuilder(),
		rectRenderer:        rectRenderer,
		labelBuilder:        labelBuilder,
		sliderChangeHandler: func(int64) {}}

	return builder
}

// Build creates a new Slider instance from the current parameters.
func (builder *SliderBuilder) Build() *Slider {
	slider := &Slider{
		rectRenderer:        builder.rectRenderer,
		sliderChangeHandler: builder.sliderChangeHandler,
		formatter:           DefaultSliderValueFormatter,
		valueMin:            builder.valueMin,
		valueMax:            builder.valueMax,
		invertedScroll:      builder.invertedScroll,
		valueUndefined:      true,
	}

	builder.areaBuilder.OnRender(slider.onRender)
	builder.areaBuilder.OnEvent(events.MouseButtonDownEventType, slider.onMouseButtonDown)
	builder.areaBuilder.OnEvent(events.MouseButtonUpEventType, slider.onMouseButtonUp)
	builder.areaBuilder.OnEvent(events.MouseMoveEventType, slider.onMouseMove)
	builder.areaBuilder.OnEvent(events.MouseScrollEventType, slider.onMouseScroll)
	builder.areaBuilder.OnEvent(events.MouseButtonClickedEventType, ui.SilentConsumer)
	slider.area = builder.areaBuilder.Build()

	builder.labelBuilder.SetParent(slider.area)
	builder.labelBuilder.SetLeft(ui.NewOffsetAnchor(slider.area.Left(), 4))
	builder.labelBuilder.SetTop(ui.NewOffsetAnchor(slider.area.Top(), 0))
	builder.labelBuilder.SetRight(ui.NewOffsetAnchor(slider.area.Right(), -4))
	builder.labelBuilder.SetBottom(ui.NewOffsetAnchor(slider.area.Bottom(), 0))
	builder.labelBuilder.AlignedHorizontallyBy(LeftAligner)

	slider.valueLabel = builder.labelBuilder.Build()

	return slider
}

// SetParent sets the parent area.
func (builder *SliderBuilder) SetParent(parent *ui.Area) *SliderBuilder {
	builder.areaBuilder.SetParent(parent)
	return builder
}

// SetLeft sets the left anchor. Default: ZeroAnchor
func (builder *SliderBuilder) SetLeft(value ui.Anchor) *SliderBuilder {
	builder.areaBuilder.SetLeft(value)
	return builder
}

// SetTop sets the top anchor. Default: ZeroAnchor
func (builder *SliderBuilder) SetTop(value ui.Anchor) *SliderBuilder {
	builder.areaBuilder.SetTop(value)
	return builder
}

// SetRight sets the right anchor. Default: ZeroAnchor
func (builder *SliderBuilder) SetRight(value ui.Anchor) *SliderBuilder {
	builder.areaBuilder.SetRight(value)
	return builder
}

// SetBottom sets the bottom anchor. Default: ZeroAnchor
func (builder *SliderBuilder) SetBottom(value ui.Anchor) *SliderBuilder {
	builder.areaBuilder.SetBottom(value)
	return builder
}

// WithSliderChangeHandler sets the handler for a value change.
func (builder *SliderBuilder) WithSliderChangeHandler(handler SliderChangeHandler) *SliderBuilder {
	builder.sliderChangeHandler = handler
	return builder
}

// WithRange sets the allowed range of the slider.
func (builder *SliderBuilder) WithRange(valueMin, valueMax int64) *SliderBuilder {
	builder.valueMin = valueMin
	builder.valueMax = valueMax
	return builder
}

// WithInvertedScroll sets whether scrolling up/down should be inverted.
func (builder *SliderBuilder) WithInvertedScroll(value bool) *SliderBuilder {
	builder.invertedScroll = value
	return builder
}
