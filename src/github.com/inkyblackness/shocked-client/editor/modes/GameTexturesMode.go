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

// GameTexturesMode is a mode for game textures.
type GameTexturesMode struct {
	context        Context
	textureAdapter *model.TextureAdapter

	area           *ui.Area
	propertiesArea *ui.Area

	selectedTextureLabel    *controls.Label
	selectedTextureSelector *controls.TextureSelector
	selectedTextureIDLabel  *controls.Label
	selectedTextureIDSlider *controls.Slider
	selectedTextureID       int

	propertiesHeader *controls.Label

	climbableLabel *controls.Label
	climbableBox   *controls.ComboBox
	climbableItems []controls.ComboBoxItem

	transparencyControlLabel *controls.Label
	transparencyControlBox   *controls.ComboBox
	transparencyControlItems []controls.ComboBoxItem

	animationGroupLabel  *controls.Label
	animationGroupSlider *controls.Slider
	animationIndexLabel  *controls.Label
	animationIndexSlider *controls.Slider

	languageLabel *controls.Label
	languageBox   *controls.ComboBox
	languageIndex int
	nameTitle     *controls.Label
	nameValue     *controls.Label
	useTextTitle  *controls.Label
	useTextValue  *controls.Label

	imageDisplayDrops map[dataModel.TextureSize]*ui.Area
	imageDisplays     map[dataModel.TextureSize]*controls.ImageDisplay
}

// NewGameTexturesMode returns a new instance.
func NewGameTexturesMode(context Context, parent *ui.Area) *GameTexturesMode {
	mode := &GameTexturesMode{
		context:           context,
		textureAdapter:    context.ModelAdapter().TextureAdapter(),
		selectedTextureID: -1,
		imageDisplayDrops: make(map[dataModel.TextureSize]*ui.Area),
		imageDisplays:     make(map[dataModel.TextureSize]*controls.ImageDisplay)}

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
			mode.selectedTextureIDLabel, mode.selectedTextureIDSlider = panelBuilder.addSliderProperty("Selected Texture ID",
				func(newValue int64) {
					mode.selectedTextureSelector.SetSelectedIndex(int(newValue))
					mode.selectedTextureSelector.DisplaySelectedIndex()
					mode.onTextureSelected(int(newValue))
				})
			mode.selectedTextureLabel, mode.selectedTextureSelector = panelBuilder.addTextureProperty("Selected Texture",
				mode.worldTextures, func(newValue int) {
					mode.selectedTextureIDSlider.SetValue(int64(newValue))
					mode.onTextureSelected(newValue)
				})
		}
		mode.propertiesHeader = panelBuilder.addTitle("Properties")
		{
			mode.languageLabel, mode.languageBox = panelBuilder.addComboProperty("Language", mode.onLanguageChanged)
			items := []controls.ComboBoxItem{&enumItem{0, "STD"}, &enumItem{1, "FRN"}, &enumItem{2, "GER"}}
			mode.languageBox.SetItems(items)
			mode.languageBox.SetSelectedItem(items[0])

			mode.nameTitle, mode.nameValue = panelBuilder.addInfo("Name")
			mode.nameValue.AllowTextChange(mode.onNameChangeRequested)
			mode.useTextTitle, mode.useTextValue = panelBuilder.addInfo("Use Text")
			mode.useTextValue.AllowTextChange(mode.onUseTextChangeRequested)
		}
		{
			mode.climbableLabel, mode.climbableBox = panelBuilder.addComboProperty("Climbable", mode.onClimbableChanged)
			mode.climbableItems = []controls.ComboBoxItem{&enumItem{0, "No"}, &enumItem{1, "Yes"}}
			mode.climbableBox.SetItems(mode.climbableItems)
		}
		{
			mode.transparencyControlLabel, mode.transparencyControlBox = panelBuilder.addComboProperty("Transparency Control", mode.onTransparencyControlChanged)
			mode.transparencyControlItems = []controls.ComboBoxItem{&enumItem{0, "Opaque"}, &enumItem{1, "Space"}, &enumItem{2, "Transparent"}}
			mode.transparencyControlBox.SetItems(mode.transparencyControlItems)
		}
		{
			mode.animationGroupLabel, mode.animationGroupSlider = panelBuilder.addSliderProperty("Animation Group", mode.onAnimationGroupChanged)
			mode.animationGroupSlider.SetRange(0, 3)
			mode.animationIndexLabel, mode.animationIndexSlider = panelBuilder.addSliderProperty("Animation Index", mode.onAnimationIndexChanged)
			mode.animationIndexSlider.SetRange(0, 3)
		}
	}
	{
		padding := float32(5.0)
		runningLeft := mode.propertiesArea.Right()
		pixelSizes := map[dataModel.TextureSize]float32{
			dataModel.TextureLarge:  128,
			dataModel.TextureMedium: 64,
			dataModel.TextureSmall:  32,
			dataModel.TextureIcon:   16}

		for _, textureSize := range dataModel.TextureSizes() {
			dropBuilder := ui.NewAreaBuilder()
			displayBuilder := mode.context.ControlFactory().ForImageDisplay()
			left := ui.NewOffsetAnchor(runningLeft, padding)
			right := ui.NewOffsetAnchor(left, pixelSizes[textureSize])
			top := ui.NewOffsetAnchor(mode.area.Top(), padding)

			dropBuilder.SetParent(mode.area)
			displayBuilder.SetParent(mode.area)
			dropBuilder.SetLeft(left)
			displayBuilder.SetLeft(left)
			dropBuilder.SetRight(right)
			displayBuilder.SetRight(right)
			dropBuilder.SetTop(top)
			displayBuilder.SetTop(top)
			dropBuilder.SetBottom(ui.NewOffsetAnchor(top, pixelSizes[dataModel.TextureLarge]))
			displayBuilder.SetBottom(ui.NewOffsetAnchor(top, pixelSizes[dataModel.TextureLarge]))
			dropBuilder.OnEvent(events.FileDropEventType, mode.textureDropHandler(textureSize))
			displayBuilder.WithProvider(mode.imageProvider(textureSize))
			mode.imageDisplayDrops[textureSize] = dropBuilder.Build()
			mode.imageDisplays[textureSize] = displayBuilder.Build()
			runningLeft = right
		}
	}
	mode.textureAdapter.OnGameTexturesChanged(mode.onGameTexturesChanged)

	return mode
}

