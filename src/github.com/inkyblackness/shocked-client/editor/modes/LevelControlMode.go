package modes

import (
	"fmt"

	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/shocked-client/editor/display"
	"github.com/inkyblackness/shocked-client/editor/model"
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/graphics/controls"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"

	dataModel "github.com/inkyblackness/shocked-model"
)

type levelPropertyItem struct {
	displayString string
	modifier      func(properties *dataModel.LevelProperties)
	formatter     controls.SliderValueFormatter
}

func (item *levelPropertyItem) String() string {
	return item.displayString
}

func lbpValueFormatter(value int64) string {
	return fmt.Sprintf("%v.%v LBP", value/2, (value%2)*5)
}

// LevelControlMode is a mode for archive level control.
type LevelControlMode struct {
	context      Context
	levelAdapter *model.LevelAdapter

	mapDisplay *display.MapDisplay

	area *ui.Area

	activeLevelLabel *controls.Label
	activeLevelBox   *controls.ComboBox

	heightShiftLabel *controls.Label
	heightShiftBox   *controls.ComboBox

	realWorldProperties *ui.Area

	levelTexturesLabel       *controls.Label
	levelTexturesSelector    *controls.TextureSelector
	currentLevelTextureIndex int
	worldTexturesLabel       *controls.Label
	worldTexturesSelector    *controls.TextureSelector
	worldTexturesIDLabel     *controls.Label
	worldTexturesIDSlider    *controls.Slider

	selectedSurveillanceIndex    int
	surveillanceIndexLabel       *controls.Label
	surveillanceIndexBox         *controls.ComboBox
	surveillanceSourceLabel      *controls.Label
	surveillanceSourceSlider     *controls.Slider
	surveillanceDeathwatchLabel  *controls.Label
	surveillanceDeathwatchSlider *controls.Slider

	ceilingEffectLabel       *controls.Label
	ceilingEffectBox         *controls.ComboBox
	ceilingEffectLevelLabel  *controls.Label
	ceilingEffectLevelSlider *controls.Slider

	floorEffectLabel       *controls.Label
	floorEffectBox         *controls.ComboBox
	floorEffectLevelLabel  *controls.Label
	floorEffectLevelSlider *controls.Slider

	selectedAnimationGroupIndex int
	animationGroupIndexLabel    *controls.Label
	animationGroupIndexBox      *controls.ComboBox
	animationGroupTimeLabel     *controls.Label
	animationGroupTimeSlider    *controls.Slider
	animationGroupFramesLabel   *controls.Label
	animationGroupFramesSlider  *controls.Slider
	animationGroupTypeLabel     *controls.Label
	animationGroupTypeBox       *controls.ComboBox
	animationGroupTypeItems     map[int]controls.ComboBoxItem
}

