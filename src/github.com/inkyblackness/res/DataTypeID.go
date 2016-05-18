package res

// DataTypeID is a basic identification of a data block
type DataTypeID byte

const (
	// Palette refers to palette data blocks
	Palette = DataTypeID(0x00)
	// Text refers to texts
	Text = DataTypeID(0x01)
	// Bitmap refers to images
	Bitmap = DataTypeID(0x02)
	// Font refers to fonts (text and icons)
	Font = DataTypeID(0x03)
	// VideoClip refers to movies
	VideoClip = DataTypeID(0x04)
	// Sound refers to audio samples
	Sound = DataTypeID(0x07)
	// Geometry refers to 3D models
	Geometry = DataTypeID(0x0F)
	// Media refers to audio logs/cutscenes
	Media = DataTypeID(0x11)
	// Map refers to level data
	Map = DataTypeID(0x30)
)
