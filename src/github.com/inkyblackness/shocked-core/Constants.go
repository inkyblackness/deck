package core

import (
	model "github.com/inkyblackness/shocked-model"
)

// MaximumLevelsPerArchive is the constant of how many levels are possible.
const MaximumLevelsPerArchive = 16

type localizedFiles struct {
	cybstrng string
}

var localized = [model.LanguageCount]localizedFiles{
	{
		cybstrng: "cybstrng.res"},
	{
		cybstrng: "frnstrng.res"},
	{
		cybstrng: "gerstrng.res"}}
