package controls

import (
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"
)

// LabelBuilder creates new label controls.
type LabelBuilder struct {
	areaBuilder *ui.AreaBuilder

	textPainter     graphics.TextPainter
	texturizer      BitmapTexturizer
	textureRenderer graphics.TextureRenderer

	fitToWidth        bool
	scale             float32
	horizontalAligner Aligner
	verticalAligner   Aligner
}

// NewLabelBuilder returns a new instance of a LabelBuilder.
func NewLabelBuilder(textPainter graphics.TextPainter, texturizer BitmapTexturizer,
	textureRenderer graphics.TextureRenderer) *LabelBuilder {
	builder := &LabelBuilder{
		areaBuilder:       ui.NewAreaBuilder(),
		textPainter:       textPainter,
		texturizer:        texturizer,
		textureRenderer:   textureRenderer,
		scale:             1.0,
		horizontalAligner: CenterAligner,
		verticalAligner:   CenterAligner}

	return builder
}

// Build creates a new Label instance from the current parameters
func (builder *LabelBuilder) Build() *Label {
	label := &Label{
		textPainter:       builder.textPainter,
		texturizer:        builder.texturizer,
		textureRenderer:   builder.textureRenderer,
		fitToWidth:        builder.fitToWidth,
		scale:             builder.scale,
		horizontalAligner: builder.horizontalAligner,
		verticalAligner:   builder.verticalAligner}

	builder.areaBuilder.OnRender(label.onRender)
	builder.areaBuilder.OnEvent(events.ClipboardCopyEventType, label.onClipboardCopy)
	builder.areaBuilder.OnEvent(events.ClipboardPasteEventType, label.onClipboardPaste)
	builder.areaBuilder.OnEvent(events.FileDropEventType, label.onFileDrop)
	label.area = builder.areaBuilder.Build()
	label.SetText("")

	return label
}

// SetParent sets the parent area.
func (builder *LabelBuilder) SetParent(parent *ui.Area) *LabelBuilder {
	builder.areaBuilder.SetParent(parent)
	return builder
}

// SetLeft sets the left anchor. Default: ZeroAnchor
func (builder *LabelBuilder) SetLeft(value ui.Anchor) *LabelBuilder {
	builder.areaBuilder.SetLeft(value)
	return builder
}

// SetTop sets the top anchor. Default: ZeroAnchor
func (builder *LabelBuilder) SetTop(value ui.Anchor) *LabelBuilder {
	builder.areaBuilder.SetTop(value)
	return builder
}

// SetRight sets the right anchor. Default: ZeroAnchor
func (builder *LabelBuilder) SetRight(value ui.Anchor) *LabelBuilder {
	builder.areaBuilder.SetRight(value)
	return builder
}

// SetBottom sets the bottom anchor. Default: ZeroAnchor
func (builder *LabelBuilder) SetBottom(value ui.Anchor) *LabelBuilder {
	builder.areaBuilder.SetBottom(value)
	return builder
}

// SetScale sets the scaling factor of the text. Default: 1.0
func (builder *LabelBuilder) SetScale(value float32) *LabelBuilder {
	builder.scale = value
	return builder
}

// AlignedHorizontallyBy sets the aligner for the horizontal axis. Default: Center.
func (builder *LabelBuilder) AlignedHorizontallyBy(aligner Aligner) *LabelBuilder {
	builder.horizontalAligner = aligner
	return builder
}

// AlignedVerticallyBy sets the aligner for the vertical axis. Default: Center.
func (builder *LabelBuilder) AlignedVerticallyBy(aligner Aligner) *LabelBuilder {
	builder.verticalAligner = aligner
	return builder
}

// SetFitToWidth has the label always fit its text into the width.
func (builder *LabelBuilder) SetFitToWidth() *LabelBuilder {
	builder.fitToWidth = true
	return builder
}
