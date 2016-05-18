package main

import (
	"github.com/docopt/docopt-go"

	"github.com/inkyblackness/shocked-client/app"
	"github.com/inkyblackness/shocked-client/editor"
	"github.com/inkyblackness/shocked-client/env/native"
	"github.com/inkyblackness/shocked-client/env/native/console"
)

func usage() string {
	return app.Title + `

Usage:
   shocked-client-console [--address=<addr>]
   shocked-client-console -h | --help
   shocked-client-console --version

Options:
   -h --help             Show this screen.
   --version             Show version.
   --address=<addr>      The ip:port combination to connect to. Default: "localhost:8080".
`
}

func main() {
	arguments, _ := docopt.Parse(usage(), nil, true, app.Title, false)
	addressArg := arguments["--address"]
	address := "localhost:8080"

	if addressArg != nil {
		address = addressArg.(string)
	}

	deferrer := make(chan func(), 100)
	defer close(deferrer)

	transport := native.NewRestTransport("http://"+address, deferrer)
	store := editor.NewRestDataStore(transport)
	app := editor.NewMainApplication(store)

	console.Run(app, deferrer)
}
