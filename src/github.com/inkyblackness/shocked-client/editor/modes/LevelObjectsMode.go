package modes

import (
	"fmt"
	"math"
	"sort"
	"strings"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data/interpreters"
	"github.com/inkyblackness/res/data/levelobj"

	dataModel "github.com/inkyblackness/shocked-model"

	"github.com/inkyblackness/shocked-client/editor/display"
	"github.com/inkyblackness/shocked-client/editor/model"
	"github.com/inkyblackness/shocked-client/env"
	"github.com/inkyblackness/shocked-client/env/keys"
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/graphics/controls"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"
	"github.com/inkyblackness/shocked-client/util"
)

// LevelObjectsMode is a mode for level objects.
type LevelObjectsMode struct {
	context        Context
	levelAdapter   *model.LevelAdapter
	objectsAdapter *model.ObjectsAdapter

	displayFilter    func(*model.LevelObject) bool
	displayedObjects []*model.LevelObject

	mapDisplay *display.MapDisplay

	area *ui.Area

	limitsPanel   *ui.Area
	limitsHeader1 *controls.Label
	limitsHeader2 *controls.Label
	limitTitles   []*controls.Label
	limitValues   []*controls.Label

	propertiesPanel      *ui.Area
	propertiesPanelRight ui.Anchor

	closestObjects              []*model.LevelObject
	closestObjectHighlightIndex int
	selectedObjects             []*model.LevelObject

	newObjectID model.ObjectID

	newObjectClassLabel *controls.Label
	newObjectClassBox   *controls.ComboBox
	newObjectTypeLabel  *controls.Label
	newObjectTypeBox    *controls.ComboBox

	highlightedObjectInfoTitle *controls.Label
	highlightedObjectInfoValue *controls.Label

	selectedObjectsTitleLabel   *controls.Label
	selectedObjectsDeleteLabel  *controls.Label
	selectedObjectsDeleteButton *controls.TextButton
	selectedObjectsIDTitleLabel *controls.Label
	selectedObjectsIDInfoLabel  *controls.Label
	selectedObjectsTypeLabel    *controls.Label
	selectedObjectsTypeBox      *controls.ComboBox

	selectedObjectsPropertiesTitle *controls.Label
	selectedObjectsPropertiesBox   *controls.ComboBox

	selectedObjectsBasePropertiesItem *tabItem
	selectedObjectsBasePropertiesArea *ui.Area
	selectedObjectsZTitle             *controls.Label
	selectedObjectsZValue             *controls.Slider
	selectedObjectsTileXTitle         *controls.Label
	selectedObjectsTileXValue         *controls.Slider
	selectedObjectsFineXTitle         *controls.Label
	selectedObjectsFineXValue         *controls.Slider
	selectedObjectsTileYTitle         *controls.Label
	selectedObjectsTileYValue         *controls.Slider
	selectedObjectsFineYTitle         *controls.Label
	selectedObjectsFineYValue         *controls.Slider
	selectedObjectsRotationXTitle     *controls.Label
	selectedObjectsRotationXValue     *controls.Slider
	selectedObjectsRotationYTitle     *controls.Label
	selectedObjectsRotationYValue     *controls.Slider
	selectedObjectsRotationZTitle     *controls.Label
	selectedObjectsRotationZValue     *controls.Slider
	selectedObjectsHitpointsTitle     *controls.Label
	selectedObjectsHitpointsValue     *controls.Slider

	selectedObjectsClassPropertiesItem  *tabItem
	selectedObjectsPropertiesMainArea   *ui.Area
	selectedObjectsPropertiesHeaderArea *ui.Area
	selectedObjectsPropertiesPanel      *propertyPanel
}

