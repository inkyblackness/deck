package modes

import (
	"fmt"

	"github.com/inkyblackness/res/data"
	"github.com/inkyblackness/shocked-client/editor/cmd"
	"github.com/inkyblackness/shocked-client/editor/display"
	"github.com/inkyblackness/shocked-client/editor/model"
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/graphics/controls"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"

	dataModel "github.com/inkyblackness/shocked-model"
)

func lbpValueFormatter(value int64) string {
	return fmt.Sprintf("%v.%v LBP", value/2, (value%2)*5)
}

type ceilingEffect uint32

const (
	ceilingEffectNone      = 0
	ceilingEffectRadiation = 1
)

func (effect ceilingEffect) formatter() (f func(int64) string) {
	f = controls.DefaultSliderValueFormatter
	if effect == ceilingEffectRadiation {
		f = lbpValueFormatter
	}
	return
}

type floorEffect uint32

const (
	floorEffectNone      = 0
	floorEffectGravity   = 1
	floorEffectBiohazard = 2
)

func (effect floorEffect) formatter() (f func(int64) string) {
	f = controls.DefaultSliderValueFormatter
	if effect == floorEffectGravity {
		f = func(value int64) string { return fmt.Sprintf("%v%%", value*25) }
	} else if effect == floorEffectBiohazard {
		f = lbpValueFormatter
	}
	return
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
	heightShiftItems enumItems

	realWorldProperties *ui.Area

	levelGenericTexturesLabel    *controls.Label
	levelGenericTexturesSelector *controls.TextureSelector
	levelWallTexturesLabel       *controls.Label
	levelWallTexturesSelector    *controls.TextureSelector
	currentLevelTextureIndex     int
	worldTexturesLabel           *controls.Label
	worldTexturesSelector        *controls.TextureSelector
	worldTexturesIDLabel         *controls.Label
	worldTexturesIDSlider        *controls.Slider

	selectedSurveillanceIndex    int
	surveillanceIndexLabel       *controls.Label
	surveillanceIndexBox         *controls.ComboBox
	surveillanceIndexItems       enumItems
	surveillanceSourceLabel      *controls.Label
	surveillanceSourceSlider     *controls.Slider
	surveillanceDeathwatchLabel  *controls.Label
	surveillanceDeathwatchSlider *controls.Slider

	ceilingEffectLabel       *controls.Label
	ceilingEffectBox         *controls.ComboBox
	ceilingEffectItems       enumItems
	ceilingEffectLevelLabel  *controls.Label
	ceilingEffectLevelSlider *controls.Slider

	floorEffectLabel       *controls.Label
	floorEffectBox         *controls.ComboBox
	floorEffectItems       enumItems
	floorEffectLevelLabel  *controls.Label
	floorEffectLevelSlider *controls.Slider

	selectedAnimationGroupIndex int
	animationGroupIndexLabel    *controls.Label
	animationGroupIndexBox      *controls.ComboBox
	animationGroupItems         enumItems
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
		selectedAnimationGroupIndex: 1}

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
				selectedLevelID := item.(int)
				context.Perform(&cmd.SetActiveLevelCommand{
					Setter: func(levelID int) error {
						context.ModelAdapter().RequestActiveLevel(levelID)
						return nil
					},
					OldValue: context.ModelAdapter().ActiveLevel().ID(),
					NewValue: selectedLevelID})
			})

			adapter := context.ModelAdapter()
			activeLevelAdapter := adapter.ActiveLevel()
			activeLevelAdapter.OnIDChanged(func() {
				levelID := mode.levelAdapter.ID()
				isProperLevel := levelID >= 0

				if isProperLevel {
					mode.activeLevelBox.SetSelectedItem(levelID)
				} else {
					mode.activeLevelBox.SetSelectedItem(nil)
				}
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
			mode.heightShiftItems = []*enumItem{
				{0, "32 Tiles"},
				{1, "16 Tiles"},
				{2, "8 Tiles"},
				{3, "4 Tiles"},
				{4, "2 Tiles"},
				{5, "1 Tile"},
				{6, "1/2 Tile"},
				{7, "1/4 Tile"}}
			mode.heightShiftBox.SetItems(mode.heightShiftItems.forComboBox())
			mode.levelAdapter.OnLevelPropertiesChanged(func() {
				heightShift := mode.levelAdapter.HeightShift()
				if (heightShift >= 0) && (heightShift < len(mode.heightShiftItems)) {
					mode.heightShiftBox.SetSelectedItem(mode.heightShiftItems[heightShift])
				} else {
					mode.heightShiftBox.SetSelectedItem(nil)
				}
			})
		}

		{
			var realWorldBuilder *controlPanelBuilder
			mode.realWorldProperties, realWorldBuilder = panelBuilder.addSection(false)

			{
				mode.levelGenericTexturesLabel, mode.levelGenericTexturesSelector = realWorldBuilder.addTextureProperty("Level Textures (floors, ceilings, walls)",
					mode.genericLevelTextures, mode.onSelectedGenericLevelTextureChanged)
				mode.levelWallTexturesLabel, mode.levelWallTexturesSelector = realWorldBuilder.addTextureProperty("Level Textures (walls only)",
					mode.wallLevelTextures, mode.onSelectedWallLevelTextureChanged)
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

				mode.ceilingEffectItems = []*enumItem{
					{ceilingEffectNone, "None"},
					{ceilingEffectRadiation, "Radiation"}}
				mode.ceilingEffectBox.SetItems(mode.ceilingEffectItems.forComboBox())
				mode.ceilingEffectLevelSlider.SetRange(0, 255)

				mode.levelAdapter.OnLevelPropertiesChanged(func() {
					radiation, level := mode.levelAdapter.CeilingEffect()
					effect := ceilingEffect(ceilingEffectNone)

					if radiation {
						effect = ceilingEffectRadiation
					}
					mode.ceilingEffectBox.SetSelectedItem(mode.ceilingEffectItems[effect])
					mode.ceilingEffectLevelSlider.SetValue(int64(level))
					mode.ceilingEffectLevelSlider.SetValueFormatter(effect.formatter())
				})
			}
			{
				mode.floorEffectLabel, mode.floorEffectBox =
					realWorldBuilder.addComboProperty("Floor Effect", mode.onLevelFloorPropertyBoxChanged)
				mode.floorEffectLevelLabel, mode.floorEffectLevelSlider =
					realWorldBuilder.addSliderProperty("Floor Effect Level", mode.onFloorEffectLevelChanged)

				mode.floorEffectItems = []*enumItem{
					{floorEffectNone, "None"},
					{floorEffectGravity, "Gravity"},
					{floorEffectBiohazard, "Biohazard"}}
				mode.floorEffectBox.SetItems(mode.floorEffectItems.forComboBox())
				mode.floorEffectLevelSlider.SetRange(0, 255)

				mode.levelAdapter.OnLevelPropertiesChanged(func() {
					effect, level := mode.currentFloorEffect()
					mode.floorEffectBox.SetSelectedItem(mode.floorEffectItems[effect])
					mode.floorEffectLevelSlider.SetValue(int64(level))
					mode.floorEffectLevelSlider.SetValueFormatter(effect.formatter())
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
			isProperLevel := mode.levelAdapter.ID() >= 0
			isRealWorld := !mode.levelAdapter.IsCyberspace()
			mode.realWorldProperties.SetVisible(isProperLevel && isRealWorld)
		})
	}

	return mode
}

// SetActive implements the Mode interface.
func (mode *LevelControlMode) SetActive(active bool) {
	mode.area.SetVisible(active)
	mode.mapDisplay.SetVisible(active)
}

func (mode *LevelControlMode) currentFloorEffect() (floorEffect, int) {
	biohazard, gravity, level := mode.levelAdapter.FloorEffect()
	effect := floorEffect(floorEffectNone)
	if gravity {
		effect = floorEffectGravity
	} else if biohazard {
		effect = floorEffectBiohazard
	}
	return effect, level
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

func (mode *LevelControlMode) genericLevelTextures() []*graphics.BitmapTexture {
	textures := mode.levelTextures()
	if len(textures) > 32 {
		return textures[:32]
	}
	return textures
}

func (mode *LevelControlMode) wallLevelTextures() []*graphics.BitmapTexture {
	textures := mode.levelTextures()
	if len(textures) > 32 {
		return textures[32:]
	}
	return nil
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
	newValue := int(item.value)

	mode.context.Perform(&cmd.SetIntPropertyCommand{
		Setter: func(value int) error {
			mode.levelAdapter.RequestLevelPropertiesChange(func(properties *dataModel.LevelProperties) {
				properties.HeightShift = intAsPointer(value)
			})
			return nil
		},
		NewValue: newValue,
		OldValue: mode.levelAdapter.HeightShift()})
}

func (mode *LevelControlMode) onSelectedGenericLevelTextureChanged(index int) {
	mode.levelWallTexturesSelector.SetSelectedIndex(-1)
	mode.onSelectedLevelTextureChanged(index)
}

func (mode *LevelControlMode) onSelectedWallLevelTextureChanged(index int) {
	mode.levelGenericTexturesSelector.SetSelectedIndex(-1)
	mode.onSelectedLevelTextureChanged(32 + index)
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
	oldIDs := levelAdapter.LevelTextureIDs()

	if (mode.currentLevelTextureIndex >= 0) && (mode.currentLevelTextureIndex < len(oldIDs)) {
		newIDs := make([]int, len(oldIDs))
		copy(newIDs, oldIDs)
		newIDs[mode.currentLevelTextureIndex] = id

		mode.context.Perform(&cmd.SetLevelTexturesCommand{
			Setter: func(textureIDs []int) error {
				mode.worldTexturesSelector.SetSelectedIndex(textureIDs[mode.currentLevelTextureIndex])
				mode.worldTexturesIDSlider.SetValue(int64(textureIDs[mode.currentLevelTextureIndex]))
				if mode.currentLevelTextureIndex < 32 {
					mode.levelGenericTexturesSelector.SetSelectedIndex(mode.currentLevelTextureIndex)
				} else {
					mode.levelGenericTexturesSelector.SetSelectedIndex(-1)
				}
				if mode.currentLevelTextureIndex >= 32 {
					mode.levelWallTexturesSelector.SetSelectedIndex(mode.currentLevelTextureIndex - 32)
				} else {
					mode.levelWallTexturesSelector.SetSelectedIndex(-1)
				}

				levelAdapter.RequestLevelTexturesChange(textureIDs)
				return nil
			},
			OldTextureIDs: oldIDs,
			NewTextureIDs: newIDs})
	}
}

func (mode *LevelControlMode) onLevelSurveillanceChanged() {
	surveillanceCount := mode.levelAdapter.ObjectSurveillanceCount()

	mode.surveillanceIndexItems = make([]*enumItem, surveillanceCount)
	for index := 0; index < surveillanceCount; index++ {
		item := &enumItem{uint32(index), fmt.Sprintf("Object %v", index)}
		mode.surveillanceIndexItems[index] = item
	}

	mode.surveillanceIndexBox.SetItems(mode.surveillanceIndexItems.forComboBox())
	mode.setSurveillanceState(mode.selectedSurveillanceIndex)
}

func (mode *LevelControlMode) onSurveillanceIndexChanged(boxItem controls.ComboBoxItem) {
	if boxItem != nil {
		item := boxItem.(*enumItem)
		mode.setSurveillanceState(int(item.value))
	} else {
		mode.setSurveillanceState(-1)
	}
}

func (mode *LevelControlMode) onSurveillanceSourceChanged(newValue int64) {
	oldValue, _ := mode.levelAdapter.ObjectSurveillanceInfo(mode.selectedSurveillanceIndex)
	mode.requestSurveillanceChange(func(value int) error {
		mode.levelAdapter.RequestObjectSurveillance(mode.selectedSurveillanceIndex, &value, nil)
		return nil
	}, int(newValue), oldValue)
}

func (mode *LevelControlMode) onSurveillanceDeathwatchChanged(newValue int64) {
	_, oldValue := mode.levelAdapter.ObjectSurveillanceInfo(mode.selectedSurveillanceIndex)
	mode.requestSurveillanceChange(func(value int) error {
		mode.levelAdapter.RequestObjectSurveillance(mode.selectedSurveillanceIndex, nil, &value)
		return nil
	}, int(newValue), oldValue)
}

func (mode *LevelControlMode) requestSurveillanceChange(executor func(value int) error, newValue, oldValue int) {
	currentSurveillanceIndex := mode.selectedSurveillanceIndex

	if currentSurveillanceIndex >= 0 {
		mode.context.Perform(&cmd.SetIntPropertyCommand{
			Setter: func(value int) error {
				mode.setSurveillanceState(currentSurveillanceIndex)
				return executor(value)
			},
			NewValue: newValue,
			OldValue: oldValue})
	}
}

func (mode *LevelControlMode) setSurveillanceState(objectIndex int) {
	mode.selectedSurveillanceIndex = objectIndex
	if (mode.selectedSurveillanceIndex >= 0) && (mode.selectedSurveillanceIndex < len(mode.surveillanceIndexItems)) {
		mode.surveillanceIndexBox.SetSelectedItem(mode.surveillanceIndexItems[mode.selectedSurveillanceIndex])
		sourceIndex, deathwatchIndex := mode.levelAdapter.ObjectSurveillanceInfo(mode.selectedSurveillanceIndex)
		mode.surveillanceSourceSlider.SetValue(int64(sourceIndex))
		mode.surveillanceDeathwatchSlider.SetValue(int64(deathwatchIndex))
	} else {
		mode.surveillanceIndexBox.SetSelectedItem(nil)
		mode.surveillanceSourceSlider.SetValueUndefined()
		mode.surveillanceDeathwatchSlider.SetValueUndefined()
	}
}

func (mode *LevelControlMode) onLevelFloorPropertyBoxChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	oldValue, _ := mode.currentFloorEffect()

	mode.context.Perform(&cmd.SetIntPropertyCommand{
		Setter: func(value int) error {
			mode.levelAdapter.RequestLevelPropertiesChange(func(properties *dataModel.LevelProperties) {
				biohazard := value == floorEffectBiohazard
				gravity := value == floorEffectGravity

				properties.FloorHasBiohazard = &biohazard
				properties.FloorHasGravity = &gravity
			})
			return nil
		},
		NewValue: int(item.value),
		OldValue: int(oldValue)})
}

func (mode *LevelControlMode) onLevelCeilingPropertyBoxChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	oldValue, _ := mode.levelAdapter.CeilingEffect()

	mode.context.Perform(&cmd.SetBooleanPropertyCommand{
		Setter: func(value bool) error {
			mode.levelAdapter.RequestLevelPropertiesChange(func(properties *dataModel.LevelProperties) {
				properties.CeilingHasRadiation = &value
			})
			return nil
		},
		NewValue: item.value == ceilingEffectRadiation,
		OldValue: oldValue})
}

