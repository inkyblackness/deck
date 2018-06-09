package modes

import (
	"fmt"
	"math"
	"sort"
	"strings"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data"
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

type interpreterFactoryFunc func(res.ObjectID, []byte) *interpreters.Instance

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
	posInfoTileX  *controls.Label
	posInfoFineX  *controls.Label
	posInfoTileY  *controls.Label
	posInfoFineY  *controls.Label
	posInfoZ      *controls.Label

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

	selectedObjectsExtraPropertiesItem  *tabItem
	selectedObjectsPropertiesExtraArea  *ui.Area
	selectedObjectsPropertiesExtraPanel *propertyPanel

	selectedObjectsClassPropertiesItem  *tabItem
	selectedObjectsPropertiesMainArea   *ui.Area
	selectedObjectsPropertiesHeaderArea *ui.Area
	selectedObjectsPropertiesPanel      *propertyPanel

	blockPuzzleArea  *ui.Area
	blockPuzzleCells [][]*controls.Label
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
		builder.OnEvent(events.MouseMoveEventType, mode.onMouseMoved)
		builder.OnEvent(events.MouseScrollEventType, mode.onMouseScrolled)
		builder.OnEvent(events.MouseButtonClickedEventType, mode.onMouseButtonClicked)
		mode.area = builder.Build()
	}
	{
		minRight := ui.NewOffsetAnchor(mode.area.Left(), scaled(100))
		maxRight := ui.NewRelativeAnchor(mode.area.Left(), mode.area.Right(), 0.5)
		mode.propertiesPanelRight = ui.NewLimitedAnchor(minRight, maxRight, ui.NewOffsetAnchor(mode.area.Left(), scaled(400)))
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
			mode.updateSelectedObjectsClassPropertiesRaw(func(objectID model.ObjectID, classData []byte) {
				for index := range classData {
					classData[index] = 0x00
				}
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

		{
			extraPropertiesBottomResolver := func() ui.Anchor { return mode.selectedObjectsPropertiesExtraPanel.Bottom() }
			var extraPanelBuilder *controlPanelBuilder
			mode.selectedObjectsPropertiesExtraArea, extraPanelBuilder =
				panelBuilder.addDynamicSection(false, extraPropertiesBottomResolver)

			mode.selectedObjectsPropertiesExtraPanel = newPropertyPanel(extraPanelBuilder,
				mode.updateSelectedObjectsExtraProperties, mode.objectItemsForClass)
		}

		{
			classPropertiesBottomResolver := func() ui.Anchor { return mode.selectedObjectsPropertiesPanel.Bottom() }
			var mainClassPanelBuilder *controlPanelBuilder
			mode.selectedObjectsPropertiesMainArea, mainClassPanelBuilder =
				panelBuilder.addDynamicSection(false, classPropertiesBottomResolver)
			mode.selectedObjectsPropertiesHeaderArea, _ = mainClassPanelBuilder.addSection(true)

			mode.selectedObjectsPropertiesPanel = newPropertyPanel(mainClassPanelBuilder,
				mode.updateSelectedObjectsClassProperties, mode.objectItemsForClass)
		}

		mode.selectedObjectsBasePropertiesItem = &tabItem{mode.selectedObjectsBasePropertiesArea, "Base Properties"}
		mode.selectedObjectsExtraPropertiesItem = &tabItem{mode.selectedObjectsPropertiesExtraArea, "Extra Properties"}
		mode.selectedObjectsClassPropertiesItem = &tabItem{mode.selectedObjectsPropertiesMainArea, "Class Properties"}
		propertiesTabItems := []controls.ComboBoxItem{
			mode.selectedObjectsBasePropertiesItem,
			mode.selectedObjectsExtraPropertiesItem,
			mode.selectedObjectsClassPropertiesItem}
		mode.selectedObjectsPropertiesBox.SetItems(propertiesTabItems)
		initialItem := mode.selectedObjectsBasePropertiesItem
		mode.selectedObjectsPropertiesBox.SetSelectedItem(initialItem)
		mode.onSelectedPropertiesDisplayChanged(initialItem)
	}
	{
		cellSize := scaled(25)
		cellsPerSide := 9
		areaSize := cellSize * float32(cellsPerSide)
		areaLeft := ui.NewOffsetAnchor(mode.propertiesPanelRight, scaled(5))
		areaBottom := ui.NewOffsetAnchor(mode.area.Bottom(), 0)
		areaTop := ui.NewOffsetAnchor(areaBottom, -areaSize)
		{

			builder := ui.NewAreaBuilder()
			builder.SetParent(mode.area)
			builder.SetLeft(areaLeft)
			builder.SetTop(areaTop)
			builder.SetRight(ui.NewOffsetAnchor(areaLeft, areaSize))
			builder.SetBottom(areaBottom)
			builder.SetVisible(false)
			builder.OnRender(func(area *ui.Area) {
				context.ForGraphics().RectangleRenderer().Fill(
					area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
					graphics.RGBA(0.5, 0.0, 0.5, 0.7))
			})
			builder.OnEvent(events.MouseMoveEventType, ui.SilentConsumer)
			builder.OnEvent(events.MouseButtonUpEventType, ui.SilentConsumer)
			builder.OnEvent(events.MouseButtonDownEventType, ui.SilentConsumer)
			builder.OnEvent(events.MouseButtonClickedEventType, ui.SilentConsumer)
			builder.OnEvent(events.MouseScrollEventType, mode.onBlockPuzzleScrolled)
			mode.blockPuzzleArea = builder.Build()
		}

		mode.blockPuzzleCells = make([][]*controls.Label, cellsPerSide)
		for cellRow := 0; cellRow < cellsPerSide; cellRow++ {
			mode.blockPuzzleCells[cellRow] = make([]*controls.Label, cellsPerSide)
			for cellColumn := 0; cellColumn < cellsPerSide; cellColumn++ {
				labelBuilder := mode.context.ControlFactory().ForLabel()
				cellTop := ui.NewOffsetAnchor(areaTop, cellSize*float32(cellRow))
				cellLeft := ui.NewOffsetAnchor(areaLeft, cellSize*float32(cellColumn))
				labelBuilder.SetParent(mode.blockPuzzleArea)
				labelBuilder.SetTop(cellTop)
				labelBuilder.SetLeft(cellLeft)
				labelBuilder.SetRight(ui.NewOffsetAnchor(cellLeft, cellSize))
				labelBuilder.SetBottom(ui.NewOffsetAnchor(cellTop, cellSize))
				cell := labelBuilder.Build()
				mode.blockPuzzleCells[cellRow][cellColumn] = cell
			}
		}
	}
	{
		builder := ui.NewAreaBuilder()
		builder.SetParent(mode.area)
		builder.SetLeft(ui.NewOffsetAnchor(mode.area.Right(), scaled(-126)))
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
		{
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
		{
			panelBuilder.addInfo("")
			_, mode.posInfoTileX = panelBuilder.addInfo("TileX")
			_, mode.posInfoFineX = panelBuilder.addInfo("FineX")
			_, mode.posInfoTileY = panelBuilder.addInfo("TileY")
			_, mode.posInfoFineY = panelBuilder.addInfo("FineY")
			_, mode.posInfoZ = panelBuilder.addInfo("Z")
		}
	}

	mode.levelAdapter.OnIDChanged(func() {
		mode.setSelectedObjects(nil)
	})
	mode.levelAdapter.OnLevelPropertiesChanged(func() {
		mode.selectedObjectsZValue.SetValueFormatter(mode.objectZToString)
	})
	mode.levelAdapter.OnLevelObjectsChanged(mode.onLevelObjectsChanged)
	mode.context.ModelAdapter().ObjectsAdapter().OnObjectsChanged(mode.onGameObjectsChanged)

	return mode
}

func (mode *LevelObjectsMode) objectZToString(value int64) string {
	return heightToString(mode.levelAdapter.HeightShift(), value, 256.0)
}

func (mode *LevelObjectsMode) objectZToShortString(value int64) string {
	return fmt.Sprintf("%.3f", heightToValue(mode.levelAdapter.HeightShift(), value, 256.0))
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
	mode.selectedObjectsExtraPropertiesItem.page.SetVisible(false)
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
		atGrid := false
		if keys.Modifier(mouseEvent.Modifier()) == keys.ModShift {
			atGrid = true
		}
		mode.createNewObject(worldX, worldY, atGrid)
	}

	return
}

func (mode *LevelObjectsMode) updateClosestDisplayedObjects(worldX, worldY float32) {
	type resultEntry struct {
		distance float32
		object   *model.LevelObject
	}
	var entries []*resultEntry
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
	if len(mode.closestObjects) == 0 {
		mode.showTileLocation(worldX, worldY)
	}
}

func (mode *LevelObjectsMode) updateClosestObjectHighlight() {
	if len(mode.closestObjects) > 0 {
		object := mode.closestObjects[mode.closestObjectHighlightIndex]
		mode.highlightedObjectInfoValue.SetText(fmt.Sprintf("%v: %v (%v)", object.Index(), object.ID(), mode.objectDisplayName(object.ID())))
		mode.mapDisplay.SetHighlightedObject(object)

		mode.posInfoTileX.SetText(fmt.Sprintf("%v", object.TileX()))
		mode.posInfoFineX.SetText(fmt.Sprintf("%v", object.FineX()))
		mode.posInfoTileY.SetText(fmt.Sprintf("%v", object.TileY()))
		mode.posInfoFineY.SetText(fmt.Sprintf("%v", object.FineY()))
		mode.posInfoZ.SetText(mode.objectZToShortString(int64(object.Z())))
	} else {
		mode.highlightedObjectInfoValue.SetText("")
		mode.mapDisplay.SetHighlightedObject(nil)
	}
}

func (mode *LevelObjectsMode) showTileLocation(worldX, worldY float32) {
	integerX, integerY := int(worldX), int(worldY)
	tileX, fineX := integerX>>8, integerX&0xFF
	tileY, fineY := integerY>>8, integerY&0xFF
	tile := mode.levelAdapter.TileMap().Tile(model.TileCoordinateOf(tileX, tileY))
	if tile != nil {
		mode.posInfoTileX.SetText(fmt.Sprintf("%v", tileX))
		mode.posInfoFineX.SetText(fmt.Sprintf("%v", fineX))
		mode.posInfoTileY.SetText(fmt.Sprintf("%v", tileY))
		mode.posInfoFineY.SetText(fmt.Sprintf("%v", fineY))
		mode.posInfoZ.SetText(mode.tileHeightUnitToShortString(int64(*tile.Properties().FloorHeight)))
	} else {
		mode.posInfoTileX.SetText("")
		mode.posInfoFineX.SetText("")
		mode.posInfoTileY.SetText("")
		mode.posInfoFineY.SetText("")
		mode.posInfoZ.SetText("")
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

	mode.recreateLevelObjectExtraProperties()
	mode.recreateLevelObjectClassProperties()
}

func (mode *LevelObjectsMode) recreateLevelObjectExtraProperties() {
	interpreterFactory := mode.extraInterpreterFactory()

	mode.recreateLevelObjectProperties(mode.selectedObjectsPropertiesExtraPanel,
		func(objID res.ObjectID, object *model.LevelObject) *interpreters.Instance {
			return interpreterFactory(objID, object.ExtraData())
		})
}

func (mode *LevelObjectsMode) recreateLevelObjectClassProperties() {
	interpreterFactory := mode.classInterpreterFactory()

	mode.recreateLevelObjectProperties(mode.selectedObjectsPropertiesPanel,
		func(objID res.ObjectID, object *model.LevelObject) *interpreters.Instance {
			return interpreterFactory(objID, object.ClassData())
		})
	mode.updateBlockPuzzleArea(func(objID res.ObjectID, object *model.LevelObject) *interpreters.Instance {
		return interpreterFactory(objID, object.ClassData())
	})
}

func (mode *LevelObjectsMode) recreateLevelObjectProperties(panel *propertyPanel,
	interpreterFactory func(res.ObjectID, *model.LevelObject) *interpreters.Instance) {
	panel.Reset()

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

		for index, object := range mode.selectedObjects {
			objID := object.ID()
			resID := res.MakeObjectID(res.ObjectClass(objID.Class()), res.ObjectSubclass(objID.Subclass()), res.ObjectType(objID.Type()))
			interpreter := interpreterFactory(resID, object)
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
				mode.createPropertyControls(panel, key, unifier.Value().(int64), propertyDescribers[key])
			}
		}
	}
}

