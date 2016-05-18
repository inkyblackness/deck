package geometry

// DynamicModel is a Model implementation that allows addition of content
type DynamicModel struct {
	DynamicNode

	vertices []Vertex
}

// NewDynamicModel returns a new DynamicModel instance.
func NewDynamicModel() *DynamicModel {
	return new(DynamicModel)
}

// VertexCount is the VertexContainer implementation
func (model *DynamicModel) VertexCount() int {
	return len(model.vertices)
}

// Vertex is the VertexContainer implementation
func (model *DynamicModel) Vertex(index int) Vertex {
	return model.vertices[index]
}

// AddVertex adds the given vertex to the model.
func (model *DynamicModel) AddVertex(vertex Vertex) {
	model.vertices = append(model.vertices, vertex)
}
