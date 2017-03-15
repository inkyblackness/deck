// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sfnt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

func moveTo(xa, ya fixed.Int26_6) Segment {
	return Segment{
		Op:   SegmentOpMoveTo,
		Args: [6]fixed.Int26_6{xa, ya},
	}
}

func lineTo(xa, ya fixed.Int26_6) Segment {
	return Segment{
		Op:   SegmentOpLineTo,
		Args: [6]fixed.Int26_6{xa, ya},
	}
}

func quadTo(xa, ya, xb, yb fixed.Int26_6) Segment {
	return Segment{
		Op:   SegmentOpQuadTo,
		Args: [6]fixed.Int26_6{xa, ya, xb, yb},
	}
}

func cubeTo(xa, ya, xb, yb, xc, yc fixed.Int26_6) Segment {
	return Segment{
		Op:   SegmentOpCubeTo,
		Args: [6]fixed.Int26_6{xa, ya, xb, yb, xc, yc},
	}
}

func checkSegmentsEqual(got, want []Segment) error {
	if len(got) != len(want) {
		return fmt.Errorf("got %d elements, want %d\noverall:\ngot  %v\nwant %v",
			len(got), len(want), got, want)
	}
	for i, g := range got {
		if w := want[i]; g != w {
			return fmt.Errorf("element %d:\ngot  %v\nwant %v\noverall:\ngot  %v\nwant %v",
				i, g, w, got, want)
		}
	}
	return nil
}

func TestTrueTypeParse(t *testing.T) {
	f, err := Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	testTrueType(t, f)
}

func TestTrueTypeParseReaderAt(t *testing.T) {
	f, err := ParseReaderAt(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatalf("ParseReaderAt: %v", err)
	}
	testTrueType(t, f)
}

func testTrueType(t *testing.T, f *Font) {
	if got, want := f.UnitsPerEm(), Units(2048); got != want {
		t.Errorf("UnitsPerEm: got %d, want %d", got, want)
	}
	// The exact number of glyphs in goregular.TTF can vary, and future
	// versions may add more glyphs, but https://blog.golang.org/go-fonts says
	// that "The WGL4 character set... [has] more than 650 characters in all.
	if got, want := f.NumGlyphs(), 650; got <= want {
		t.Errorf("NumGlyphs: got %d, want > %d", got, want)
	}
}

func TestGoRegularGlyphIndex(t *testing.T) {
	f, err := Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	testCases := []struct {
		r    rune
		want GlyphIndex
	}{
		// Glyphs that aren't present in Go Regular.
		{'\u001f', 0}, // U+001F <control>
		{'\u0200', 0}, // U+0200 LATIN CAPITAL LETTER A WITH DOUBLE GRAVE
		{'\u2000', 0}, // U+2000 EN QUAD

		// The want values below can be verified by running the ttx tool on
		// Go-Regular.ttf.
		//
		// The actual values are ad hoc, and result from whatever tools the
		// Bigelow & Holmes type foundry used and the order in which they
		// crafted the glyphs. They may change over time as newer versions of
		// the font are released. In practice, though, running this test with
		// coverage analysis suggests that it covers both the zero and non-zero
		// cmapEntry16.offset cases for a format-4 cmap table.

		{'\u0020', 3},   // U+0020 SPACE
		{'\u0021', 4},   // U+0021 EXCLAMATION MARK
		{'\u0022', 5},   // U+0022 QUOTATION MARK
		{'\u0023', 6},   // U+0023 NUMBER SIGN
		{'\u0024', 223}, // U+0024 DOLLAR SIGN
		{'\u0025', 7},   // U+0025 PERCENT SIGN
		{'\u0026', 8},   // U+0026 AMPERSAND
		{'\u0027', 9},   // U+0027 APOSTROPHE

		{'\u03bd', 423}, // U+03BD GREEK SMALL LETTER NU
		{'\u03be', 424}, // U+03BE GREEK SMALL LETTER XI
		{'\u03bf', 438}, // U+03BF GREEK SMALL LETTER OMICRON
		{'\u03c0', 208}, // U+03C0 GREEK SMALL LETTER PI
		{'\u03c1', 425}, // U+03C1 GREEK SMALL LETTER RHO
		{'\u03c2', 426}, // U+03C2 GREEK SMALL LETTER FINAL SIGMA
	}

	var b Buffer
	for _, tc := range testCases {
		got, err := f.GlyphIndex(&b, tc.r)
		if err != nil {
			t.Errorf("r=%q: %v", tc.r, err)
			continue
		}
		if got != tc.want {
			t.Errorf("r=%q: got %d, want %d", tc.r, got, tc.want)
			continue
		}
	}
}