// NewLevelControlMode returns a new instance.
func NewLevelControlMode(context Context, parent *ui.Area, mapDisplay *display.MapDisplay) *LevelControlMode {
	mode := &LevelControlMode{
		context:                     context,
		levelAdapter:                context.ModelAdapter().ActiveLevel(),
		mapDisplay:                  mapDisplay,
		currentLevelTextureIndex:    -1,
		selectedSurveillanceIndex:   -1,
		selectedAnimationGroupIndex: -1}

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
			mode.activeLevelLabel, mode.activeLevelBox = panelBuilder.addComboProperty("Active Level", func(item controls.ComboBoxItem) {
				context.ModelAdapter().RequestActiveLevel(item.(int))
			})

			adapter := context.ModelAdapter()
			activeLevelAdapter := adapter.ActiveLevel()
			activeLevelAdapter.OnIDChanged(func() {
				mode.activeLevelBox.SetSelectedItem(activeLevelAdapter.ID())
			})
			adapter.OnAvailableLevelsChanged(func() {
				ids := adapter.AvailableLevelIDs()
				items := make([]controls.ComboBoxItem, len(ids))
				for index, id := range ids {
					items[index] = id
				}
				mode.activeLevelBox.SetItems(items)
			})
		}
		{
			mode.heightShiftLabel, mode.heightShiftBox = panelBuilder.addComboProperty("Tile Height", mode.onHeightShiftChanged)
			heightShiftItems := make([]controls.ComboBoxItem, 8)

			heightShiftItems[0] = &enumItem{0, "32 Tiles"}
			heightShiftItems[1] = &enumItem{1, "16 Tiles"}
			heightShiftItems[2] = &enumItem{2, "8 Tiles"}
			heightShiftItems[3] = &enumItem{3, "4 Tiles"}
			heightShiftItems[4] = &enumItem{4, "2 Tiles"}
			heightShiftItems[5] = &enumItem{5, "1 Tile"}
			heightShiftItems[6] = &enumItem{6, "1/2 Tile"}
			heightShiftItems[7] = &enumItem{7, "1/4 Tile"}
			mode.heightShiftBox.SetItems(heightShiftItems)
			mode.levelAdapter.OnLevelPropertiesChanged(func() {
				heightShift := mode.levelAdapter.HeightShift()
				if (heightShift >= 0) && (heightShift < len(heightShiftItems)) {
					mode.heightShiftBox.SetSelectedItem(heightShiftItems[heightShift])
				} else {
					mode.heightShiftBox.SetSelectedItem(nil)
				}
			})
		}

		{
			var realWorldBuilder *controlPanelBuilder
			mode.realWorldProperties, realWorldBuilder = panelBuilder.addSection(false)

			{
				mode.levelTexturesLabel, mode.levelTexturesSelector = realWorldBuilder.addTextureProperty("Level Textures",
					mode.levelTextures, mode.onSelectedLevelTextureChanged)
				mode.worldTexturesLabel, mode.worldTexturesSelector = realWorldBuilder.addTextureProperty("World Textures",
					mode.worldTextures, mode.onSelectedWorldTextureChanged)
				mode.worldTexturesIDLabel, mode.worldTexturesIDSlider = realWorldBuilder.addSliderProperty("World Texture ID",
					mode.onSelectedWorldTextureIDChanged)

				textureAdapter := mode.context.ModelAdapter().TextureAdapter()
				textureAdapter.OnGameTexturesChanged(func() {
					mode.worldTexturesIDSlider.SetRange(0, int64(textureAdapter.WorldTextureCount()-1))
				})
			}
			{
				mode.surveillanceIndexLabel, mode.surveillanceIndexBox =
					realWorldBuilder.addComboProperty("Surveillance Object", mode.onSurveillanceIndexChanged)
				mode.surveillanceSourceLabel, mode.surveillanceSourceSlider =
					realWorldBuilder.addSliderProperty("Surveillance Source", mode.onSurveillanceSourceChanged)
				mode.surveillanceDeathwatchLabel, mode.surveillanceDeathwatchSlider =
					realWorldBuilder.addSliderProperty("Surveillance Deathwatch", mode.onSurveillanceDeathwatchChanged)
				mode.surveillanceSourceSlider.SetRange(0, 871)
				mode.surveillanceDeathwatchSlider.SetRange(0, 871)

				mode.levelAdapter.OnLevelSurveillanceChanged(mode.onLevelSurveillanceChanged)
			}
			{
				mode.ceilingEffectLabel, mode.ceilingEffectBox =
					realWorldBuilder.addComboProperty("Ceiling Effect", mode.onLevelCeilingPropertyBoxChanged)
				mode.ceilingEffectLevelLabel, mode.ceilingEffectLevelSlider =
					realWorldBuilder.addSliderProperty("Ceiling Effect Level", mode.onCeilingEffectLevelChanged)

				noEffectItem := &levelPropertyItem{"None",
					func(properties *dataModel.LevelProperties) { properties.CeilingHasRadiation = boolAsPointer(false) },
					controls.DefaultSliderValueFormatter}
				radiationEffectItem := &levelPropertyItem{"Radiation",
					func(properties *dataModel.LevelProperties) { properties.CeilingHasRadiation = boolAsPointer(true) },
					lbpValueFormatter}
				ceilingItems := []controls.ComboBoxItem{noEffectItem, radiationEffectItem}

				mode.ceilingEffectBox.SetItems(ceilingItems)
				mode.ceilingEffectLevelSlider.SetRange(0, 255)

				mode.levelAdapter.OnLevelPropertiesChanged(func() {
					radiation, level := mode.levelAdapter.CeilingEffect()
					item := noEffectItem

					if radiation {
						item = radiationEffectItem
					}
					mode.ceilingEffectBox.SetSelectedItem(item)
					mode.ceilingEffectLevelSlider.SetValue(int64(level))
					mode.ceilingEffectLevelSlider.SetValueFormatter(item.formatter)
				})
			}
			{
				mode.floorEffectLabel, mode.floorEffectBox =
					realWorldBuilder.addComboProperty("Floor Effect", mode.onLevelFloorPropertyBoxChanged)
				mode.floorEffectLevelLabel, mode.floorEffectLevelSlider =
					realWorldBuilder.addSliderProperty("Floor Effect Level", mode.onFloorEffectLevelChanged)

				noEffectItem := &levelPropertyItem{"None", func(properties *dataModel.LevelProperties) {
					properties.FloorHasBiohazard = boolAsPointer(false)
					properties.FloorHasGravity = boolAsPointer(false)
				}, controls.DefaultSliderValueFormatter}
				gravityEffectItem := &levelPropertyItem{"Gravity", func(properties *dataModel.LevelProperties) {
					properties.FloorHasBiohazard = boolAsPointer(false)
					properties.FloorHasGravity = boolAsPointer(true)
				}, func(value int64) string { return fmt.Sprintf("%v%%", value*25) }}
				biohazardEffectItem := &levelPropertyItem{"Biohazard", func(properties *dataModel.LevelProperties) {
					properties.FloorHasBiohazard = boolAsPointer(true)
					properties.FloorHasGravity = boolAsPointer(false)
				}, lbpValueFormatter}
				floorItems := []controls.ComboBoxItem{noEffectItem, gravityEffectItem, biohazardEffectItem}

				mode.floorEffectBox.SetItems(floorItems)
				mode.floorEffectLevelSlider.SetRange(0, 255)

				mode.levelAdapter.OnLevelPropertiesChanged(func() {
					biohazard, gravity, level := mode.levelAdapter.FloorEffect()
					item := noEffectItem

					if gravity {
						item = gravityEffectItem
					} else if biohazard {
						item = biohazardEffectItem
					}
					mode.floorEffectBox.SetSelectedItem(item)
					mode.floorEffectLevelSlider.SetValue(int64(level))
					mode.floorEffectLevelSlider.SetValueFormatter(item.formatter)
				})
			}
			{
				mode.animationGroupIndexLabel, mode.animationGroupIndexBox =
					realWorldBuilder.addComboProperty("Texture Animation Group", mode.onAnimationGroupIndexChanged)
				mode.animationGroupTimeLabel, mode.animationGroupTimeSlider =
					realWorldBuilder.addSliderProperty("Texture Animation Time", mode.onAnimationGroupTimeChanged)
				mode.animationGroupTimeSlider.SetRange(0, 1000)
				mode.animationGroupTimeSlider.SetValueFormatter(func(value int64) string {
					return fmt.Sprintf("%v msec", value)
				})
				mode.animationGroupFramesLabel, mode.animationGroupFramesSlider =
					realWorldBuilder.addSliderProperty("Texture Animation Frame Count", mode.onAnimationGroupFramesChanged)
				mode.animationGroupFramesSlider.SetRange(0, 10)
				mode.animationGroupTypeLabel, mode.animationGroupTypeBox =
					realWorldBuilder.addComboProperty("Texture Animation Type", mode.onAnimationGroupTypeChanged)
				items := []controls.ComboBoxItem{
					&enumItem{uint32(data.TextureAnimationForward), "Forward"},
					&enumItem{uint32(data.TextureAnimationForthAndBack), "Forth-And-Back"},
					&enumItem{uint32(data.TextureAnimationBackAndForth), "Back-And-Forth"}}
				mode.animationGroupTypeItems = make(map[int]controls.ComboBoxItem)
				for _, boxItem := range items {
					item := boxItem.(*enumItem)
					mode.animationGroupTypeItems[int(item.value)] = item
				}
				mode.animationGroupTypeBox.SetItems(items)

				mode.levelAdapter.OnLevelTextureAnimationsChanged(mode.onLevelTextureAnimationsChanged)
			}
		}
		mode.levelAdapter.OnLevelPropertiesChanged(func() {
			mode.realWorldProperties.SetVisible(!mode.levelAdapter.IsCyberspace())
		})
	}

	return mode
}

