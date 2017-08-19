package modes

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/inkyblackness/shocked-client/editor/model"
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/graphics/controls"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"

	dataModel "github.com/inkyblackness/shocked-model"
)

var bitmapCount = map[dataModel.ResourceType]int{
	dataModel.ResourceTypeMfdDataImages: 64}

// GameBitmapsMode is a mode for arbitrary game bitmaps.
type GameBitmapsMode struct {
	context        Context
	bitmapsAdapter *model.BitmapsAdapter

	area           *ui.Area
	propertiesArea *ui.Area

	typeLabel            *controls.Label
	typeBox              *controls.ComboBox
	selectedResourceType dataModel.ResourceType

	languageLabel *controls.Label
	languageBox   *controls.ComboBox
	language      dataModel.ResourceLanguage

	selectedBitmapLabel    *controls.Label
	selectedBitmapSelector *controls.TextureSelector
	selectedBitmapIDLabel  *controls.Label
	selectedBitmapIDSlider *controls.Slider
	selectedBitmapID       int

	imageDisplayDrop *ui.Area
	imageDisplay     *controls.ImageDisplay
}

// NewGameBitmapsMode returns a new instance.
func NewGameBitmapsMode(context Context, parent *ui.Area) *GameBitmapsMode {
	mode := &GameBitmapsMode{
		context:          context,
		bitmapsAdapter:   context.ModelAdapter().BitmapsAdapter(),
		selectedBitmapID: -1}

	{
		builder := ui.NewAreaBuilder()
		builder.SetParent(parent)
		builder.SetLeft(ui.NewOffsetAnchor(parent.Left(), 0))
		builder.SetTop(ui.NewOffsetAnchor(parent.Top(), 0))
		builder.SetRight(ui.NewOffsetAnchor(parent.Right(), 0))
		builder.SetBottom(ui.NewOffsetAnchor(parent.Bottom(), 0))
		builder.SetVisible(false)
		mode.area = builder.Build()
	}
	{
		builder := ui.NewAreaBuilder()
		builder.SetParent(mode.area)
		builder.SetLeft(ui.NewOffsetAnchor(parent.Left(), 0))
		builder.SetTop(ui.NewOffsetAnchor(parent.Top(), 0))
		builder.SetRight(ui.NewRelativeAnchor(parent.Left(), parent.Right(), 0.5))
		builder.SetBottom(ui.NewOffsetAnchor(parent.Bottom(), 0))
		builder.SetVisible(true)
		builder.OnRender(func(area *ui.Area) {
			context.ForGraphics().RectangleRenderer().Fill(
				area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
				graphics.RGBA(0.7, 0.0, 0.7, 0.3))
		})
		builder.OnEvent(events.MouseMoveEventType, ui.SilentConsumer)
		builder.OnEvent(events.MouseButtonUpEventType, ui.SilentConsumer)
		builder.OnEvent(events.MouseButtonDownEventType, ui.SilentConsumer)
		builder.OnEvent(events.MouseButtonClickedEventType, ui.SilentConsumer)
		builder.OnEvent(events.MouseScrollEventType, ui.SilentConsumer)
		mode.propertiesArea = builder.Build()
	}
	{
		panelBuilder := newControlPanelBuilder(mode.propertiesArea, context.ControlFactory())
		var initialTypeItem controls.ComboBoxItem
		var initialLanguageItem controls.ComboBoxItem

		{
			mode.typeLabel, mode.typeBox = panelBuilder.addComboProperty("Bitmap Type", mode.onBitmapTypeChanged)
			items := []controls.ComboBoxItem{
				&enumItem{uint32(dataModel.ResourceTypeMfdDataImages), "MFD Data Images"}}

			mode.typeBox.SetItems(items)
			initialTypeItem = items[0]
		}
		{
			mode.languageLabel, mode.languageBox = panelBuilder.addComboProperty("Language", mode.onLanguageChanged)
			items := []controls.ComboBoxItem{
				&enumItem{uint32(dataModel.ResourceLanguageStandard), "STD"},
				&enumItem{uint32(dataModel.ResourceLanguageFrench), "FRN"},
				&enumItem{uint32(dataModel.ResourceLanguageGerman), "GER"}}
			mode.languageBox.SetItems(items)
			initialLanguageItem = items[0]
		}
		{
			mode.selectedBitmapIDLabel, mode.selectedBitmapIDSlider = panelBuilder.addSliderProperty("Selected Bitmap ID",
				func(newValue int64) {
					mode.selectedBitmapSelector.SetSelectedIndex(int(newValue))
					mode.onBitmapSelected(int(newValue))
				})
			mode.selectedBitmapLabel, mode.selectedBitmapSelector = panelBuilder.addTextureProperty("Selected Bitmap",
				mode.bitmapTextures, func(newValue int) {
					mode.selectedBitmapIDSlider.SetValue(int64(newValue))
					mode.onBitmapSelected(newValue)
				})
		}
		mode.languageBox.SetSelectedItem(initialLanguageItem)
		mode.onLanguageChanged(initialLanguageItem)
		mode.typeBox.SetSelectedItem(initialTypeItem)
		mode.onBitmapTypeChanged(initialTypeItem)
	}
	{
		padding := float32(5.0)
		runningLeft := mode.propertiesArea.Right()
		displayWidth := float32(256)

		{
			dropBuilder := ui.NewAreaBuilder()
			displayBuilder := mode.context.ControlFactory().ForImageDisplay()
			left := ui.NewOffsetAnchor(runningLeft, padding)
			right := ui.NewOffsetAnchor(left, displayWidth)
			top := ui.NewOffsetAnchor(mode.area.Top(), padding)

			dropBuilder.SetParent(mode.area)
			displayBuilder.SetParent(mode.area)
			dropBuilder.SetLeft(left)
			displayBuilder.SetLeft(left)
			dropBuilder.SetRight(right)
			displayBuilder.SetRight(right)
			dropBuilder.SetTop(top)
			displayBuilder.SetTop(top)
			dropBuilder.SetBottom(ui.NewOffsetAnchor(top, displayWidth))
			displayBuilder.SetBottom(ui.NewOffsetAnchor(top, displayWidth))
			dropBuilder.OnEvent(events.FileDropEventType, mode.bitmapDropHandler)
			displayBuilder.WithProvider(mode.imageProvider)
			mode.imageDisplayDrop = dropBuilder.Build()
			mode.imageDisplay = displayBuilder.Build()
			runningLeft = right
		}
	}

	return mode
}

