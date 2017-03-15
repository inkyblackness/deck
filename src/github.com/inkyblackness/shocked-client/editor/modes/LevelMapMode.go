package modes

import (
	"fmt"

	"github.com/inkyblackness/shocked-client/editor/display"
	"github.com/inkyblackness/shocked-client/editor/model"
	"github.com/inkyblackness/shocked-client/env"
	"github.com/inkyblackness/shocked-client/env/keys"
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/graphics/controls"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"
	"github.com/inkyblackness/shocked-client/util"
	dataModel "github.com/inkyblackness/shocked-model"
)

type tilePropertySetter func(properties *dataModel.TileProperties, value interface{})
type tilePropertyItem struct {
	value  interface{}
	setter tilePropertySetter
}

func (item *tilePropertyItem) String() string {
	return fmt.Sprintf("%v", item.value)
}

// LevelMapMode is a mode for level maps.
type LevelMapMode struct {
	context      Context
	levelAdapter *model.LevelAdapter

	mapDisplay *display.MapDisplay

	area       *ui.Area
	panel      *ui.Area
	panelRight ui.Anchor

	selectedTiles []model.TileCoordinate

	tileTypeLabel      *controls.Label
	tileTypeBox        *controls.ComboBox
	tileTypeItems      map[dataModel.TileType]*tilePropertyItem
	floorHeightLabel   *controls.Label
	floorHeightBox     *controls.ComboBox
	floorHeightItems   map[dataModel.HeightUnit]*tilePropertyItem
	ceilingHeightLabel *controls.Label
	ceilingHeightBox   *controls.ComboBox
	ceilingHeightItems map[dataModel.HeightUnit]*tilePropertyItem
	slopeHeightLabel   *controls.Label
	slopeHeightBox     *controls.ComboBox
	slopeHeightItems   map[dataModel.HeightUnit]*tilePropertyItem
	slopeControlLabel  *controls.Label
	slopeControlBox    *controls.ComboBox
	slopeControlItems  map[dataModel.SlopeControl]*tilePropertyItem

	realWorldArea                *ui.Area
	floorTextureLabel            *controls.Label
	floorTextureSelector         *controls.TextureSelector
	floorTextureRotationsLabel   *controls.Label
	floorTextureRotationsBox     *controls.ComboBox
	floorTextureRotationsItems   map[int]*tilePropertyItem
	ceilingTextureLabel          *controls.Label
	ceilingTextureSelector       *controls.TextureSelector
	ceilingTextureRotationsLabel *controls.Label
	ceilingTextureRotationsBox   *controls.ComboBox
	ceilingTextureRotationsItems map[int]*tilePropertyItem
	wallTextureLabel             *controls.Label
	wallTextureSelector          *controls.TextureSelector
	wallTextureOffsetLabel       *controls.Label
	wallTextureOffsetBox         *controls.ComboBox
	wallTextureOffsetItems       map[dataModel.HeightUnit]*tilePropertyItem
	useAdjacentWallTextureLabel  *controls.Label
	useAdjacentWallTextureBox    *controls.ComboBox
	useAdjacentWallTextureItems  map[string]*tilePropertyItem

	cyberspaceArea *ui.Area
}

