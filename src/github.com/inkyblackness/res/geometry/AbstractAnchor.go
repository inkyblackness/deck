package geometry

// abstractAnchor is the abstract base type for an Anchor implementation.
type abstractAnchor struct {
	normal    Vector
	reference Vector
}

// Normal returns the normal vector of the anchor.
func (anchor *abstractAnchor) Normal() Vector {
	return anchor.normal
}

// Reference returns the position of the anchor.
func (anchor *abstractAnchor) Reference() Vector {
	return anchor.reference
}