// NewLevelObjectsMode returns a new instance.
func NewLevelObjectsMode(context Context, parent *ui.Area, mapDisplay *display.MapDisplay) *LevelObjectsMode {
	mode := &LevelObjectsMode{
		context:        context,
		levelAdapter:   context.ModelAdapter().ActiveLevel(),
		objectsAdapter: context.ModelAdapter().ObjectsAdapter(),
		displayFilter:  func(*model.LevelObject) bool { return true },

		newObjectID: model.MakeObjectID(0, 0, 0),

		mapDisplay: mapDisplay}

	{
		builder := ui.NewAreaBuilder()
		builder.SetParent(parent)
		builder.SetLeft(ui.NewOffsetAnchor(parent.Left(), 0))
		builder.SetTop(ui.NewOffsetAnchor(parent.Top(), 0))
		builder.SetRight(ui.NewOffsetAnchor(parent.Right(), 0))
		builder.SetBottom(ui.NewOffsetAnchor(parent.Bottom(), 0))
		builder.SetVisible(false)
		builder.OnEvent(events.MouseMoveEventType, mode.onMouseMoved)
		builder.OnEvent(events.MouseScrollEventType, mode.onMouseScrolled)
		builder.OnEvent(events.MouseButtonClickedEventType, mode.onMouseButtonClicked)
		mode.area = builder.Build()
	}
	{
		minRight := ui.NewOffsetAnchor(mode.area.Left(), 100)
		maxRight := ui.NewRelativeAnchor(mode.area.Left(), mode.area.Right(), 0.5)
		mode.propertiesPanelRight = ui.NewLimitedAnchor(minRight, maxRight, ui.NewOffsetAnchor(mode.area.Left(), 400))
		builder := ui.NewAreaBuilder()
		builder.SetParent(mode.area)
		builder.SetLeft(ui.NewOffsetAnchor(mode.area.Left(), 0))
		builder.SetTop(ui.NewOffsetAnchor(mode.area.Top(), 0))
		builder.SetRight(mode.propertiesPanelRight)
		builder.SetBottom(ui.NewOffsetAnchor(mode.area.Bottom(), 0))
		builder.OnRender(func(area *ui.Area) {
			context.ForGraphics().RectangleRenderer().Fill(
				area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
				graphics.RGBA(0.7, 0.0, 0.7, 0.3))
		})

		builder.OnEvent(events.MouseButtonClickedEventType, ui.SilentConsumer)
		builder.OnEvent(events.MouseScrollEventType, ui.SilentConsumer)

		lastGrabX := float32(0.0)
		grabbing := false
		builder.OnEvent(events.MouseButtonDownEventType, func(area *ui.Area, event events.Event) bool {
			buttonEvent := event.(*events.MouseButtonEvent)
			if buttonEvent.Buttons() == env.MousePrimary {
				area.RequestFocus()
				lastGrabX, _ = buttonEvent.Position()
				grabbing = true
			}
			return true
		})
		builder.OnEvent(events.MouseButtonUpEventType, func(area *ui.Area, event events.Event) bool {
			buttonEvent := event.(*events.MouseButtonEvent)
			if buttonEvent.AffectedButtons() == env.MousePrimary {
				area.ReleaseFocus()
				grabbing = false
			}
			return true
		})
		builder.OnEvent(events.MouseMoveEventType, func(area *ui.Area, event events.Event) bool {
			moveEvent := event.(*events.MouseMoveEvent)
			if grabbing {
				newX, _ := moveEvent.Position()
				mode.propertiesPanelRight.RequestValue(mode.propertiesPanelRight.Value() + (newX - lastGrabX))
				lastGrabX = newX
			}
			return true
		})

		mode.propertiesPanel = builder.Build()
	}
	{
		panelBuilder := newControlPanelBuilder(mode.propertiesPanel, context.ControlFactory())

		mode.newObjectClassLabel, mode.newObjectClassBox = panelBuilder.addComboProperty("New Object Class", mode.onNewObjectClassChanged)
		mode.newObjectTypeLabel, mode.newObjectTypeBox = panelBuilder.addComboProperty("New Object Type", mode.onNewObjectTypeChanged)

		mode.highlightedObjectInfoTitle, mode.highlightedObjectInfoValue = panelBuilder.addInfo("Highlighted Object")

		mode.selectedObjectsTitleLabel = panelBuilder.addTitle("Selected Object(s)")
		mode.selectedObjectsDeleteLabel, mode.selectedObjectsDeleteButton = panelBuilder.addTextButton("Delete Selected", "Delete", mode.deleteSelectedObjects)
		mode.selectedObjectsIDTitleLabel, mode.selectedObjectsIDInfoLabel = panelBuilder.addInfo("Object ID")
		mode.selectedObjectsTypeLabel, mode.selectedObjectsTypeBox = panelBuilder.addComboProperty("Type", func(item controls.ComboBoxItem) {
			typeItem := item.(*objectTypeItem)
			mode.updateSelectedObjectsBaseProperties(func(properties *dataModel.LevelObjectProperties) {
				properties.Subclass = intAsPointer(typeItem.id.Subclass())
				properties.Type = intAsPointer(typeItem.id.Type())
			})
		})

		mode.selectedObjectsPropertiesTitle, mode.selectedObjectsPropertiesBox = panelBuilder.addComboProperty("Show Properties", mode.onSelectedPropertiesDisplayChanged)

		var basePropertiesPanelBuilder *controlPanelBuilder
		mode.selectedObjectsBasePropertiesArea, basePropertiesPanelBuilder = panelBuilder.addSection(false)
		mode.selectedObjectsZTitle, mode.selectedObjectsZValue = basePropertiesPanelBuilder.addSliderProperty("Z", func(newValue int64) {
			mode.updateSelectedObjectsBaseProperties(func(properties *dataModel.LevelObjectProperties) {
				properties.Z = intAsPointer(int(newValue))
			})
		})
		mode.selectedObjectsZValue.SetRange(0, 255)

		mode.selectedObjectsTileXTitle, mode.selectedObjectsTileXValue = basePropertiesPanelBuilder.addSliderProperty("TileX", func(newValue int64) {
			mode.updateSelectedObjectsBaseProperties(func(properties *dataModel.LevelObjectProperties) {
				properties.TileX = intAsPointer(int(newValue))
			})
		})
		mode.selectedObjectsTileXValue.SetRange(0, 63)
		mode.selectedObjectsFineXTitle, mode.selectedObjectsFineXValue = basePropertiesPanelBuilder.addSliderProperty("FineX", func(newValue int64) {
			mode.updateSelectedObjectsBaseProperties(func(properties *dataModel.LevelObjectProperties) {
				properties.FineX = intAsPointer(int(newValue))
			})
		})
		mode.selectedObjectsFineXValue.SetRange(0, 255)

		mode.selectedObjectsTileYTitle, mode.selectedObjectsTileYValue = basePropertiesPanelBuilder.addSliderProperty("TileY", func(newValue int64) {
			mode.updateSelectedObjectsBaseProperties(func(properties *dataModel.LevelObjectProperties) {
				properties.TileY = intAsPointer(int(newValue))
			})
		})
		mode.selectedObjectsTileYValue.SetRange(0, 63)
		mode.selectedObjectsFineYTitle, mode.selectedObjectsFineYValue = basePropertiesPanelBuilder.addSliderProperty("FineY", func(newValue int64) {
			mode.updateSelectedObjectsBaseProperties(func(properties *dataModel.LevelObjectProperties) {
				properties.FineY = intAsPointer(int(newValue))
			})
		})
		mode.selectedObjectsFineYValue.SetRange(0, 255)

		mode.selectedObjectsRotationXTitle, mode.selectedObjectsRotationXValue = basePropertiesPanelBuilder.addSliderProperty("RotationX", func(newValue int64) {
			mode.updateSelectedObjectsBaseProperties(func(properties *dataModel.LevelObjectProperties) {
				properties.RotationX = intAsPointer(int(newValue))
			})
		})
		mode.selectedObjectsRotationXValue.SetRange(0, 255)
		mode.selectedObjectsRotationXValue.SetValueFormatter(mode.rotationToString)
		mode.selectedObjectsRotationYTitle, mode.selectedObjectsRotationYValue = basePropertiesPanelBuilder.addSliderProperty("RotationY", func(newValue int64) {
			mode.updateSelectedObjectsBaseProperties(func(properties *dataModel.LevelObjectProperties) {
				properties.RotationY = intAsPointer(int(newValue))
			})
		})
		mode.selectedObjectsRotationYValue.SetRange(0, 255)
		mode.selectedObjectsRotationYValue.SetValueFormatter(mode.rotationToString)
		mode.selectedObjectsRotationZTitle, mode.selectedObjectsRotationZValue = basePropertiesPanelBuilder.addSliderProperty("RotationZ", func(newValue int64) {
			mode.updateSelectedObjectsBaseProperties(func(properties *dataModel.LevelObjectProperties) {
				properties.RotationZ = intAsPointer(int(newValue))
			})
		})
		mode.selectedObjectsRotationZValue.SetRange(0, 255)
		mode.selectedObjectsRotationZValue.SetValueFormatter(mode.rotationToString)

		mode.selectedObjectsHitpointsTitle, mode.selectedObjectsHitpointsValue = basePropertiesPanelBuilder.addSliderProperty("Hitpoints", func(newValue int64) {
			mode.updateSelectedObjectsBaseProperties(func(properties *dataModel.LevelObjectProperties) {
				properties.Hitpoints = intAsPointer(int(newValue))
			})
		})
		mode.selectedObjectsHitpointsValue.SetRange(0, 10000)

		classPropertiesBottomResolver := func() ui.Anchor { return mode.selectedObjectsPropertiesPanel.Bottom() }
		var mainClassPanelBuilder *controlPanelBuilder
		mode.selectedObjectsPropertiesMainArea, mainClassPanelBuilder =
			panelBuilder.addDynamicSection(true, classPropertiesBottomResolver)
		mode.selectedObjectsPropertiesHeaderArea, _ = mainClassPanelBuilder.addSection(true)

		mode.selectedObjectsPropertiesPanel = newPropertyPanel(mainClassPanelBuilder,
			mode.updateSelectedObjectsClassPropertiesFiltered, mode.objectItemsForClass)

		mode.selectedObjectsBasePropertiesItem = &tabItem{mode.selectedObjectsBasePropertiesArea, "Base Properties"}
		mode.selectedObjectsClassPropertiesItem = &tabItem{mode.selectedObjectsPropertiesMainArea, "Class Properties"}
		propertiesTabItems := []controls.ComboBoxItem{mode.selectedObjectsBasePropertiesItem, mode.selectedObjectsClassPropertiesItem}
		mode.selectedObjectsPropertiesBox.SetItems(propertiesTabItems)
		mode.selectedObjectsPropertiesBox.SetSelectedItem(mode.selectedObjectsClassPropertiesItem)
	}
	{
		builder := ui.NewAreaBuilder()
		builder.SetParent(mode.area)
		builder.SetLeft(ui.NewOffsetAnchor(mode.area.Right(), -126))
		builder.SetTop(ui.NewOffsetAnchor(mode.area.Top(), 0))
		builder.SetRight(ui.NewOffsetAnchor(mode.area.Right(), 0))
		builder.SetBottom(ui.NewOffsetAnchor(mode.area.Bottom(), 0))
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
		mode.limitsPanel = builder.Build()
	}
	{
		panelBuilder := newControlPanelBuilder(mode.limitsPanel, context.ControlFactory())
		mode.limitsHeader1, mode.limitsHeader2 = panelBuilder.addInfo("Class")
		mode.limitsHeader2.SetText("Usage")
		classes := len(maxObjectsPerClass)
		mode.limitTitles = make([]*controls.Label, classes+1)
		mode.limitValues = make([]*controls.Label, classes+1)
		for class := 0; class < classes; class++ {
			mode.limitTitles[class], mode.limitValues[class] = panelBuilder.addInfo(fmt.Sprintf("%d", class))
		}
		mode.limitTitles[classes], mode.limitValues[classes] = panelBuilder.addInfo("Total")
	}

	mode.levelAdapter.OnLevelPropertiesChanged(func() {
		mode.selectedObjectsZValue.SetValueFormatter(mode.objectZToString)
	})
	mode.levelAdapter.OnLevelObjectsChanged(mode.onLevelObjectsChanged)
	mode.context.ModelAdapter().ObjectsAdapter().OnObjectsChanged(mode.onGameObjectsChanged)

	return mode
}