// NewLevelMapMode returns a new instance.
func NewLevelMapMode(context Context, parent *ui.Area, mapDisplay *display.MapDisplay) *LevelMapMode {
	mode := &LevelMapMode{
		context:      context,
		levelAdapter: context.ModelAdapter().ActiveLevel(),
		mapDisplay:   mapDisplay}

	{
		builder := ui.NewAreaBuilder()
		builder.SetParent(parent)
		builder.SetLeft(ui.NewOffsetAnchor(parent.Left(), 0))
		builder.SetTop(ui.NewOffsetAnchor(parent.Top(), 0))
		builder.SetRight(ui.NewOffsetAnchor(parent.Right(), 0))
		builder.SetBottom(ui.NewOffsetAnchor(parent.Bottom(), 0))
		builder.SetVisible(false)
		builder.OnEvent(events.MouseMoveEventType, mode.onMouseMoved)
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

		mode.tileTypeLabel, mode.tileTypeBox = panelBuilder.addComboProperty("Tile Type", mode.onTilePropertyChangeRequested)
		{
			setter := func(properties *dataModel.TileProperties, value interface{}) {
				tileType := value.(dataModel.TileType)
				properties.Type = &tileType
			}
			tileTypes := dataModel.TileTypes()
			tileTypeItems := make([]controls.ComboBoxItem, len(tileTypes))
			mode.tileTypeItems = make(map[dataModel.TileType]*tilePropertyItem)
			for index, tileType := range tileTypes {
				item := &tilePropertyItem{tileType, setter}
				tileTypeItems[index] = item
				mode.tileTypeItems[tileType] = item
			}
			mode.tileTypeItems[dataModel.TileType("")] = &tilePropertyItem{"", nil}
			mode.tileTypeBox.SetItems(tileTypeItems)
		}

		mode.floorHeightLabel, mode.floorHeightBox = panelBuilder.addComboProperty("Floor Height", mode.onTilePropertyChangeRequested)
		mode.ceilingHeightLabel, mode.ceilingHeightBox = panelBuilder.addComboProperty("Ceiling Height", mode.onTilePropertyChangeRequested)
		mode.slopeHeightLabel, mode.slopeHeightBox = panelBuilder.addComboProperty("Slope Height", mode.onTilePropertyChangeRequested)

		positiveHeight := func(value int) dataModel.HeightUnit { return dataModel.HeightUnit(value) }
		negativeHeight := func(value int) dataModel.HeightUnit { return dataModel.HeightUnit(31 - value) }
		setupHeightCollections := func(mapper func(int) dataModel.HeightUnit,
			setter func(*dataModel.TileProperties, dataModel.HeightUnit)) ([]controls.ComboBoxItem, map[dataModel.HeightUnit]*tilePropertyItem) {
			mappingSetter := func(properties *dataModel.TileProperties, value interface{}) {
				height := mapper(int(value.(dataModel.HeightUnit)))
				setter(properties, height)
			}
			itemsSlice := make([]controls.ComboBoxItem, 32)
			heightItems := make(map[dataModel.HeightUnit]*tilePropertyItem)
			for height := 0; height < 32; height++ {
				heightUnit := mapper(height)
				item := &tilePropertyItem{heightUnit, mappingSetter}
				itemsSlice[height] = item
				heightItems[heightUnit] = item
			}
			heightItems[dataModel.HeightUnit(-1)] = &tilePropertyItem{"", nil}
			return itemsSlice, heightItems
		}

		setupBooleanCollections := func(
			setter func(*dataModel.TileProperties, bool)) ([]controls.ComboBoxItem, map[string]*tilePropertyItem) {
			mappingSetter := func(properties *dataModel.TileProperties, value interface{}) {
				mappedValue := value.(bool)
				setter(properties, mappedValue)
			}
			itemsSlice := make([]controls.ComboBoxItem, 2)
			keyedItems := make(map[string]*tilePropertyItem)
			falseItem := &tilePropertyItem{false, mappingSetter}
			itemsSlice[0] = falseItem
			keyedItems["false"] = falseItem
			trueItem := &tilePropertyItem{true, mappingSetter}
			itemsSlice[1] = trueItem
			keyedItems["true"] = trueItem
			keyedItems[""] = &tilePropertyItem{"", nil}
			return itemsSlice, keyedItems
		}

		{
			var floorHeightItemsSlice []controls.ComboBoxItem
			var ceilingHeightItemsSlice []controls.ComboBoxItem
			var slopeHeightItemsSlice []controls.ComboBoxItem
			floorHeightItemsSlice, mode.floorHeightItems = setupHeightCollections(positiveHeight,
				func(properties *dataModel.TileProperties, height dataModel.HeightUnit) {
					properties.FloorHeight = &height
				})
			ceilingHeightItemsSlice, mode.ceilingHeightItems = setupHeightCollections(negativeHeight,
				func(properties *dataModel.TileProperties, height dataModel.HeightUnit) {
					properties.CeilingHeight = &height
				})
			slopeHeightItemsSlice, mode.slopeHeightItems = setupHeightCollections(positiveHeight,
				func(properties *dataModel.TileProperties, height dataModel.HeightUnit) {
					properties.SlopeHeight = &height
				})

			mode.floorHeightBox.SetItems(floorHeightItemsSlice)
			mode.ceilingHeightBox.SetItems(ceilingHeightItemsSlice)
			mode.slopeHeightBox.SetItems(slopeHeightItemsSlice)
		}

		mode.slopeControlLabel, mode.slopeControlBox = panelBuilder.addComboProperty("Slope Control", mode.onTilePropertyChangeRequested)
		{
			setter := func(properties *dataModel.TileProperties, value interface{}) {
				slopeControl := value.(dataModel.SlopeControl)
				properties.SlopeControl = &slopeControl
			}
			slopeControls := dataModel.SlopeControls()
			slopeControlItems := make([]controls.ComboBoxItem, len(slopeControls))
			mode.slopeControlItems = make(map[dataModel.SlopeControl]*tilePropertyItem)
			for index, slopeControl := range slopeControls {
				item := &tilePropertyItem{slopeControl, setter}
				slopeControlItems[index] = item
				mode.slopeControlItems[slopeControl] = item
			}
			mode.slopeControlItems[dataModel.SlopeControl("")] = &tilePropertyItem{"", nil}
			mode.slopeControlBox.SetItems(slopeControlItems)
		}
		{
			var realWorldPanelBuilder *controlPanelBuilder
			mode.realWorldArea, realWorldPanelBuilder = panelBuilder.addSection(false)

			{
				setupRotations := func(setter func(*dataModel.TileProperties, int)) ([]controls.ComboBoxItem, map[int]*tilePropertyItem) {
					mappingSetter := func(properties *dataModel.TileProperties, value interface{}) {
						setter(properties, value.(int))
					}
					itemsSlice := make([]controls.ComboBoxItem, 4)
					propertyItems := make(map[int]*tilePropertyItem)
					for rotations := 0; rotations < 4; rotations++ {
						item := &tilePropertyItem{rotations, mappingSetter}
						itemsSlice[rotations] = item
						propertyItems[rotations] = item
					}
					propertyItems[-1] = &tilePropertyItem{"", nil}
					return itemsSlice, propertyItems
				}

				mode.floorTextureLabel, mode.floorTextureSelector = realWorldPanelBuilder.addTextureProperty("Floor Texture",
					mode.levelTextures, mode.onFloorTextureChanged)
				mode.floorTextureRotationsLabel, mode.floorTextureRotationsBox = realWorldPanelBuilder.addComboProperty("Floor Texture Rotations", mode.onTilePropertyChangeRequested)
				mode.ceilingTextureLabel, mode.ceilingTextureSelector = realWorldPanelBuilder.addTextureProperty("Ceiling Texture",
					mode.levelTextures, mode.onCeilingTextureChanged)
				mode.ceilingTextureRotationsLabel, mode.ceilingTextureRotationsBox = realWorldPanelBuilder.addComboProperty("Ceiling Texture Rotations", mode.onTilePropertyChangeRequested)

				var floorTextureRotationsItemsSlice []controls.ComboBoxItem
				var ceilingTextureRotationsItemsSlice []controls.ComboBoxItem
				floorTextureRotationsItemsSlice, mode.floorTextureRotationsItems = setupRotations(
					func(properties *dataModel.TileProperties, value int) {
						properties.RealWorld.FloorTextureRotations = &value
					})
				ceilingTextureRotationsItemsSlice, mode.ceilingTextureRotationsItems = setupRotations(
					func(properties *dataModel.TileProperties, value int) {
						properties.RealWorld.CeilingTextureRotations = &value
					})
				mode.floorTextureRotationsBox.SetItems(floorTextureRotationsItemsSlice)
				mode.ceilingTextureRotationsBox.SetItems(ceilingTextureRotationsItemsSlice)
			}
			{
				mode.wallTextureLabel, mode.wallTextureSelector = realWorldPanelBuilder.addTextureProperty("Wall Texture",
					mode.levelTextures, mode.onWallTextureChanged)

				var wallTextureOffsetItemsSlice []controls.ComboBoxItem
				wallTextureOffsetItemsSlice, mode.wallTextureOffsetItems = setupHeightCollections(positiveHeight,
					func(properties *dataModel.TileProperties, height dataModel.HeightUnit) {
						properties.RealWorld.WallTextureOffset = &height
					})

				mode.wallTextureOffsetLabel, mode.wallTextureOffsetBox = realWorldPanelBuilder.addComboProperty("Wall Texture Offset", mode.onTilePropertyChangeRequested)
				mode.wallTextureOffsetBox.SetItems(wallTextureOffsetItemsSlice)
			}
			{
				var useAdjacentWallTextureSlice []controls.ComboBoxItem
				useAdjacentWallTextureSlice, mode.useAdjacentWallTextureItems = setupBooleanCollections(
					func(properties *dataModel.TileProperties, value bool) {
						properties.RealWorld.UseAdjacentWallTexture = &value
					})

				mode.useAdjacentWallTextureLabel, mode.useAdjacentWallTextureBox = realWorldPanelBuilder.addComboProperty("Use Adjacent Wall Texture", mode.onTilePropertyChangeRequested)
				mode.useAdjacentWallTextureBox.SetItems(useAdjacentWallTextureSlice)
			}
		}
		{
			mode.cyberspaceArea, _ = panelBuilder.addSection(false)

		}
		mode.levelAdapter.OnIDChanged(func() {
			mode.realWorldArea.SetVisible(!mode.levelAdapter.IsCyberspace())
			mode.cyberspaceArea.SetVisible(mode.levelAdapter.IsCyberspace())
		})
	}

	return mode
}

