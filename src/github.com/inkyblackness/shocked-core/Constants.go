package core

import (
	model "github.com/inkyblackness/shocked-model"
)

// MaximumLevelsPerArchive is the constant of how many levels are possible.
const MaximumLevelsPerArchive = 16

type localizedFiles struct {
	cybstrng string
	mfdart   string
}

var localized = [model.LanguageCount]localizedFiles{
	{
		cybstrng: "cybstrng.res",
		mfdart:   "mfdart.res"},
	{
		cybstrng: "frnstrng.res",
		mfdart:   "mfdfrn.res"},
	{
		cybstrng: "gerstrng.res",
		mfdart:   "mfdger.res"}}