func (mode *LevelObjectsMode) updateBlockPuzzleArea(interpreterFactory func(res.ObjectID, *model.LevelObject) *interpreters.Instance) {
	isBlockPuzzle := false

	for _, row := range mode.blockPuzzleCells {
		for _, cell := range row {
			cell.SetText("")
		}
	}

	if len(mode.selectedObjects) == 1 {
		singleObject := mode.selectedObjects[0]
		singleObjID := singleObject.ID()
		singleResID := res.MakeObjectID(res.ObjectClass(singleObjID.Class()), res.ObjectSubclass(singleObjID.Subclass()), res.ObjectType(singleObjID.Type()))
		singleInterpreter := interpreterFactory(singleResID, singleObject)
		isBlockPuzzle = singleInterpreter.Refined("Puzzle").Get("Type") == 0x10

		if isBlockPuzzle {
			blockInfo := singleInterpreter.Refined("Puzzle").Refined("Block")
			blockLayout := blockInfo.Get("Layout")
			blockPuzzleDataIndex := blockInfo.Get("StateStoreObjectIndex")
			dataObject := mode.levelAdapter.LevelObject(int(blockPuzzleDataIndex))

			if (dataObject != nil) && (dataObject.ID() == model.MakeObjectID(12, 0, 1)) {
				raw := dataObject.ClassData()[6 : 6+16]
				blockWidth := int((blockLayout >> 20) & 7)
				blockHeight := int((blockLayout >> 24) & 7)
				state := data.NewBlockPuzzleState(raw, blockHeight, blockWidth)
				startRow := 1 + (7-blockHeight)/2
				startColumn := 1 + (7-blockWidth)/2
				placeConnector := func(side, offset int, text string) {
					xOffsets := []int{offset, offset, -1, blockWidth}
					yOffsets := []int{-1, blockHeight, offset, offset}
					y := startRow + yOffsets[side]
					x := startColumn + xOffsets[side]

					if (x >= 0) && (x < 9) && (y >= 0) && (y < 9) {
						mode.blockPuzzleCells[y][x].SetText(text)
					}
				}
				stateMapping := []string{".", "X", "+", "(+)", "F", "(F)", "H", "(H)"}

				for row := 0; row < blockHeight; row++ {
					for col := 0; col < blockWidth; col++ {
						value := state.CellValue(row, col)
						mode.blockPuzzleCells[startRow+row][startColumn+col].SetText(stateMapping[value])
					}
				}
				placeConnector(int((blockLayout>>7)&3), int((blockLayout>>4)&7), "S")
				placeConnector(int((blockLayout>>15)&3), int((blockLayout>>12)&7), "D")

			} else {
				isBlockPuzzle = false
			}
		}
	}
	mode.blockPuzzleArea.SetVisible(isBlockPuzzle)
}