func (mode *LevelObjectsMode) objectZToString(value int64) string {
	tileHeights := []float64{32.0, 16.0, 8.0, 4.0, 2.0, 1.0, 0.5, 0.25}
	heightShift := mode.levelAdapter.HeightShift()

	return fmt.Sprintf("%.3f tile(s)  - raw: %v", (float64(value)*tileHeights[heightShift])/256.0, value)
}

func (mode *LevelObjectsMode) rotationToString(value int64) string {
	return fmt.Sprintf("%.3f degrees  - raw: %v", (float64(value)*360.0)/256.0, value)
}

// SetActive implements the Mode interface.
func (mode *LevelObjectsMode) SetActive(active bool) {
	if active {
		mode.updateDisplayedObjects()
		mode.mapDisplay.SetSelectedObjects(mode.selectedObjects)
	} else {
		mode.mapDisplay.SetDisplayedObjects(nil)
		mode.mapDisplay.SetHighlightedObject(nil)
		mode.mapDisplay.SetSelectedObjects(nil)
	}
	mode.area.SetVisible(active)
	mode.mapDisplay.SetVisible(active)
}

func (mode *LevelObjectsMode) onLevelObjectsChanged() {
	mode.updateNewObjectClassQuota()
	if mode.area.IsVisible() {
		mode.updateDisplayedObjects()
	}
}

