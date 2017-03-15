package main

import (
	"log"

	"github.com/docopt/docopt-go"

	"github.com/inkyblackness/shocked-client/editor"
	"github.com/inkyblackness/shocked-client/env/native"
	"github.com/inkyblackness/shocked-core"
	"github.com/inkyblackness/shocked-core/release"
)

func usage() string {
	return Title + `

Usage:
   shocked-client-console --path=<datadir>...
   shocked-client-console -h | --help
   shocked-client-console --version

Options:
   -h --help             Show this screen.
   --version             Show version.
   --path=<datadir>      A path to data directory for inplace modifications. Repeat option for multiple directories.
`
}

func main() {
	arguments, _ := docopt.Parse(usage(), nil, true, Title, false)

	pathArg := arguments["--path"]

	source, srcErr := release.FromAbsolutePaths(pathArg.([]string))
	if srcErr != nil {
		log.Fatalf("Source is not available: %v", srcErr)
		return
	}

	deferrer := make(chan func(), 100)
	defer close(deferrer)

	store := core.NewInplaceDataStore(source, deferrer)
	app := editor.NewMainApplication(store)

	native.Run(app, deferrer)
}
