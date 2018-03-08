package display

import (
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/editor/camera"
	"github.com/inkyblackness/shocked-client/editor/model"
	"github.com/inkyblackness/shocked-client/env"
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"
	dataModel "github.com/inkyblackness/shocked-model"
)

// MapDisplay is a display for a level map
type MapDisplay struct {
	context      Context
	levelAdapter *model.LevelAdapter

	area *ui.Area

	camera        *camera.LimitedCamera
	viewMatrix    mgl.Mat4
	renderContext *graphics.RenderContext

	paletteTexture *graphics.PaletteTexture

	highlighter *BasicHighlighter
	background  *GridRenderable
	mapGrid     *TileGridMapRenderable
	textures    *TileTextureMapRenderable
	colors      *TileColorMapRenderable
	slopeGrid   *TileSlopeMapRenderable
	objects     *PlacedIconsRenderable

	selectedTileAreas   []Area
	highlightedTileArea Area

	displayedObjectAreas  []Area
	displayedObjectIcons  []PlacedIcon
	selectedObjectAreas   []Area
	highlightedObjectArea Area
	highlightedObjectIcon PlacedIcon

	moveCapture func(pixelX, pixelY float32)
}

// NewMapDisplay returns a new instance.
func NewMapDisplay(context Context, parent *ui.Area, scale float32) *MapDisplay {
	tileBaseLength := fineCoordinatesPerTileSide
	tileBaseHalf := tileBaseLength / 2.0
	camLimit := tilesPerMapSide*tileBaseLength - tileBaseHalf
	zoomShift := scale - 1.0
	zoomLevelMin := float32(-5) + zoomShift
	zoomLevelMax := float32(1) + zoomShift

	display := &MapDisplay{
		context:      context,
		levelAdapter: context.ModelAdapter().ActiveLevel(),
		camera:       camera.NewLimited(zoomLevelMin, zoomLevelMax, -tileBaseHalf, camLimit),
		moveCapture:  func(float32, float32) {}}

	centerX, centerY := float32(tilesPerMapSide*tileBaseLength)/-2.0, float32(tilesPerMapSide*tileBaseLength)/-2.0
	display.camera.ZoomAt(-3+zoomShift, centerX, centerY)
	display.camera.MoveTo(centerX, centerY)

	{
		builder := ui.NewAreaBuilder()
		builder.SetParent(parent)
		builder.SetLeft(ui.NewOffsetAnchor(parent.Left(), 0))
		builder.SetTop(ui.NewOffsetAnchor(parent.Top(), 0))
		builder.SetRight(ui.NewOffsetAnchor(parent.Right(), 0))
		builder.SetBottom(ui.NewOffsetAnchor(parent.Bottom(), 0))
		builder.SetVisible(false)
		builder.OnRender(func(area *ui.Area) { display.render() })
		builder.OnEvent(events.MouseScrollEventType, display.onMouseScroll)
		builder.OnEvent(events.MouseMoveEventType, display.onMouseMove)
		builder.OnEvent(events.MouseButtonDownEventType, display.onMouseButtonDown)
		builder.OnEvent(events.MouseButtonUpEventType, display.onMouseButtonUp)
		display.area = builder.Build()
	}

	display.paletteTexture = context.ForGraphics().NewPaletteTexture(display.paletteEntry)
	display.context.ModelAdapter().OnGamePaletteChanged(func() {
		display.paletteTexture.Update()
	})

	display.renderContext = context.NewRenderContext(display.camera.ViewMatrix())
	display.highlighter = NewBasicHighlighter(display.renderContext)
	display.background = NewGridRenderable(display.renderContext)
	display.mapGrid = NewTileGridMapRenderable(display.renderContext)
	display.textures = NewTileTextureMapRenderable(display.renderContext, display.paletteTexture, func(index int) *graphics.BitmapTexture {
		id := display.levelAdapter.LevelTextureID(index)
		return display.context.ForGraphics().WorldTextureStore(dataModel.TextureLarge).Texture(graphics.TextureKeyFromInt(id))
	})
	display.colors = NewTileColorMapRenderable(display.renderContext)
	display.slopeGrid = NewTileSlopeMapRenderable(display.renderContext)
	display.objects = NewPlacedIconsRenderable(display.renderContext, display.paletteTexture)

	linkTileProperties := func(coord model.TileCoordinate) {
		tile := display.levelAdapter.TileMap().Tile(coord)
		tile.OnPropertiesChanged(func() {
			x, y := coord.XY()
			properties := tile.Properties()
			display.mapGrid.SetTile(x, y, properties)
			display.textures.SetTile(x, y, properties)
			display.colors.SetTile(x, y, properties)
			display.slopeGrid.SetTile(x, y, properties)
		})
	}

	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			linkTileProperties(model.TileCoordinateOf(x, y))
		}
	}

	return display
}