func (mode *LevelControlMode) onCeilingEffectLevelChanged(newValue int64) {
	_, oldValue := mode.levelAdapter.CeilingEffect()
	mode.context.Perform(&cmd.SetIntPropertyCommand{
		Setter: func(value int) error {
			mode.levelAdapter.RequestLevelPropertiesChange(func(properties *dataModel.LevelProperties) {
				properties.CeilingEffectLevel = &value
			})
			return nil
		},
		NewValue: int(newValue),
		OldValue: oldValue})
}

func (mode *LevelControlMode) onFloorEffectLevelChanged(newValue int64) {
	_, _, oldValue := mode.levelAdapter.FloorEffect()
	mode.context.Perform(&cmd.SetIntPropertyCommand{
		Setter: func(value int) error {
			mode.levelAdapter.RequestLevelPropertiesChange(func(properties *dataModel.LevelProperties) {
				properties.FloorEffectLevel = &value
			})
			return nil
		},
		NewValue: int(newValue),
		OldValue: oldValue})
}

func (mode *LevelControlMode) onLevelTextureAnimationsChanged() {
	groupCount := mode.levelAdapter.TextureAnimationGroupCount()
	mode.animationGroupItems = nil
	if groupCount > 0 {
		mode.animationGroupItems = make([]*enumItem, 0, groupCount-1)
		for index := 1; index < groupCount; index++ {
			item := &enumItem{uint32(index), fmt.Sprintf("Group %v", index)}
			mode.animationGroupItems = append(mode.animationGroupItems, item)
		}
	}
	mode.animationGroupIndexBox.SetItems(mode.animationGroupItems.forComboBox())
	mode.setAnimationGroupState(mode.selectedAnimationGroupIndex)
}

