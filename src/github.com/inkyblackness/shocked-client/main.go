package main

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"fmt"
	"log"
	"os"

	"github.com/docopt/docopt-go"

	"github.com/inkyblackness/shocked-client/editor"
	"github.com/inkyblackness/shocked-client/env/native"
	"github.com/inkyblackness/shocked-core"
	"github.com/inkyblackness/shocked-core/release"
)

func usage() string {
	return Title + `

Usage:
   shocked-client --path=<datadir>... [--autosave=<sec>] [--scale=<scale>] [--invertedSliderScroll]
   shocked-client -h | --help
   shocked-client --version

Options:
   -h --help               Show this screen.
   --version               Show version.
   --path=<datadir>        A path to data directory for inplace modifications. Repeat option for multiple directories.
   --autosave=<sec>        A duration, in seconds (1..1800), after which changed files are automatically saved. Default: 5.
   --scale=<scale>         A factor for scaling the UI (0.5 .. 1.0). 1080p displays should use default. 4K most likely 2.0. Default: 1.0.
   --invertedSliderScroll  Specify to have sliders go "down" if scrolling "up" (= old behaviour)
`
}

func main() {
	opts, _ := docopt.ParseArgs(usage(), nil, Title)
	autoSaveTimeoutMSec := 5000
	scale := 1.0
	invertedSliderScroll := false

	autoSaveValue, err := opts.Int("--autosave")
	if err == nil {
		if (autoSaveValue > 0) && (autoSaveValue <= 1800) {
			autoSaveTimeoutMSec = int(autoSaveValue) * 1000
		} else {
			fmt.Fprintf(os.Stderr, "--autosave is supported only between 1 and 1800 -- value ignored: <%v>\n", autoSaveValue)
		}
	}
	scaleValue, err := opts.Float64("--scale")
	if err == nil {
		if (scaleValue >= 0.5) && (scaleValue <= 10.0) {
			scale = scaleValue
		} else {
			fmt.Fprintf(os.Stderr, "--scale is supported only between 0.5 and 10.0 -- value ignored: <%v>\n", scaleValue)
		}
	}
	invertedSliderScrollArg, err := opts.Bool("--invertedSliderScroll")
	if err == nil {
		invertedSliderScroll = invertedSliderScrollArg
	}
	pathArg := opts["--path"]

	source, srcErr := release.FromAbsolutePaths(pathArg.([]string))
	if srcErr != nil {
		log.Fatalf("Source is not available: %v", srcErr)
		return
	}

	deferrer := make(chan func(), 100)
	defer close(deferrer)

	store := core.NewInplaceDataStore(source, deferrer, autoSaveTimeoutMSec)
	app := editor.NewMainApplication(store, float32(scale), invertedSliderScroll)

	native.Run(app, deferrer)
}