func (display *MapDisplay) paletteEntry(index int) (r, g, b, a byte) {
	pal := display.context.ModelAdapter().GamePalette()
	color := &pal[index]

	r = byte(color.Red)
	g = byte(color.Green)
	b = byte(color.Blue)
	if index > 0 {
		a = 0xFF
	}

	return
}

// SetVisible sets the display visibility state.
func (display *MapDisplay) SetVisible(visible bool) {
	display.area.SetVisible(visible)
}

// SetTextureIndexQuery sets which texture shall be shown.
func (display *MapDisplay) SetTextureIndexQuery(query TextureIndexQuery) {
	display.textures.SetTextureIndexQuery(query)
}

// SetHighlightedTile requests to highlight the identified tile.
func (display *MapDisplay) SetHighlightedTile(coord model.TileCoordinate) {
	tileX, tileY := coord.XY()
	display.highlightedTileArea = NewSimpleArea(float32(tileX<<8+128), float32(tileY<<8+128), 256.0, 256.0)
}

// ClearHighlightedTile requests to remove any tile highlight.
func (display *MapDisplay) ClearHighlightedTile() {
	display.highlightedTileArea = nil
}

// SetSelectedTiles requests to show the given set of tiles as selected.
func (display *MapDisplay) SetSelectedTiles(tiles []model.TileCoordinate) {
	display.selectedTileAreas = make([]Area, len(tiles))

	for index, coord := range tiles {
		tileX, tileY := coord.XY()
		display.selectedTileAreas[index] = NewSimpleArea(float32(tileX<<8+128), float32(tileY<<8+128), 256.0, 256.0)
	}
}

// SetDisplayedObjects requests to show the given set of objects.
func (display *MapDisplay) SetDisplayedObjects(objects []*model.LevelObject) {
	display.displayedObjectIcons = make([]PlacedIcon, len(objects))
	display.displayedObjectAreas = make([]Area, len(objects))

	for index, object := range objects {
		icon := display.iconForObject(object)
		display.displayedObjectAreas[index] = icon
		display.displayedObjectIcons[index] = icon
	}
}

// SetSelectedObjects requests to show the given set of objects as selected.
func (display *MapDisplay) SetSelectedObjects(objects []*model.LevelObject) {
	display.selectedObjectAreas = make([]Area, len(objects))

	for index, object := range objects {
		icon := display.iconForObject(object)
		display.selectedObjectAreas[index] = icon
	}
}

// SetHighlightedObject registers an object that shall be highlighted.
func (display *MapDisplay) SetHighlightedObject(object *model.LevelObject) {
	// TODO: In some future update (past 1.8), see if it still crashes if
	// these two assignments are done from a variable of type *referringPlacedIcon.
	if object != nil {
		icon := display.iconForObject(object)
		display.highlightedObjectArea = icon
		display.highlightedObjectIcon = icon
	} else {
		display.highlightedObjectArea = nil
		display.highlightedObjectIcon = nil
	}
}

// SetTileColoring sets the query function for coloring tiles.
func (display *MapDisplay) SetTileColoring(colorQuery ColorQuery) {
	display.colors.SetColorQuery(colorQuery)
}

