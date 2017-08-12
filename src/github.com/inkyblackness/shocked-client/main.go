package main

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"log"
	"strconv"

	"github.com/docopt/docopt-go"

	"github.com/inkyblackness/shocked-client/editor"
	"github.com/inkyblackness/shocked-client/env/native"
	"github.com/inkyblackness/shocked-core"
	"github.com/inkyblackness/shocked-core/release"
)

func usage() string {
	return Title + `

Usage:
   shocked-client --path=<datadir>... [--autosave=<sec>]
   shocked-client -h | --help
   shocked-client --version

Options:
   -h --help             Show this screen.
   --version             Show version.
   --path=<datadir>      A path to data directory for inplace modifications. Repeat option for multiple directories.
   --autosave=<sec>      A duration, in seconds (1..1800), after which changed files are automatically saved. Default: 5.
`
}

func main() {
	arguments, _ := docopt.Parse(usage(), nil, true, Title, false)
	autoSaveTimeoutMSec := 5000

	autoSaveArg := arguments["--autosave"]
	if autoSaveArg != nil {
		autoSaveValue, autoSaveErr := strconv.ParseInt(autoSaveArg.(string), 10, 16)
		if autoSaveErr == nil {
			if (autoSaveValue > 0) && (autoSaveValue <= 1800) {
				autoSaveTimeoutMSec = int(autoSaveValue) * 1000
			}
		}
	}
	pathArg := arguments["--path"]

	source, srcErr := release.FromAbsolutePaths(pathArg.([]string))
	if srcErr != nil {
		log.Fatalf("Source is not available: %v", srcErr)
		return
	}

	deferrer := make(chan func(), 100)
	defer close(deferrer)

	store := core.NewInplaceDataStore(source, deferrer, autoSaveTimeoutMSec)
	app := editor.NewMainApplication(store)

	native.Run(app, deferrer)
}
