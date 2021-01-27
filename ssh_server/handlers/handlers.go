package handlers

import "github.com/tbotnz/cisshgo/fakedevices"

// PlatformHandler defines a default type for all platform handlers
type PlatformHandler func(*fakedevices.FakeDevice)
