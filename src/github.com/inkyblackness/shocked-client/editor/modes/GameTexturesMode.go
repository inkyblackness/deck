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

// GameTexturesMode is a mode for game textures.
type GameTexturesMode struct {
	context        Context
	textureAdapter *model.TextureAdapter

	area           *ui.Area
	propertiesArea *ui.Area

	textureLabel      *controls.Label
	textureSelector   *controls.TextureSelector
	textureIDLabel    *controls.Label
	textureIDSlider   *controls.Slider
	selectedTextureID int

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

	languageLabel    *controls.Label
	languageBox      *controls.ComboBox
	languageItems    enumItems
	selectedLanguage dataModel.ResourceLanguage

	nameTitle    *controls.Label
	nameValue    *controls.Label
	useTextTitle *controls.Label
	useTextValue *controls.Label

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
			mode.textureIDLabel, mode.textureIDSlider = panelBuilder.addSliderProperty("Selected Texture ID",
				func(newValue int64) {
					mode.textureSelector.SetSelectedIndex(int(newValue))
					mode.textureSelector.DisplaySelectedIndex()
					mode.onTextureSelected(int(newValue))
				})
			mode.textureLabel, mode.textureSelector = panelBuilder.addTextureProperty("Selected Texture",
				mode.worldTextures, func(newValue int) {
					mode.textureIDSlider.SetValue(int64(newValue))
					mode.onTextureSelected(newValue)
				})
		}
		mode.propertiesHeader = panelBuilder.addTitle("Properties")
		{
			mode.languageLabel, mode.languageBox = panelBuilder.addComboProperty("Language", mode.onLanguageChanged)
			mode.languageItems = []*enumItem{
				{uint32(dataModel.ResourceLanguageStandard), "STD"},
				{uint32(dataModel.ResourceLanguageFrench), "FRN"},
				{uint32(dataModel.ResourceLanguageGerman), "GER"}}
			mode.languageBox.SetItems(mode.languageItems.forComboBox())

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
		mode.setState(dataModel.ResourceLanguageStandard, 0)
	}
	{
		padding := scaled(5.0)
		runningLeft := mode.propertiesArea.Right()
		pixelSizes := map[dataModel.TextureSize]float32{
			dataModel.TextureLarge:  scaled(128),
			dataModel.TextureMedium: scaled(64),
			dataModel.TextureSmall:  scaled(32),
			dataModel.TextureIcon:   scaled(16)}

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

	mode.textureIDSlider.SetRange(0, int64(textureCount)-1)
	mode.updateData()
}

func (mode *GameTexturesMode) onTextureSelected(id int) {
	mode.selectedTextureID = id
	mode.updateData()
}

func (mode *GameTexturesMode) onLanguageChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.selectedLanguage = dataModel.ResourceLanguage(item.value)
	mode.updateTextureText()
}

func (mode *GameTexturesMode) updateData() {
	climbable := false
	transparencyControl := 0
	animationGroup := 0
	animationIndex := 0

	if texture := mode.textureAdapter.GameTexture(mode.selectedTextureID); texture != nil {
		climbable = texture.Climbable()
		transparencyControl = texture.TransparencyControl()
		animationGroup = texture.AnimationGroup()
		animationIndex = texture.AnimationIndex()
	}

	if climbable {
		mode.climbableBox.SetSelectedItem(mode.climbableItems[1])
	} else {
		mode.climbableBox.SetSelectedItem(mode.climbableItems[0])
	}
	mode.transparencyControlBox.SetSelectedItem(mode.transparencyControlItems[transparencyControl])
	mode.animationGroupSlider.SetValue(int64(animationGroup))
	mode.animationIndexSlider.SetValue(int64(animationIndex))
	mode.updateTextureText()
}

func (mode *GameTexturesMode) updateTextureText() {
	name := ""
	useText := ""

	if texture := mode.textureAdapter.GameTexture(mode.selectedTextureID); texture != nil {
		name = texture.Name(mode.selectedLanguage)
		useText = texture.UseText(mode.selectedLanguage)
	}
	mode.nameValue.SetText(name)
	mode.useTextValue.SetText(useText)
}

func (mode *GameTexturesMode) onNameChangeRequested(newValue string) {
	mode.requestStringPropertyChange(func(properties *dataModel.TextureProperties, value *string) {
		properties.Name[mode.selectedLanguage.ToIndex()] = value
	}, newValue, mode.textureAdapter.GameTexture(mode.selectedTextureID).Name(mode.selectedLanguage))
}

func (mode *GameTexturesMode) onUseTextChangeRequested(newValue string) {
	mode.requestStringPropertyChange(func(properties *dataModel.TextureProperties, value *string) {
		properties.CantBeUsed[mode.selectedLanguage.ToIndex()] = value
	}, newValue, mode.textureAdapter.GameTexture(mode.selectedTextureID).UseText(mode.selectedLanguage))
}

func (mode *GameTexturesMode) onClimbableChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	newValue := item.value != 0
	mode.requestBooleanPropertyChange(func(properties *dataModel.TextureProperties, value *bool) {
		properties.Climbable = value
	}, newValue, mode.textureAdapter.GameTexture(mode.selectedTextureID).Climbable())
}

func (mode *GameTexturesMode) onTransparencyControlChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	newValue := int(item.value)
	mode.requestIntPropertyChange(func(properties *dataModel.TextureProperties, value *int) {
		properties.TransparencyControl = value
	}, newValue, mode.textureAdapter.GameTexture(mode.selectedTextureID).TransparencyControl())
}