func TestGlyphIndex(t *testing.T) {
	data, err := ioutil.ReadFile(filepath.FromSlash("../testdata/cmapTest.ttf"))
	if err != nil {
		t.Fatal(err)
	}

	for _, format := range []int{-1, 0, 4, 12} {
		testGlyphIndex(t, data, format)
	}
}

func testGlyphIndex(t *testing.T, data []byte, cmapFormat int) {
	if cmapFormat >= 0 {
		originalSupportedCmapFormat := supportedCmapFormat
		defer func() {
			supportedCmapFormat = originalSupportedCmapFormat
		}()
		supportedCmapFormat = func(format, pid, psid uint16) bool {
			return int(format) == cmapFormat && originalSupportedCmapFormat(format, pid, psid)
		}
	}

	f, err := Parse(data)
	if err != nil {
		t.Errorf("cmapFormat=%d: %v", cmapFormat, err)
		return
	}

	testCases := []struct {
		r    rune
		want GlyphIndex
	}{
		// Glyphs that aren't present in cmapTest.ttf.
		{'?', 0},
		{'\ufffd', 0},
		{'\U0001f4a9', 0},

		// For a .TTF file, FontForge maps:
		//	- ".notdef"          to glyph index 0.
		//	- ".null"            to glyph index 1.
		//	- "nonmarkingreturn" to glyph index 2.

		{'/', 0},
		{'0', 3},
		{'1', 4},
		{'2', 5},
		{'3', 0},

		{'@', 0},
		{'A', 6},
		{'B', 7},
		{'C', 0},

		{'`', 0},
		{'a', 8},
		{'b', 0},

		// Of the remaining runes, only U+00FF LATIN SMALL LETTER Y WITH
		// DIAERESIS is in both the Mac Roman encoding and the cmapTest.ttf
		// font file.
		{'\u00fe', 0},
		{'\u00ff', 9},
		{'\u0100', 10},
		{'\u0101', 11},
		{'\u0102', 0},

		{'\u4e2c', 0},
		{'\u4e2d', 12},
		{'\u4e2e', 0},

		{'\U0001f0a0', 0},
		{'\U0001f0a1', 13},
		{'\U0001f0a2', 0},

		{'\U0001f0b0', 0},
		{'\U0001f0b1', 14},
		{'\U0001f0b2', 15},
		{'\U0001f0b3', 0},
	}

	var b Buffer
	for _, tc := range testCases {
		want := tc.want
		switch {
		case cmapFormat == 0 && tc.r > '\u007f' && tc.r != '\u00ff':
			// cmap format 0, with the Macintosh Roman encoding, can only
			// represent a limited set of non-ASCII runes, e.g. U+00FF.
			want = 0
		case cmapFormat == 4 && tc.r > '\uffff':
			// cmap format 4 only supports the Basic Multilingual Plane (BMP).
			want = 0
		}

		got, err := f.GlyphIndex(&b, tc.r)
		if err != nil {
			t.Errorf("cmapFormat=%d, r=%q: %v", cmapFormat, tc.r, err)
			continue
		}
		if got != want {
			t.Errorf("cmapFormat=%d, r=%q: got %d, want %d", cmapFormat, tc.r, got, want)
			continue
		}
	}
}

