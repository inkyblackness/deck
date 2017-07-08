package modes

import (
	"math"
	"sort"

	"github.com/inkyblackness/res/data/interpreters"
	"github.com/inkyblackness/shocked-client/editor/model"
	"github.com/inkyblackness/shocked-client/graphics/controls"
	"github.com/inkyblackness/shocked-client/ui"
)

type disposableControl interface {
	Dispose()
}

type enumItem struct {
	value       uint32
	displayName string
}

func (item *enumItem) String() string {
	return item.displayName
}

type bitfieldItem struct {
	displayName string
	shift       uint32
	mask        uint32
	maxValue    uint32
}

func (item *bitfieldItem) String() string {
	return item.displayName
}

type propertyEntry struct {
	title   *controls.Label
	control disposableControl
}

type propertyUpdateFunction func(currentValue, parameter uint32) uint32

func setUpdate() propertyUpdateFunction {
	return func(currentValue, parameter uint32) uint32 { return parameter }
}

func maskedUpdate(shift, mask uint32) propertyUpdateFunction {
	return func(currentValue, parameter uint32) uint32 {
		return (parameter << shift) | (currentValue & ^mask)
	}
}

type propertyChangeHandler func(key string, parameter uint32, update propertyUpdateFunction)

type objectTypeItemsRetriever func(objectClass int) []controls.ComboBoxItem

type propertyPanel struct {
	area          *ui.Area
	builder       *controlPanelBuilder
	changeHandler propertyChangeHandler
	entries       []*propertyEntry

	objectItemsForClass objectTypeItemsRetriever

	// selectionCache keeps the most recently selected mask of a bitfield, per key.
	// This helps re-selecting the same entry on re-creation.
	selectionCache map[string]uint32
}

func newPropertyPanel(parentBuilder *controlPanelBuilder, changeHandler propertyChangeHandler,
	objectItemsForClass objectTypeItemsRetriever) *propertyPanel {
	panel := &propertyPanel{
		changeHandler:       changeHandler,
		objectItemsForClass: objectItemsForClass,
		selectionCache:      make(map[string]uint32)}

	panel.area, panel.builder = parentBuilder.addDynamicSection(true, panel.Bottom)

	return panel
}

func (panel *propertyPanel) SetVisible(visible bool) {
	panel.area.SetVisible(visible)
}

func (panel *propertyPanel) Bottom() ui.Anchor {
	return panel.builder.bottom()
}

func (panel *propertyPanel) Reset() {
	for _, entry := range panel.entries {
		entry.title.Dispose()
		entry.control.Dispose()
	}
	panel.entries = []*propertyEntry{}
	panel.builder.reset()
}