func (mode *LevelControlMode) onAnimationGroupIndexChanged(boxItem controls.ComboBoxItem) {
	if boxItem != nil {
		item := boxItem.(*enumItem)
		mode.setAnimationGroupState(int(item.value))
	} else {
		mode.setAnimationGroupState(-1)
	}
}

func (mode *LevelControlMode) onAnimationGroupTimeChanged(newValue int64) {
	mode.requestAnimationGroupChange(func(properties *dataModel.TextureAnimation, value *int) {
		properties.FrameTime = value
	}, int(newValue), mode.levelAdapter.TextureAnimationGroup(mode.selectedAnimationGroupIndex).FrameTime())
}

func (mode *LevelControlMode) onAnimationGroupFramesChanged(newValue int64) {
	mode.requestAnimationGroupChange(func(properties *dataModel.TextureAnimation, value *int) {
		properties.FrameCount = value
	}, int(newValue), mode.levelAdapter.TextureAnimationGroup(mode.selectedAnimationGroupIndex).FrameCount())
}

func (mode *LevelControlMode) onAnimationGroupTypeChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.requestAnimationGroupChange(func(properties *dataModel.TextureAnimation, value *int) {
		properties.LoopType = value
	}, int(item.value), mode.levelAdapter.TextureAnimationGroup(mode.selectedAnimationGroupIndex).LoopType())
}