func (mode *LevelObjectsMode) onBlockPuzzleScrolled(area *ui.Area, event events.Event) bool {
	interpreterFactory := mode.classInterpreterFactory()
	mouseEvent := event.(*events.MouseScrollEvent)

	mouseX, mouseY := mouseEvent.Position()
	_, scrollY := mouseEvent.Deltas()
	areaTop := area.Top().Value()
	areaLeft := area.Left().Value()
	cellClickX := int((mouseX - areaLeft) / (area.Right().Value() - areaLeft) * float32(9))
	cellClickY := int((mouseY - areaTop) / (area.Bottom().Value() - areaTop) * float32(9))
	scrollOffset := 1

	if scrollY < 0 {
		scrollOffset = -1
	}

	singleObject := mode.selectedObjects[0]
	singleObjID := singleObject.ID()
	singleResID := res.MakeObjectID(res.ObjectClass(singleObjID.Class()), res.ObjectSubclass(singleObjID.Subclass()), res.ObjectType(singleObjID.Type()))
	singleInterpreter := interpreterFactory(singleResID, singleObject.ClassData())
	isBlockPuzzle := singleInterpreter.Refined("Puzzle").Get("Type") == 0x10

	if isBlockPuzzle {
		blockInfo := singleInterpreter.Refined("Puzzle").Refined("Block")
		blockLayout := blockInfo.Get("Layout")
		blockPuzzleDataIndex := blockInfo.Get("StateStoreObjectIndex")
		dataObject := mode.levelAdapter.LevelObject(int(blockPuzzleDataIndex))

		if (dataObject != nil) && (dataObject.ID() == model.MakeObjectID(12, 0, 1)) {
			blockWidth := int((blockLayout >> 20) & 7)
			blockHeight := int((blockLayout >> 24) & 7)
			startRow := 1 + (7-blockHeight)/2
			startColumn := 1 + (7-blockWidth)/2

			clickRow := cellClickY - startRow
			clickCol := cellClickX - startColumn
			if (clickRow >= 0) && (clickRow < blockHeight) && (clickCol >= 0) && (clickCol < blockWidth) {
				mode.updateObjectClassPropertiesRaw(dataObject, func(classData []byte) {
					state := data.NewBlockPuzzleState(classData[6:6+16], blockHeight, blockWidth)
					oldValue := state.CellValue(clickRow, clickCol)
					state.SetCellValue(clickRow, clickCol, (8+oldValue+scrollOffset)%8)
				})
			}
		}
	}

	return true
}

