package controls

import (
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"
)

// ImageDisplayBuilder creates new instances of an ImageDisplay.
type ImageDisplayBuilder struct {
	areaBuilder *ui.AreaBuilder

	textureRenderer graphics.TextureRenderer
	provider        ImageProvider
}

// NewImageDisplayBuilder returns a new instance of an ImageDisplayBuilder.
func NewImageDisplayBuilder(textureRenderer graphics.TextureRenderer) *ImageDisplayBuilder {
	builder := &ImageDisplayBuilder{
		areaBuilder:     ui.NewAreaBuilder(),
		textureRenderer: textureRenderer,
		provider:        func() *graphics.BitmapTexture { return nil }}

	return builder
}

// Build creates a new ImageDisplay instance from the current parameters.
func (builder *ImageDisplayBuilder) Build() *ImageDisplay {
	control := &ImageDisplay{
		textureRenderer: builder.textureRenderer,
		provider:        builder.provider}

	builder.areaBuilder.OnRender(control.onRender)
	builder.areaBuilder.OnEvent(events.MouseMoveEventType, ui.SilentConsumer)
	builder.areaBuilder.OnEvent(events.MouseButtonDownEventType, ui.SilentConsumer)
	builder.areaBuilder.OnEvent(events.MouseButtonUpEventType, ui.SilentConsumer)
	control.area = builder.areaBuilder.Build()

	return control
}

// SetParent sets the parent area.
func (builder *ImageDisplayBuilder) SetParent(parent *ui.Area) *ImageDisplayBuilder {
	builder.areaBuilder.SetParent(parent)
	return builder
}

// SetLeft sets the left anchor. Default: ZeroAnchor
func (builder *ImageDisplayBuilder) SetLeft(value ui.Anchor) *ImageDisplayBuilder {
	builder.areaBuilder.SetLeft(value)
	return builder
}

// SetTop sets the top anchor. Default: ZeroAnchor
func (builder *ImageDisplayBuilder) SetTop(value ui.Anchor) *ImageDisplayBuilder {
	builder.areaBuilder.SetTop(value)
	return builder
}

// SetRight sets the right anchor. Default: ZeroAnchor
func (builder *ImageDisplayBuilder) SetRight(value ui.Anchor) *ImageDisplayBuilder {
	builder.areaBuilder.SetRight(value)
	return builder
}

// SetBottom sets the bottom anchor. Default: ZeroAnchor
func (builder *ImageDisplayBuilder) SetBottom(value ui.Anchor) *ImageDisplayBuilder {
	builder.areaBuilder.SetBottom(value)
	return builder
}

// WithProvider registers the image source
func (builder *ImageDisplayBuilder) WithProvider(provider ImageProvider) *ImageDisplayBuilder {
	builder.provider = provider
	return builder
}
