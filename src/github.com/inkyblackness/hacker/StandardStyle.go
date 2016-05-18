package main

import (
	"github.com/inkyblackness/hacker/styling"

	"github.com/fatih/color"
)

type standardStyle struct {
	prompt  styling.StyleFunc
	err     styling.StyleFunc
	status  styling.StyleFunc
	added   styling.StyleFunc
	removed styling.StyleFunc
}

func newStandardStyle() *standardStyle {
	style := &standardStyle{
		prompt:  color.New(color.FgGreen).SprintFunc(),
		err:     color.New(color.FgRed, color.Bold).SprintFunc(),
		status:  color.New(color.FgCyan).SprintFunc(),
		added:   color.New(color.FgMagenta, color.Bold).SprintFunc(),
		removed: color.New(color.FgCyan, color.Bold).SprintFunc()}

	return style
}

func (style *standardStyle) Prompt() styling.StyleFunc {
	return style.prompt
}

func (style *standardStyle) Error() styling.StyleFunc {
	return style.err
}

func (style *standardStyle) Status() styling.StyleFunc {
	return style.status
}

func (style *standardStyle) Added() styling.StyleFunc {
	return style.added
}

func (style *standardStyle) Removed() styling.StyleFunc {
	return style.removed
}
