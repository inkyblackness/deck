package core

import (
	"strings"
)

// DataLocation specifies where data files are stored
type DataLocation string

func (location DataLocation) String() string {
	return string(location)
}

const (
	// HD is the data location for often-accessed files
	HD = DataLocation("hd")
	// CD is the data location for remaining files
	CD = DataLocation("cd")
)

// FileDesc is a description of a specific data file
type FileDesc struct {
	name string
}

func newFileDesc(name string) FileDesc {
	return FileDesc{name: name}
}

// ReleaseDesc is a description of a certain release
type ReleaseDesc struct {
	name string

	dataFiles map[DataLocation][]FileDesc
}

var dosHdDemo = ReleaseDesc{
	name: "DOS HD Demo",
	dataFiles: map[DataLocation][]FileDesc{
		HD: []FileDesc{
			newFileDesc("archive.dat"),
			newFileDesc("citmat.res"),
			newFileDesc("cybstrng.res"),
			newFileDesc("digifx.res"),
			newFileDesc("frnstrng.res"),
			newFileDesc("gamepal.res"),
			newFileDesc("gamescr.res"),
			newFileDesc("gerstrng.res"),
			newFileDesc("handart.res"),
			newFileDesc("intro.res"),
			newFileDesc("mfdart.res"),
			newFileDesc("obj3d.res"),
			newFileDesc("objart.res"),
			newFileDesc("objart2.res"),
			newFileDesc("objart3.res"),
			newFileDesc("objprop.dat"),
			newFileDesc("sideart.res"),
			newFileDesc("splash.res"),
			newFileDesc("textprop.dat"),
			newFileDesc("texture.res")}}}

var dosCdDemo = ReleaseDesc{
	name: "DOS CD Demo",
	dataFiles: map[DataLocation][]FileDesc{
		HD: []FileDesc{
			newFileDesc("citmat.res"),
			newFileDesc("cybstrng.res"),
			newFileDesc("digifx.res"),
			newFileDesc("frnstrng.res"),
			newFileDesc("gamescr.res"),
			newFileDesc("gerstrng.res"),
			newFileDesc("handart.res"),
			newFileDesc("mfdart.res"),
			newFileDesc("obj3d.res"),
			newFileDesc("objart2.res"),
			newFileDesc("objart3.res"),
			newFileDesc("sideart.res"),
			newFileDesc("texture.res")},
		CD: []FileDesc{
			newFileDesc("archive.dat"),
			newFileDesc("citalog.res"),
			newFileDesc("citbark.res"),
			newFileDesc("gamepal.res"),
			newFileDesc("intro.res"),
			newFileDesc("objart.res"),
			newFileDesc("objprop.dat"),
			newFileDesc("splash.res"),
			newFileDesc("textprop.dat")}}}

var dosHdRelease = ReleaseDesc{
	name: "DOS HD Release",
	dataFiles: map[DataLocation][]FileDesc{
		HD: []FileDesc{
			newFileDesc("archive.dat"),
			newFileDesc("citmat.res"),
			newFileDesc("cutspal.res"),
			newFileDesc("cybstrng.res"),
			newFileDesc("death.res"),
			newFileDesc("digifx.res"),
			newFileDesc("frnstrng.res"),
			newFileDesc("gamepal.res"),
			newFileDesc("gamescr.res"),
			newFileDesc("gerstrng.res"),
			newFileDesc("handart.res"),
			newFileDesc("intro.res"),
			newFileDesc("mfdart.res"),
			newFileDesc("mfdfrn.res"),
			newFileDesc("obj3d.res"),
			newFileDesc("mfdger.res"),
			newFileDesc("objart.res"),
			newFileDesc("objart2.res"),
			newFileDesc("objart3.res"),
			newFileDesc("objprop.dat"),
			newFileDesc("sideart.res"),
			newFileDesc("splash.res"),
			newFileDesc("splshpal.res"),
			newFileDesc("start1.res"),
			newFileDesc("textprop.dat"),
			newFileDesc("texture.res"),
			newFileDesc("vidmail.res"),
			newFileDesc("win1.res")}}}

var dosCdRelease = ReleaseDesc{
	name: "DOS CD Release",
	dataFiles: map[DataLocation][]FileDesc{
		HD: []FileDesc{
			newFileDesc("citmat.res"),
			newFileDesc("cybstrng.res"),
			newFileDesc("digifx.res"),
			newFileDesc("frnstrng.res"),
			newFileDesc("gamescr.res"),
			newFileDesc("gerstrng.res"),
			newFileDesc("handart.res"),
			newFileDesc("intro.res"),
			newFileDesc("mfdart.res"),
			newFileDesc("mfdfrn.res"),
			newFileDesc("mfdger.res"),
			newFileDesc("obj3d.res"),
			newFileDesc("objart2.res"),
			newFileDesc("objart3.res"),
			newFileDesc("sideart.res"),
			newFileDesc("texture.res"),
			newFileDesc("objprop.dat")},
		CD: []FileDesc{
			newFileDesc("archive.dat"),
			newFileDesc("citalog.res"),
			newFileDesc("citbark.res"),
			newFileDesc("cutspal.res"),
			newFileDesc("death.res"),
			newFileDesc("frnalog.res"),
			newFileDesc("frnbark.res"),
			newFileDesc("gamepal.res"),
			newFileDesc("geralog.res"),
			newFileDesc("gerbark.res"),
			newFileDesc("intro.res"),
			newFileDesc("lofrintr.res"),
			newFileDesc("logeintr.res"),
			newFileDesc("lowdeth.res"),
			newFileDesc("lowend.res"),
			newFileDesc("lowintr.res"),
			newFileDesc("objart.res"),
			newFileDesc("objprop.dat"),
			newFileDesc("splash.res"),
			newFileDesc("splshpal.res"),
			newFileDesc("start1.res"),
			newFileDesc("svfrintr.res"),
			newFileDesc("svgadeth.res"),
			newFileDesc("svgaend.res"),
			newFileDesc("svgaintr.res"),
			newFileDesc("svgeintr.res"),
			newFileDesc("vidmail.res"),
			newFileDesc("win1.res"),
			newFileDesc("textprop.dat")}}}

// Releases contains all the release descriptions known to Hacker.
var Releases = []*ReleaseDesc{&dosHdDemo, &dosCdDemo, &dosHdRelease, &dosCdRelease}

// DataFiles maps over the provided release description and returns the names
// of the data files per file location
func DataFiles(release *ReleaseDesc) (hdFiles, cdFiles []string) {
	for _, fileDesc := range release.dataFiles[HD] {
		hdFiles = append(hdFiles, fileDesc.name)
	}
	for _, fileDesc := range release.dataFiles[CD] {
		cdFiles = append(cdFiles, fileDesc.name)
	}
	return
}

// FindRelease tries to determine which release the two sets of files represent. Returns nil if none found.
func FindRelease(hdFiles, cdFiles []string) (release *ReleaseDesc) {
	for _, testRelease := range Releases {
		if allFilesExist(hdFiles, testRelease.dataFiles[HD]) && allFilesExist(cdFiles, testRelease.dataFiles[CD]) {
			release = testRelease
		}
	}

	return
}

func allFilesExist(files []string, descriptions []FileDesc) (ok bool) {
	ok = true
	for _, desc := range descriptions {
		found := false
		for _, file := range files {
			if strings.ToLower(file) == desc.name {
				found = true
			}
		}
		if !found {
			ok = false
		}
	}

	return ok
}
