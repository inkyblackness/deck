package modes

import (
	"strings"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data/gameobj"
	"github.com/inkyblackness/res/data/interpreters"
	"github.com/inkyblackness/shocked-client/editor/model"
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/graphics/controls"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"

	datamodel "github.com/inkyblackness/shocked-model"
)

// GameObjectsMode is a mode for game object properties.
type GameObjectsMode struct {
	context        Context
	objectsAdapter *model.ObjectsAdapter

	area *ui.Area

	selectedObjectClassLabel *controls.Label
	selectedObjectClassBox   *controls.ComboBox
	selectedObjectLabel      *controls.Label
	selectedObjectBox        *controls.ComboBox
	selectedObjectID         model.ObjectID

	selectedPropertiesTitle *controls.Label
	selectedPropertiesBox   *controls.ComboBox

	commonPropertiesItem    *tabItem
	commonPropertiesPanel   *propertyPanel
	genericPropertiesItem   *tabItem
	genericPropertiesPanel  *propertyPanel
	specificPropertiesItem  *tabItem
	specificPropertiesPanel *propertyPanel
}

// NewGameObjectsMode returns a new instance.
func NewGameObjectsMode(context Context, parent *ui.Area) *GameObjectsMode {
	mode := &GameObjectsMode{
		context:        context,
		objectsAdapter: context.ModelAdapter().ObjectsAdapter()}

	{
		builder := ui.NewAreaBuilder()
		builder.SetParent(parent)
		builder.SetLeft(ui.NewOffsetAnchor(parent.Left(), 0))
		builder.SetTop(ui.NewOffsetAnchor(parent.Top(), 0))
		builder.SetRight(ui.NewRelativeAnchor(parent.Left(), parent.Right(), 0.66))
		builder.SetBottom(ui.NewOffsetAnchor(parent.Bottom(), 0))
		builder.SetVisible(false)
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
		mode.area = builder.Build()
	}
	{
		panelBuilder := newControlPanelBuilder(mode.area, context.ControlFactory())

		{
			mode.selectedObjectClassLabel, mode.selectedObjectClassBox = panelBuilder.addComboProperty("Object Class", mode.onSelectedObjectClassChanged)
			mode.selectedObjectLabel, mode.selectedObjectBox = panelBuilder.addComboProperty("Selected Object", func(item controls.ComboBoxItem) {
				typeItem := item.(*objectTypeItem)
				mode.onSelectedObjectTypeChanged(typeItem.id)
			})

			mode.objectsAdapter.OnObjectsChanged(mode.onObjectsChanged)
		}

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

	return mode
}

// SetActive implements the Mode interface.
func (mode *GameObjectsMode) SetActive(active bool) {
	mode.area.SetVisible(active)
}

func (mode *GameObjectsMode) onObjectsChanged() {
	classItems := make([]controls.ComboBoxItem, len(classNames))

	for index := range classNames {
		classItems[index] = &objectClassItem{index}
	}
	mode.selectedObjectClassBox.SetItems(classItems)
	mode.selectedObjectClassBox.SetSelectedItem(classItems[mode.selectedObjectID.Class()])
	mode.updateSelectedObjectClass(mode.selectedObjectID.Class())
}

func (mode *GameObjectsMode) onSelectedObjectClassChanged(item controls.ComboBoxItem) {
	classItem := item.(*objectClassItem)
	mode.updateSelectedObjectClass(classItem.class)
}

func (mode *GameObjectsMode) updateSelectedObjectClass(objectClass int) {
	typeItems := mode.objectItemsForClass(objectClass)
	selectedIndex := -1

	mode.selectedObjectBox.SetItems(typeItems)
	for index, item := range typeItems {
		typeItem := item.(*objectTypeItem)
		if typeItem.id == mode.selectedObjectID {
			selectedIndex = index
		}
	}

	if selectedIndex >= 0 {
		mode.selectedObjectBox.SetSelectedItem(typeItems[selectedIndex])
		mode.recreatePropertyControls()
	} else {
		mode.selectedObjectBox.SetSelectedItem(nil)
		mode.commonPropertiesPanel.Reset()
		mode.genericPropertiesPanel.Reset()
		mode.specificPropertiesPanel.Reset()
	}
}

func (mode *GameObjectsMode) objectItemsForClass(objectClass int) []controls.ComboBoxItem {
	availableGameObjects := mode.objectsAdapter.ObjectsOfClass(objectClass)
	typeItems := make([]controls.ComboBoxItem, len(availableGameObjects))

	for index, gameObject := range availableGameObjects {
		typeItems[index] = &objectTypeItem{gameObject.ID(), gameObject.DisplayName()}
	}

	return typeItems
}

func (mode *GameObjectsMode) onSelectedObjectTypeChanged(id model.ObjectID) {
	mode.selectedObjectID = id
	mode.recreatePropertyControls()
}

func (mode *GameObjectsMode) recreatePropertyControls() {
	object := mode.objectsAdapter.Object(mode.selectedObjectID)

	mode.commonPropertiesPanel.Reset()
	mode.genericPropertiesPanel.Reset()
	mode.specificPropertiesPanel.Reset()
	if object != nil {
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
	mode.requestObjectPropertiesChange(func(object *model.GameObject, properties *datamodel.GameObjectProperties) {
		properties.Data.Common = cloneBytes(object.CommonData())
		interpreter := gameobj.CommonProperties(properties.Data.Common)
		mode.updateObjectProperty(interpreter, fullPath, parameter, update)
	})
}

func (mode *GameObjectsMode) updateGenericProperty(fullPath string, parameter uint32, update propertyUpdateFunction) {
	mode.requestObjectPropertiesChange(func(object *model.GameObject, properties *datamodel.GameObjectProperties) {
		properties.Data.Generic = cloneBytes(object.GenericData())
		interpreter := gameobj.GenericProperties(res.ObjectClass(object.ID().Class()), properties.Data.Generic)
		mode.updateObjectProperty(interpreter, fullPath, parameter, update)
	})
}

func (mode *GameObjectsMode) updateSpecificProperty(fullPath string, parameter uint32, update propertyUpdateFunction) {
	mode.requestObjectPropertiesChange(func(object *model.GameObject, properties *datamodel.GameObjectProperties) {
		properties.Data.Specific = cloneBytes(object.SpecificData())
		interpreter := gameobj.SpecificProperties(
			res.MakeObjectID(res.ObjectClass(object.ID().Class()), res.ObjectSubclass(object.ID().Subclass()), res.ObjectType(object.ID().Type())),
			properties.Data.Specific)
		mode.updateObjectProperty(interpreter, fullPath, parameter, update)
	})
}

func (mode *GameObjectsMode) requestObjectPropertiesChange(modifier func(*model.GameObject, *datamodel.GameObjectProperties)) {
	object := mode.objectsAdapter.Object(mode.selectedObjectID)
	var properties datamodel.GameObjectProperties

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
