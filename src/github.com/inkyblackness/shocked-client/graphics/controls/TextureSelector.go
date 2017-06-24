package controls

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"
)

// TextureSelectionChangeHandler is the type of the callback for changed selections.
type TextureSelectionChangeHandler func(index int)

// TextureProvider returns the textures available to the selector.
type TextureProvider func() []*graphics.BitmapTexture

// TextureSelector is a graphical selection tool for textures.
type TextureSelector struct {
	area        *ui.Area
	displayArea *ui.Area
	listArea    *ui.Area

	rectangleRenderer *graphics.RectangleRenderer
	textureRenderer   graphics.TextureRenderer

	provider TextureProvider

	firstDisplayedIndex int

	selectedIndex          int
	selectionChangeHandler TextureSelectionChangeHandler
}

// Dispose releases all current resources.
func (selector *TextureSelector) Dispose() {
	selector.listArea.Remove()
	selector.displayArea.Remove()
	selector.area.Remove()
}

// SetSelectedIndex updates the currently selected item.
// The change handler will no be called.
func (selector *TextureSelector) SetSelectedIndex(index int) {
	selector.selectedIndex = index
}

// DisplaySelectedIndex sets the first displayed index to the currently selected item.
func (selector *TextureSelector) DisplaySelectedIndex() {
	if selector.selectedIndex >= 0 {
		selector.firstDisplayedIndex = selector.selectedIndex
	}
}

func (selector *TextureSelector) onRender(area *ui.Area) {
	areaTop := area.Top().Value()
	areaBottom := area.Bottom().Value()
	areaLeft := area.Left().Value()
	areaRight := area.Right().Value()
	areaHeight := areaBottom - areaTop
	padding := float32(4.0)
	cellSize := areaHeight
	iconSize := cellSize - (padding * 2)
	textures := selector.textureListForRender()
	availableTextures := len(textures)
	selectedTexture := textures[0]

	selector.rectangleRenderer.Fill(areaLeft, areaTop, areaRight, areaBottom, graphics.RGBA(0.0, 0.0, 0.0, 0.9))

	runningLeft := areaLeft
	for index := 0; (index < availableTextures) && (runningLeft < areaRight); index++ {
		texture := textures[index]
		if (selectedTexture != nil) && (texture == selectedTexture) {
			toRight := runningLeft + cellSize
			if toRight > areaRight {
				toRight = areaRight
			}
			selector.rectangleRenderer.Fill(runningLeft, areaTop, toRight, areaBottom, graphics.RGBA(0.31, 0.56, 0.34, 0.8))
		}
		if texture != nil {
			u, v := texture.UV()
			imageWidth, imageHeight := texture.Size()
			fromLeft := float32(0.0)
			fromTop := float32(0.0)
			fromRight := u
			fromBottom := v
			widthFitting := iconSize / imageWidth
			heightFitting := iconSize / imageHeight
			scale := float32(math.Min(float64(widthFitting), float64(heightFitting)))
			toLeft := runningLeft + padding + (iconSize-(imageWidth*scale))/2.0
			toTop := areaTop + padding + (iconSize-(imageHeight*scale))/2.0
			toRight := toLeft + iconSize
			if toRight > areaRight {
				fromRight -= (u / iconSize) * (toRight - areaRight)
				toRight = areaRight
			}

			modelMatrix := mgl.Ident4().Mul4(mgl.Translate3D(toLeft, toTop, 0.0)).Mul4(mgl.Scale3D(toRight-toLeft, imageHeight*scale, 1.0))
			selector.textureRenderer.Render(&modelMatrix, texture, graphics.RectByCoord(fromLeft, fromTop, fromRight, fromBottom))
		}
		runningLeft += cellSize
	}
}

func (selector *TextureSelector) textureListForRender() []*graphics.BitmapTexture {
	allTextures := selector.provider()
	totalCount := len(allTextures)
	result := make([]*graphics.BitmapTexture, totalCount+1)
	var selectedTexture *graphics.BitmapTexture

	copy(result[1:], allTextures)
	if (selector.selectedIndex >= 0) && (selector.selectedIndex < totalCount) {
		selectedTexture = allTextures[selector.selectedIndex]
	}
	if selector.firstDisplayedIndex >= totalCount {
		selector.firstDisplayedIndex = totalCount - 1
	}
	if selector.firstDisplayedIndex < 0 {
		selector.firstDisplayedIndex = 0
	}
	result[selector.firstDisplayedIndex] = selectedTexture

	return result[selector.firstDisplayedIndex:]
}

func (selector *TextureSelector) onListMouseScroll(area *ui.Area, event events.Event) bool {
	mouseEvent := event.(*events.MouseScrollEvent)
	_, dy := mouseEvent.Deltas()

	if dy < 0 {
		selector.firstDisplayedIndex--
	} else {
		selector.firstDisplayedIndex++
	}

	return true
}

func (selector *TextureSelector) onListMouseButtonClicked(area *ui.Area, event events.Event) bool {
	mouseEvent := event.(*events.MouseButtonEvent)
	mouseX, _ := mouseEvent.Position()
	areaLeft := area.Left().Value()
	cellSize := area.Bottom().Value() - area.Top().Value()
	relativeIndex := (mouseX - areaLeft) / cellSize

	selector.selectedIndex = selector.firstDisplayedIndex + int(relativeIndex)
	selector.selectionChangeHandler(selector.selectedIndex)

	return true
}

func (selector *TextureSelector) onDisplayMouseButtonClicked(area *ui.Area, event events.Event) bool {
	selector.DisplaySelectedIndex()
	return true
}
