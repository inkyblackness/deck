package core

// DataNode represents a container with data.
type DataNode interface {
	// Parent returns the parent node or nil if none known.
	Parent() DataNode
	// Children returns all currently available DataNodes from this node.
	Children() []DataNode
	// Info returns human readable information about this node.
	Info() string
	// ID returns the identification for this node. The returned value must be
	// the same by which the parent resolves this node.
	ID() string
	// Resolve returns a DataNode this node knows for the given ID.
	Resolve(string) DataNode
	// Data returns the data of the given node or nil if no data available.
	Data() []byte
	// UnknownData returns a byte array similar to Data(), but with known data
	// cleared to be 0x00.
	UnknownData() []byte
}
