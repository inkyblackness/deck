package chunk

// TypeID is a numerical presentation of a chunk type
type TypeID byte

const (
	compressionFlag = byte(0x01)
	directoryFlag   = byte(0x02)
)

// BasicChunkType specifies an uncompressed, flat chunk
const BasicChunkType = TypeID(0x00)

func (id TypeID) String() (result string) {
	if id.IsCompressed() {
		result += "Compressed"
	}
	if id.HasDirectory() {
		result += "Directory"
	}

	if result == "" {
		result = "Basic"
	}

	return
}

func (id TypeID) hasFlag(flag byte) bool {
	return (byte(id) & flag) != 0
}

func (id TypeID) setFlag(flag byte) TypeID {
	return TypeID(byte(id) | flag)
}

func (id TypeID) clearFlag(flag byte) TypeID {
	return TypeID(byte(id) & ^flag)
}

// IsCompressed returns true if the type specifies compression
func (id TypeID) IsCompressed() bool {
	return id.hasFlag(compressionFlag)
}

// WithCompression returns a TypeID marked with compression
func (id TypeID) WithCompression() TypeID {
	return id.setFlag(compressionFlag)
}

// WithoutCompression returns a TypeID without compression
func (id TypeID) WithoutCompression() TypeID {
	return id.clearFlag(compressionFlag)
}

// HasDirectory returns true if the type has a directory
func (id TypeID) HasDirectory() bool {
	return id.hasFlag(directoryFlag)
}

// WithDirectory returns a TypeID that specifies a directory chunk
func (id TypeID) WithDirectory() TypeID {
	return id.setFlag(directoryFlag)
}

// WithoutDirectory returns a TypeID that specifies a flat chunk
func (id TypeID) WithoutDirectory() TypeID {
	return id.clearFlag(directoryFlag)
}