// SetActive implements the Mode interface.
func (mode *LevelControlMode) SetActive(active bool) {
	mode.area.SetVisible(active)
	mode.mapDisplay.SetVisible(active)
}

func (mode *LevelControlMode) levelTextures() []*graphics.BitmapTexture {
	ids := mode.context.ModelAdapter().ActiveLevel().LevelTextureIDs()
	textures := make([]*graphics.BitmapTexture, len(ids))
	store := mode.context.ForGraphics().WorldTextureStore(dataModel.TextureLarge)

	for index, id := range ids {
		textures[index] = store.Texture(graphics.TextureKeyFromInt(id))
	}

	return textures
}

func (mode *LevelControlMode) worldTextures() []*graphics.BitmapTexture {
	textureCount := mode.context.ModelAdapter().TextureAdapter().WorldTextureCount()
	textures := make([]*graphics.BitmapTexture, textureCount)
	store := mode.context.ForGraphics().WorldTextureStore(dataModel.TextureLarge)

	for index := 0; index < textureCount; index++ {
		textures[index] = store.Texture(graphics.TextureKeyFromInt(index))
	}

	return textures
}

func (mode *LevelControlMode) onHeightShiftChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.levelAdapter.RequestLevelPropertiesChange(func(properties *dataModel.LevelProperties) {
		properties.HeightShift = intAsPointer(int(item.value))
	})
}

