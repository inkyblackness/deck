package modes

import (
	"fmt"
	"sort"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/data/interpreters"
	"github.com/inkyblackness/res/data/levelobj"

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

var classNames = []string{
	" 0: Weapons",
	" 1: AmmoClips",
	" 2: Projectiles",
	" 3: Explosives",
	" 4: Patches",
	" 5: Hardware",
	" 6: Software",
	" 7: Scenery",
	" 8: Items",
	" 9: Panels",
	"10: Barriers",
	"11: Animations",
	"12: Markers",
	"13: Containers",
	"14: Critters"}

type levelObjectProperty struct {
	title *controls.Label
	value *controls.Label
}

// LevelObjectsMode is a mode for level objects.
type LevelObjectsMode struct {
	context        Context
	levelAdapter   *model.LevelAdapter
	objectsAdapter *model.ObjectsAdapter

	displayFilter    func(*model.LevelObject) bool
	displayedObjects []*model.LevelObject

	mapDisplay *display.MapDisplay

	area       *ui.Area
	panel      *ui.Area
	panelRight ui.Anchor

	closestObjects              []*model.LevelObject
	closestObjectHighlightIndex int
	selectedObjects             []*model.LevelObject

	highlightedObjectIndexTitle *controls.Label
	highlightedObjectIndexValue *controls.Label

	selectedObjectsTitleLabel         *controls.Label
	selectedObjectsClassTitleLabel    *controls.Label
	selectedObjectsClassInfoLabel     *controls.Label
	selectedObjectsSubclassTitleLabel *controls.Label
	selectedObjectsSubclassInfoLabel  *controls.Label
	selectedObjectsTypeTitleLabel     *controls.Label
	selectedObjectsTypeInfoLabel      *controls.Label

	selectedObjectsPropertiesArea         *ui.Area
	selectedObjectsPropertiesPanelBuilder *controlPanelBuilder
	selectedObjectsPropertiesBottom       ui.Anchor
	selectedObjectsProperties             []*levelObjectProperty
}

// NewLevelObjectsMode returns a new instance.
func NewLevelObjectsMode(context Context, parent *ui.Area, mapDisplay *display.MapDisplay) *LevelObjectsMode {
	mode := &LevelObjectsMode{
		context:        context,
		levelAdapter:   context.ModelAdapter().ActiveLevel(),
		objectsAdapter: context.ModelAdapter().ObjectsAdapter(),
		displayFilter:  func(*model.LevelObject) bool { return true },

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
		mode.panelRight = ui.NewLimitedAnchor(minRight, maxRight, ui.NewOffsetAnchor(mode.area.Left(), 400))
		builder := ui.NewAreaBuilder()
		builder.SetParent(mode.area)
		builder.SetLeft(ui.NewOffsetAnchor(mode.area.Left(), 0))
		builder.SetTop(ui.NewOffsetAnchor(mode.area.Top(), 0))
		builder.SetRight(mode.panelRight)
		builder.SetBottom(ui.NewOffsetAnchor(mode.area.Bottom(), 0))
		builder.OnRender(func(area *ui.Area) {
			context.ForGraphics().RectangleRenderer().Fill(
				area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
				graphics.RGBA(0.7, 0.0, 0.7, 0.3))
		})

		builder.OnEvent(events.MouseButtonClickedEventType, ui.SilentConsumer)
		builder.OnEvent(events.MouseScrollEventType, ui.SilentConsumer)

		lastGrabX := float32(0.0)
		builder.OnEvent(events.MouseButtonDownEventType, func(area *ui.Area, event events.Event) bool {
			buttonEvent := event.(*events.MouseButtonEvent)
			if buttonEvent.Buttons() == env.MousePrimary {
				area.RequestFocus()
				lastGrabX, _ = buttonEvent.Position()
			}
			return true
		})
		builder.OnEvent(events.MouseButtonUpEventType, func(area *ui.Area, event events.Event) bool {
			buttonEvent := event.(*events.MouseButtonEvent)
			if buttonEvent.AffectedButtons() == env.MousePrimary {
				area.ReleaseFocus()
			}
			return true
		})
		builder.OnEvent(events.MouseMoveEventType, func(area *ui.Area, event events.Event) bool {
			moveEvent := event.(*events.MouseMoveEvent)
			if area.HasFocus() {
				newX, _ := moveEvent.Position()
				mode.panelRight.RequestValue(mode.panelRight.Value() + (newX - lastGrabX))
				lastGrabX = newX
			}
			return true
		})

		mode.panel = builder.Build()
	}
	{
		panelBuilder := newControlPanelBuilder(mode.panel, context.ControlFactory())

		mode.highlightedObjectIndexTitle, mode.highlightedObjectIndexValue = panelBuilder.addInfo("Highlighted Index")

		mode.selectedObjectsTitleLabel = panelBuilder.addTitle("Selected Object(s)")
		mode.selectedObjectsClassTitleLabel, mode.selectedObjectsClassInfoLabel = panelBuilder.addInfo("Class")
		mode.selectedObjectsSubclassTitleLabel, mode.selectedObjectsSubclassInfoLabel = panelBuilder.addInfo("Subclass")
		mode.selectedObjectsTypeTitleLabel, mode.selectedObjectsTypeInfoLabel = panelBuilder.addInfo("Type")

		mode.selectedObjectsPropertiesArea, mode.selectedObjectsPropertiesPanelBuilder =
			panelBuilder.addDynamicSection(true, func() ui.Anchor { return mode.selectedObjectsPropertiesBottom })
		mode.selectedObjectsPropertiesBottom = mode.selectedObjectsPropertiesArea.Top()
	}

	mode.levelAdapter.OnLevelObjectsChanged(mode.onLevelObjectsChanged)

	return mode
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
	if mode.area.IsVisible() {
		mode.updateDisplayedObjects()
	}
}

func (mode *LevelObjectsMode) updateDisplayedObjects() {
	mode.displayedObjects = mode.levelAdapter.LevelObjects(mode.displayFilter)
	mode.mapDisplay.SetDisplayedObjects(mode.displayedObjects)

	mode.closestObjects = nil
	mode.closestObjectHighlightIndex = 0
	mode.updateClosestObjectHighlight()
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
		mode.highlightedObjectIndexValue.SetText(fmt.Sprintf("%v", object.Index()))
		mode.mapDisplay.SetHighlightedObject(object)
	} else {
		mode.highlightedObjectIndexValue.SetText("")
		mode.mapDisplay.SetHighlightedObject(nil)
	}
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

	for _, object := range mode.selectedObjects {
		classUnifier.Add(object.ID().Class())
		subclassUnifier.Add(object.ID().Subclass())
		typeUnifier.Add(object.ID().Type())
	}
	unifiedClass := classUnifier.Value().(int)
	unifiedSubclass := subclassUnifier.Value().(int)
	unifiedType := typeUnifier.Value().(int)
	if unifiedClass != -1 {
		mode.selectedObjectsClassInfoLabel.SetText(classNames[unifiedClass])
	} else {
		mode.selectedObjectsClassInfoLabel.SetText("")
		unifiedSubclass = -1
	}
	if unifiedSubclass != -1 {
		mode.selectedObjectsSubclassInfoLabel.SetText(fmt.Sprintf("%v", unifiedSubclass))
	} else {
		mode.selectedObjectsSubclassInfoLabel.SetText("")
		unifiedType = -1
	}
	if unifiedType != -1 {
		mode.selectedObjectsTypeInfoLabel.SetText(fmt.Sprintf("%v", unifiedType))
	} else {
		mode.selectedObjectsTypeInfoLabel.SetText("")
	}
	mode.recreateLevelObjectProperties()
}

