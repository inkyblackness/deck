// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go

// Package sfnt implements a decoder for SFNT font file formats, including
// TrueType and OpenType.
package sfnt // import "golang.org/x/image/font/sfnt"

// This implementation was written primarily to the
// https://www.microsoft.com/en-us/Typography/OpenTypeSpecification.aspx
// specification. Additional documentation is at
// http://developer.apple.com/fonts/TTRefMan/
//
// The pyftinspect tool from https://github.com/fonttools/fonttools is useful
// for inspecting SFNT fonts.
//
// The ttfdump tool is also useful. For example:
//	ttfdump -t cmap ../testdata/CFFTest.otf dump.txt

import (
	"errors"
	"io"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/encoding/charmap"
)

// These constants are not part of the specifications, but are limitations used
// by this implementation.
const (
	// This value is arbitrary, but defends against parsing malicious font
	// files causing excessive memory allocations. For reference, Adobe's
	// SourceHanSansSC-Regular.otf has 65535 glyphs and:
	//	- its format-4  cmap table has  1581 segments.
	//	- its format-12 cmap table has 16498 segments.
	//
	// TODO: eliminate this constraint? If the cmap table is very large, load
	// some or all of it lazily (at the time Font.GlyphIndex is called) instead
	// of all of it eagerly (at the time Font.initialize is called), while
	// keeping an upper bound on the memory used? This will make the code in
	// cmap.go more complicated, considering that all of the Font methods are
	// safe to call concurrently, as long as each call has a different *Buffer.
	maxCmapSegments = 20000

	maxGlyphDataLength  = 64 * 1024
	maxHintBits         = 256
	maxNumTables        = 256
	maxRealNumberStrLen = 64 // Maximum length in bytes of the "-123.456E-7" representation.

	// (maxTableOffset + maxTableLength) will not overflow an int32.
	maxTableLength = 1 << 29
	maxTableOffset = 1 << 29
)

var (
	// ErrNotFound indicates that the requested value was not found.
	ErrNotFound = errors.New("sfnt: not found")

	errInvalidBounds        = errors.New("sfnt: invalid bounds")
	errInvalidCFFTable      = errors.New("sfnt: invalid CFF table")
	errInvalidCmapTable     = errors.New("sfnt: invalid cmap table")
	errInvalidGlyphData     = errors.New("sfnt: invalid glyph data")
	errInvalidHeadTable     = errors.New("sfnt: invalid head table")
	errInvalidKernTable     = errors.New("sfnt: invalid kern table")
	errInvalidLocaTable     = errors.New("sfnt: invalid loca table")
	errInvalidLocationData  = errors.New("sfnt: invalid location data")
	errInvalidMaxpTable     = errors.New("sfnt: invalid maxp table")
	errInvalidNameTable     = errors.New("sfnt: invalid name table")
	errInvalidPostTable     = errors.New("sfnt: invalid post table")
	errInvalidSourceData    = errors.New("sfnt: invalid source data")
	errInvalidTableOffset   = errors.New("sfnt: invalid table offset")
	errInvalidTableTagOrder = errors.New("sfnt: invalid table tag order")
	errInvalidUCS2String    = errors.New("sfnt: invalid UCS-2 string")
	errInvalidVersion       = errors.New("sfnt: invalid version")

	errUnsupportedCFFVersion           = errors.New("sfnt: unsupported CFF version")
	errUnsupportedCmapEncodings        = errors.New("sfnt: unsupported cmap encodings")
	errUnsupportedCompoundGlyph        = errors.New("sfnt: unsupported compound glyph")
	errUnsupportedGlyphDataLength      = errors.New("sfnt: unsupported glyph data length")
	errUnsupportedKernTable            = errors.New("sfnt: unsupported kern table")
	errUnsupportedRealNumberEncoding   = errors.New("sfnt: unsupported real number encoding")
	errUnsupportedNumberOfCmapSegments = errors.New("sfnt: unsupported number of cmap segments")
	errUnsupportedNumberOfHints        = errors.New("sfnt: unsupported number of hints")
	errUnsupportedNumberOfTables       = errors.New("sfnt: unsupported number of tables")
	errUnsupportedPlatformEncoding     = errors.New("sfnt: unsupported platform encoding")
	errUnsupportedPostTable            = errors.New("sfnt: unsupported post table")
	errUnsupportedTableOffsetLength    = errors.New("sfnt: unsupported table offset or length")
	errUnsupportedType2Charstring      = errors.New("sfnt: unsupported Type 2 Charstring")
)

// GlyphIndex is a glyph index in a Font.
type GlyphIndex uint16

