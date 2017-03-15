package controls

import (
	"fmt"

	"github.com/inkyblackness/shocked-client/env"
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"
)

// ComboBoxItem is the interface type for items within a combo box.
type ComboBoxItem interface{}

// SelectionChangeHandler is a callback for notifying the current selection.
type SelectionChangeHandler func(item ComboBoxItem)

// ComboBox provides the ability to select one item from a list.
type ComboBox struct {
	labelBuilder *LabelBuilder

	area         *ui.Area
	rectRenderer *graphics.RectangleRenderer

	selectedLabel *Label
	hintLabel     *Label

	selectionChangeHandler SelectionChangeHandler

	items        []ComboBoxItem
	selectedItem ComboBoxItem

	listArea       *ui.Area
	listItemCount  int
	listItemLabels []*Label
	listStartIndex int
}

// Dispose releases the resources.
func (box *ComboBox) Dispose() {
	box.hideList()
	box.selectedLabel.Dispose()
	box.hintLabel.Dispose()
	box.area.Remove()
}

// SetItems sets the lits of available items.
func (box *ComboBox) SetItems(items []ComboBoxItem) {
	box.hideList()
	box.items = items
}

// SetSelectedItem changes what is currently selected. Does not fire change handler.
func (box *ComboBox) SetSelectedItem(item ComboBoxItem) {
	if box.selectedItem != item {
		box.selectedItem = item
		if box.selectedItem != nil {
			box.selectedLabel.SetText(fmt.Sprintf("%v", item))
		} else {
			box.selectedLabel.SetText("")
		}
	}
}

func (box *ComboBox) onRender(area *ui.Area) {
	box.rectRenderer.Fill(area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
		graphics.RGBA(0.31, 0.56, 0.34, 0.8))
}

func (box *ComboBox) onMouseDown(area *ui.Area, event events.Event) (consumed bool) {
	mouseEvent := event.(*events.MouseButtonEvent)

	if mouseEvent.Buttons() == env.MousePrimary {
		if box.listArea == nil {
			box.showList()
		} else {
			box.hideList()
		}
		consumed = true
	}

	return
}

func (box *ComboBox) contains(area *ui.Area, event events.PositionalEvent) bool {
	x, y := event.Position()

	return (x >= area.Left().Value()) && (x < area.Right().Value()) &&
		(y >= area.Top().Value()) && (y < area.Bottom().Value())
}

func (box *ComboBox) showList() {
	if box.listArea == nil {
		listAreaBuilder := ui.NewAreaBuilder()
		root := box.area.Root()
		boxTop := box.area.Top().Value()
		boxBottom := box.area.Bottom().Value()
		boxHeight := boxBottom - boxTop
		box.listItemCount = len(box.items)
		if box.listItemCount > 6 {
			box.listItemCount = 6
		}
		listHeight := boxHeight * float32(box.listItemCount)
		var listTop ui.Anchor

		if (listHeight > (root.Bottom().Value() - boxBottom)) && (listHeight <= (boxTop - root.Top().Value())) {
			listTop = ui.NewOffsetAnchor(box.area.Top(), -listHeight)
		} else {
			listTop = ui.NewOffsetAnchor(box.area.Bottom(), 0)
		}

		listAreaBuilder.SetParent(root)
		listAreaBuilder.SetLeft(box.area.Left())
		listAreaBuilder.SetRight(box.area.Right())
		listAreaBuilder.SetTop(listTop)
		listAreaBuilder.SetBottom(ui.NewOffsetAnchor(listTop, listHeight))
		listAreaBuilder.OnRender(box.onListRender)
		listAreaBuilder.OnEvent(events.MouseButtonDownEventType, box.onListMouseDown)
		listAreaBuilder.OnEvent(events.MouseButtonUpEventType, box.onListMouseUp)
		listAreaBuilder.OnEvent(events.MouseScrollEventType, box.onListScroll)
		listAreaBuilder.OnEvent(events.MouseButtonClickedEventType, ui.SilentConsumer)

		box.listArea = listAreaBuilder.Build()
		box.listArea.RequestFocus()

		box.listItemLabels = make([]*Label, box.listItemCount)
		lastBottom := listTop

		box.labelBuilder.SetParent(box.listArea)
		box.labelBuilder.AlignedHorizontallyBy(LeftAligner)
		box.labelBuilder.SetLeft(ui.NewOffsetAnchor(box.area.Left(), 4))
		box.labelBuilder.SetRight(ui.NewOffsetAnchor(box.area.Right(), -4))
		for listIndex := 0; listIndex < box.listItemCount; listIndex++ {
			nextBottom := ui.NewOffsetAnchor(lastBottom, boxHeight)
			box.labelBuilder.SetTop(lastBottom)
			box.labelBuilder.SetBottom(nextBottom)
			box.listItemLabels[listIndex] = box.labelBuilder.Build()
			lastBottom = nextBottom
		}
		box.updateListItemLabels()
	}
}

func (box *ComboBox) hideList() {
	if box.listArea != nil {
		box.listArea.Remove()
		box.listArea = nil
		for _, label := range box.listItemLabels {
			label.Dispose()
		}
		box.listItemLabels = nil
	}
}

func (box *ComboBox) updateListItemLabels() {
	for listIndex, label := range box.listItemLabels {
		label.SetText(fmt.Sprintf("%v", box.items[box.listStartIndex+listIndex]))
	}
}

func (box *ComboBox) onListRender(area *ui.Area) {
	box.rectRenderer.Fill(area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
		graphics.RGBA(0.31, 0.56, 0.34, 0.7))
}

func (box *ComboBox) onListMouseDown(area *ui.Area, event events.Event) (consumed bool) {
	mouseEvent := event.(*events.MouseButtonEvent)

	if mouseEvent.Buttons() == env.MousePrimary {
		box.listArea.RequestFocus()
		consumed = true
	}

	return
}

func (box *ComboBox) onListMouseUp(area *ui.Area, event events.Event) (consumed bool) {
	mouseEvent := event.(*events.MouseButtonEvent)

	if mouseEvent.AffectedButtons() == env.MousePrimary {
		if box.listArea != nil {
			if box.listArea.HasFocus() {
				box.listArea.ReleaseFocus()
			}
			if box.contains(box.listArea, mouseEvent) {
				_, mouseY := mouseEvent.Position()
				chosenItem := ((mouseY - box.listArea.Top().Value()) * float32(box.listItemCount)) /
					(box.listArea.Bottom().Value() - box.listArea.Top().Value())
				box.hideList()
				box.onItemChosen(box.items[box.listStartIndex+int(chosenItem)])
			}
		}
		consumed = true
	}

	return
}

func (box *ComboBox) onItemChosen(item ComboBoxItem) {
	if item != box.selectedItem {
		box.SetSelectedItem(item)
		box.selectionChangeHandler(item)
	}
}

func (box *ComboBox) onListScroll(area *ui.Area, event events.Event) (consumed bool) {
	mouseEvent := event.(*events.MouseScrollEvent)
	_, dy := mouseEvent.Deltas()
	toScroll := func(available int) int {
		result := 1
		if result > available {
			result = available
		}
		return result
	}

	if dy < 0 {
		available := box.listStartIndex
		box.listStartIndex -= toScroll(available)
	} else if dy > 0 {
		available := len(box.items) - (box.listStartIndex + box.listItemCount)
		box.listStartIndex += toScroll(available)
	}
	box.updateListItemLabels()
	consumed = true

	return
}