// SetActive implements the Mode interface.
func (mode *GameTexturesMode) SetActive(active bool) {
	mode.area.SetVisible(active)
}

func (mode *GameTexturesMode) worldTextures() []*graphics.BitmapTexture {
	textureCount := mode.context.ModelAdapter().TextureAdapter().WorldTextureCount()
	textures := make([]*graphics.BitmapTexture, textureCount)
	store := mode.context.ForGraphics().WorldTextureStore(dataModel.TextureLarge)

	for index := 0; index < textureCount; index++ {
		textures[index] = store.Texture(graphics.TextureKeyFromInt(index))
	}

	return textures
}

func (mode *GameTexturesMode) imageProvider(size dataModel.TextureSize) controls.ImageProvider {
	store := mode.context.ForGraphics().WorldTextureStore(size)
	return func() (texture *graphics.BitmapTexture) {
		if mode.selectedTextureID >= 0 {
			texture = store.Texture(graphics.TextureKeyFromInt(mode.selectedTextureID))
		}
		return
	}
}

func (mode *GameTexturesMode) onGameTexturesChanged() {
	textureCount := mode.textureAdapter.WorldTextureCount()

	mode.selectedTextureIDSlider.SetRange(0, int64(textureCount)-1)
	if mode.selectedTextureID >= 0 {
		mode.onTextureSelected(mode.selectedTextureID)
	}
}

func (mode *GameTexturesMode) onTextureSelected(id int) {
	texture := mode.textureAdapter.GameTexture(id)

	mode.selectedTextureID = id
	if texture.Climbable() {
		mode.climbableBox.SetSelectedItem(mode.climbableItems[1])
	} else {
		mode.climbableBox.SetSelectedItem(mode.climbableItems[0])
	}
	mode.transparencyControlBox.SetSelectedItem(mode.transparencyControlItems[texture.TransparencyControl()])
	mode.animationGroupSlider.SetValue(int64(texture.AnimationGroup()))
	mode.animationIndexSlider.SetValue(int64(texture.AnimationIndex()))
	mode.updateTextureText()
}

func (mode *GameTexturesMode) requestTexturePropertiesChange(modifier func(*dataModel.TextureProperties)) {
	if mode.selectedTextureID >= 0 {
		var properties dataModel.TextureProperties

		modifier(&properties)
		mode.textureAdapter.RequestTexturePropertiesChange(mode.selectedTextureID, &properties)
	}
}

func (mode *GameTexturesMode) onLanguageChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.languageIndex = int(item.value)
	mode.updateTextureText()
}

func (mode *GameTexturesMode) updateTextureText() {
	name := ""
	useText := ""

	if mode.selectedTextureID >= 0 {
		texture := mode.textureAdapter.GameTexture(mode.selectedTextureID)

		name = texture.Name(mode.languageIndex)
		useText = texture.UseText(mode.languageIndex)
	}
	mode.nameValue.SetText(name)
	mode.useTextValue.SetText(useText)
}

func (mode *GameTexturesMode) onNameChangeRequested(newValue string) {
	mode.requestTexturePropertiesChange(func(properties *dataModel.TextureProperties) {
		properties.Name[mode.languageIndex] = stringAsPointer(newValue)
	})
}

func (mode *GameTexturesMode) onUseTextChangeRequested(newValue string) {
	mode.requestTexturePropertiesChange(func(properties *dataModel.TextureProperties) {
		properties.CantBeUsed[mode.languageIndex] = stringAsPointer(newValue)
	})
}

func (mode *GameTexturesMode) onClimbableChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.requestTexturePropertiesChange(func(properties *dataModel.TextureProperties) {
		properties.Climbable = boolAsPointer(item.value != 0)
	})
}

func (mode *GameTexturesMode) onTransparencyControlChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.requestTexturePropertiesChange(func(properties *dataModel.TextureProperties) {
		properties.TransparencyControl = intAsPointer(int(item.value))
	})
}

func (mode *GameTexturesMode) onAnimationGroupChanged(newValue int64) {
	mode.requestTexturePropertiesChange(func(properties *dataModel.TextureProperties) {
		properties.AnimationGroup = intAsPointer(int(newValue))
	})
}

func (mode *GameTexturesMode) onAnimationIndexChanged(newValue int64) {
	mode.requestTexturePropertiesChange(func(properties *dataModel.TextureProperties) {
		properties.AnimationIndex = intAsPointer(int(newValue))
	})
}

func (mode *GameTexturesMode) textureDropHandler(textureSize dataModel.TextureSize) ui.EventHandler {
	return func(area *ui.Area, event events.Event) (result bool) {
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
				mode.setTextureBitmap(textureSize, img)
			}
		}
		return
	}
}

func (mode *GameTexturesMode) setTextureBitmap(textureSize dataModel.TextureSize, img image.Image) {
	if mode.selectedTextureID >= 0 {
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

		mode.textureAdapter.RequestTextureBitmapChange(mode.selectedTextureID, textureSize, &rawBitmap)
	}
}
