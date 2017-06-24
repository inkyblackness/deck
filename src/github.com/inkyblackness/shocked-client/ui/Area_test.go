package ui

import (
	"github.com/inkyblackness/shocked-client/ui/events"

	check "gopkg.in/check.v1"
)

type testingEvent struct {
	eventType events.EventType
}

func (event *testingEvent) EventType() events.EventType {
	return event.eventType
}

type AreaSuite struct {
	builder *AreaBuilder
}

var _ = check.Suite(&AreaSuite{})

func (suite *AreaSuite) aPositionalEvent(x, y float32) *events.MouseEvent {
	event := events.InitMouseEvent(events.EventType("test.positional"), x, y, 0, 0)

	return &event
}

func (suite *AreaSuite) SetUpTest(c *check.C) {
	suite.builder = NewAreaBuilder()
	suite.builder.SetRight(NewAbsoluteAnchor(100.0))
	suite.builder.SetBottom(NewAbsoluteAnchor(100.0))
}

func (suite *AreaSuite) TestRenderCallsOtherAreas(c *check.C) {
	renderCounter := 0
	renderCalls := make(map[int]int)
	renderFunc := func(index int) func(*Area) {
		return func(*Area) {
			renderCalls[index] = renderCounter
			renderCounter++
		}
	}
	parent := suite.builder.Build()
	NewAreaBuilder().SetParent(parent).OnRender(renderFunc(0)).Build()
	NewAreaBuilder().SetParent(parent).OnRender(renderFunc(1)).Build()

	parent.Render()

	c.Check(renderCalls, check.DeepEquals, map[int]int{0: 0, 1: 1})
}

func (suite *AreaSuite) TestHandleEventReturnsFalseForUnhandledEvent(c *check.C) {
	area := suite.builder.Build()
	event1 := &testingEvent{events.EventType("testEvent")}

	c.Check(area.HandleEvent(event1), check.Equals, false)
}

func (suite *AreaSuite) TestHandleEventCallsRegisteredHandler(c *check.C) {
	eventType := events.EventType("registeredEvent")
	called := false
	handler := func(area *Area, event events.Event) bool {
		called = true
		return false
	}
	suite.builder.OnEvent(eventType, handler)
	area := suite.builder.Build()
	event1 := &testingEvent{eventType}

	area.HandleEvent(event1)

	c.Check(called, check.Equals, true)
}

func (suite *AreaSuite) TestHandleEventReturnsResultFromRegisteredHandler_A(c *check.C) {
	eventType := events.EventType("registeredEvent")
	handler := func(area *Area, event events.Event) bool {
		return false
	}
	suite.builder.OnEvent(eventType, handler)
	area := suite.builder.Build()
	event1 := &testingEvent{eventType}

	c.Check(area.HandleEvent(event1), check.Equals, false)
}

func (suite *AreaSuite) TestHandleEventReturnsResultFromRegisteredHandler_B(c *check.C) {
	eventType := events.EventType("registeredEvent")
	handler := func(area *Area, event events.Event) bool {
		return true
	}
	suite.builder.OnEvent(eventType, handler)
	area := suite.builder.Build()
	event1 := &testingEvent{eventType}

	c.Check(area.HandleEvent(event1), check.Equals, true)
}

func (suite *AreaSuite) TestDispatchPositionalEventCallsHandleEventIfNoChildMatches(c *check.C) {
	testEvent := suite.aPositionalEvent(10.0, 10.0)
	called := false
	handler := func(area *Area, event events.Event) bool {
		called = true
		return false
	}
	suite.builder.OnEvent(testEvent.EventType(), handler)
	area := suite.builder.Build()

	area.DispatchPositionalEvent(testEvent)

	c.Check(called, check.Equals, true)
}