// NameID identifies a name table entry.
//
// See the "Name IDs" section of
// https://www.microsoft.com/typography/otspec/name.htm
type NameID uint16

const (
	NameIDCopyright                  NameID = 0
	NameIDFamily                            = 1
	NameIDSubfamily                         = 2
	NameIDUniqueIdentifier                  = 3
	NameIDFull                              = 4
	NameIDVersion                           = 5
	NameIDPostScript                        = 6
	NameIDTrademark                         = 7
	NameIDManufacturer                      = 8
	NameIDDesigner                          = 9
	NameIDDescription                       = 10
	NameIDVendorURL                         = 11
	NameIDDesignerURL                       = 12
	NameIDLicense                           = 13
	NameIDLicenseURL                        = 14
	NameIDTypographicFamily                 = 16
	NameIDTypographicSubfamily              = 17
	NameIDCompatibleFull                    = 18
	NameIDSampleText                        = 19
	NameIDPostScriptCID                     = 20
	NameIDWWSFamily                         = 21
	NameIDWWSSubfamily                      = 22
	NameIDLightBackgroundPalette            = 23
	NameIDDarkBackgroundPalette             = 24
	NameIDVariationsPostScriptPrefix        = 25
)

// Units are an integral number of abstract, scalable "font units". The em
// square is typically 1000 or 2048 "font units". This would map to a certain
// number (e.g. 30 pixels) of physical pixels, depending on things like the
// display resolution (DPI) and font size (e.g. a 12 point font).
type Units int32

// scale returns x divided by unitsPerEm, rounded to the nearest fixed.Int26_6
// value (1/64th of a pixel).
func scale(x fixed.Int26_6, unitsPerEm Units) fixed.Int26_6 {
	if x >= 0 {
		x += fixed.Int26_6(unitsPerEm) / 2
	} else {
		x -= fixed.Int26_6(unitsPerEm) / 2
	}
	return x / fixed.Int26_6(unitsPerEm)
}

func u16(b []byte) uint16 {
	_ = b[1] // Bounds check hint to compiler.
	return uint16(b[0])<<8 | uint16(b[1])<<0
}

func u32(b []byte) uint32 {
	_ = b[3] // Bounds check hint to compiler.
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])<<0
}

// source is a source of byte data. Conceptually, it is like an io.ReaderAt,
// except that a common source of SFNT font data is in-memory instead of
// on-disk: a []byte containing the entire data, either as a global variable
// (e.g. "goregular.TTF") or the result of an ioutil.ReadFile call. In such
// cases, as an optimization, we skip the io.Reader / io.ReaderAt model of
// copying from the source to a caller-supplied buffer, and instead provide
// direct access to the underlying []byte data.
type source struct {
	b []byte
	r io.ReaderAt

	// TODO: add a caching layer, if we're using the io.ReaderAt? Note that
	// this might make a source no longer safe to use concurrently.
}

// valid returns whether exactly one of s.b and s.r is nil.
func (s *source) valid() bool {
	return (s.b == nil) != (s.r == nil)
}

// viewBufferWritable returns whether the []byte returned by source.view can be
// written to by the caller, including by passing it to the same method
// (source.view) on other receivers (i.e. different sources).
//
// In other words, it returns whether the source's underlying data is an
// io.ReaderAt, not a []byte.
func (s *source) viewBufferWritable() bool {
	return s.b == nil
}

// view returns the length bytes at the given offset. buf is an optional
// scratch buffer to reduce allocations when calling view multiple times. A nil
// buf is valid. The []byte returned may be a sub-slice of buf[:cap(buf)], or
// it may be an unrelated slice. In any case, the caller should not modify the
// contents of the returned []byte, other than passing that []byte back to this
// method on the same source s.
func (s *source) view(buf []byte, offset, length int) ([]byte, error) {
	if 0 > offset || offset > offset+length {
		return nil, errInvalidBounds
	}

	// Try reading from the []byte.
	if s.b != nil {
		if offset+length > len(s.b) {
			return nil, errInvalidBounds
		}
		return s.b[offset : offset+length], nil
	}

	// Read from the io.ReaderAt.
	if length <= cap(buf) {
		buf = buf[:length]
	} else {
		// Round length up to the nearest KiB. The slack can lead to fewer
		// allocations if the buffer is re-used for multiple source.view calls.
		n := length
		n += 1023
		n &^= 1023
		buf = make([]byte, length, n)
	}
	if n, err := s.r.ReadAt(buf, int64(offset)); n != length {
		return nil, err
	}
	return buf, nil
}