func (mode *LevelObjectsMode) classInterpreterFactory() interpreterFactoryFunc {
	factory := levelobj.ForRealWorld
	if mode.levelAdapter.IsCyberspace() {
		factory = levelobj.ForCyberspace
	}
	return factory
}

func (mode *LevelObjectsMode) extraInterpreterFactory() interpreterFactoryFunc {
	factory := levelobj.RealWorldExtra
	if mode.levelAdapter.IsCyberspace() {
		factory = levelobj.CyberspaceExtra
	}
	return factory
}

func (mode *LevelObjectsMode) createPropertyControls(panel *propertyPanel, key string, unifiedValue int64, describer func(*interpreters.Simplifier)) {
	simplifier := panel.NewSimplifier(key, unifiedValue)

	simplifier.SetObjectIndexHandler(func() {
		slider := panel.NewSlider(key, "", setUpdate())
		slider.SetRange(0, 871)
		if unifiedValue != math.MinInt64 {
			slider.SetValue(unifiedValue)
		}
	})

	addVariableKey := func() {
		typeBox := panel.NewComboBox(key, "Type", maskedUpdate(0, 0x1000))
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

		indexSlider := panel.NewSlider(key, "Index", maskedUpdate(0, 0x1FF))
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

		comparisonBox := panel.NewComboBox(key, "Check", maskedUpdate(13, 0xE000))
		var selectedItem controls.ComboBoxItem
		items := []controls.ComboBoxItem{
			&enumItem{0, "Var == Val"},
			&enumItem{1, "Var < Val"},
			&enumItem{2, "Var <= Val"},
			&enumItem{3, "Var > Val"},
			&enumItem{4, "Var >= Val"},
			&enumItem{5, "Var != Val"},
			&enumItem{6, "Val > Random(0..254)"}}

		if unifiedValue != math.MinInt64 {
			selectedItem = items[unifiedValue>>13]
		}
		comparisonBox.SetItems(items)
		comparisonBox.SetSelectedItem(selectedItem)
	})

	simplifier.SetSpecialHandler("BinaryCodedDecimal", func() {
		slider := panel.NewSlider(key, "", func(currentValue, parameter uint32) uint32 {
			return uint32(util.ToBinaryCodedDecimal(uint16(parameter)))
		})
		slider.SetRange(0, 999)
		slider.SetValueFormatter(func(value int64) string {
			return fmt.Sprintf("%03d", value)
		})
		if unifiedValue != math.MinInt64 {
			slider.SetValue(int64(util.FromBinaryCodedDecimal(uint16(unifiedValue))))
		}
	})

	simplifier.SetSpecialHandler("LevelTexture", func() {
		selector := panel.NewTextureSelector(key, "", setUpdate(), mode.levelTextures)
		if unifiedValue != math.MinInt64 {
			selector.SetSelectedIndex(int(unifiedValue))
		}
	})

	simplifier.SetSpecialHandler("MaterialOrLevelTexture", func() {
		selectionBox := panel.NewComboBox(key, "Type", maskedUpdate(7, 0xFF))
		var selectedItem controls.ComboBoxItem
		selectedType := -1
		selectedIndex := 0
		items := []controls.ComboBoxItem{
			&enumItem{0, "Material"},
			&enumItem{1, "Level texture"}}

		if unifiedValue != math.MinInt64 {
			selectedType = int(unifiedValue >> 7)
			selectedIndex = int(unifiedValue & 0x7F)
			selectedItem = items[selectedType]
		}
		selectionBox.SetItems(items)
		selectionBox.SetSelectedItem(selectedItem)

		if selectedType == 0 {
			slider := panel.NewSlider(key, "Material", maskedUpdate(0, 0x7F))
			slider.SetRange(0, 127)
			slider.SetValue(int64(selectedIndex))
		} else if selectedType == 1 {
			selector := panel.NewTextureSelector(key, "", maskedUpdate(0, 0x7F), mode.levelTextures)
			selector.SetSelectedIndex(selectedIndex)
		}
	})

	simplifier.SetSpecialHandler("ObjectHeight", func() {
		slider := panel.NewSlider(key, "", func(currentValue, parameter uint32) uint32 {
			return parameter
		})
		slider.SetRange(0, 255)
		slider.SetValueFormatter(mode.objectZToString)
		if unifiedValue != math.MinInt64 {
			slider.SetValue(int64(unifiedValue))
		}
	})
	simplifier.SetSpecialHandler("MoveTileHeight", func() {
		slider := panel.NewSlider(key, "", func(currentValue, parameter uint32) uint32 {
			return parameter
		})
		slider.SetRange(0, 0x0FFF)
		slider.SetValueFormatter(mode.moveTileHeightUnitToString)
		if unifiedValue != math.MinInt64 {
			slider.SetValue(int64(unifiedValue))
		}
	})

	describer(simplifier)
}

