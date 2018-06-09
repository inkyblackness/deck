package modes

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path"
	"strings"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data/gameobj"
	"github.com/inkyblackness/res/data/interpreters"
	"github.com/inkyblackness/shocked-client/editor/cmd"
	"github.com/inkyblackness/shocked-client/editor/model"
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/graphics/controls"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"

	dataModel "github.com/inkyblackness/shocked-model"
)

// GameObjectsMode is a mode for game object properties.
type GameObjectsMode struct {
	context        Context
	objectsAdapter *model.ObjectsAdapter

	area           *ui.Area
	propertiesArea *ui.Area

	objectClassLabel *controls.Label
	objectClassBox   *controls.ComboBox
	objectClassItems enumItems
	objectTypeLabel  *controls.Label
	objectTypeBox    *controls.ComboBox
	objectTypeItems  []controls.ComboBoxItem
	selectedObjectID model.ObjectID

	bitmapIndexLabel    *controls.Label
	bitmapIndexSlider   *controls.Slider
	selectedBitmapIndex int

	selectedPropertiesTitle *controls.Label
	selectedPropertiesBox   *controls.ComboBox

	commonPropertiesItem    *tabItem
	commonPropertiesPanel   *propertyPanel
	genericPropertiesItem   *tabItem
	genericPropertiesPanel  *propertyPanel
	specificPropertiesItem  *tabItem
	specificPropertiesPanel *propertyPanel

	imageDisplayDrop *ui.Area
	imageDisplay     *controls.ImageDisplay
}

// NewGameObjectsMode returns a new instance.
func NewGameObjectsMode(context Context, parent *ui.Area) *GameObjectsMode {
	mode := &GameObjectsMode{
		context:        context,
		objectsAdapter: context.ModelAdapter().ObjectsAdapter(),

		selectedBitmapIndex: -1}

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
		builder.SetRight(ui.NewRelativeAnchor(parent.Left(), parent.Right(), 0.66))
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
			mode.objectClassLabel, mode.objectClassBox = panelBuilder.addComboProperty("Object Class", mode.onSelectedObjectClassChanged)
			mode.objectTypeLabel, mode.objectTypeBox = panelBuilder.addComboProperty("Selected Object", func(item controls.ComboBoxItem) {
				typeItem := item.(*objectTypeItem)
				mode.onSelectedObjectTypeChanged(typeItem.id)
			})

			mode.objectClassItems = objectClassItems()
			mode.objectClassBox.SetItems(mode.objectClassItems.forComboBox())
			mode.objectsAdapter.OnObjectsChanged(mode.onObjectsChanged)
		}

		mode.bitmapIndexLabel, mode.bitmapIndexSlider = panelBuilder.addSliderProperty("Bitmap", mode.onSelectedBitmapChanged)

		mode.selectedPropertiesTitle, mode.selectedPropertiesBox = panelBuilder.addComboProperty("Show Properties", mode.onSelectedPropertiesDisplayChanged)

		mode.commonPropertiesPanel = newPropertyPanel(panelBuilder, mode.updateCommonProperty, mode.objectItemsForClass)
		mode.genericPropertiesPanel = newPropertyPanel(panelBuilder, mode.updateGenericProperty, mode.objectItemsForClass)
		mode.specificPropertiesPanel = newPropertyPanel(panelBuilder, mode.updateSpecificProperty, mode.objectItemsForClass)

		mode.commonPropertiesItem = &tabItem{mode.commonPropertiesPanel, "Common Properties"}
		mode.genericPropertiesItem = &tabItem{mode.genericPropertiesPanel, "Generic Properties"}
		mode.specificPropertiesItem = &tabItem{mode.specificPropertiesPanel, "Specific Properties"}
		propertiesTabItems := []controls.ComboBoxItem{mode.commonPropertiesItem, mode.genericPropertiesItem, mode.specificPropertiesItem}
		mode.selectedPropertiesBox.SetItems(propertiesTabItems)
		mode.selectedPropertiesBox.SetSelectedItem(mode.commonPropertiesItem)
		mode.onSelectedPropertiesDisplayChanged(mode.commonPropertiesItem)
	}
	{
		padding := scaled(5.0)
		displayWidth := scaled(256)

		{
			dropBuilder := ui.NewAreaBuilder()
			displayBuilder := mode.context.ControlFactory().ForImageDisplay()
			left := ui.NewOffsetAnchor(mode.propertiesArea.Right(), padding)
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
		}
	}

	return mode
}