func (mode *LevelControlMode) onSelectedLevelTextureChanged(index int) {
	ids := mode.context.ModelAdapter().ActiveLevel().LevelTextureIDs()

	if (index >= 0) && (index < len(ids)) {
		textureID := ids[index]
		mode.worldTexturesSelector.SetSelectedIndex(textureID)
		mode.worldTexturesIDSlider.SetValue(int64(textureID))
		mode.currentLevelTextureIndex = index
	} else {
		mode.worldTexturesSelector.SetSelectedIndex(-1)
		mode.worldTexturesIDSlider.SetValueUndefined()
		mode.currentLevelTextureIndex = -1
	}
}

func (mode *LevelControlMode) onSelectedWorldTextureChanged(index int) {
	mode.worldTexturesIDSlider.SetValue(int64(index))
	mode.setLevelTextureID(index)
}

func (mode *LevelControlMode) onSelectedWorldTextureIDChanged(newValue int64) {
	mode.worldTexturesSelector.SetSelectedIndex(int(newValue))
	mode.worldTexturesSelector.DisplaySelectedIndex()
	mode.setLevelTextureID(int(newValue))
}

func (mode *LevelControlMode) setLevelTextureID(id int) {
	levelAdapter := mode.context.ModelAdapter().ActiveLevel()
	ids := levelAdapter.LevelTextureIDs()

	if (mode.currentLevelTextureIndex >= 0) && (mode.currentLevelTextureIndex < len(ids)) {
		newIDs := make([]int, len(ids))
		copy(newIDs, ids)
		newIDs[mode.currentLevelTextureIndex] = id
		levelAdapter.RequestLevelTexturesChange(newIDs)
	}
}

func (mode *LevelControlMode) onLevelSurveillanceChanged() {
	surveillanceCount := mode.levelAdapter.ObjectSurveillanceCount()
	items := make([]controls.ComboBoxItem, surveillanceCount)
	var selectedItem controls.ComboBoxItem

	for index := 0; index < surveillanceCount; index++ {
		item := &enumItem{uint32(index), fmt.Sprintf("Object %v", index)}
		items[index] = item
		if index == mode.selectedSurveillanceIndex {
			selectedItem = item
		}
	}

	mode.surveillanceIndexBox.SetItems(items)
	mode.surveillanceIndexBox.SetSelectedItem(selectedItem)
	mode.onSurveillanceIndexChanged(selectedItem)
}

func (mode *LevelControlMode) onSurveillanceIndexChanged(boxItem controls.ComboBoxItem) {
	if boxItem != nil {
		item := boxItem.(*enumItem)
		mode.selectedSurveillanceIndex = int(item.value)
		sourceIndex, deathwatchIndex := mode.levelAdapter.ObjectSurveillanceInfo(mode.selectedSurveillanceIndex)
		mode.surveillanceSourceSlider.SetValue(int64(sourceIndex))
		mode.surveillanceDeathwatchSlider.SetValue(int64(deathwatchIndex))
	} else {
		mode.surveillanceSourceSlider.SetValueUndefined()
		mode.surveillanceDeathwatchSlider.SetValueUndefined()
	}
}