func TestPostScriptSegments(t *testing.T) {
	// wants' vectors correspond 1-to-1 to what's in the CFFTest.sfd file,
	// although OpenType/CFF and FontForge's SFD have reversed orders.
	// https://fontforge.github.io/validation.html says that "All paths must be
	// drawn in a consistent direction. Clockwise for external paths,
	// anti-clockwise for internal paths. (Actually PostScript requires the
	// exact opposite, but FontForge reverses PostScript contours when it loads
	// them so that everything is consistant internally -- and reverses them
	// again when it saves them, of course)."
	//
	// The .notdef glyph isn't explicitly in the SFD file, but for some unknown
	// reason, FontForge generates it in the OpenType/CFF file.
	wants := [][]Segment{{
		// .notdef
		// - contour #0
		moveTo(50, 0),
		lineTo(450, 0),
		lineTo(450, 533),
		lineTo(50, 533),
		// - contour #1
		moveTo(100, 50),
		lineTo(100, 483),
		lineTo(400, 483),
		lineTo(400, 50),
	}, {
		// zero
		// - contour #0
		moveTo(300, 700),
		cubeTo(380, 700, 420, 580, 420, 500),
		cubeTo(420, 350, 390, 100, 300, 100),
		cubeTo(220, 100, 180, 220, 180, 300),
		cubeTo(180, 450, 210, 700, 300, 700),
		// - contour #1
		moveTo(300, 800),
		cubeTo(200, 800, 100, 580, 100, 400),
		cubeTo(100, 220, 200, 0, 300, 0),
		cubeTo(400, 0, 500, 220, 500, 400),
		cubeTo(500, 580, 400, 800, 300, 800),
	}, {
		// one
		// - contour #0
		moveTo(100, 0),
		lineTo(300, 0),
		lineTo(300, 800),
		lineTo(100, 800),
	}, {
		// Q
		// - contour #0
		moveTo(657, 237),
		lineTo(289, 387),
		lineTo(519, 615),
		// - contour #1
		moveTo(792, 169),
		cubeTo(867, 263, 926, 502, 791, 665),
		cubeTo(645, 840, 380, 831, 228, 673),
		cubeTo(71, 509, 110, 231, 242, 93),
		cubeTo(369, -39, 641, 18, 722, 93),
		lineTo(802, 3),
		lineTo(864, 83),
	}, {
		// uni4E2D
		// - contour #0
		moveTo(141, 520),
		lineTo(137, 356),
		lineTo(245, 400),
		lineTo(331, 26),
		lineTo(355, 414),
		lineTo(463, 434),
		lineTo(453, 620),
		lineTo(341, 592),
		lineTo(331, 758),
		lineTo(243, 752),
		lineTo(235, 562),
		// TODO: explicitly (not implicitly) close these contours?
	}}

	testSegments(t, "CFFTest.otf", wants)
}

func TestTrueTypeSegments(t *testing.T) {
	// wants' vectors correspond 1-to-1 to what's in the glyfTest.sfd file,
	// although FontForge's SFD format stores quadratic Bézier curves as cubics
	// with duplicated off-curve points. quadTo(bx, by, cx, cy) is stored as
	// "bx by bx by cx cy".
	//
	// The .notdef, .null and nonmarkingreturn glyphs aren't explicitly in the
	// SFD file, but for some unknown reason, FontForge generates them in the
	// TrueType file.
	wants := [][]Segment{{
		// .notdef
		// - contour #0
		moveTo(68, 0),
		lineTo(68, 1365),
		lineTo(612, 1365),
		lineTo(612, 0),
		lineTo(68, 0),
		// - contour #1
		moveTo(136, 68),
		lineTo(544, 68),
		lineTo(544, 1297),
		lineTo(136, 1297),
		lineTo(136, 68),
	}, {
	// .null
	// Empty glyph.
	}, {
	// nonmarkingreturn
	// Empty glyph.
	}, {
		// zero
		// - contour #0
		moveTo(614, 1434),
		quadTo(369, 1434, 369, 614),
		quadTo(369, 471, 435, 338),
		quadTo(502, 205, 614, 205),
		quadTo(860, 205, 860, 1024),
		quadTo(860, 1167, 793, 1300),
		quadTo(727, 1434, 614, 1434),
		// - contour #1
		moveTo(614, 1638),
		quadTo(1024, 1638, 1024, 819),
		quadTo(1024, 0, 614, 0),
		quadTo(205, 0, 205, 819),
		quadTo(205, 1638, 614, 1638),
	}, {
		// one
		// - contour #0
		moveTo(205, 0),
		lineTo(205, 1638),
		lineTo(614, 1638),
		lineTo(614, 0),
		lineTo(205, 0),
	}}

	testSegments(t, "glyfTest.ttf", wants)
}