// u16 returns the uint16 in the table t at the relative offset i.
//
// buf is an optional scratch buffer as per the source.view method.
func (s *source) u16(buf []byte, t table, i int) (uint16, error) {
	if i < 0 || uint(t.length) < uint(i+2) {
		return 0, errInvalidBounds
	}
	buf, err := s.view(buf, int(t.offset)+i, 2)
	if err != nil {
		return 0, err
	}
	return u16(buf), nil
}

// u32 returns the uint32 in the table t at the relative offset i.
//
// buf is an optional scratch buffer as per the source.view method.
func (s *source) u32(buf []byte, t table, i int) (uint32, error) {
	if i < 0 || uint(t.length) < uint(i+4) {
		return 0, errInvalidBounds
	}
	buf, err := s.view(buf, int(t.offset)+i, 4)
	if err != nil {
		return 0, err
	}
	return u32(buf), nil
}

// table is a section of the font data.
type table struct {
	offset, length uint32
}

// Parse parses an SFNT font from a []byte data source.
func Parse(src []byte) (*Font, error) {
	f := &Font{src: source{b: src}}
	if err := f.initialize(); err != nil {
		return nil, err
	}
	return f, nil
}

// ParseReaderAt parses an SFNT font from an io.ReaderAt data source.
func ParseReaderAt(src io.ReaderAt) (*Font, error) {
	f := &Font{src: source{r: src}}
	if err := f.initialize(); err != nil {
		return nil, err
	}
	return f, nil
}

// Font is an SFNT font.
//
// Many of its methods take a *Buffer argument, as re-using buffers can reduce
// the total memory allocation of repeated Font method calls, such as measuring
// and rasterizing every unique glyph in a string of text. If efficiency is not
// a concern, passing a nil *Buffer is valid, and implies using a temporary
// buffer for a single call.
//
// It is valid to re-use a *Buffer with multiple Font method calls, even with
// different *Font receivers, as long as they are not concurrent calls.
//
// All of the Font methods are safe to call concurrently, as long as each call
// has a different *Buffer (or nil).
//
// The Font methods that don't take a *Buffer argument are always safe to call
// concurrently.
//
// Some methods provide lengths or coordinates, e.g. bounds, font metrics and
// control points. All of these methods take a ppem parameter, which is the
// number of pixels in 1 em, expressed as a 26.6 fixed point value. For
// example, if 1 em is 10 pixels then ppem is fixed.I(10), which equals
// fixed.Int26_6(10 << 6).
//
// To get those lengths or coordinates in terms of font units instead of
// pixels, use ppem = fixed.Int26_6(f.UnitsPerEm()) and if those methods take a
// font.Hinting parameter, use font.HintingNone. The return values will have
// type fixed.Int26_6, but those numbers can be converted back to Units with no
// further scaling necessary.
type Font struct {
	src source

	// https://www.microsoft.com/typography/otspec/otff.htm#otttables
	// "Required Tables".
	cmap table
	head table
	hhea table
	hmtx table
	maxp table
	name table
	os2  table
	post table

	// https://www.microsoft.com/typography/otspec/otff.htm#otttables
	// "Tables Related to TrueType Outlines".
	//
	// This implementation does not support hinting, so it does not read the
	// cvt, fpgm gasp or prep tables.
	glyf table
	loca table

	// https://www.microsoft.com/typography/otspec/otff.htm#otttables
	// "Tables Related to PostScript Outlines".
	//
	// TODO: cff2, vorg?
	cff table

	// https://www.microsoft.com/typography/otspec/otff.htm#otttables
	// "Advanced Typographic Tables".
	//
	// TODO: base, gdef, gpos, gsub, jstf, math?

	// https://www.microsoft.com/typography/otspec/otff.htm#otttables
	// "Other OpenType Tables".
	//
	// TODO: hdmx, vmtx? Others?
	kern table

	cached struct {
		glyphIndex       glyphIndexFunc
		indexToLocFormat bool // false means short, true means long.
		isPostScript     bool
		kernNumPairs     int32
		kernOffset       int32
		postTableVersion uint32
		unitsPerEm       Units

		// The glyph data for the glyph index i is in
		// src[locations[i+0]:locations[i+1]].
		locations []uint32
	}
}

// NumGlyphs returns the number of glyphs in f.
func (f *Font) NumGlyphs() int { return len(f.cached.locations) - 1 }

// UnitsPerEm returns the number of units per em for f.
func (f *Font) UnitsPerEm() Units { return f.cached.unitsPerEm }

