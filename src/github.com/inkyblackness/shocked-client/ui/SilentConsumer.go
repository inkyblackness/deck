package ui

import (
	"github.com/inkyblackness/shocked-client/ui/events"
)

// SilentConsumer is an event consumer always returning true.
func SilentConsumer(*Area, events.Event) bool {
	return true
}