func (suite *AreaSuite) TestDispatchPositionalEventCallsChildrenAtPositionHighestFirst(c *check.C) {
	testEvent := suite.aPositionalEvent(50.0, 50.0)
	handleSequence := []int{}
	aHandler := func(index int) EventHandler {
		return func(*Area, events.Event) bool {
			handleSequence = append(handleSequence, index)
			return false
		}
	}

	suite.builder.OnEvent(testEvent.EventType(), aHandler(0))
	area := suite.builder.Build()

	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.OnEvent(testEvent.EventType(), aHandler(1))
		subAreaBuilder.SetRight(area.Right())
		subAreaBuilder.SetBottom(area.Bottom())
		subAreaBuilder.SetParent(area)
		subAreaBuilder.Build()
	}
	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.OnEvent(testEvent.EventType(), aHandler(2))
		subAreaBuilder.SetRight(area.Right())
		subAreaBuilder.SetBottom(area.Bottom())
		subAreaBuilder.SetParent(area)
		subAreaBuilder.Build()
	}
	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.OnEvent(testEvent.EventType(), aHandler(3))
		subAreaBuilder.SetRight(NewAbsoluteAnchor(10.0))
		subAreaBuilder.SetBottom(NewAbsoluteAnchor(10.0))
		subAreaBuilder.SetParent(area)
		subAreaBuilder.Build()
	}

	area.DispatchPositionalEvent(testEvent)

	c.Check(handleSequence, check.DeepEquals, []int{2, 1, 0})
}

func (suite *AreaSuite) TestRootReturnsRootArea(c *check.C) {
	area := suite.builder.Build()
	var subArea *Area
	var subSubArea *Area

	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.SetParent(area)
		subArea = subAreaBuilder.Build()
	}
	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.SetParent(subArea)
		subSubArea = subAreaBuilder.Build()
	}

	c.Check(area.Root(), check.Equals, area)
	c.Check(subArea.Root(), check.Equals, area)
	c.Check(subSubArea.Root(), check.Equals, area)
}

func (suite *AreaSuite) TestHasFocusReturnsTrueForRoot(c *check.C) {
	area := suite.builder.Build()

	c.Check(area.HasFocus(), check.Equals, true)
}

func (suite *AreaSuite) TestHasFocusReturnsFalseForNonFocusedArea(c *check.C) {
	area := suite.builder.Build()
	var subArea *Area

	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.SetParent(area)
		subArea = subAreaBuilder.Build()
	}

	c.Check(subArea.HasFocus(), check.Equals, false)
}

func (suite *AreaSuite) TestRequestFocusSetsFocusForAreaRecursively(c *check.C) {
	area := suite.builder.Build()
	var subArea *Area
	var subSubArea *Area

	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.SetParent(area)
		subArea = subAreaBuilder.Build()
	}
	{
		subSubAreaBuilder := NewAreaBuilder()
		subSubAreaBuilder.SetParent(subArea)
		subSubArea = subSubAreaBuilder.Build()
	}

	subSubArea.RequestFocus()

	c.Check(subArea.HasFocus(), check.Equals, true)
	c.Check(subSubArea.HasFocus(), check.Equals, true)
}

func (suite *AreaSuite) TestReleaseFocusClearsFocusForAreaRecursively(c *check.C) {
	area := suite.builder.Build()
	var subArea *Area
	var subSubArea *Area

	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.SetParent(area)
		subArea = subAreaBuilder.Build()
	}
	{
		subSubAreaBuilder := NewAreaBuilder()
		subSubAreaBuilder.SetParent(subArea)
		subSubArea = subSubAreaBuilder.Build()
	}

	subSubArea.RequestFocus()
	subSubArea.ReleaseFocus()

	c.Check(subArea.HasFocus(), check.Equals, false)
	c.Check(subSubArea.HasFocus(), check.Equals, false)
}

func (suite *AreaSuite) TestRequestFocusClearsFocusForLostBranch(c *check.C) {
	area := suite.builder.Build()
	var subArea *Area
	var leftArea *Area
	var leftChildArea *Area
	var rightArea *Area
	var rightChildArea *Area

	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.SetParent(area)
		subArea = subAreaBuilder.Build()
	}
	{
		leftAreaBuilder := NewAreaBuilder()
		leftAreaBuilder.SetParent(subArea)
		leftArea = leftAreaBuilder.Build()
	}
	{
		leftChildAreaBuilder := NewAreaBuilder()
		leftChildAreaBuilder.SetParent(leftArea)
		leftChildArea = leftChildAreaBuilder.Build()
	}
	{
		rightAreaBuilder := NewAreaBuilder()
		rightAreaBuilder.SetParent(subArea)
		rightArea = rightAreaBuilder.Build()
	}
	{
		rightChildAreaBuilder := NewAreaBuilder()
		rightChildAreaBuilder.SetParent(rightArea)
		rightChildArea = rightChildAreaBuilder.Build()
	}

	leftChildArea.RequestFocus()
	rightChildArea.RequestFocus()

	c.Check(subArea.HasFocus(), check.Equals, true)
	c.Check(leftArea.HasFocus(), check.Equals, false)
	c.Check(leftChildArea.HasFocus(), check.Equals, false)
	c.Check(rightArea.HasFocus(), check.Equals, true)
	c.Check(rightChildArea.HasFocus(), check.Equals, true)
}