func (f *Font) initialize() error {
	if !f.src.valid() {
		return errInvalidSourceData
	}
	buf, isPostScript, err := f.initializeTables(nil)
	if err != nil {
		return err
	}

	// The order of these parseXxx calls matters. Later calls may depend on
	// information parsed by earlier calls, such as the maxp table's numGlyphs.
	// To enforce these dependencies, such information is passed and returned
	// explicitly, and the f.cached fields are only set afterwards.
	//
	// When implementing new parseXxx methods, take care not to call methods
	// such as Font.NumGlyphs that implicitly depend on f.cached fields.

	buf, indexToLocFormat, unitsPerEm, err := f.parseHead(buf)
	if err != nil {
		return err
	}
	buf, numGlyphs, locations, err := f.parseMaxp(buf, indexToLocFormat, isPostScript)
	if err != nil {
		return err
	}
	buf, glyphIndex, err := f.parseCmap(buf)
	if err != nil {
		return err
	}
	buf, kernNumPairs, kernOffset, err := f.parseKern(buf)
	if err != nil {
		return err
	}
	buf, postTableVersion, err := f.parsePost(buf, numGlyphs)
	if err != nil {
		return err
	}

	f.cached.glyphIndex = glyphIndex
	f.cached.indexToLocFormat = indexToLocFormat
	f.cached.isPostScript = isPostScript
	f.cached.kernNumPairs = kernNumPairs
	f.cached.kernOffset = kernOffset
	f.cached.postTableVersion = postTableVersion
	f.cached.unitsPerEm = unitsPerEm
	f.cached.locations = locations

	return nil
}

func (f *Font) initializeTables(buf []byte) (buf1 []byte, isPostScript bool, err error) {
	// https://www.microsoft.com/typography/otspec/otff.htm "Organization of an
	// OpenType Font" says that "The OpenType font starts with the Offset
	// Table", which is 12 bytes.
	buf, err = f.src.view(buf, 0, 12)
	if err != nil {
		return nil, false, err
	}
	switch u32(buf) {
	default:
		return nil, false, errInvalidVersion
	case 0x00010000:
		// No-op.
	case 0x4f54544f: // "OTTO".
		isPostScript = true
	}
	numTables := int(u16(buf[4:]))
	if numTables > maxNumTables {
		return nil, false, errUnsupportedNumberOfTables
	}

	// "The Offset Table is followed immediately by the Table Record entries...
	// sorted in ascending order by tag", 16 bytes each.
	buf, err = f.src.view(buf, 12, 16*numTables)
	if err != nil {
		return nil, false, err
	}
	for b, first, prevTag := buf, true, uint32(0); len(b) > 0; b = b[16:] {
		tag := u32(b)
		if first {
			first = false
		} else if tag <= prevTag {
			return nil, false, errInvalidTableTagOrder
		}
		prevTag = tag

		o, n := u32(b[8:12]), u32(b[12:16])
		if o > maxTableOffset || n > maxTableLength {
			return nil, false, errUnsupportedTableOffsetLength
		}
		// We ignore the checksums, but "all tables must begin on four byte
		// boundries [sic]".
		if o&3 != 0 {
			return nil, false, errInvalidTableOffset
		}

		// Match the 4-byte tag as a uint32. For example, "OS/2" is 0x4f532f32.
		switch tag {
		case 0x43464620:
			f.cff = table{o, n}
		case 0x4f532f32:
			f.os2 = table{o, n}
		case 0x636d6170:
			f.cmap = table{o, n}
		case 0x676c7966:
			f.glyf = table{o, n}
		case 0x68656164:
			f.head = table{o, n}
		case 0x68686561:
			f.hhea = table{o, n}
		case 0x686d7478:
			f.hmtx = table{o, n}
		case 0x6b65726e:
			f.kern = table{o, n}
		case 0x6c6f6361:
			f.loca = table{o, n}
		case 0x6d617870:
			f.maxp = table{o, n}
		case 0x6e616d65:
			f.name = table{o, n}
		case 0x706f7374:
			f.post = table{o, n}
		}
	}
	return buf, isPostScript, nil
}