// SetActive implements the Mode interface.
func (mode *GameBitmapsMode) SetActive(active bool) {
	mode.area.SetVisible(active)
}

func (mode *GameBitmapsMode) bitmapTextures() []*graphics.BitmapTexture {
	textureCount := bitmapCount[mode.selectedResourceType]
	textures := make([]*graphics.BitmapTexture, textureCount)
	store := mode.context.ForGraphics().BitmapsStore()

	for index := 0; index < textureCount; index++ {
		key := dataModel.MakeLocalizedResourceKey(mode.selectedResourceType, mode.language, uint16(index))
		textures[index] = store.Texture(graphics.TextureKeyFromInt(key.ToInt()))
	}

	return textures
}

func (mode *GameBitmapsMode) imageProvider() (texture *graphics.BitmapTexture) {
	store := mode.context.ForGraphics().BitmapsStore()

	if mode.selectedBitmapID >= 0 {
		key := dataModel.MakeLocalizedResourceKey(mode.selectedResourceType, mode.language, uint16(mode.selectedBitmapID))
		texture = store.Texture(graphics.TextureKeyFromInt(key.ToInt()))
	}
	return
}

func (mode *GameBitmapsMode) onBitmapTypeChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.selectedResourceType = dataModel.ResourceType(item.value)

	mode.onBitmapSelected(0)
	mode.selectedBitmapIDSlider.SetRange(0, int64(bitmapCount[mode.selectedResourceType]-1))
	mode.selectedBitmapIDSlider.SetValue(0)
	mode.selectedBitmapSelector.SetSelectedIndex(0)
}

func (mode *GameBitmapsMode) onLanguageChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.language = dataModel.ResourceLanguage(item.value)
}

func (mode *GameBitmapsMode) onBitmapSelected(id int) {
	mode.selectedBitmapID = id
}

func (mode *GameBitmapsMode) bitmapDropHandler(area *ui.Area, event events.Event) (result bool) {
	fileDropEvent := event.(*events.FileDropEvent)
	filePaths := fileDropEvent.FilePaths()
	if len(filePaths) == 1 {
		file, err := os.Open(filePaths[0])
		var img image.Image

		if err == nil {
			defer file.Close()
			img, _, err = image.Decode(file)
			if err != nil {
				mode.context.ModelAdapter().SetMessage(fmt.Sprintf("File <%v> has unknown image format", filePaths[0]))
			}
		} else {
			mode.context.ModelAdapter().SetMessage(fmt.Sprintf("Could not open file <%v>", filePaths[0]))
		}
		if err == nil {
			mode.setBitmap(img)
		}
	}
	return
}

func (mode *GameBitmapsMode) setBitmap(img image.Image) {
	if mode.selectedBitmapID >= 0 {
		rawPalette := mode.context.ModelAdapter().GamePalette()
		palette := make([]color.Color, len(rawPalette))
		for index, clr := range rawPalette {
			palette[index] = clr
		}
		bitmapper := graphics.NewStandardBitmapper(palette)
		gfxBitmap := bitmapper.Map(img)
		var rawBitmap dataModel.RawBitmap

		rawBitmap.Width = gfxBitmap.Width
		rawBitmap.Height = gfxBitmap.Height
		rawBitmap.Pixels = base64.StdEncoding.EncodeToString(gfxBitmap.Pixels)

		key := dataModel.MakeLocalizedResourceKey(mode.selectedResourceType, mode.language, uint16(mode.selectedBitmapID))
		mode.bitmapsAdapter.RequestBitmapChange(key, &rawBitmap)
	}
}