func (display *MapDisplay) iconForObject(object *model.LevelObject) *referringPlacedIcon {
	return &referringPlacedIcon{
		center: func() (float32, float32) { return object.Center() },
		texture: func() *graphics.BitmapTexture {
			return display.context.ForGraphics().GameObjectIconsStore().Texture(graphics.TextureKeyFromInt(object.ID().ToInt()))
		}}
}

func (display *MapDisplay) render() {
	root := display.area.Root()
	display.camera.SetViewportSize(root.Right().Value(), root.Bottom().Value())
	display.background.Render()
	if !display.levelAdapter.IsCyberspace() {
		display.textures.Render()
	}
	display.colors.Render()
	display.slopeGrid.Render()
	display.highlighter.Render(display.selectedTileAreas, graphics.RGBA(0.0, 0.8, 0.2, 0.5))
	if display.highlightedTileArea != nil {
		display.highlighter.Render([]Area{display.highlightedTileArea}, graphics.RGBA(0.0, 0.2, 0.8, 0.3))
	}
	display.mapGrid.Render()
	display.highlighter.Render(display.displayedObjectAreas, graphics.RGBA(1.0, 1.0, 1.0, 0.3))
	display.highlighter.Render(display.selectedObjectAreas, graphics.RGBA(0.0, 0.8, 0.2, 0.5))
	if display.highlightedObjectArea != nil {
		display.highlighter.Render([]Area{display.highlightedObjectArea}, graphics.RGBA(0.0, 0.2, 0.8, 0.3))
	}
	display.objects.Render(display.displayedObjectIcons)
	if display.highlightedObjectIcon != nil {
		display.objects.Render([]PlacedIcon{display.highlightedObjectIcon})
	}
}

// WorldCoordinatesForPixel returns the world coordinates at the given pixel position.
func (display *MapDisplay) WorldCoordinatesForPixel(pixelX, pixelY float32) (x, y float32) {
	return display.unprojectPixel(pixelX, pixelY)
}

func (display *MapDisplay) unprojectPixel(pixelX, pixelY float32) (x, y float32) {
	pixelVec := mgl.Vec4{pixelX, pixelY, 0.0, 1.0}
	invertedView := display.camera.ViewMatrix().Inv()
	result := invertedView.Mul4x1(pixelVec)

	return result[0], result[1]
}

func (display *MapDisplay) onMouseScroll(area *ui.Area, event events.Event) bool {
	mouseEvent := event.(*events.MouseScrollEvent)
	mouseX, mouseY := mouseEvent.Position()
	worldX, worldY := display.unprojectPixel(mouseX, mouseY)
	_, dy := mouseEvent.Deltas()

	if dy > 0 {
		display.camera.ZoomAt(-0.5, worldX, worldY)
	}
	if dy < 0 {
		display.camera.ZoomAt(0.5, worldX, worldY)
	}

	return true
}

func (display *MapDisplay) onMouseMove(area *ui.Area, event events.Event) bool {
	mouseEvent := event.(*events.MouseMoveEvent)
	display.moveCapture(mouseEvent.Position())

	return true
}

func (display *MapDisplay) onMouseButtonDown(area *ui.Area, event events.Event) bool {
	mouseEvent := event.(*events.MouseButtonEvent)

	if mouseEvent.Buttons() == env.MousePrimary {
		lastPixelX, lastPixelY := mouseEvent.Position()

		display.area.RequestFocus()
		display.moveCapture = func(pixelX, pixelY float32) {
			lastWorldX, lastWorldY := display.unprojectPixel(lastPixelX, lastPixelY)
			worldX, worldY := display.unprojectPixel(pixelX, pixelY)

			display.camera.MoveBy(worldX-lastWorldX, worldY-lastWorldY)
			lastPixelX, lastPixelY = pixelX, pixelY
		}
	}

	return true
}

func (display *MapDisplay) onMouseButtonUp(area *ui.Area, event events.Event) bool {
	mouseEvent := event.(*events.MouseButtonEvent)

	if mouseEvent.AffectedButtons() == env.MousePrimary {
		if display.area.HasFocus() {
			display.area.ReleaseFocus()
		}
		display.moveCapture = func(float32, float32) {}
	}

	return true
}
