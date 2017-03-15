package ui

import (
	"github.com/inkyblackness/shocked-client/ui/events"
)

// RenderFunction is called when an area wants to render its content.
type RenderFunction func(*Area)

// EventHandler is called for events dispatched to the area.
type EventHandler func(*Area, events.Event) bool

// Area specifies one rectangular area within the user-interface stack.
type Area struct {
	parent   *Area
	children []*Area

	focusedArea *Area

	left   Anchor
	top    Anchor
	right  Anchor
	bottom Anchor

	visible bool

	onRender     RenderFunction
	eventHandler map[events.EventType]EventHandler
}

// Remove removes the area from the parent.
func (area *Area) Remove() {
	area.ReleaseFocus()
	if area.parent != nil {
		area.parent.removeChild(area)
		area.parent = nil
	}
}

// IsVisible returns true if the area is currently visible.
func (area *Area) IsVisible() bool {
	return area.visible
}

// SetVisible determines whether the area (and all of its children) shall
// be visible and target for events.
// Invisible areas are not rendered and will not handle any events.
func (area *Area) SetVisible(visible bool) {
	area.visible = visible
	if !area.IsVisible() {
		area.ReleaseFocus()
	}
}

// Root returns the area at the base of the UI tree.
func (area *Area) Root() (root *Area) {
	root = area
	if root.parent != nil {
		root = root.parent.Root()
	}
	return
}

func (area *Area) isChild(other *Area) (result bool) {
	for _, child := range area.children {
		if child == other {
			result = true
		}
	}
	return
}

func (area *Area) removeChild(child *Area) {
	newChildren := []*Area{}

	for _, other := range area.children {
		if other != child {
			newChildren = append(newChildren, other)
		}
	}
	area.children = newChildren
}

func (area *Area) currentChildren() []*Area {
	children := make([]*Area, len(area.children))
	copy(children[:], area.children[:])
	return children
}

func (area *Area) isRoot() bool {
	return area.parent == nil
}

// Left returns the left anchor.
func (area *Area) Left() Anchor {
	return area.left
}

// Top returns the top anchor.
func (area *Area) Top() Anchor {
	return area.top
}

// Right returns the right anchor.
func (area *Area) Right() Anchor {
	return area.right
}

// Bottom returns the bottom anchor.
func (area *Area) Bottom() Anchor {
	return area.bottom
}

// Render first renders this area, then sequentially all children.
func (area *Area) Render() {
	if area.IsVisible() {
		area.onRender(area)
		for _, child := range area.children {
			child.Render()
		}
	}
}

// HandleEvent tries to process the given event.
// It returns true if the area consumed the event.
func (area *Area) HandleEvent(event events.Event) (consumed bool) {
	if area.IsVisible() {
		if area.focusedArea != nil {
			consumed = area.focusedArea.HandleEvent(event)
		}
		if !consumed {
			consumed = area.tryEventHandlerFor(event)
		}
	}

	return
}

// DispatchPositionalEvent tries to find an event handler in this areas
// UI tree at the position of the event. The event is tried depth-first,
// before trying to handle it within this area.
func (area *Area) DispatchPositionalEvent(event events.PositionalEvent) (consumed bool) {
	if area.IsVisible() {
		if area.focusedArea != nil {
			consumed = area.focusedArea.DispatchPositionalEvent(event)
		}
		if !consumed {
			children := area.currentChildren()
			x, y := event.Position()

			for childIndex := len(children) - 1; !consumed && (childIndex >= 0); childIndex-- {
				child := children[childIndex]
				if area.isChild(child) && (child != area.focusedArea) &&
					(x >= child.Left().Value()) && (x < child.Right().Value()) &&
					(y >= child.Top().Value()) && (y < child.Bottom().Value()) {
					consumed = child.DispatchPositionalEvent(event)
				}
			}
		}
		if !consumed {
			consumed = area.tryEventHandlerFor(event)
		}
	}

	return
}

func (area *Area) tryEventHandlerFor(event events.Event) (consumed bool) {
	handler, existing := area.eventHandler[event.EventType()]

	if existing {
		consumed = handler(area, event)
	}

	return
}

// HasFocus returns true if this area (or any child) currently has the focus.
// The root area always has focus.
func (area *Area) HasFocus() bool {
	return area.isRoot() || area.parent.hasAreaFocus(area)
}

func (area *Area) hasAreaFocus(child *Area) bool {
	return area.focusedArea == child
}

// RequestFocus sets this area (and all of its parents) first in receiving events.
// Any previously focused area not in the parent list of this area will lose its focus.
func (area *Area) RequestFocus() {
	area.loseFocus()
	if area.parent != nil {
		area.parent.requestFocusFor(area)
	}
}

func (area *Area) requestFocusFor(child *Area) {
	area.RequestFocus()
	if (area.focusedArea != child) && (area.focusedArea != nil) {
		area.focusedArea.loseFocus()
	}
	area.focusedArea = child
}

// ReleaseFocus lets this area (and any of its children) lose focus.
func (area *Area) ReleaseFocus() {
	area.loseFocus()
	if area.HasFocus() && (area.parent != nil) {
		area.parent.ReleaseFocus()
	}
}

func (area *Area) loseFocus() {
	if area.focusedArea != nil {
		area.focusedArea.loseFocus()
		area.focusedArea = nil
	}
}
