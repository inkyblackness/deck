package data

import "github.com/inkyblackness/res"

// LevelObjectClassMeta contains meta information about a level object class in an archive.
type LevelObjectClassMeta struct {
	EntrySize  int
	EntryCount int
}

var levelObjectClassMetaList = []LevelObjectClassMeta{
	{LevelObjectPrefixSize + 2, 16},
	{LevelObjectPrefixSize + 0, 32},
	{LevelObjectPrefixSize + 34, 32},
	{LevelObjectPrefixSize + 6, 32},
	{LevelObjectPrefixSize + 0, 32},
	{LevelObjectPrefixSize + 1, 8},
	{LevelObjectPrefixSize + 3, 16},
	{LevelObjectPrefixSize + 10, 176},
	{LevelObjectPrefixSize + 10, 128},
	{LevelObjectPrefixSize + 24, 64},
	{LevelObjectPrefixSize + 8, 64},
	{LevelObjectPrefixSize + 4, 32},
	{LevelObjectPrefixSize + 22, 160},
	{LevelObjectPrefixSize + 15, 64},
	{LevelObjectPrefixSize + 40, 64}}

// LevelObjectClassMetaEntry returns the meta entry for the corresponding object class.
func LevelObjectClassMetaEntry(class res.ObjectClass) LevelObjectClassMeta {
	return levelObjectClassMetaList[int(class)]
}
