package data

import (
	"github.com/inkyblackness/res/serial"
	"github.com/inkyblackness/res/text"
)

// TODO this appears to be only used for the archive name. Consider
// creating a dedicated type for that and drop this generic one.

// String wraps an encoded string type.
type String struct {
	encoded []byte
}

// NewString returns a new object prepared to hold given amount of bytes.
func NewString(length int) *String {
	return &String{encoded: make([]byte, length)}
}

// Code serializes the string with given coder.
func (str *String) Code(coder serial.Coder) {
	coder.Code(str.encoded)
}

func (str *String) String() string {
	cp := text.DefaultCodepage()

	return "\"" + cp.Decode(str.encoded) + "\""
}