// SetActive implements the Mode interface.
func (mode *LevelMapMode) SetActive(active bool) {
	if active {
		mode.mapDisplay.SetSelectedTiles(mode.selectedTiles)
	} else {
		mode.mapDisplay.ClearHighlightedTile()
		mode.mapDisplay.SetSelectedTiles(nil)
	}
	mode.area.SetVisible(active)
	mode.mapDisplay.SetVisible(active)
}

func (mode *LevelMapMode) levelTextures() []*graphics.BitmapTexture {
	ids := mode.levelAdapter.LevelTextureIDs()
	textures := make([]*graphics.BitmapTexture, len(ids))
	store := mode.context.ForGraphics().WorldTextureStore(dataModel.TextureLarge)

	for index, id := range ids {
		textures[index] = store.Texture(graphics.TextureKeyFromInt(id))
	}

	return textures
}

func (mode *LevelMapMode) onMouseMoved(area *ui.Area, event events.Event) (consumed bool) {
	mouseEvent := event.(*events.MouseMoveEvent)

	if mouseEvent.Buttons() == 0 {
		worldX, worldY := mode.mapDisplay.WorldCoordinatesForPixel(mouseEvent.Position())
		coord := model.TileCoordinateOf(int(worldX)>>8, int(worldY)>>8)
		tileX, tileY := coord.XY()

		if tileX >= 0 && tileX < 64 && tileY >= 0 && tileY < 64 {
			mode.mapDisplay.SetHighlightedTile(coord)
		} else {
			mode.mapDisplay.ClearHighlightedTile()
		}
		consumed = true
	}

	return
}