func (mode *GameTexturesMode) onAnimationGroupChanged(newValue int64) {
	mode.requestIntPropertyChange(func(properties *dataModel.TextureProperties, value *int) {
		properties.AnimationGroup = value
	}, int(newValue), mode.textureAdapter.GameTexture(mode.selectedTextureID).AnimationGroup())
}

func (mode *GameTexturesMode) onAnimationIndexChanged(newValue int64) {
	mode.requestIntPropertyChange(func(properties *dataModel.TextureProperties, value *int) {
		properties.AnimationIndex = value
	}, int(newValue), mode.textureAdapter.GameTexture(mode.selectedTextureID).AnimationIndex())
}

func (mode *GameTexturesMode) textureDropHandler(textureSize dataModel.TextureSize) ui.EventHandler {
	return func(area *ui.Area, event events.Event) (result bool) {
		fileDropEvent := event.(*events.FileDropEvent)
		filePaths := fileDropEvent.FilePaths()
		if len(filePaths) == 1 {
			file, err := os.Open(filePaths[0])
			var img image.Image

			if err == nil {
				defer func() {
					_ = file.Close()
				}()
				img, _, err = image.Decode(file)
				if err != nil {
					mode.context.ModelAdapter().SetMessage(fmt.Sprintf("File <%v> has unknown image format", filePaths[0]))
				}
			} else {
				mode.context.ModelAdapter().SetMessage(fmt.Sprintf("Could not open file <%v>", filePaths[0]))
			}
			if err == nil {
				mode.importTextureBitmap(textureSize, img)
			}
		}
		return
	}
}

func (mode *GameTexturesMode) importTextureBitmap(textureSize dataModel.TextureSize, img image.Image) {
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

		mode.requestTextureBitmapChange(textureSize, &rawBitmap)
	}
}

func (mode *GameTexturesMode) requestTextureBitmapChange(textureSize dataModel.TextureSize, newBitmap *dataModel.RawBitmap) {
	restoreState := mode.stateSnapshot()

	mode.context.Perform(&cmd.SetBitmapCommand{
		Setter: func(bmp *dataModel.RawBitmap) error {
			restoreState()
			mode.textureAdapter.RequestTextureBitmapChange(mode.selectedTextureID, textureSize, bmp)
			return nil
		},
		NewValue: newBitmap,
		OldValue: mode.textureAdapter.TextureBitmap(mode.selectedTextureID, textureSize)})
}

func (mode *GameTexturesMode) requestStringPropertyChange(modifier func(*dataModel.TextureProperties, *string),
	newValue, oldValue string) {
	if mode.existingTextureSelected() {
		restoreState := mode.stateSnapshot()

		mode.context.Perform(&cmd.SetStringPropertyCommand{
			Setter: func(value string) error {
				return mode.requestPropertyChange(restoreState, func(properties *dataModel.TextureProperties) { modifier(properties, &value) })
			},
			NewValue: newValue,
			OldValue: oldValue})
	}
}

func (mode *GameTexturesMode) requestBooleanPropertyChange(modifier func(*dataModel.TextureProperties, *bool),
	newValue, oldValue bool) {
	if mode.existingTextureSelected() {
		restoreState := mode.stateSnapshot()

		mode.context.Perform(&cmd.SetBooleanPropertyCommand{
			Setter: func(value bool) error {
				return mode.requestPropertyChange(restoreState, func(properties *dataModel.TextureProperties) { modifier(properties, &value) })
			},
			NewValue: newValue,
			OldValue: oldValue})
	}
}

func (mode *GameTexturesMode) requestIntPropertyChange(modifier func(*dataModel.TextureProperties, *int),
	newValue, oldValue int) {
	if mode.existingTextureSelected() {
		restoreState := mode.stateSnapshot()

		mode.context.Perform(&cmd.SetIntPropertyCommand{
			Setter: func(value int) error {
				return mode.requestPropertyChange(restoreState, func(properties *dataModel.TextureProperties) { modifier(properties, &value) })
			},
			NewValue: newValue,
			OldValue: oldValue})
	}
}

func (mode *GameTexturesMode) requestPropertyChange(restoreState func(), modifier func(*dataModel.TextureProperties)) error {
	restoreState()
	var properties dataModel.TextureProperties
	modifier(&properties)
	mode.textureAdapter.RequestTexturePropertiesChange(mode.selectedTextureID, &properties)
	return nil
}

func (mode *GameTexturesMode) existingTextureSelected() bool {
	return (mode.selectedTextureID >= 0) && (mode.selectedTextureID < mode.textureAdapter.WorldTextureCount())
}

func (mode *GameTexturesMode) stateSnapshot() func() {
	currentLanguage := mode.selectedLanguage
	currentID := mode.selectedTextureID
	return func() {
		mode.setState(currentLanguage, currentID)
	}
}

func (mode *GameTexturesMode) setState(language dataModel.ResourceLanguage, id int) {
	{
		mode.selectedLanguage = language
		for _, item := range mode.languageItems {
			if item.value == uint32(language) {
				mode.languageBox.SetSelectedItem(item)
			}
		}
	}
	mode.selectedTextureID = id
	mode.textureIDSlider.SetValue(int64(mode.selectedTextureID))
	mode.textureSelector.SetSelectedIndex(mode.selectedTextureID)
	mode.textureSelector.DisplaySelectedIndex()
	mode.updateData()
}