func (mode *LevelObjectsMode) onSelectedPropertiesDisplayChanged(item controls.ComboBoxItem) {
	tabItem := item.(*tabItem)

	mode.selectedObjectsBasePropertiesItem.page.SetVisible(false)
	mode.selectedObjectsClassPropertiesItem.page.SetVisible(false)
	tabItem.page.SetVisible(true)
}

func (mode *LevelObjectsMode) updateDisplayedObjects() {
	mode.displayedObjects = mode.levelAdapter.LevelObjects(mode.displayFilter)
	mode.mapDisplay.SetDisplayedObjects(mode.displayedObjects)

	mode.closestObjects = nil
	mode.closestObjectHighlightIndex = 0
	mode.updateClosestObjectHighlight()

	mode.updateSelectedFromDisplayedObjects()
}

func (mode *LevelObjectsMode) updateSelectedFromDisplayedObjects() {
	displayedIndices := make(map[int]bool)
	for _, displayedObject := range mode.displayedObjects {
		displayedIndices[displayedObject.Index()] = true
	}

	selectedObjects := make([]*model.LevelObject, 0, len(mode.selectedObjects))
	for _, selectedObject := range mode.selectedObjects {
		if displayedIndices[selectedObject.Index()] {
			selectedObjects = append(selectedObjects, selectedObject)
		}
	}
	mode.setSelectedObjects(selectedObjects)
}

func (mode *LevelObjectsMode) onMouseMoved(area *ui.Area, event events.Event) (consumed bool) {
	mouseEvent := event.(*events.MouseMoveEvent)

	if mouseEvent.Buttons() == 0 {
		worldX, worldY := mode.mapDisplay.WorldCoordinatesForPixel(mouseEvent.Position())
		mode.updateClosestDisplayedObjects(worldX, worldY)
		consumed = true
	}

	return
}