// SetActive implements the Mode interface.
func (mode *GameObjectsMode) SetActive(active bool) {
	mode.area.SetVisible(active)
}

func (mode *GameObjectsMode) objectItemsForClass(objectClass int) []controls.ComboBoxItem {
	availableGameObjects := mode.objectsAdapter.ObjectsOfClass(objectClass)
	typeItems := make([]controls.ComboBoxItem, len(availableGameObjects))

	for index, gameObject := range availableGameObjects {
		typeItems[index] = &objectTypeItem{gameObject.ID(), gameObject.DisplayName()}
	}

	return typeItems
}

func (mode *GameObjectsMode) onObjectsChanged() {
	mode.setState(mode.selectedObjectID, mode.selectedBitmapIndex)
}

func (mode *GameObjectsMode) onSelectedObjectClassChanged(item controls.ComboBoxItem) {
	classItem := item.(*enumItem)
	mode.setState(model.MakeObjectID(int(classItem.value), 0, 0), 0)
}

func (mode *GameObjectsMode) onSelectedObjectTypeChanged(id model.ObjectID) {
	mode.setState(id, mode.selectedBitmapIndex)
}

func (mode *GameObjectsMode) onSelectedBitmapChanged(newValue int64) {
	mode.setState(mode.selectedObjectID, int(newValue))
}

func (mode *GameObjectsMode) recreatePropertyControls(bitmapCount int) {
	object := mode.objectsAdapter.Object(mode.selectedObjectID)

	mode.commonPropertiesPanel.Reset()
	mode.genericPropertiesPanel.Reset()
	mode.specificPropertiesPanel.Reset()
	mode.bitmapIndexSlider.SetRange(0, 2)
	if object != nil {
		mode.bitmapIndexSlider.SetRange(0, int64(bitmapCount)-1)
		mode.createPropertyControls(gameobj.CommonProperties(object.CommonData()), mode.commonPropertiesPanel)
		mode.createPropertyControls(gameobj.GenericProperties(res.ObjectClass(mode.selectedObjectID.Class()),
			object.GenericData()), mode.genericPropertiesPanel)
		mode.createPropertyControls(gameobj.SpecificProperties(
			res.MakeObjectID(res.ObjectClass(mode.selectedObjectID.Class()), res.ObjectSubclass(mode.selectedObjectID.Subclass()), res.ObjectType(mode.selectedObjectID.Type())),
			object.SpecificData()), mode.specificPropertiesPanel)
	}
}

func (mode *GameObjectsMode) createPropertyControls(rootInterpreter *interpreters.Instance, panel *propertyPanel) {
	var processInterpreter func(string, *interpreters.Instance)
	processInterpreter = func(path string, interpreter *interpreters.Instance) {
		for _, key := range interpreter.Keys() {
			fullPath := path + key
			simplifier := panel.NewSimplifier(fullPath, int64(interpreter.Get(key)))

			interpreter.Describe(key, simplifier)
		}
		for _, key := range interpreter.ActiveRefinements() {
			processInterpreter(path+key+".", interpreter.Refined(key))
		}
	}
	processInterpreter("", rootInterpreter)
}

func (mode *GameObjectsMode) onSelectedPropertiesDisplayChanged(item controls.ComboBoxItem) {
	tabItem := item.(*tabItem)

	mode.commonPropertiesItem.page.SetVisible(false)
	mode.genericPropertiesItem.page.SetVisible(false)
	mode.specificPropertiesItem.page.SetVisible(false)
	tabItem.page.SetVisible(true)
}

func (mode *GameObjectsMode) updateCommonProperty(fullPath string, parameter uint32, update propertyUpdateFunction) {
	mode.requestObjectPropertiesChange(func(object *model.GameObject, properties *dataModel.GameObjectProperties) {
		properties.Data.Common = cloneBytes(object.CommonData())
		interpreter := gameobj.CommonProperties(properties.Data.Common)
		mode.updateObjectProperty(interpreter, fullPath, parameter, update)
	})
}

