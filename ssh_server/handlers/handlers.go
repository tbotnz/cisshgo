package handlers

import (
	"github.com/gliderlabs/ssh"

	"github.com/tbotnz/cisshgo/fakedevices"
)

// PlatformHandler defines a default type for all platform handlers
type PlatformHandler func(*fakedevices.FakeDevice) ssh.Handler