func (f *Font) parseCmap(buf []byte) (buf1 []byte, glyphIndex glyphIndexFunc, err error) {
	// https://www.microsoft.com/typography/OTSPEC/cmap.htm

	const headerSize, entrySize = 4, 8
	if f.cmap.length < headerSize {
		return nil, nil, errInvalidCmapTable
	}
	u, err := f.src.u16(buf, f.cmap, 2)
	if err != nil {
		return nil, nil, err
	}
	numSubtables := int(u)
	if f.cmap.length < headerSize+entrySize*uint32(numSubtables) {
		return nil, nil, errInvalidCmapTable
	}

	var (
		bestWidth  int
		bestOffset uint32
		bestLength uint32
		bestFormat uint16
	)

	// Scan all of the subtables, picking the widest supported one. See the
	// platformEncodingWidth comment for more discussion of width.
	for i := 0; i < numSubtables; i++ {
		buf, err = f.src.view(buf, int(f.cmap.offset)+headerSize+entrySize*i, entrySize)
		if err != nil {
			return nil, nil, err
		}
		pid := u16(buf)
		psid := u16(buf[2:])
		width := platformEncodingWidth(pid, psid)
		if width <= bestWidth {
			continue
		}
		offset := u32(buf[4:])

		if offset > f.cmap.length-4 {
			return nil, nil, errInvalidCmapTable
		}
		buf, err = f.src.view(buf, int(f.cmap.offset+offset), 4)
		if err != nil {
			return nil, nil, err
		}
		format := u16(buf)
		if !supportedCmapFormat(format, pid, psid) {
			continue
		}
		length := uint32(u16(buf[2:]))

		bestWidth = width
		bestOffset = offset
		bestLength = length
		bestFormat = format
	}

	if bestWidth == 0 {
		return nil, nil, errUnsupportedCmapEncodings
	}
	return f.makeCachedGlyphIndex(buf, bestOffset, bestLength, bestFormat)
}

func (f *Font) parseHead(buf []byte) (buf1 []byte, indexToLocFormat bool, unitsPerEm Units, err error) {
	// https://www.microsoft.com/typography/otspec/head.htm

	if f.head.length != 54 {
		return nil, false, 0, errInvalidHeadTable
	}
	u, err := f.src.u16(buf, f.head, 18)
	if err != nil {
		return nil, false, 0, err
	}
	if u == 0 {
		return nil, false, 0, errInvalidHeadTable
	}
	unitsPerEm = Units(u)
	u, err = f.src.u16(buf, f.head, 50)
	if err != nil {
		return nil, false, 0, err
	}
	indexToLocFormat = u != 0
	return buf, indexToLocFormat, unitsPerEm, nil
}

func (f *Font) parseKern(buf []byte) (buf1 []byte, kernNumPairs, kernOffset int32, err error) {
	// https://www.microsoft.com/typography/otspec/kern.htm

	if f.kern.length == 0 {
		return buf, 0, 0, nil
	}
	const headerSize = 4
	if f.kern.length < headerSize {
		return nil, 0, 0, errInvalidKernTable
	}
	buf, err = f.src.view(buf, int(f.kern.offset), headerSize)
	if err != nil {
		return nil, 0, 0, err
	}
	offset := int(f.kern.offset) + headerSize
	length := int(f.kern.length) - headerSize

	switch version := u16(buf); version {
	case 0:
		// TODO: support numTables != 1. Testing that requires finding such a font.
		if numTables := int(u16(buf[2:])); numTables != 1 {
			return nil, 0, 0, errUnsupportedKernTable
		}
		return f.parseKernVersion0(buf, offset, length)
	case 1:
		// TODO: find such a (proprietary?) font, and support it. Both of
		// https://www.microsoft.com/typography/otspec/kern.htm
		// https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6kern.html
		// say that such fonts work on Mac OS but not on Windows.
	}
	return nil, 0, 0, errUnsupportedKernTable
}

func (f *Font) parseKernVersion0(buf []byte, offset, length int) (buf1 []byte, kernNumPairs, kernOffset int32, err error) {
	const headerSize = 6
	if length < headerSize {
		return nil, 0, 0, errInvalidKernTable
	}
	buf, err = f.src.view(buf, offset, headerSize)
	if err != nil {
		return nil, 0, 0, err
	}
	if version := u16(buf); version != 0 {
		return nil, 0, 0, errUnsupportedKernTable
	}
	subtableLength := int(u16(buf[2:]))
	if subtableLength < headerSize || length < subtableLength {
		return nil, 0, 0, errInvalidKernTable
	}
	if coverageBits := buf[5]; coverageBits != 0x01 {
		// We only support horizontal kerning.
		return nil, 0, 0, errUnsupportedKernTable
	}
	offset += headerSize
	length -= headerSize
	subtableLength -= headerSize

	switch format := buf[4]; format {
	case 0:
		return f.parseKernFormat0(buf, offset, subtableLength)
	case 2:
		// TODO: find such a (proprietary?) font, and support it.
	}
	return nil, 0, 0, errUnsupportedKernTable
}

