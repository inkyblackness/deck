package controls

import (
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"
)

// ComboBoxBuilder is a builder for ComboBox instances.
type ComboBoxBuilder struct {
	areaBuilder  *ui.AreaBuilder
	rectRenderer *graphics.RectangleRenderer
	labelBuilder *LabelBuilder

	selectionChangeHandler SelectionChangeHandler

	items []ComboBoxItem
}

// NewComboBoxBuilder returns a new ComboBoxBuilder instance.
func NewComboBoxBuilder(labelBuilder *LabelBuilder, rectRenderer *graphics.RectangleRenderer) *ComboBoxBuilder {
	builder := &ComboBoxBuilder{
		areaBuilder:            ui.NewAreaBuilder(),
		rectRenderer:           rectRenderer,
		labelBuilder:           labelBuilder,
		selectionChangeHandler: func(ComboBoxItem) {}}

	return builder
}

// Build creates a new ComboBox instance from the current parameters.
func (builder *ComboBoxBuilder) Build() *ComboBox {
	box := &ComboBox{
		labelBuilder:           builder.labelBuilder,
		rectRenderer:           builder.rectRenderer,
		selectionChangeHandler: builder.selectionChangeHandler,
		items: builder.items}

	builder.areaBuilder.OnRender(box.onRender)
	builder.areaBuilder.OnEvent(events.MouseButtonDownEventType, box.onMouseDown)
	builder.areaBuilder.OnEvent(events.MouseButtonUpEventType, ui.SilentConsumer)
	builder.areaBuilder.OnEvent(events.MouseButtonClickedEventType, ui.SilentConsumer)
	builder.areaBuilder.OnEvent(events.MouseScrollEventType, ui.SilentConsumer)
	box.area = builder.areaBuilder.Build()

	builder.labelBuilder.SetParent(box.area)
	builder.labelBuilder.SetTop(ui.NewOffsetAnchor(box.area.Top(), 0))
	builder.labelBuilder.SetBottom(ui.NewOffsetAnchor(box.area.Bottom(), 0))

	hintRight := ui.NewOffsetAnchor(box.area.Right(), 0)
	hintLeft := ui.NewResolvingAnchor(func() ui.Anchor {
		return ui.NewOffsetAnchor(hintRight, box.area.Top().Value()-box.area.Bottom().Value())
	})
	builder.labelBuilder.SetLeft(hintLeft)
	builder.labelBuilder.SetRight(hintRight)
	box.hintLabel = builder.labelBuilder.Build()
	box.hintLabel.SetText("...")

	builder.labelBuilder.SetLeft(ui.NewOffsetAnchor(box.area.Left(), 4))
	builder.labelBuilder.SetRight(ui.NewOffsetAnchor(hintLeft, -4))
	builder.labelBuilder.AlignedHorizontallyBy(LeftAligner)
	box.selectedLabel = builder.labelBuilder.Build()

	return box
}

// SetParent sets the parent area.
func (builder *ComboBoxBuilder) SetParent(parent *ui.Area) *ComboBoxBuilder {
	builder.areaBuilder.SetParent(parent)
	return builder
}

// SetLeft sets the left anchor. Default: ZeroAnchor
func (builder *ComboBoxBuilder) SetLeft(value ui.Anchor) *ComboBoxBuilder {
	builder.areaBuilder.SetLeft(value)
	return builder
}

// SetTop sets the top anchor. Default: ZeroAnchor
func (builder *ComboBoxBuilder) SetTop(value ui.Anchor) *ComboBoxBuilder {
	builder.areaBuilder.SetTop(value)
	return builder
}

// SetRight sets the right anchor. Default: ZeroAnchor
func (builder *ComboBoxBuilder) SetRight(value ui.Anchor) *ComboBoxBuilder {
	builder.areaBuilder.SetRight(value)
	return builder
}

// SetBottom sets the bottom anchor. Default: ZeroAnchor
func (builder *ComboBoxBuilder) SetBottom(value ui.Anchor) *ComboBoxBuilder {
	builder.areaBuilder.SetBottom(value)
	return builder
}

// WithItems sets the list of contained items.
func (builder *ComboBoxBuilder) WithItems(items []ComboBoxItem) *ComboBoxBuilder {
	builder.items = make([]ComboBoxItem, len(items))
	copy(builder.items, items)
	return builder
}

// WithSelectionChangeHandler sets the handler for a selection change.
func (builder *ComboBoxBuilder) WithSelectionChangeHandler(handler SelectionChangeHandler) *ComboBoxBuilder {
	builder.selectionChangeHandler = handler
	return builder
}
