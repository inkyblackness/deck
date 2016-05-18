package main

import (
	"github.com/inkyblackness/shocked-client/editor"
	"github.com/inkyblackness/shocked-client/env/browser"
)

func main() {
	transport := browser.NewRestTransport()
	store := editor.NewRestDataStore(transport)
	app := editor.NewMainApplication(store)

	browser.Run(app)
}
