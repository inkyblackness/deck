package viewmodel

type TestingNodeVisitor struct {
	sectionSelectionNodes []Node
	sectionNodes          []Node
	valueSelectionNodes   []Node
	containerNodes        []Node
	tableNodes            []Node
	stringValueNodes      []Node
	boolValueNodes        []Node
	actionNodes           []Node
}

func NewTestingNodeVisitor() *TestingNodeVisitor {
	return &TestingNodeVisitor{}
}

func (visitor *TestingNodeVisitor) SectionSelection(node *SectionSelectionNode) {
	visitor.sectionSelectionNodes = append(visitor.sectionSelectionNodes, node)
}

func (visitor *TestingNodeVisitor) Section(node *SectionNode) {
	visitor.sectionNodes = append(visitor.sectionNodes, node)
}

func (visitor *TestingNodeVisitor) ValueSelection(node *ValueSelectionNode) {
	visitor.valueSelectionNodes = append(visitor.valueSelectionNodes, node)
}

func (visitor *TestingNodeVisitor) StringValue(node *StringValueNode) {
	visitor.stringValueNodes = append(visitor.stringValueNodes, node)
}

func (visitor *TestingNodeVisitor) BoolValue(node *BoolValueNode) {
	visitor.boolValueNodes = append(visitor.boolValueNodes, node)
}

func (visitor *TestingNodeVisitor) Container(node *ContainerNode) {
	visitor.containerNodes = append(visitor.containerNodes, node)
}

func (visitor *TestingNodeVisitor) Table(node *TableNode) {
	visitor.tableNodes = append(visitor.tableNodes, node)
}

func (visitor *TestingNodeVisitor) Action(node *ActionNode) {
	visitor.actionNodes = append(visitor.actionNodes, node)
}