func (mode *LevelObjectsMode) recreateLevelObjectProperties() {
	for _, oldProperty := range mode.selectedObjectsProperties {
		oldProperty.title.Dispose()
		oldProperty.value.Dispose()
	}
	mode.selectedObjectsPropertiesPanelBuilder.reset()
	mode.selectedObjectsPropertiesBottom = mode.selectedObjectsPropertiesArea.Top()

	var newProperties = []*levelObjectProperty{}
	if len(mode.selectedObjects) > 0 {
		propertyUnifier := make(map[string]*util.ValueUnifier)
		propertyOrder := []string{}

		var unifyInterpreter func(string, *interpreters.Instance, bool, map[string]bool)
		unifyInterpreter = func(path string, interpreter *interpreters.Instance, first bool, thisKeys map[string]bool) {
			for _, key := range interpreter.Keys() {
				fullPath := path + key
				thisKeys[fullPath] = true
				if unifier, existing := propertyUnifier[fullPath]; existing || first {
					if !existing {
						unifier = util.NewValueUnifier(uint32(0xFFFFFFFF))
						propertyUnifier[fullPath] = unifier
						propertyOrder = append(propertyOrder, fullPath)
					}
					unifier.Add(interpreter.Get(key))
				}
			}
			for _, key := range interpreter.ActiveRefinements() {
				unifyInterpreter(path+key+".", interpreter.Refined(key), first, thisKeys)
			}
		}

		interpreterFactory := levelobj.ForRealWorld
		if mode.levelAdapter.IsCyberspace() {
			interpreterFactory = levelobj.ForCyberspace
		}

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
				property := &levelObjectProperty{}
				newProperties = append(newProperties, property)
				property.title, property.value = mode.selectedObjectsPropertiesPanelBuilder.addInfo(key)
				property.value.SetText(fmt.Sprintf("%v", unifier.Value().(uint32)))
				mode.selectedObjectsPropertiesBottom = mode.selectedObjectsPropertiesPanelBuilder.bottom()
			}
		}
	}
	mode.selectedObjectsProperties = newProperties
}