func (suite *AreaSuite) TestRequestFocusClearsFocusForPreviouslyFocusedDescendant(c *check.C) {
	area := suite.builder.Build()
	var subArea *Area
	var leftArea *Area
	var leftChildArea *Area

	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.SetParent(area)
		subArea = subAreaBuilder.Build()
	}
	{
		leftAreaBuilder := NewAreaBuilder()
		leftAreaBuilder.SetParent(subArea)
		leftArea = leftAreaBuilder.Build()
	}
	{
		leftChildAreaBuilder := NewAreaBuilder()
		leftChildAreaBuilder.SetParent(leftArea)
		leftChildArea = leftChildAreaBuilder.Build()
	}

	leftChildArea.RequestFocus()
	leftArea.RequestFocus()

	c.Check(subArea.HasFocus(), check.Equals, true)
	c.Check(leftArea.HasFocus(), check.Equals, true)
	c.Check(leftChildArea.HasFocus(), check.Equals, false)
}

func (suite *AreaSuite) TestHandleEventForwardsEventToAreaWithFocusBeforeHandlingItself(c *check.C) {
	event1 := &testingEvent{events.EventType("TestingEvent")}
	var subArea *Area
	var leftArea *Area
	handleCounter := 0
	handleCalls := make(map[string]int)
	aHandler := func(id string, consume bool) EventHandler {
		return func(*Area, events.Event) bool {
			handleCalls[id] = handleCounter
			handleCounter++
			return consume
		}
	}

	suite.builder.OnEvent(event1.EventType(), aHandler("root", false))
	area := suite.builder.Build()

	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.SetParent(area)
		subAreaBuilder.OnEvent(event1.EventType(), aHandler("subArea", true))
		subArea = subAreaBuilder.Build()
	}
	{
		leftAreaBuilder := NewAreaBuilder()
		leftAreaBuilder.SetParent(subArea)
		leftAreaBuilder.OnEvent(event1.EventType(), aHandler("left", false))
		leftArea = leftAreaBuilder.Build()
	}
	{
		rightAreaBuilder := NewAreaBuilder()
		rightAreaBuilder.SetParent(subArea)
		rightAreaBuilder.OnEvent(event1.EventType(), aHandler("right", false))
		rightAreaBuilder.Build()
	}

	leftArea.RequestFocus()
	area.HandleEvent(event1)

	c.Check(handleCalls, check.DeepEquals, map[string]int{"left": 0, "subArea": 1})
}

func (suite *AreaSuite) TestDispatchPositionalEventCallsFocusedItemBeforeChildren(c *check.C) {
	testEvent := suite.aPositionalEvent(50.0, 50.0)
	handleCounter := 0
	handleCalls := make(map[string]int)
	aHandler := func(id string, consume bool) EventHandler {
		return func(*Area, events.Event) bool {
			handleCalls[id] = handleCounter
			handleCounter++
			return consume
		}
	}
	var testedArea *Area

	suite.builder.OnEvent(testEvent.EventType(), aHandler("root", false))
	area := suite.builder.Build()

	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.OnEvent(testEvent.EventType(), aHandler("area1", true))
		subAreaBuilder.SetRight(area.Right())
		subAreaBuilder.SetBottom(area.Bottom())
		subAreaBuilder.SetParent(area)
		subAreaBuilder.Build()
	}
	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.OnEvent(testEvent.EventType(), aHandler("area2", false))
		subAreaBuilder.SetRight(area.Right())
		subAreaBuilder.SetBottom(area.Bottom())
		subAreaBuilder.SetParent(area)
		subAreaBuilder.Build()
	}
	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.OnEvent(testEvent.EventType(), aHandler("area3", false))
		subAreaBuilder.SetRight(NewAbsoluteAnchor(10.0))
		subAreaBuilder.SetBottom(NewAbsoluteAnchor(10.0))
		subAreaBuilder.SetParent(area)
		testedArea = subAreaBuilder.Build()
	}

	testedArea.RequestFocus()
	area.DispatchPositionalEvent(testEvent)

	c.Check(handleCalls, check.DeepEquals, map[string]int{"area3": 0, "area2": 1, "area1": 2})
}