func (mode *LevelMapMode) onMouseButtonClicked(area *ui.Area, event events.Event) (consumed bool) {
	mouseEvent := event.(*events.MouseButtonEvent)

	if mouseEvent.AffectedButtons() == env.MousePrimary {
		worldX, worldY := mode.mapDisplay.WorldCoordinatesForPixel(mouseEvent.Position())
		coord := model.TileCoordinateOf(int(worldX)>>8, int(worldY)>>8)
		tileX, tileY := coord.XY()

		if tileX >= 0 && tileX < 64 && tileY >= 0 && tileY < 64 {
			if keys.Modifier(mouseEvent.Modifier()) == keys.ModControl {
				mode.toggleSelectedTile(coord)
			} else {
				mode.setSelectedTiles([]model.TileCoordinate{coord})
			}
			consumed = true
		}
	}

	return
}

func (mode *LevelMapMode) setSelectedTiles(tiles []model.TileCoordinate) {
	mode.selectedTiles = tiles
	mode.onSelectedTilesChanged()
}

func (mode *LevelMapMode) toggleSelectedTile(coord model.TileCoordinate) {
	newList := []model.TileCoordinate{}
	wasSelected := false

	for _, other := range mode.selectedTiles {
		if other != coord {
			newList = append(newList, other)
		} else {
			wasSelected = true
		}
	}
	if !wasSelected {
		newList = append(newList, coord)
	}

	mode.selectedTiles = newList
	mode.onSelectedTilesChanged()
}