func (f *Font) parseKernFormat0(buf []byte, offset, length int) (buf1 []byte, kernNumPairs, kernOffset int32, err error) {
	const headerSize, entrySize = 8, 6
	if length < headerSize {
		return nil, 0, 0, errInvalidKernTable
	}
	buf, err = f.src.view(buf, offset, headerSize)
	if err != nil {
		return nil, 0, 0, err
	}
	kernNumPairs = int32(u16(buf))
	if length != headerSize+entrySize*int(kernNumPairs) {
		return nil, 0, 0, errInvalidKernTable
	}
	return buf, kernNumPairs, int32(offset) + headerSize, nil
}

func (f *Font) parseMaxp(buf []byte, indexToLocFormat, isPostScript bool) (buf1 []byte, numGlyphs int, locations []uint32, err error) {
	// https://www.microsoft.com/typography/otspec/maxp.htm

	if isPostScript {
		if f.maxp.length != 6 {
			return nil, 0, nil, errInvalidMaxpTable
		}
	} else {
		if f.maxp.length != 32 {
			return nil, 0, nil, errInvalidMaxpTable
		}
	}
	u, err := f.src.u16(buf, f.maxp, 4)
	if err != nil {
		return nil, 0, nil, err
	}
	numGlyphs = int(u)

	if isPostScript {
		p := cffParser{
			src:    &f.src,
			base:   int(f.cff.offset),
			offset: int(f.cff.offset),
			end:    int(f.cff.offset + f.cff.length),
		}
		locations, err = p.parse()
		if err != nil {
			return nil, 0, nil, err
		}
	} else {
		locations, err = parseLoca(&f.src, f.loca, f.glyf.offset, indexToLocFormat, numGlyphs)
		if err != nil {
			return nil, 0, nil, err
		}
	}
	if len(locations) != numGlyphs+1 {
		return nil, 0, nil, errInvalidLocationData
	}

	return buf, numGlyphs, locations, nil
}

func (f *Font) parsePost(buf []byte, numGlyphs int) (buf1 []byte, postTableVersion uint32, err error) {
	// https://www.microsoft.com/typography/otspec/post.htm

	const headerSize = 32
	if f.post.length < headerSize {
		return nil, 0, errInvalidPostTable
	}
	u, err := f.src.u32(buf, f.post, 0)
	if err != nil {
		return nil, 0, err
	}
	switch u {
	case 0x20000:
		if f.post.length < headerSize+2+2*uint32(numGlyphs) {
			return nil, 0, errInvalidPostTable
		}
	case 0x30000:
		// No-op.
	default:
		return nil, 0, errUnsupportedPostTable
	}
	return buf, u, nil
}

// TODO: API for looking up glyph variants?? For example, some fonts may
// provide both slashed and dotted zero glyphs ('0'), or regular and 'old
// style' numerals, and users can direct software to choose a variant.

type glyphIndexFunc func(f *Font, b *Buffer, r rune) (GlyphIndex, error)

// GlyphIndex returns the glyph index for the given rune.
//
// It returns (0, nil) if there is no glyph for r.
// https://www.microsoft.com/typography/OTSPEC/cmap.htm says that "Character
// codes that do not correspond to any glyph in the font should be mapped to
// glyph index 0. The glyph at this location must be a special glyph
// representing a missing character, commonly known as .notdef."
func (f *Font) GlyphIndex(b *Buffer, r rune) (GlyphIndex, error) {
	return f.cached.glyphIndex(f, b, r)
}

func (f *Font) viewGlyphData(b *Buffer, x GlyphIndex) ([]byte, error) {
	xx := int(x)
	if f.NumGlyphs() <= xx {
		return nil, ErrNotFound
	}
	i := f.cached.locations[xx+0]
	j := f.cached.locations[xx+1]
	if j-i > maxGlyphDataLength {
		return nil, errUnsupportedGlyphDataLength
	}
	return b.view(&f.src, int(i), int(j-i))
}

// LoadGlyphOptions are the options to the Font.LoadGlyph method.
type LoadGlyphOptions struct {
	// TODO: transform / hinting.
}