func (mode *LevelObjectsMode) moveTileHeightUnitToString(value int64) (result string) {
	if (value >= 0) && (value < 32) {
		result = mode.tileHeightUnitToString(value)
	} else {
		result = fmt.Sprintf("Don't change  - raw: 0x%04X", value)
	}
	return
}

func (mode *LevelObjectsMode) tileHeightUnitToString(value int64) string {
	return heightToString(mode.levelAdapter.HeightShift(), value, 32.0)
}

func (mode *LevelObjectsMode) tileHeightUnitToShortString(value int64) string {
	return fmt.Sprintf("%.3f", heightToValue(mode.levelAdapter.HeightShift(), value, 32.0))
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

func (mode *LevelObjectsMode) updateSelectedObjectsClassProperties(key string, value uint32, update propertyUpdateFunction) {
	interpreterFactory := mode.classInterpreterFactory()

	mode.updateSelectedObjectsClassPropertiesRaw(func(objectID model.ObjectID, classData []byte) {
		mode.setInterpreterValue(interpreterFactory, objectID, classData, key, value, update)
	})
}

func (mode *LevelObjectsMode) updateSelectedObjectsClassPropertiesRaw(modifier func(objectID model.ObjectID, classData []byte)) {
	for _, object := range mode.selectedObjects {
		mode.updateObjectClassPropertiesRaw(object, func(classData []byte) {
			modifier(object.ID(), classData)
		})
	}
}

func (mode *LevelObjectsMode) updateObjectClassPropertiesRaw(object *model.LevelObject, modifier func(classData []byte)) {
	var properties dataModel.LevelObjectProperties

	properties.ClassData = object.ClassData()
	modifier(properties.ClassData)
	mode.levelAdapter.RequestObjectPropertiesChange([]int{object.Index()}, &properties)
}

func (mode *LevelObjectsMode) updateSelectedObjectsExtraProperties(key string, value uint32, update propertyUpdateFunction) {
	interpreterFactory := mode.extraInterpreterFactory()

	for _, object := range mode.selectedObjects {
		var properties dataModel.LevelObjectProperties

		properties.ExtraData = object.ExtraData()
		mode.setInterpreterValue(interpreterFactory, object.ID(), properties.ExtraData, key, value, update)
		mode.levelAdapter.RequestObjectPropertiesChange([]int{object.Index()}, &properties)
	}
}

func (mode *LevelObjectsMode) setInterpreterValue(interpreterFactory interpreterFactoryFunc, objID model.ObjectID, data []byte,
	key string, value uint32, update propertyUpdateFunction) {
	resID := res.MakeObjectID(res.ObjectClass(objID.Class()), res.ObjectSubclass(objID.Subclass()), res.ObjectType(objID.Type()))

	interpreter := interpreterFactory(resID, data)
	subKeys := strings.Split(key, ".")
	valueIndex := len(subKeys) - 1
	for subIndex := 0; subIndex < valueIndex; subIndex++ {
		interpreter = interpreter.Refined(subKeys[subIndex])
	}
	subKey := subKeys[valueIndex]
	interpreter.Set(subKey, update(interpreter.Get(subKey), value))
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

func (mode *LevelObjectsMode) createNewObject(worldX, worldY float32, atGrid bool) {
	mode.levelAdapter.RequestNewObject(worldX, worldY, mode.newObjectID, atGrid)
}

func (mode *LevelObjectsMode) objectItemsForClass(objectClass int) []controls.ComboBoxItem {
	availableGameObjects := mode.context.ModelAdapter().ObjectsAdapter().ObjectsOfClass(objectClass)
	typeItems := make([]controls.ComboBoxItem, len(availableGameObjects))

	for index, gameObject := range availableGameObjects {
		typeItems[index] = &objectTypeItem{gameObject.ID(), gameObject.DisplayName()}
	}

	return typeItems
}
