// Package handlers contains SSH session handlers for specific device types
// in order to best emulate their actual behavior.
package handlers

import (
	"strings"

	"github.com/gliderlabs/ssh"

	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/transcript"
)

// PlatformHandler defines a default type for all platform handlers
type PlatformHandler func(*fakedevices.FakeDevice) ssh.Handler

// ScenarioHandler defines a handler type that includes a sequence of steps
type ScenarioHandler func(*fakedevices.FakeDevice, []transcript.SequenceStep) ssh.Handler

// ResolvePlatformHandlers returns platform-appropriate interactive and scenario handlers.
func ResolvePlatformHandlers(platform string) (PlatformHandler, ScenarioHandler) {
	switch strings.ToLower(platform) {
	case "fortios":
		return FortiOSHandler, FortiOSScenarioHandler
	default:
		return GenericCiscoHandler, GenericCiscoScenarioHandler
	}
}