// LoadGlyph returns the vector segments for the x'th glyph. ppem is the number
// of pixels in 1 em.
//
// If b is non-nil, the segments become invalid to use once b is re-used.
//
// It returns ErrNotFound if the glyph index is out of range.
func (f *Font) LoadGlyph(b *Buffer, x GlyphIndex, ppem fixed.Int26_6, opts *LoadGlyphOptions) ([]Segment, error) {
	if b == nil {
		b = &Buffer{}
	}

	buf, err := f.viewGlyphData(b, x)
	if err != nil {
		return nil, err
	}

	b.segments = b.segments[:0]
	if f.cached.isPostScript {
		b.psi.type2Charstrings.initialize(b.segments)
		if err := b.psi.run(psContextType2Charstring, buf); err != nil {
			return nil, err
		}
		b.segments = b.psi.type2Charstrings.segments
	} else {
		segments, err := appendGlyfSegments(b.segments, buf)
		if err != nil {
			return nil, err
		}
		b.segments = segments
	}

	// Scale the segments. If we want to support hinting, we'll have to push
	// the scaling computation into the PostScript / TrueType specific glyph
	// loading code, such as the appendGlyfSegments body, since TrueType
	// hinting bytecode works on the scaled glyph vectors. For now, though,
	// it's simpler to scale as a post-processing step.
	for i := range b.segments {
		s := &b.segments[i]
		for j := range s.Args {
			s.Args[j] = scale(s.Args[j]*ppem, f.cached.unitsPerEm)
		}
	}

	// TODO: look at opts to transform / hint the Buffer.segments.

	return b.segments, nil
}

// GlyphName returns the name of the x'th glyph.
//
// Not every font contains glyph names. If not present, GlyphName will return
// ("", nil).
//
// If present, the glyph name, provided by the font, is assumed to follow the
// Adobe Glyph List Specification:
// https://github.com/adobe-type-tools/agl-specification/blob/master/README.md
//
// This is also known as the "Adobe Glyph Naming convention", the "Adobe
// document [for] Unicode and Glyph Names" or "PostScript glyph names".
//
// It returns ErrNotFound if the glyph index is out of range.
func (f *Font) GlyphName(b *Buffer, x GlyphIndex) (string, error) {
	if int(x) >= f.NumGlyphs() {
		return "", ErrNotFound
	}
	if f.cached.postTableVersion != 0x20000 {
		return "", nil
	}
	if b == nil {
		b = &Buffer{}
	}

	// The wire format for a Version 2 post table is documented at:
	// https://www.microsoft.com/typography/otspec/post.htm
	const glyphNameIndexOffset = 34

	buf, err := b.view(&f.src, int(f.post.offset)+glyphNameIndexOffset+2*int(x), 2)
	if err != nil {
		return "", err
	}
	u := u16(buf)
	if u < numBuiltInPostNames {
		i := builtInPostNamesOffsets[u+0]
		j := builtInPostNamesOffsets[u+1]
		return builtInPostNamesData[i:j], nil
	}
	// https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6post.html
	// says that "32768 through 65535 are reserved for future use".
	if u > 32767 {
		return "", errUnsupportedPostTable
	}
	u -= numBuiltInPostNames

	// Iterate through the list of Pascal-formatted strings. A linear scan is
	// clearly O(u), which isn't great (as the obvious loop, calling
	// Font.GlyphName, to get all of the glyph names in a font has quadratic
	// complexity), but the wire format doesn't suggest a better alternative.

	offset := glyphNameIndexOffset + 2*f.NumGlyphs()
	buf, err = b.view(&f.src, int(f.post.offset)+offset, int(f.post.length)-offset)
	if err != nil {
		return "", err
	}

	for {
		if len(buf) == 0 {
			return "", errInvalidPostTable
		}
		n := 1 + int(buf[0])
		if len(buf) < n {
			return "", errInvalidPostTable
		}
		if u == 0 {
			return string(buf[1:n]), nil
		}
		buf = buf[n:]
		u--
	}
}