func (mode *GameObjectsMode) updateGenericProperty(fullPath string, parameter uint32, update propertyUpdateFunction) {
	mode.requestObjectPropertiesChange(func(object *model.GameObject, properties *dataModel.GameObjectProperties) {
		properties.Data.Generic = cloneBytes(object.GenericData())
		interpreter := gameobj.GenericProperties(res.ObjectClass(object.ID().Class()), properties.Data.Generic)
		mode.updateObjectProperty(interpreter, fullPath, parameter, update)
	})
}

func (mode *GameObjectsMode) updateSpecificProperty(fullPath string, parameter uint32, update propertyUpdateFunction) {
	mode.requestObjectPropertiesChange(func(object *model.GameObject, properties *dataModel.GameObjectProperties) {
		properties.Data.Specific = cloneBytes(object.SpecificData())
		interpreter := gameobj.SpecificProperties(
			res.MakeObjectID(res.ObjectClass(object.ID().Class()), res.ObjectSubclass(object.ID().Subclass()), res.ObjectType(object.ID().Type())),
			properties.Data.Specific)
		mode.updateObjectProperty(interpreter, fullPath, parameter, update)
	})
}

func (mode *GameObjectsMode) requestObjectPropertiesChange(modifier func(*model.GameObject, *dataModel.GameObjectProperties)) {
	object := mode.objectsAdapter.Object(mode.selectedObjectID)
	var properties dataModel.GameObjectProperties

	modifier(object, &properties)
	mode.objectsAdapter.RequestObjectPropertiesChange(object.ID(), &properties)
}

func (mode *GameObjectsMode) updateObjectProperty(interpreter *interpreters.Instance,
	fullPath string, parameter uint32, update propertyUpdateFunction) {
	keys := strings.Split(fullPath, ".")
	valueIndex := len(keys) - 1

	for subIndex := 0; subIndex < valueIndex; subIndex++ {
		interpreter = interpreter.Refined(keys[subIndex])
	}
	subKey := keys[valueIndex]
	interpreter.Set(subKey, update(interpreter.Get(subKey), parameter))
}

func (mode *GameObjectsMode) imageProvider() (texture *graphics.BitmapTexture) {
	store := mode.context.ForGraphics().GameObjectBitmapsStore()

	if mode.selectedBitmapIndex >= 0 {
		id := model.ObjectBitmapID{ObjectID: mode.selectedObjectID, Index: mode.selectedBitmapIndex}
		texture = store.Texture(graphics.TextureKeyFromInt(id.ToInt()))
	}
	return
}

func (mode *GameObjectsMode) bitmapDropHandler(area *ui.Area, event events.Event) (consumed bool) {
	dropEvent := event.(*events.FileDropEvent)

	if len(dropEvent.FilePaths()) == 1 {
		filePath := dropEvent.FilePaths()[0]
		fileInfo, err := os.Stat(filePath)

		if err == nil {
			if fileInfo.IsDir() {
				mode.exportBitmap(filePath)
			} else {
				mode.importBitmap(filePath)
			}
		} else {
			mode.context.ModelAdapter().SetMessage(fmt.Sprintf("File is not found/recognized %s", filePath))
		}
		consumed = true
	}

	return
}