func (mode *LevelControlMode) requestAnimationGroupChange(modifier func(*dataModel.TextureAnimation, *int), newValue, oldValue int) {
	currentAnimationGroupIndex := mode.selectedAnimationGroupIndex

	if currentAnimationGroupIndex >= 0 {
		mode.context.Perform(&cmd.SetIntPropertyCommand{
			Setter: func(value int) error {
				var properties dataModel.TextureAnimation
				mode.setAnimationGroupState(currentAnimationGroupIndex)
				modifier(&properties, &value)
				mode.levelAdapter.RequestLevelTextureAnimationGroupChange(mode.selectedAnimationGroupIndex, properties)
				return nil
			},
			NewValue: newValue,
			OldValue: oldValue})
	}
}

func (mode *LevelControlMode) setAnimationGroupState(groupIndex int) {
	mode.selectedAnimationGroupIndex = groupIndex
	if (mode.selectedAnimationGroupIndex >= 1) && (mode.selectedAnimationGroupIndex < mode.levelAdapter.TextureAnimationGroupCount()) {
		group := mode.levelAdapter.TextureAnimationGroup(mode.selectedAnimationGroupIndex)
		mode.animationGroupIndexBox.SetSelectedItem(mode.animationGroupItems[mode.selectedAnimationGroupIndex-1])
		mode.animationGroupFramesSlider.SetValue(int64(group.FrameCount()))
		mode.animationGroupTimeSlider.SetValue(int64(group.FrameTime()))
		mode.animationGroupTypeBox.SetSelectedItem(mode.animationGroupTypeItems[group.LoopType()])
	} else {
		mode.animationGroupIndexBox.SetSelectedItem(nil)
		mode.animationGroupFramesSlider.SetValueUndefined()
		mode.animationGroupTimeSlider.SetValueUndefined()
		mode.animationGroupTypeBox.SetSelectedItem(nil)
	}
}