func (suite *AreaSuite) TestDispatchPositionalEventCallsFocusedAreaHandlerWithoutFocusedChildBeforeChildren(c *check.C) {
	testEvent := suite.aPositionalEvent(50.0, 50.0)
	handleCounter := 0
	handleCalls := make(map[string]int)
	aHandler := func(id string, consume bool) EventHandler {
		return func(*Area, events.Event) bool {
			handleCalls[id] = handleCounter
			handleCounter++
			return consume
		}
	}
	var testedArea *Area

	suite.builder.OnEvent(testEvent.EventType(), aHandler("root", false))
	area := suite.builder.Build()

	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.OnEvent(testEvent.EventType(), aHandler("area1", false))
		subAreaBuilder.SetRight(area.Right())
		subAreaBuilder.SetBottom(area.Bottom())
		subAreaBuilder.SetParent(area)
		testedArea = subAreaBuilder.Build()
	}
	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.OnEvent(testEvent.EventType(), aHandler("area2", false))
		subAreaBuilder.SetRight(area.Right())
		subAreaBuilder.SetBottom(area.Bottom())
		subAreaBuilder.SetParent(testedArea)
		subAreaBuilder.Build()
	}

	testedArea.RequestFocus()
	area.DispatchPositionalEvent(testEvent)

	c.Check(handleCalls, check.DeepEquals, map[string]int{"root": 2, "area2": 1, "area1": 0})
}

func (suite *AreaSuite) TestDispatchPositionalEventCallsFocusedItemOnlyOnce(c *check.C) {
	testEvent := suite.aPositionalEvent(50.0, 50.0)
	handleCounter := 0
	handleCalls := make(map[string]int)
	aHandler := func(id string, consume bool) EventHandler {
		return func(*Area, events.Event) bool {
			handleCalls[id] = handleCounter
			handleCounter++
			return consume
		}
	}
	var testedArea *Area

	suite.builder.OnEvent(testEvent.EventType(), aHandler("root", false))
	area := suite.builder.Build()

	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.OnEvent(testEvent.EventType(), aHandler("area3", false))
		subAreaBuilder.SetRight(NewAbsoluteAnchor(100.0))
		subAreaBuilder.SetBottom(NewAbsoluteAnchor(100.0))
		subAreaBuilder.SetParent(area)
		testedArea = subAreaBuilder.Build()
	}

	testedArea.RequestFocus()
	area.DispatchPositionalEvent(testEvent)

	c.Check(handleCalls, check.DeepEquals, map[string]int{"area3": 0, "root": 1})
}

func (suite *AreaSuite) TestRemoveDisassociatesAreaFromParent(c *check.C) {
	var testedArea *Area

	area := suite.builder.Build()

	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.SetParent(area)
		testedArea = subAreaBuilder.Build()
	}

	testedArea.RequestFocus()
	testedArea.Remove()

	c.Check(testedArea.parent, check.IsNil)
	c.Check(area.focusedArea, check.IsNil)
	c.Check(area.children, check.DeepEquals, []*Area{})
}

func (suite *AreaSuite) TestDispatchPositionalEventCallsOnlySurvivingChildren(c *check.C) {
	testEvent := suite.aPositionalEvent(50.0, 50.0)
	aSubArea := func(parent *Area, handler EventHandler) *Area {
		subAreaBuilder := NewAreaBuilder()
		if handler != nil {
			subAreaBuilder.OnEvent(testEvent.EventType(), handler)
		}
		subAreaBuilder.SetRight(NewAbsoluteAnchor(100.0))
		subAreaBuilder.SetBottom(NewAbsoluteAnchor(100.0))
		subAreaBuilder.SetParent(parent)
		return subAreaBuilder.Build()
	}
	sub1Called := false
	sub2Called := false

	area := suite.builder.Build()

	aSubArea(area, func(area *Area, event events.Event) bool {
		sub1Called = true
		return false
	})
	sub2 := aSubArea(area, func(area *Area, event events.Event) bool {
		sub2Called = true
		return false
	})
	var sub3 *Area
	sub3 = aSubArea(area, func(area *Area, event events.Event) bool {
		sub2.Remove()
		sub3.Remove()
		return false
	})

	area.DispatchPositionalEvent(testEvent)

	c.Check(sub1Called, check.Equals, true)
	c.Check(sub2Called, check.Equals, false)
}

