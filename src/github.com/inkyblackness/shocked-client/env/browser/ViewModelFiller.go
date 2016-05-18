package browser

import (
	"github.com/Archs/js/gopherjs-ko"
	"github.com/gopherjs/gopherjs/js"

	"github.com/inkyblackness/shocked-client/viewmodel"
)

type viewModelFiller struct {
	object *js.Object
}

func newViewModelFiller() *viewModelFiller {
	return &viewModelFiller{}
}

func (filler *viewModelFiller) Section(node *viewmodel.SectionNode) {
	filler.object = js.Global.Get("Object").New()

	filler.object.Set("type", "section")
	filler.object.Set("label", node.Label())
	{
		availableFiller := newViewModelFiller()
		node.Available().Specialize(availableFiller)
		filler.object.Set("available", availableFiller.object)
	}
	{
		nodes := node.Get()
		objNodes := make([]*js.Object, len(nodes))

		for index, subNode := range nodes {
			subFiller := newViewModelFiller()
			subNode.Specialize(subFiller)
			objNodes[index] = subFiller.object
		}
		filler.object.Set("nodes", objNodes)
	}
}

func (filler *viewModelFiller) SectionSelection(node *viewmodel.SectionSelectionNode) {
	filler.ValueSelection(node.Selection())
	filler.object.Set("type", "sectionSelection")
	filler.object.Set("label", node.Label())

	nodeSections := node.Sections()
	objSections := js.Global.Get("Object").New()
	for key, section := range nodeSections {
		sectionFiller := newViewModelFiller()
		section.Specialize(sectionFiller)
		objSections.Set(key, sectionFiller.object)
	}
	filler.object.Set("sections", objSections)
}

func (filler *viewModelFiller) ValueSelection(node *viewmodel.ValueSelectionNode) {
	filler.object = js.Global.Get("Object").New()
	filler.object.Set("type", "valueSelection")
	filler.object.Set("label", node.Label())
	{
		selectedFiller := newViewModelFiller()
		node.Selected().Specialize(selectedFiller)
		filler.object.Set("selected", selectedFiller.object)
	}
	{
		observable := ko.NewObservableArray(node.Values())

		filler.object.Set("values", observable.ToJS())
		node.Subscribe(func(newValues []string) {
			observable.Set(newValues)
		})
	}
}

func (filler *viewModelFiller) BoolValue(node *viewmodel.BoolValueNode) {
	observable := ko.NewObservable(node.Get())

	filler.object = observable.ToJS()
	filler.object.Set("type", "bool")
	filler.object.Set("label", node.Label())
	node.Subscribe(func(newValue bool) {
		if observable.Get().Bool() != newValue {
			observable.Set(newValue)
		}
	})
	observable.Subscribe(func(jsValue *js.Object) {
		newValue := jsValue.Bool()

		if node.Get() != newValue {
			node.Set(newValue)
		}
	})
}

func (filler *viewModelFiller) StringValue(node *viewmodel.StringValueNode) {
	observable := ko.NewObservable(node.Get())

	filler.object = js.Global.Get("Object").New()
	filler.object.Set("type", "string")
	filler.object.Set("label", node.Label())
	filler.object.Set("readonly", !node.Editable())
	filler.object.Set("data", observable.ToJS())
	node.Subscribe(func(newValue string) {
		if observable.Get().String() != newValue {
			observable.Set(newValue)
		}
	})
	observable.Subscribe(func(jsValue *js.Object) {
		newValue := jsValue.String()

		if jsValue == js.Undefined {
			newValue = ""
		}
		if node.Get() != newValue {
			node.Set(newValue)
		}
	})
}

func (filler *viewModelFiller) Container(node *viewmodel.ContainerNode) {
	filler.object = js.Global.Get("Object").New()
	filler.object.Set("type", "container")
	filler.object.Set("label", node.Label())

	dataObject := js.Global.Get("Object").New()
	filler.object.Set("data", dataObject)
	for name, sub := range node.Get() {
		subFiller := newViewModelFiller()
		sub.Specialize(subFiller)
		dataObject.Set(name, subFiller.object)
	}
}

func (filler *viewModelFiller) Table(node *viewmodel.TableNode) {
	observable := ko.NewObservableArray()
	setEntries := func(nodeRows []*viewmodel.ContainerNode) {
		objRows := make([]*js.Object, len(nodeRows))

		for index, nodeRow := range nodeRows {
			subFiller := newViewModelFiller()
			nodeRow.Specialize(subFiller)
			objRows[index] = subFiller.object
		}
		observable.Set(objRows)
	}

	filler.object = observable.ToJS()
	filler.object.Set("type", "table")
	filler.object.Set("label", node.Label())

	setEntries(node.Get())
	node.Subscribe(setEntries)
}

func (filler *viewModelFiller) Action(node *viewmodel.ActionNode) {
	filler.object = js.Global.Get("Object").New()
	filler.object.Set("type", "action")
	filler.object.Set("label", node.Label())
	filler.object.Set("act", node.Act)
}