func (mode *LevelObjectsMode) onMouseScrolled(area *ui.Area, event events.Event) (consumed bool) {
	mouseEvent := event.(*events.MouseScrollEvent)

	if (mouseEvent.Buttons() == 0) && (keys.Modifier(mouseEvent.Modifier()) == keys.ModControl) {
		available := len(mode.closestObjects)

		if available > 1 {
			_, dy := mouseEvent.Deltas()
			delta := 1

			if dy < 0 {
				delta = -1
			}
			mode.closestObjectHighlightIndex = (available + mode.closestObjectHighlightIndex + delta) % available
			mode.updateClosestObjectHighlight()
		}
		consumed = true
	}

	return
}

func (mode *LevelObjectsMode) onMouseButtonClicked(area *ui.Area, event events.Event) (consumed bool) {
	mouseEvent := event.(*events.MouseButtonEvent)

	if mouseEvent.AffectedButtons() == env.MousePrimary {
		if len(mode.closestObjects) > 0 {
			object := mode.closestObjects[mode.closestObjectHighlightIndex]
			if keys.Modifier(mouseEvent.Modifier()) == keys.ModControl {
				mode.toggleSelectedObject(object)
			} else {
				mode.setSelectedObjects([]*model.LevelObject{object})
			}
		}
		consumed = true
	} else if mouseEvent.AffectedButtons() == env.MouseSecondary {
		worldX, worldY := mode.mapDisplay.WorldCoordinatesForPixel(mouseEvent.Position())
		mode.createNewObject(worldX, worldY)
	}

	return
}

func (mode *LevelObjectsMode) updateClosestDisplayedObjects(worldX, worldY float32) {
	type resultEntry struct {
		distance float32
		object   *model.LevelObject
	}
	entries := []*resultEntry{}
	refPoint := mgl.Vec2{worldX, worldY}
	limit := float32(48.0)

	for _, object := range mode.displayedObjects {
		otherX, otherY := object.Center()
		otherPoint := mgl.Vec2{otherX, otherY}
		delta := refPoint.Sub(otherPoint)
		len := delta.Len()

		if len <= limit {
			entries = append(entries, &resultEntry{len, object})
		}
	}
	sort.Slice(entries, func(a int, b int) bool { return entries[a].distance < entries[b].distance })
	mode.closestObjects = make([]*model.LevelObject, len(entries))
	for index, entry := range entries {
		mode.closestObjects[index] = entry.object
	}
	mode.closestObjectHighlightIndex = 0
	mode.updateClosestObjectHighlight()
}

func (mode *LevelObjectsMode) updateClosestObjectHighlight() {
	if len(mode.closestObjects) > 0 {
		object := mode.closestObjects[mode.closestObjectHighlightIndex]
		mode.highlightedObjectInfoValue.SetText(fmt.Sprintf("%v: %v (%v)", object.Index(), object.ID(), mode.objectDisplayName(object.ID())))
		mode.mapDisplay.SetHighlightedObject(object)
	} else {
		mode.highlightedObjectInfoValue.SetText("")
		mode.mapDisplay.SetHighlightedObject(nil)
	}
}

func (mode *LevelObjectsMode) objectDisplayName(id model.ObjectID) string {
	displayName := "unknown"

	if gameObject := mode.context.ModelAdapter().ObjectsAdapter().Object(id); gameObject != nil {
		displayName = gameObject.DisplayName()
	}

	return displayName
}

func (mode *LevelObjectsMode) setSelectedObjects(objects []*model.LevelObject) {
	mode.selectedObjects = objects
	mode.onSelectedObjectsChanged()
}

func (mode *LevelObjectsMode) toggleSelectedObject(object *model.LevelObject) {
	newList := []*model.LevelObject{}
	wasSelected := false

	for _, other := range mode.selectedObjects {
		if other.Index() != object.Index() {
			newList = append(newList, other)
		} else {
			wasSelected = true
		}
	}
	if !wasSelected {
		newList = append(newList, object)
	}
	mode.setSelectedObjects(newList)
}