func testSegments(t *testing.T, filename string, wants [][]Segment) {
	data, err := ioutil.ReadFile(filepath.FromSlash("../testdata/" + filename))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	f, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	ppem := fixed.Int26_6(f.UnitsPerEm())

	if ng := f.NumGlyphs(); ng != len(wants) {
		t.Fatalf("NumGlyphs: got %d, want %d", ng, len(wants))
	}
	var b Buffer
	for i, want := range wants {
		got, err := f.LoadGlyph(&b, GlyphIndex(i), ppem, nil)
		if err != nil {
			t.Errorf("i=%d: LoadGlyph: %v", i, err)
			continue
		}
		if err := checkSegmentsEqual(got, want); err != nil {
			t.Errorf("i=%d: %v", i, err)
			continue
		}
	}
	if _, err := f.LoadGlyph(nil, 0xffff, ppem, nil); err != ErrNotFound {
		t.Errorf("LoadGlyph(..., 0xffff, ...):\ngot  %v\nwant %v", err, ErrNotFound)
	}

	name, err := f.Name(nil, NameIDFamily)
	if err != nil {
		t.Errorf("Name: %v", err)
	} else if want := filename[:len(filename)-len(".ttf")]; name != want {
		t.Errorf("Name:\ngot  %q\nwant %q", name, want)
	}
}

func TestPPEM(t *testing.T) {
	data, err := ioutil.ReadFile(filepath.FromSlash("../testdata/glyfTest.ttf"))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	f, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var b Buffer
	x, err := f.GlyphIndex(&b, '1')
	if err != nil {
		t.Fatalf("GlyphIndex: %v", err)
	}
	if x == 0 {
		t.Fatalf("GlyphIndex: no glyph index found for the rune '1'")
	}

	testCases := []struct {
		ppem fixed.Int26_6
		want []Segment
	}{{
		ppem: fixed.Int26_6(12 << 6),
		want: []Segment{
			moveTo(77, 0),
			lineTo(77, 614),
			lineTo(230, 614),
			lineTo(230, 0),
			lineTo(77, 0),
		},
	}, {
		ppem: fixed.Int26_6(2048),
		want: []Segment{
			moveTo(205, 0),
			lineTo(205, 1638),
			lineTo(614, 1638),
			lineTo(614, 0),
			lineTo(205, 0),
		},
	}}

	for i, tc := range testCases {
		got, err := f.LoadGlyph(&b, x, tc.ppem, nil)
		if err != nil {
			t.Errorf("i=%d: LoadGlyph: %v", i, err)
			continue
		}
		if err := checkSegmentsEqual(got, tc.want); err != nil {
			t.Errorf("i=%d: %v", i, err)
			continue
		}
	}
}

func TestGlyphName(t *testing.T) {
	f, err := Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	testCases := []struct {
		r    rune
		want string
	}{
		{'\x00', "NULL"},
		{'!', "exclam"},
		{'A', "A"},
		{'{', "braceleft"},
		{'\u00c4', "Adieresis"}, // U+00C4 LATIN CAPITAL LETTER A WITH DIAERESIS
		{'\u2020', "dagger"},    // U+2020 DAGGER
		{'\u2660', "spade"},     // U+2660 BLACK SPADE SUIT
		{'\uf800', "gopher"},    // U+F800 <Private Use>
		{'\ufffe', ".notdef"},   // Not in the Go Regular font, so GlyphIndex returns (0, nil).
	}

	var b Buffer
	for _, tc := range testCases {
		x, err := f.GlyphIndex(&b, tc.r)
		if err != nil {
			t.Errorf("r=%q: GlyphIndex: %v", tc.r, err)
			continue
		}
		got, err := f.GlyphName(&b, x)
		if err != nil {
			t.Errorf("r=%q: GlyphName: %v", tc.r, err)
			continue
		}
		if got != tc.want {
			t.Errorf("r=%q: got %q, want %q", tc.r, got, tc.want)
			continue
		}
	}
}

func TestBuiltInPostNames(t *testing.T) {
	testCases := []struct {
		x    GlyphIndex
		want string
	}{
		{0, ".notdef"},
		{1, ".null"},
		{2, "nonmarkingreturn"},
		{13, "asterisk"},
		{36, "A"},
		{93, "z"},
		{123, "ocircumflex"},
		{202, "Edieresis"},
		{255, "Ccaron"},
		{256, "ccaron"},
		{257, "dcroat"},
		{258, ""},
		{999, ""},
		{0xffff, ""},
	}

	for _, tc := range testCases {
		if tc.x >= numBuiltInPostNames {
			continue
		}
		i := builtInPostNamesOffsets[tc.x+0]
		j := builtInPostNamesOffsets[tc.x+1]
		got := builtInPostNamesData[i:j]
		if got != tc.want {
			t.Errorf("x=%d: got %q, want %q", tc.x, got, tc.want)
		}
	}
}
