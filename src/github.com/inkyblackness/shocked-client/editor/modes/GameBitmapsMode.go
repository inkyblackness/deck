package modes

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/inkyblackness/shocked-client/editor/cmd"
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

	resourceTypeLabel    *controls.Label
	resourceTypeBox      *controls.ComboBox
	resourceTypeItems    enumItems
	selectedResourceType dataModel.ResourceType

	languageLabel    *controls.Label
	languageBox      *controls.ComboBox
	languageItems    enumItems
	selectedLanguage dataModel.ResourceLanguage

	bitmapLabel      *controls.Label
	bitmapSelector   *controls.TextureSelector
	bitmapIDLabel    *controls.Label
	bitmapIDSlider   *controls.Slider
	selectedBitmapID int

	imageDisplayDrop *ui.Area
	imageDisplay     *controls.ImageDisplay
}

// NewGameBitmapsMode returns a new instance.
func NewGameBitmapsMode(context Context, parent *ui.Area) *GameBitmapsMode {
	mode := &GameBitmapsMode{
		context:          context,
		bitmapsAdapter:   context.ModelAdapter().BitmapsAdapter(),
		selectedBitmapID: -1}

	scaled := func(value float32) float32 {
		return value * context.ControlFactory().Scale()
	}

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

		{
			mode.resourceTypeLabel, mode.resourceTypeBox = panelBuilder.addComboProperty("Bitmap Type", mode.onResourceTypeChanged)
			mode.resourceTypeItems = []*enumItem{
				{uint32(dataModel.ResourceTypeMfdDataImages), "MFD Data Images"}}

			mode.resourceTypeBox.SetItems(mode.resourceTypeItems.forComboBox())
		}
		{
			mode.languageLabel, mode.languageBox = panelBuilder.addComboProperty("Language", mode.onLanguageChanged)
			mode.languageItems = []*enumItem{
				{uint32(dataModel.ResourceLanguageStandard), "STD"},
				{uint32(dataModel.ResourceLanguageFrench), "FRN"},
				{uint32(dataModel.ResourceLanguageGerman), "GER"}}
			mode.languageBox.SetItems(mode.languageItems.forComboBox())
		}
		{
			mode.bitmapIDLabel, mode.bitmapIDSlider = panelBuilder.addSliderProperty("Selected Bitmap ID",
				func(newValue int64) {
					mode.bitmapSelector.SetSelectedIndex(int(newValue))
					mode.onBitmapSelected(int(newValue))
				})
			mode.bitmapLabel, mode.bitmapSelector = panelBuilder.addTextureProperty("Selected Bitmap",
				mode.bitmapTextures, func(newValue int) {
					mode.bitmapIDSlider.SetValue(int64(newValue))
					mode.onBitmapSelected(newValue)
				})
		}

		mode.setState(dataModel.ResourceTypeMfdDataImages, dataModel.ResourceLanguageStandard, 0)
	}
	{
		padding := scaled(5.0)
		runningLeft := mode.propertiesArea.Right()
		displayWidth := scaled(256)

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
		key := dataModel.MakeLocalizedResourceKey(mode.selectedResourceType, mode.selectedLanguage, uint16(index))
		textures[index] = store.Texture(graphics.TextureKeyFromInt(key.ToInt()))
	}

	return textures
}

func (mode *GameBitmapsMode) imageProvider() (texture *graphics.BitmapTexture) {
	store := mode.context.ForGraphics().BitmapsStore()

	if mode.selectedBitmapID >= 0 {
		key := dataModel.MakeLocalizedResourceKey(mode.selectedResourceType, mode.selectedLanguage, uint16(mode.selectedBitmapID))
		texture = store.Texture(graphics.TextureKeyFromInt(key.ToInt()))
	}
	return
}

func (mode *GameBitmapsMode) onResourceTypeChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.setState(dataModel.ResourceType(item.value), mode.selectedLanguage, 0)
}

func (mode *GameBitmapsMode) onLanguageChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.selectedLanguage = dataModel.ResourceLanguage(item.value)
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
			mode.importImage(img)
		}
	}
	return
}

func (mode *GameBitmapsMode) importImage(img image.Image) {
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

		mode.requestBitmapChange(&rawBitmap)
	}
}

func (mode *GameBitmapsMode) requestBitmapChange(newBitmap *dataModel.RawBitmap) {
	restoreState := mode.stateSnapshot()
	key := dataModel.MakeLocalizedResourceKey(mode.selectedResourceType, mode.selectedLanguage, uint16(mode.selectedBitmapID))
	mode.context.Perform(&cmd.SetBitmapCommand{
		Setter: func(bmp *dataModel.RawBitmap) error {
			restoreState()
			mode.bitmapsAdapter.RequestBitmapChange(key, bmp)
			return nil
		},
		NewValue: newBitmap,
		OldValue: mode.bitmapsAdapter.Bitmap(key)})
}

func (mode *GameBitmapsMode) stateSnapshot() func() {
	currentType := mode.selectedResourceType
	currentLanguage := mode.selectedLanguage
	currentID := mode.selectedBitmapID
	return func() {
		mode.setState(currentType, currentLanguage, currentID)
	}
}

func (mode *GameBitmapsMode) setState(resourceType dataModel.ResourceType, language dataModel.ResourceLanguage, id int) {
	{
		mode.selectedResourceType = resourceType
		for _, item := range mode.resourceTypeItems {
			if item.value == uint32(resourceType) {
				mode.resourceTypeBox.SetSelectedItem(item)
			}
		}
		mode.bitmapIDSlider.SetRange(0, int64(bitmapCount[mode.selectedResourceType]-1))
	}
	{
		mode.selectedLanguage = language
		for _, item := range mode.languageItems {
			if item.value == uint32(language) {
				mode.languageBox.SetSelectedItem(item)
			}
		}
	}
	mode.selectedBitmapID = id
	mode.bitmapIDSlider.SetValue(int64(mode.selectedBitmapID))
	mode.bitmapSelector.SetSelectedIndex(mode.selectedBitmapID)
	mode.bitmapSelector.DisplaySelectedIndex()
}
