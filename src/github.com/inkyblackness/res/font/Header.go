package font

// Type identifies how the bits of a font are stored in the bitmap.
type Type uint16

const (
	// Monochrome fonts use single bits per output pixel.
	Monochrome = Type(0x0000)
	// Color fonts use one byte per pixel.
	Color = Type(0xCCCC)
)

// HeaderSize in bytes
const HeaderSize = 84

// Header describes the font header in the resource files
type Header struct {
	// Type of the font
	Type Type

	// Unknown0002 data
	Unknown0002 [34]byte

	// FirstCharacter index (inclusive)
	FirstCharacter uint16
	// LastCharacter index (inclusive)
	LastCharacter uint16

	// Unknown0028 data
	Unknown0028 [32]byte

	// XOffsetStart from beginning of file
	XOffsetStart uint32
	// BitmapStart from beginning of file
	BitmapStart uint32

	// Width of the bitmap in bytes
	Width uint16
	// Height of the bitmap in bytes
	Height uint16
}