func (mode *LevelObjectsMode) onSelectedObjectsChanged() {
	mode.mapDisplay.SetSelectedObjects(mode.selectedObjects)

	classUnifier := util.NewValueUnifier(-1)
	subclassUnifier := util.NewValueUnifier(-1)
	typeUnifier := util.NewValueUnifier(-1)
	tileXUnifier := util.NewValueUnifier(-1)
	fineXUnifier := util.NewValueUnifier(-1)
	tileYUnifier := util.NewValueUnifier(-1)
	fineYUnifier := util.NewValueUnifier(-1)
	zUnifier := util.NewValueUnifier(-1)
	rotationXUnifier := util.NewValueUnifier(-1)
	rotationYUnifier := util.NewValueUnifier(-1)
	rotationZUnifier := util.NewValueUnifier(-1)
	hitpointsUnifier := util.NewValueUnifier(-1)

	for _, object := range mode.selectedObjects {
		classUnifier.Add(object.ID().Class())
		subclassUnifier.Add(object.ID().Subclass())
		typeUnifier.Add(object.ID().Type())
		zUnifier.Add(object.Z())
		tileXUnifier.Add(object.TileX())
		fineXUnifier.Add(object.FineX())
		tileYUnifier.Add(object.TileY())
		fineYUnifier.Add(object.FineY())
		rotationXUnifier.Add(object.RotationX())
		rotationYUnifier.Add(object.RotationY())
		rotationZUnifier.Add(object.RotationZ())
		hitpointsUnifier.Add(object.Hitpoints())
	}
	unifiedClass := classUnifier.Value().(int)
	unifiedSubclass := subclassUnifier.Value().(int)
	unifiedType := typeUnifier.Value().(int)
	var unifiedIDString string
	if unifiedClass != -1 {
		unifiedIDString = classNames[unifiedClass]
		typeItems := mode.objectItemsForClass(unifiedClass)
		mode.selectedObjectsTypeBox.SetItems(typeItems)
	} else {
		unifiedIDString = "**"
		unifiedSubclass = -1
		mode.selectedObjectsTypeBox.SetItems(nil)
	}
	unifiedIDString += "/"
	if unifiedSubclass != -1 {
		unifiedIDString += fmt.Sprintf("%v", unifiedSubclass)
	} else {
		unifiedIDString += "*"
		unifiedType = -1
	}
	unifiedIDString += "/"
	if unifiedType != -1 {
		unifiedIDString += fmt.Sprintf("%v", unifiedType)
	} else {
		unifiedIDString += "*"
	}
	if len(mode.selectedObjects) > 0 {
		mode.selectedObjectsIDInfoLabel.SetText(unifiedIDString)
	} else {
		mode.selectedObjectsIDInfoLabel.SetText("")
	}
	if (unifiedClass != -1) && (unifiedSubclass != -1) && (unifiedType != -1) {
		objectID := model.MakeObjectID(unifiedClass, unifiedSubclass, unifiedType)
		item := &objectTypeItem{objectID, mode.objectDisplayName(objectID)}
		mode.selectedObjectsTypeBox.SetSelectedItem(item)
	} else {
		mode.selectedObjectsTypeBox.SetSelectedItem(nil)
	}

	setSliderValue := func(slider *controls.Slider, unifier *util.ValueUnifier) {
		value := unifier.Value().(int)
		if value >= 0 {
			slider.SetValue(int64(value))
		} else {
			slider.SetValueUndefined()
		}
	}

	setSliderValue(mode.selectedObjectsTileXValue, tileXUnifier)
	setSliderValue(mode.selectedObjectsFineXValue, fineXUnifier)
	setSliderValue(mode.selectedObjectsTileYValue, tileYUnifier)
	setSliderValue(mode.selectedObjectsFineYValue, fineYUnifier)
	setSliderValue(mode.selectedObjectsZValue, zUnifier)
	setSliderValue(mode.selectedObjectsRotationXValue, rotationXUnifier)
	setSliderValue(mode.selectedObjectsRotationYValue, rotationYUnifier)
	setSliderValue(mode.selectedObjectsRotationZValue, rotationZUnifier)
	setSliderValue(mode.selectedObjectsHitpointsValue, hitpointsUnifier)

	mode.recreateLevelObjectProperties()
}

func (mode *LevelObjectsMode) recreateLevelObjectProperties() {
	mode.selectedObjectsPropertiesPanel.Reset()

	if len(mode.selectedObjects) > 0 {
		propertyUnifier := make(map[string]*util.ValueUnifier)
		propertyDescribers := make(map[string]func(*interpreters.Simplifier))
		propertyOrder := []string{}
		describer := func(interpreter *interpreters.Instance, key string) func(simpl *interpreters.Simplifier) {
			return func(simpl *interpreters.Simplifier) { interpreter.Describe(key, simpl) }
		}

		var unifyInterpreter func(string, *interpreters.Instance, bool, map[string]bool)
		unifyInterpreter = func(path string, interpreter *interpreters.Instance, first bool, thisKeys map[string]bool) {
			for _, key := range interpreter.Keys() {
				fullPath := path + key
				thisKeys[fullPath] = true
				if unifier, existing := propertyUnifier[fullPath]; existing || first {
					if !existing {
						unifier = util.NewValueUnifier(int64(math.MinInt64))
						propertyUnifier[fullPath] = unifier
						propertyDescribers[fullPath] = describer(interpreter, key)
						propertyOrder = append(propertyOrder, fullPath)
					}
					unifier.Add(int64(interpreter.Get(key)))
				}
			}
			for _, key := range interpreter.ActiveRefinements() {
				unifyInterpreter(path+key+".", interpreter.Refined(key), first, thisKeys)
			}
		}

		interpreterFactory := mode.interpreterFactory()

		for index, object := range mode.selectedObjects {
			objID := object.ID()
			resID := res.MakeObjectID(res.ObjectClass(objID.Class()), res.ObjectSubclass(objID.Subclass()), res.ObjectType(objID.Type()))
			interpreter := interpreterFactory(resID, object.ClassData())
			thisKeys := make(map[string]bool)
			unifyInterpreter("", interpreter, index == 0, thisKeys)
			{
				toRemove := []string{}
				for previousKey := range propertyUnifier {
					if !thisKeys[previousKey] {
						toRemove = append(toRemove, previousKey)
					}
				}
				for _, key := range toRemove {
					delete(propertyUnifier, key)
				}
			}
		}

		for _, key := range propertyOrder {
			if unifier, existing := propertyUnifier[key]; existing {
				mode.createPropertyControls(key, unifier.Value().(int64), propertyDescribers[key])
			}
		}
	}
}