func (suite *AreaSuite) TestInvisibleAreaIsNotRendered(c *check.C) {
	renderCounter := 0
	renderCalls := make(map[int]int)
	renderFunc := func(index int) func(*Area) {
		return func(*Area) {
			renderCalls[index] = renderCounter
			renderCounter++
		}
	}
	suite.builder.OnRender(renderFunc(0))
	parent := suite.builder.Build()
	NewAreaBuilder().SetParent(parent).OnRender(renderFunc(1)).Build()

	parent.SetVisible(false)
	parent.Render()

	c.Check(renderCalls, check.DeepEquals, map[int]int{})
}

func (suite *AreaSuite) TestInvisibleAreaDoesNotHandleEvents(c *check.C) {
	event1 := &testingEvent{events.EventType("TestingEvent")}
	var subArea *Area
	var leftArea *Area
	handleCounter := 0
	handleCalls := make(map[string]int)
	aHandler := func(id string, consume bool) EventHandler {
		return func(*Area, events.Event) bool {
			handleCalls[id] = handleCounter
			handleCounter++
			return consume
		}
	}

	suite.builder.OnEvent(event1.EventType(), aHandler("root", false))
	area := suite.builder.Build()

	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.SetParent(area)
		subAreaBuilder.OnEvent(event1.EventType(), aHandler("subArea", true))
		subArea = subAreaBuilder.Build()
	}
	{
		leftAreaBuilder := NewAreaBuilder()
		leftAreaBuilder.SetParent(subArea)
		leftAreaBuilder.OnEvent(event1.EventType(), aHandler("left", false))
		leftArea = leftAreaBuilder.Build()
	}
	{
		rightAreaBuilder := NewAreaBuilder()
		rightAreaBuilder.SetParent(subArea)
		rightAreaBuilder.OnEvent(event1.EventType(), aHandler("right", false))
		rightAreaBuilder.Build()
	}

	leftArea.RequestFocus()
	area.SetVisible(false)
	area.HandleEvent(event1)

	c.Check(handleCalls, check.DeepEquals, map[string]int{})
}

func (suite *AreaSuite) TestInvisibleAreaDoesNotDispatchEvents(c *check.C) {
	testEvent := suite.aPositionalEvent(50.0, 50.0)
	handleCounter := 0
	handleCalls := make(map[string]int)
	aHandler := func(id string, consume bool) EventHandler {
		return func(*Area, events.Event) bool {
			handleCalls[id] = handleCounter
			handleCounter++
			return consume
		}
	}
	var testedArea *Area

	suite.builder.OnEvent(testEvent.EventType(), aHandler("root", false))
	area := suite.builder.Build()

	{
		subAreaBuilder := NewAreaBuilder()
		subAreaBuilder.OnEvent(testEvent.EventType(), aHandler("area3", false))
		subAreaBuilder.SetRight(NewAbsoluteAnchor(100.0))
		subAreaBuilder.SetBottom(NewAbsoluteAnchor(100.0))
		subAreaBuilder.SetParent(area)
		testedArea = subAreaBuilder.Build()
	}

	testedArea.RequestFocus()
	area.SetVisible(false)
	area.DispatchPositionalEvent(testEvent)

	c.Check(handleCalls, check.DeepEquals, map[string]int{})
}

func (suite *AreaSuite) TestInvisibleAreaLosesFocus(c *check.C) {
	parent := suite.builder.Build()
	subArea := NewAreaBuilder().SetParent(parent).Build()

	subArea.RequestFocus()
	parent.SetVisible(false)

	c.Check(subArea.HasFocus(), check.Equals, false)
}