func (panel *propertyPanel) NewSimplifier(key string, unifiedValue int64) *interpreters.Simplifier {
	simplifier := interpreters.NewSimplifier(func(minValue, maxValue int64) {
		slider := panel.NewSlider(key, "", setUpdate())
		slider.SetRange(minValue, maxValue)
		if unifiedValue != math.MinInt64 {
			slider.SetValue(unifiedValue)
		}
	})

	simplifier.SetEnumValueHandler(func(values map[uint32]string) {
		box := panel.NewComboBox(key, "", setUpdate())
		valueKeys := make([]uint32, 0, len(values))
		for valueKey := range values {
			valueKeys = append(valueKeys, valueKey)
		}
		sort.Slice(valueKeys, func(indexA, indexB int) bool { return valueKeys[indexA] < valueKeys[indexB] })
		items := make([]controls.ComboBoxItem, len(valueKeys))
		var selectedItem controls.ComboBoxItem
		for index, valueKey := range valueKeys {
			items[index] = &enumItem{valueKey, values[valueKey]}
			if int64(valueKey) == unifiedValue {
				selectedItem = items[index]
			}
		}
		box.SetItems(items)
		box.SetSelectedItem(selectedItem)
	})

	simplifier.SetBitfieldHandler(func(values map[uint32]string) {
		masks := make([]uint32, 0, len(values))
		var valueSlider *controls.Slider
		var currentUpdate propertyUpdateFunction

		onFieldSelectionChanged := func(boxItem controls.ComboBoxItem) {
			item := boxItem.(*bitfieldItem)

			currentUpdate = maskedUpdate(item.shift, item.mask)
			valueSlider.SetRange(0, int64(item.maxValue))
			if unifiedValue != math.MinInt64 {
				valueSlider.SetValue(int64((uint32(unifiedValue) & item.mask) >> item.shift))
			} else {
				valueSlider.SetValueUndefined()
			}

			panel.selectionCache[key] = item.mask
		}

		selectionTitle, selectionBox := panel.builder.addComboProperty(key+"-Part", onFieldSelectionChanged)
		panel.entries = append(panel.entries, &propertyEntry{selectionTitle, selectionBox})
		valueSlider = panel.NewSlider(key, "PartValue", func(currentValue, parameter uint32) uint32 {
			return currentUpdate(currentValue, parameter)
		})

		for mask := range values {
			masks = append(masks, mask)
		}
		sort.Slice(masks, func(indexA, indexB int) bool { return masks[indexA] < masks[indexB] })
		var items []controls.ComboBoxItem
		var selectedItem controls.ComboBoxItem
		for _, mask := range masks {
			item := &bitfieldItem{
				displayName: values[mask],
				mask:        mask,
				shift:       0,
				maxValue:    mask}

			for (item.maxValue & 1) == 0 {
				item.shift++
				item.maxValue >>= 1
			}
			items = append(items, item)
			if (selectedItem == nil) || (panel.selectionCache[key] == mask) {
				selectedItem = item
			}
		}
		selectionBox.SetItems(items)
		selectionBox.SetSelectedItem(selectedItem)
		onFieldSelectionChanged(selectedItem)
	})

	simplifier.SetSpecialHandler("ObjectType", func() {
		var typeBox *controls.ComboBox
		setTypeBox := func(objectID model.ObjectID) {
			typeItems := panel.objectItemsForClass(objectID.Class())
			selectedIndex := -1

			typeBox.SetItems(typeItems)
			for index, item := range typeItems {
				typeItem := item.(*objectTypeItem)
				if typeItem.id == objectID {
					selectedIndex = index
				}
			}

			if selectedIndex >= 0 {
				typeBox.SetSelectedItem(typeItems[selectedIndex])
			} else {
				typeBox.SetSelectedItem(nil)
			}
		}
		classTitle, classBox := panel.builder.addComboProperty(key+"-Class", func(boxItem controls.ComboBoxItem) {
			classItem := boxItem.(*objectClassItem)
			objectID := model.MakeObjectID(int(classItem.class), 0, 0)
			setTypeBox(objectID)

			panel.changeHandler(key, uint32(objectID.ToInt()), setUpdate())
		})
		panel.entries = append(panel.entries, &propertyEntry{classTitle, classBox})
		typeTitle, typeBox := panel.builder.addComboProperty(key+"-Type", func(boxItem controls.ComboBoxItem) {
			typeItem := boxItem.(*objectTypeItem)
			panel.changeHandler(key, uint32(typeItem.id.ToInt()), setUpdate())
		})
		panel.entries = append(panel.entries, &propertyEntry{typeTitle, typeBox})

		classItems := make([]controls.ComboBoxItem, len(classNames))
		for index := range classNames {
			classItems[index] = &objectClassItem{index}
		}
		classBox.SetItems(classItems)
		if unifiedValue != math.MinInt64 {
			objectID := model.ObjectID(unifiedValue)
			classBox.SetSelectedItem(classItems[objectID.Class()])
			setTypeBox(objectID)
		}
	})

	simplifier.SetSpecialHandler("Mistake", func() {})
	simplifier.SetSpecialHandler("Ignored", func() {})

	return simplifier
}

func (panel *propertyPanel) fullName(key, nameSuffix string) (fullName string) {
	fullName = key
	if len(nameSuffix) > 0 {
		fullName += "-" + nameSuffix
	}
	return
}

func (panel *propertyPanel) NewSlider(key string, nameSuffix string, update propertyUpdateFunction) *controls.Slider {
	fullName := panel.fullName(key, nameSuffix)
	title, control := panel.builder.addSliderProperty(fullName, func(newValue int64) {
		panel.changeHandler(key, uint32(newValue), update)
	})

	panel.entries = append(panel.entries, &propertyEntry{title, control})

	return control
}

func (panel *propertyPanel) NewComboBox(key string, nameSuffix string, update propertyUpdateFunction) *controls.ComboBox {
	fullName := panel.fullName(key, nameSuffix)
	title, control := panel.builder.addComboProperty(fullName, func(item controls.ComboBoxItem) {
		enumItem := item.(*enumItem)
		panel.changeHandler(key, enumItem.value, update)
	})

	panel.entries = append(panel.entries, &propertyEntry{title, control})

	return control
}

func (panel *propertyPanel) NewTextureSelector(key string, nameSuffix string,
	update propertyUpdateFunction, textureProvider controls.TextureProvider) *controls.TextureSelector {
	fullName := panel.fullName(key, nameSuffix)
	title, control := panel.builder.addTextureProperty(fullName, textureProvider, func(newIndex int) {
		panel.changeHandler(key, uint32(newIndex), update)
	})

	panel.entries = append(panel.entries, &propertyEntry{title, control})

	return control
}
