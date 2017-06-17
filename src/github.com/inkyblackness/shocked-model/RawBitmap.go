package model

// RawBitmap is a simple palette based image.
type RawBitmap struct {
	// Width of the image in pixel
	Width int
	// Height of the image in pixel
	Height int
	// Pixel data is provided as base64 encoded byte string, with the stride equal the width.
	Pixels string
}