func (mode *LevelObjectsMode) interpreterFactory() func(resID res.ObjectID, classData []byte) *interpreters.Instance {
	factory := levelobj.ForRealWorld
	if mode.levelAdapter.IsCyberspace() {
		factory = levelobj.ForCyberspace
	}
	return factory
}

func (mode *LevelObjectsMode) createPropertyControls(key string, unifiedValue int64, describer func(*interpreters.Simplifier)) {
	simplifier := mode.selectedObjectsPropertiesPanel.NewSimplifier(key, unifiedValue)

	simplifier.SetObjectIndexHandler(func() {
		slider := mode.selectedObjectsPropertiesPanel.NewSlider(key, "", setUpdate())
		slider.SetRange(0, 871)
		if unifiedValue != math.MinInt64 {
			slider.SetValue(unifiedValue)
		}
	})

	addVariableKey := func() {
		typeBox := mode.selectedObjectsPropertiesPanel.NewComboBox(key, "Type", maskedUpdate(0, 0x1000))
		items := make([]controls.ComboBoxItem, 2)
		items[0] = &enumItem{0, "Boolean"}
		items[1] = &enumItem{0x1000, "Integer"}
		var selectedItem controls.ComboBoxItem
		if unifiedValue != math.MinInt64 {
			if (unifiedValue & 0x1000) != 0 {
				selectedItem = items[1]
			} else {
				selectedItem = items[0]
			}
		}
		typeBox.SetItems(items)
		typeBox.SetSelectedItem(selectedItem)

		indexSlider := mode.selectedObjectsPropertiesPanel.NewSlider(key, "Index", maskedUpdate(0, 0x1FF))
		indexSlider.SetRange(0, 0x1FF)
		if unifiedValue != math.MinInt64 {
			indexSlider.SetValue(unifiedValue & 0x1FF)
			if (unifiedValue & 0x1000) != 0 {
				indexSlider.SetRange(0, 0x3F)
			}
		}
	}

	simplifier.SetSpecialHandler("VariableKey", addVariableKey)
	simplifier.SetSpecialHandler("VariableCondition", func() {
		addVariableKey()

		comparisonBox := mode.selectedObjectsPropertiesPanel.NewComboBox(key, "Check", maskedUpdate(13, 0xE000))
		var selectedItem controls.ComboBoxItem
		items := []controls.ComboBoxItem{
			&enumItem{0, "Var == Val"},
			&enumItem{1, "Var < Val"},
			&enumItem{2, "Var <= Val"},
			&enumItem{3, "Var > Val"},
			&enumItem{4, "Var >= Val"},
			&enumItem{5, "Var != Val"}}

		if unifiedValue != math.MinInt64 {
			selectedItem = items[unifiedValue>>13]
		}
		comparisonBox.SetItems(items)
		comparisonBox.SetSelectedItem(selectedItem)
	})

	simplifier.SetSpecialHandler("BinaryCodedDecimal", func() {
		slider := mode.selectedObjectsPropertiesPanel.NewSlider(key, "", func(currentValue, parameter uint32) uint32 {
			return uint32(util.ToBinaryCodedDecimal(uint16(parameter)))
		})
		slider.SetRange(0, 999)
		if unifiedValue != math.MinInt64 {
			slider.SetValue(int64(util.FromBinaryCodedDecimal(uint16(unifiedValue))))
		}
	})

	simplifier.SetSpecialHandler("LevelTexture", func() {
		selector := mode.selectedObjectsPropertiesPanel.NewTextureSelector(key, "", setUpdate(), mode.levelTextures)
		if unifiedValue != math.MinInt64 {
			selector.SetSelectedIndex(int(unifiedValue))
		}
	})

	describer(simplifier)
}

