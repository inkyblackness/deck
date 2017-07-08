package core

import (
	model "github.com/inkyblackness/shocked-model"
)

// MaximumLevelsPerArchive is the constant of how many levels are possible.
const MaximumLevelsPerArchive = 16

type localizedFiles struct {
	cybstrng string
	mfdart   string
	citalog  string
	citbark  string
}

var localized = [model.LanguageCount]localizedFiles{
	{
		cybstrng: "cybstrng.res",
		mfdart:   "mfdart.res",
		citalog:  "citalog.res",
		citbark:  "citbark.res"},
	{
		cybstrng: "frnstrng.res",
		mfdart:   "mfdfrn.res",
		citalog:  "frnalog.res",
		citbark:  "frnbark.res"},
	{
		cybstrng: "gerstrng.res",
		mfdart:   "mfdger.res",
		citalog:  "geralog.res",
		citbark:  "gerbark.res"}}