func (mode *LevelControlMode) onSurveillanceSourceChanged(newValue int64) {
	newIndex := int(newValue)
	mode.levelAdapter.RequestObjectSurveillance(mode.selectedSurveillanceIndex, &newIndex, nil)
}

func (mode *LevelControlMode) onSurveillanceDeathwatchChanged(newValue int64) {
	newIndex := int(newValue)
	mode.levelAdapter.RequestObjectSurveillance(mode.selectedSurveillanceIndex, nil, &newIndex)
}

func (mode *LevelControlMode) onLevelFloorPropertyBoxChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*levelPropertyItem)
	mode.levelAdapter.RequestLevelPropertiesChange(item.modifier)
	mode.floorEffectLevelSlider.SetValueFormatter(item.formatter)
}

func (mode *LevelControlMode) onLevelCeilingPropertyBoxChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*levelPropertyItem)
	mode.levelAdapter.RequestLevelPropertiesChange(item.modifier)
	mode.ceilingEffectLevelSlider.SetValueFormatter(item.formatter)
}

func (mode *LevelControlMode) onCeilingEffectLevelChanged(newValue int64) {
	mode.levelAdapter.RequestLevelPropertiesChange(func(properties *dataModel.LevelProperties) {
		properties.CeilingEffectLevel = intAsPointer(int(newValue))
	})
}

func (mode *LevelControlMode) onFloorEffectLevelChanged(newValue int64) {
	mode.levelAdapter.RequestLevelPropertiesChange(func(properties *dataModel.LevelProperties) {
		properties.FloorEffectLevel = intAsPointer(int(newValue))
	})
}

func (mode *LevelControlMode) onLevelTextureAnimationsChanged() {
	groupCount := mode.levelAdapter.TextureAnimationGroupCount()
	items := []controls.ComboBoxItem{}
	var selectedItem controls.ComboBoxItem

	for index := 1; index < groupCount; index++ {
		item := &enumItem{uint32(index), fmt.Sprintf("Group %v", index)}
		items = append(items, item)
		if index == mode.selectedAnimationGroupIndex {
			selectedItem = item
		}
	}
	mode.animationGroupIndexBox.SetItems(items)
	mode.animationGroupIndexBox.SetSelectedItem(selectedItem)
	mode.onAnimationGroupIndexChanged(selectedItem)
}

func (mode *LevelControlMode) onAnimationGroupIndexChanged(boxItem controls.ComboBoxItem) {
	if boxItem != nil {
		item := boxItem.(*enumItem)
		group := mode.levelAdapter.TextureAnimationGroup(int(item.value))
		mode.selectedAnimationGroupIndex = int(item.value)

		mode.animationGroupFramesSlider.SetValue(int64(group.FrameCount()))
		mode.animationGroupTimeSlider.SetValue(int64(group.FrameTime()))
		mode.animationGroupTypeBox.SetSelectedItem(mode.animationGroupTypeItems[group.LoopType()])
	} else {
		mode.animationGroupFramesSlider.SetValueUndefined()
		mode.animationGroupTimeSlider.SetValueUndefined()
		mode.animationGroupTypeBox.SetSelectedItem(nil)
	}
}

func (mode *LevelControlMode) requestAnimationGroupChange(modifier func(*dataModel.TextureAnimation)) {
	if mode.selectedAnimationGroupIndex > 0 {
		var properties dataModel.TextureAnimation

		modifier(&properties)
		mode.levelAdapter.RequestLevelTextureAnimationGroupChange(mode.selectedAnimationGroupIndex, properties)
	}
}

func (mode *LevelControlMode) onAnimationGroupTimeChanged(newValue int64) {
	mode.requestAnimationGroupChange(func(properties *dataModel.TextureAnimation) {
		properties.FrameTime = intAsPointer(int(newValue))
	})
}

func (mode *LevelControlMode) onAnimationGroupFramesChanged(newValue int64) {
	mode.requestAnimationGroupChange(func(properties *dataModel.TextureAnimation) {
		properties.FrameCount = intAsPointer(int(newValue))
	})
}

func (mode *LevelControlMode) onAnimationGroupTypeChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.requestAnimationGroupChange(func(properties *dataModel.TextureAnimation) {
		properties.LoopType = intAsPointer(int(item.value))
	})
}