func (mode *LevelObjectsMode) levelTextures() []*graphics.BitmapTexture {
	ids := mode.levelAdapter.LevelTextureIDs()
	textures := make([]*graphics.BitmapTexture, len(ids))
	store := mode.context.ForGraphics().WorldTextureStore(dataModel.TextureLarge)

	for index, id := range ids {
		textures[index] = store.Texture(graphics.TextureKeyFromInt(id))
	}

	return textures
}

func (mode *LevelObjectsMode) selectedObjectIndices() []int {
	objectIndices := make([]int, len(mode.selectedObjects))
	for index, object := range mode.selectedObjects {
		objectIndices[index] = object.Index()
	}
	return objectIndices
}

func (mode *LevelObjectsMode) updateSelectedObjectsBaseProperties(modifier func(properties *dataModel.LevelObjectProperties)) {
	var properties dataModel.LevelObjectProperties
	modifier(&properties)
	mode.levelAdapter.RequestObjectPropertiesChange(mode.selectedObjectIndices(), &properties)
}

func (mode *LevelObjectsMode) updateSelectedObjectsClassPropertiesFiltered(key string, value uint32, update propertyUpdateFunction) {
	interpreterFactory := mode.interpreterFactory()

	for _, object := range mode.selectedObjects {
		objID := object.ID()
		resID := res.MakeObjectID(res.ObjectClass(objID.Class()), res.ObjectSubclass(objID.Subclass()), res.ObjectType(objID.Type()))
		var properties dataModel.LevelObjectProperties

		properties.ClassData = object.ClassData()
		interpreter := interpreterFactory(resID, properties.ClassData)
		subKeys := strings.Split(key, ".")
		valueIndex := len(subKeys) - 1
		for subIndex := 0; subIndex < valueIndex; subIndex++ {
			interpreter = interpreter.Refined(subKeys[subIndex])
		}
		subKey := subKeys[valueIndex]
		interpreter.Set(subKey, update(interpreter.Get(subKey), value))
		mode.levelAdapter.RequestObjectPropertiesChange([]int{object.Index()}, &properties)
	}
}

func (mode *LevelObjectsMode) deleteSelectedObjects() {
	mode.levelAdapter.RequestRemoveObjects(mode.selectedObjectIndices())
}

func (mode *LevelObjectsMode) onGameObjectsChanged() {
	newClassItems := make([]controls.ComboBoxItem, len(classNames))

	for index := range classNames {
		newClassItems[index] = &objectClassItem{index}
	}
	mode.newObjectClassBox.SetItems(newClassItems)
	mode.newObjectClassBox.SetSelectedItem(newClassItems[mode.newObjectID.Class()])
	mode.updateNewObjectClass(mode.newObjectID.Class())
}

func (mode *LevelObjectsMode) onNewObjectClassChanged(item controls.ComboBoxItem) {
	classItem := item.(*objectClassItem)
	mode.updateNewObjectClass(classItem.class)
	mode.updateNewObjectClassQuota()
}

func (mode *LevelObjectsMode) updateNewObjectClass(objectClass int) {
	typeItems := mode.objectItemsForClass(objectClass)

	mode.newObjectTypeBox.SetItems(typeItems)
	if len(typeItems) > 0 {
		mode.newObjectTypeBox.SetSelectedItem(typeItems[0])
		mode.onNewObjectTypeChanged(typeItems[0])
	} else {
		mode.newObjectTypeBox.SetSelectedItem(nil)
		mode.newObjectID = model.MakeObjectID(objectClass, 0, 0)
	}
}

func (mode *LevelObjectsMode) updateNewObjectClassQuota() {
	classes := len(maxObjectsPerClass)
	objectsPerClass := make(map[int]int)
	allObjects := mode.levelAdapter.LevelObjects(func(object *model.LevelObject) bool {
		objectsPerClass[object.ID().Class()]++
		return true
	})
	for class := 0; class < classes; class++ {
		mode.limitValues[class].SetText(fmt.Sprintf("%03d/%03d", objectsPerClass[class], maxObjectsPerClass[class]-1))
	}
	mode.limitValues[classes].SetText(fmt.Sprintf("%03d/%03d", len(allObjects), 871))
}

func (mode *LevelObjectsMode) onNewObjectTypeChanged(item controls.ComboBoxItem) {
	typeItem := item.(*objectTypeItem)
	mode.newObjectID = typeItem.id
}

func (mode *LevelObjectsMode) createNewObject(worldX, worldY float32) {
	mode.levelAdapter.RequestNewObject(worldX, worldY, mode.newObjectID)
}

func (mode *LevelObjectsMode) objectItemsForClass(objectClass int) []controls.ComboBoxItem {
	availableGameObjects := mode.context.ModelAdapter().ObjectsAdapter().ObjectsOfClass(objectClass)
	typeItems := make([]controls.ComboBoxItem, len(availableGameObjects))

	for index, gameObject := range availableGameObjects {
		typeItems[index] = &objectTypeItem{gameObject.ID(), gameObject.DisplayName()}
	}

	return typeItems
}