// Kern returns the horizontal adjustment for the kerning pair (x0, x1). A
// positive kern means to move the glyphs further apart. ppem is the number of
// pixels in 1 em.
//
// It returns ErrNotFound if either glyph index is out of range.
func (f *Font) Kern(b *Buffer, x0, x1 GlyphIndex, ppem fixed.Int26_6, h font.Hinting) (fixed.Int26_6, error) {
	// TODO: how should this work with the GPOS table and CFF fonts?
	// https://www.microsoft.com/typography/otspec/kern.htm says that
	// "OpenType™ fonts containing CFF outlines are not supported by the 'kern'
	// table and must use the 'GPOS' OpenType Layout table."

	if n := f.NumGlyphs(); int(x0) >= n || int(x1) >= n {
		return 0, ErrNotFound
	}
	// Not every font has a kern table. If it doesn't, there's no need to
	// allocate a Buffer.
	if f.kern.length == 0 {
		return 0, nil
	}
	if b == nil {
		b = &Buffer{}
	}

	key := uint32(x0)<<16 | uint32(x1)
	lo, hi := int32(0), f.cached.kernNumPairs
	for lo < hi {
		i := (lo + hi) / 2

		// TODO: this view call inside the inner loop can lead to many small
		// reads instead of fewer larger reads, which can be expensive. We
		// should be able to do better, although we don't want to make (one)
		// arbitrarily large read. Perhaps we should round up reads to 4K or 8K
		// chunks. For reference, Arial.ttf's kern table is 5472 bytes.
		// Times_New_Roman.ttf's kern table is 5220 bytes.
		const entrySize = 6
		buf, err := b.view(&f.src, int(f.cached.kernOffset+i*entrySize), entrySize)
		if err != nil {
			return 0, err
		}

		k := u32(buf)
		if k < key {
			lo = i + 1
		} else if k > key {
			hi = i
		} else {
			kern := fixed.Int26_6(int16(u16(buf[4:])))
			kern = scale(kern*ppem, f.cached.unitsPerEm)
			if h == font.HintingFull {
				// Quantize the fixed.Int26_6 value to the nearest pixel.
				kern = (kern + 32) &^ 63
			}
			return kern, nil
		}
	}
	return 0, nil
}

// Name returns the name value keyed by the given NameID.
//
// It returns ErrNotFound if there is no value for that key.
func (f *Font) Name(b *Buffer, id NameID) (string, error) {
	if b == nil {
		b = &Buffer{}
	}

	const headerSize, entrySize = 6, 12
	if f.name.length < headerSize {
		return "", errInvalidNameTable
	}
	buf, err := b.view(&f.src, int(f.name.offset), headerSize)
	if err != nil {
		return "", err
	}
	numSubtables := u16(buf[2:])
	if f.name.length < headerSize+entrySize*uint32(numSubtables) {
		return "", errInvalidNameTable
	}
	stringOffset := u16(buf[4:])

	seen := false
	for i, n := 0, int(numSubtables); i < n; i++ {
		buf, err := b.view(&f.src, int(f.name.offset)+headerSize+entrySize*i, entrySize)
		if err != nil {
			return "", err
		}
		if u16(buf[6:]) != uint16(id) {
			continue
		}
		seen = true

		var stringify func([]byte) (string, error)
		switch u32(buf) {
		default:
			continue
		case pidMacintosh<<16 | psidMacintoshRoman:
			stringify = stringifyMacintosh
		case pidWindows<<16 | psidWindowsUCS2:
			stringify = stringifyUCS2
		}

		nameLength := u16(buf[8:])
		nameOffset := u16(buf[10:])
		buf, err = b.view(&f.src, int(f.name.offset)+int(nameOffset)+int(stringOffset), int(nameLength))
		if err != nil {
			return "", err
		}
		return stringify(buf)
	}

	if seen {
		return "", errUnsupportedPlatformEncoding
	}
	return "", ErrNotFound
}

func stringifyMacintosh(b []byte) (string, error) {
	for _, c := range b {
		if c >= 0x80 {
			// b contains some non-ASCII bytes.
			s, _ := charmap.Macintosh.NewDecoder().Bytes(b)
			return string(s), nil
		}
	}
	// b contains only ASCII bytes.
	return string(b), nil
}

func stringifyUCS2(b []byte) (string, error) {
	if len(b)&1 != 0 {
		return "", errInvalidUCS2String
	}
	r := make([]rune, len(b)/2)
	for i := range r {
		r[i] = rune(u16(b))
		b = b[2:]
	}
	return string(r), nil
}

// Buffer holds re-usable buffers that can reduce the total memory allocation
// of repeated Font method calls.
//
// See the Font type's documentation comment for more details.
type Buffer struct {
	// buf is a byte buffer for when a Font's source is an io.ReaderAt.
	buf []byte
	// segments holds glyph vector path segments.
	segments []Segment
	// psi is a PostScript interpreter for when the Font is an OpenType/CFF
	// font.
	psi psInterpreter
}

func (b *Buffer) view(src *source, offset, length int) ([]byte, error) {
	buf, err := src.view(b.buf, offset, length)
	if err != nil {
		return nil, err
	}
	// Only update b.buf if it is safe to re-use buf.
	if src.viewBufferWritable() {
		b.buf = buf
	}
	return buf, nil
}

// Segment is a segment of a vector path.
type Segment struct {
	Op   SegmentOp
	Args [6]fixed.Int26_6
}

// SegmentOp is a vector path segment's operator.
type SegmentOp uint32

const (
	SegmentOpMoveTo SegmentOp = iota
	SegmentOpLineTo
	SegmentOpQuadTo
	SegmentOpCubeTo
)