func (mode *GameObjectsMode) exportBitmap(filePath string) {
	if mode.selectedBitmapIndex >= 0 {
		fileName := path.Join(filePath, fmt.Sprintf("gameobj_%02d-%02d-%02d_%d.png",
			mode.selectedObjectID.Class(), mode.selectedObjectID.Subclass(), mode.selectedObjectID.Type(), mode.selectedBitmapIndex))
		file, err := os.Create(fileName)

		if err == nil {
			defer file.Close()
			key := model.ObjectBitmapID{ObjectID: mode.selectedObjectID, Index: mode.selectedBitmapIndex}
			rawBitmap := mode.objectsAdapter.Bitmaps().RawBitmap(key.ToInt())
			pixBitmap := graphics.BitmapFromRaw(*rawBitmap)
			gamePalette := mode.context.ModelAdapter().GamePalette()
			imgPalette := make([]color.Color, len(gamePalette))

			for index, paletteColor := range gamePalette {
				imgPalette[index] = paletteColor
			}

			img := image.NewPaletted(image.Rect(0, 0, pixBitmap.Width, pixBitmap.Height), imgPalette)
			for row := 0; row < pixBitmap.Height; row++ {
				start := row * pixBitmap.Width
				copy(img.Pix[row*img.Stride:], pixBitmap.Pixels[start:start+pixBitmap.Width])
			}
			png.Encode(file, img)
			mode.context.ModelAdapter().SetMessage(fmt.Sprintf("Exported %s", fileName))
		} else {
			mode.context.ModelAdapter().SetMessage("Could not create file for export.")
		}
	}
}

func (mode *GameObjectsMode) importBitmap(filePath string) {
	file, err := os.Open(filePath)
	var img image.Image

	if err == nil {
		defer file.Close()
		img, _, err = image.Decode(file)
		if err != nil {
			mode.context.ModelAdapter().SetMessage(fmt.Sprintf("File <%v> has unknown image format", filePath))
		}
	} else {
		mode.context.ModelAdapter().SetMessage(fmt.Sprintf("Could not open file <%v>", filePath))
	}
	if err == nil {
		mode.importBitmapImage(img)
	}
}

func (mode *GameObjectsMode) importBitmapImage(img image.Image) {
	if mode.selectedBitmapIndex >= 0 {
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

func (mode *GameObjectsMode) requestBitmapChange(newBitmap *dataModel.RawBitmap) {
	restoreState := mode.stateSnapshot()
	key := model.ObjectBitmapID{ObjectID: mode.selectedObjectID, Index: mode.selectedBitmapIndex}
	mode.context.Perform(&cmd.SetBitmapCommand{
		Setter: func(bmp *dataModel.RawBitmap) error {
			restoreState()
			mode.objectsAdapter.RequestBitmapChange(key, bmp)
			return nil
		},
		NewValue: newBitmap,
		OldValue: mode.objectsAdapter.Bitmap(key)})
}

func (mode *GameObjectsMode) stateSnapshot() func() {
	currentObjectID := mode.selectedObjectID
	currentBitmapIndex := mode.selectedBitmapIndex
	return func() {
		mode.setState(currentObjectID, currentBitmapIndex)
	}
}

func (mode *GameObjectsMode) setState(objectID model.ObjectID, bitmapIndex int) {
	var selectedTypeItem controls.ComboBoxItem

	mode.selectedObjectID = objectID
	mode.objectTypeItems = mode.objectItemsForClass(objectID.Class())
	mode.objectTypeBox.SetItems(mode.objectTypeItems)
	for _, item := range mode.objectTypeItems {
		typeItem := item.(*objectTypeItem)
		if typeItem.id == mode.selectedObjectID {
			selectedTypeItem = typeItem
		}
	}

	mode.objectClassBox.SetSelectedItem(mode.objectClassItems[mode.selectedObjectID.Class()])
	if selectedTypeItem != nil {
		object := mode.objectsAdapter.Object(mode.selectedObjectID)
		bitmapCount := 2

		if object != nil {
			commonProperties := gameobj.CommonProperties(object.CommonData())
			bitmapCount = 3 + int(commonProperties.Get("Extra")>>4)
		}

		mode.objectTypeBox.SetSelectedItem(selectedTypeItem)
		mode.recreatePropertyControls(bitmapCount)
		if (bitmapIndex >= 0) && (bitmapIndex < bitmapCount) {
			mode.selectedBitmapIndex = bitmapIndex
		} else {
			mode.selectedBitmapIndex = 0
		}
		mode.bitmapIndexSlider.SetValue(int64(mode.selectedBitmapIndex))
	} else {
		mode.objectTypeBox.SetSelectedItem(nil)
		mode.bitmapIndexSlider.SetValueUndefined()
		mode.selectedBitmapIndex = -1
		mode.commonPropertiesPanel.Reset()
		mode.genericPropertiesPanel.Reset()
		mode.specificPropertiesPanel.Reset()
	}
}
