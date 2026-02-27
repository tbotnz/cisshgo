package handlers

import (
	"github.com/gliderlabs/ssh"

	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/transcript"
)

// PlatformHandler defines a default type for all platform handlers
type PlatformHandler func(*fakedevices.FakeDevice) ssh.Handler

// ScenarioHandler defines a handler type that includes a sequence of steps
type ScenarioHandler func(*fakedevices.FakeDevice, []transcript.SequenceStep) ssh.Handler
