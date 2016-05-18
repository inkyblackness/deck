package dos

const (
	// MagicHeader is the header value in a texture properties file.
	MagicHeader = uint32(0x09)
	// Size of the magic header
	MagicHeaderSize = 4
	// How big the table is within the file
	TableSize = uint32(4000)
)