func (mode *LevelMapMode) onSelectedTilesChanged() {
	mode.mapDisplay.SetSelectedTiles(mode.selectedTiles)
	tileMap := mode.levelAdapter.TileMap()
	typeUnifier := util.NewValueUnifier(dataModel.TileType(""))
	floorHeightUnifier := util.NewValueUnifier(dataModel.HeightUnit(-1))
	ceilingHeightUnifier := util.NewValueUnifier(dataModel.HeightUnit(-1))
	slopeHeightUnifier := util.NewValueUnifier(dataModel.HeightUnit(-1))
	slopeControlUnifier := util.NewValueUnifier(dataModel.SlopeControl(""))
	floorTextureUnifier := util.NewValueUnifier(-1)
	floorTextureRotationsUnifier := util.NewValueUnifier(-1)
	ceilingTextureUnifier := util.NewValueUnifier(-1)
	ceilingTextureRotationsUnifier := util.NewValueUnifier(-1)
	wallTextureUnifier := util.NewValueUnifier(-1)
	wallTextureOffsetUnifier := util.NewValueUnifier(dataModel.HeightUnit(-1))
	useAdjacentWallTextureUnifier := util.NewValueUnifier("")

	for _, coord := range mode.selectedTiles {
		tile := tileMap.Tile(coord)
		properties := tile.Properties()
		if properties != nil {
			typeUnifier.Add(*properties.Type)
			floorHeightUnifier.Add(*properties.FloorHeight)
			ceilingHeightUnifier.Add(*properties.CeilingHeight)
			slopeHeightUnifier.Add(*properties.SlopeHeight)
			slopeControlUnifier.Add(*properties.SlopeControl)
			if properties.RealWorld != nil {
				floorTextureUnifier.Add(*properties.RealWorld.FloorTexture)
				floorTextureRotationsUnifier.Add(*properties.RealWorld.FloorTextureRotations)
				ceilingTextureUnifier.Add(*properties.RealWorld.CeilingTexture)
				ceilingTextureRotationsUnifier.Add(*properties.RealWorld.CeilingTextureRotations)
				wallTextureUnifier.Add(*properties.RealWorld.WallTexture)
				wallTextureOffsetUnifier.Add(*properties.RealWorld.WallTextureOffset)
				useAdjacentWallTextureUnifier.Add(fmt.Sprintf("%v", *properties.RealWorld.UseAdjacentWallTexture))
			}
		}
	}
	mode.tileTypeBox.SetSelectedItem(mode.tileTypeItems[typeUnifier.Value().(dataModel.TileType)])
	mode.floorHeightBox.SetSelectedItem(mode.floorHeightItems[floorHeightUnifier.Value().(dataModel.HeightUnit)])
	mode.ceilingHeightBox.SetSelectedItem(mode.ceilingHeightItems[ceilingHeightUnifier.Value().(dataModel.HeightUnit)])
	mode.slopeHeightBox.SetSelectedItem(mode.slopeHeightItems[slopeHeightUnifier.Value().(dataModel.HeightUnit)])
	mode.slopeControlBox.SetSelectedItem(mode.slopeControlItems[slopeControlUnifier.Value().(dataModel.SlopeControl)])
	mode.floorTextureSelector.SetSelectedIndex(floorTextureUnifier.Value().(int))
	mode.floorTextureRotationsBox.SetSelectedItem(mode.floorTextureRotationsItems[floorTextureRotationsUnifier.Value().(int)])
	mode.ceilingTextureSelector.SetSelectedIndex(ceilingTextureUnifier.Value().(int))
	mode.ceilingTextureRotationsBox.SetSelectedItem(mode.floorTextureRotationsItems[ceilingTextureRotationsUnifier.Value().(int)])
	mode.wallTextureSelector.SetSelectedIndex(wallTextureUnifier.Value().(int))
	mode.wallTextureOffsetBox.SetSelectedItem(mode.wallTextureOffsetItems[wallTextureOffsetUnifier.Value().(dataModel.HeightUnit)])
	mode.useAdjacentWallTextureBox.SetSelectedItem(mode.useAdjacentWallTextureItems[useAdjacentWallTextureUnifier.Value().(string)])
}

func (mode *LevelMapMode) changeSelectedTileProperties(modifier func(*dataModel.TileProperties)) {
	properties := &dataModel.TileProperties{}

	if !mode.levelAdapter.IsCyberspace() {
		properties.RealWorld = &dataModel.RealWorldTileProperties{}
	}
	modifier(properties)
	mode.levelAdapter.RequestTilePropertyChange(mode.selectedTiles, properties)
}

func (mode *LevelMapMode) onTilePropertyChangeRequested(item controls.ComboBoxItem) {
	propertyItem := item.(*tilePropertyItem)
	mode.changeSelectedTileProperties(func(properties *dataModel.TileProperties) {
		propertyItem.setter(properties, propertyItem.value)
	})
}

func (mode *LevelMapMode) onFloorTextureChanged(index int) {
	mode.changeSelectedTileProperties(func(properties *dataModel.TileProperties) {
		properties.RealWorld.FloorTexture = &index
	})
}

func (mode *LevelMapMode) onCeilingTextureChanged(index int) {
	mode.changeSelectedTileProperties(func(properties *dataModel.TileProperties) {
		properties.RealWorld.CeilingTexture = &index
	})
}

func (mode *LevelMapMode) onWallTextureChanged(index int) {
	mode.changeSelectedTileProperties(func(properties *dataModel.TileProperties) {
		properties.RealWorld.WallTexture = &index
	})
}
