// Package fakedevices provides types and initialization for simulated
// network devices used by cisshgo.
package fakedevices

import (
	"fmt"
	"os"

	"github.com/tbotnz/cisshgo/utils"
)

// SupportedCommands is a map of the commands a FakeDevice supports and it's corresponding output
type SupportedCommands map[string]string

// FakeDevice Struct for the device we will be simulating
type FakeDevice struct {
	Vendor            string            // Vendor of this fake device
	Platform          string            // Platform of this fake device
	Hostname          string            // Hostname of the fake device
	DefaultHostname   string            // Default Hostname of the fake device (for resetting)
	Password          string            // Password of the fake device
	SupportedCommands SupportedCommands // What commands this fake device supports
	ContextSearch     map[string]string // The available CLI prompt/contexts on this fake device
	ContextHierarchy  map[string]string // The hierarchy of the available contexts
}

// Copy returns a deep copy of the FakeDevice, safe for use in a separate goroutine.
func (fd *FakeDevice) Copy() *FakeDevice {
	c := *fd
	c.SupportedCommands = make(SupportedCommands, len(fd.SupportedCommands))
	for k, v := range fd.SupportedCommands {
		c.SupportedCommands[k] = v
	}
	c.ContextSearch = make(map[string]string, len(fd.ContextSearch))
	for k, v := range fd.ContextSearch {
		c.ContextSearch[k] = v
	}
	c.ContextHierarchy = make(map[string]string, len(fd.ContextHierarchy))
	for k, v := range fd.ContextHierarchy {
		c.ContextHierarchy[k] = v
	}
	return &c
}
func readFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", filename, err)
	}
	return string(content), nil
}

// InitGeneric builds a FakeDevice struct for use with cisshgo
func InitGeneric(
	vendor string,
	platform string,
	myTranscriptMap utils.TranscriptMap,
) (*FakeDevice, error) {

	p := myTranscriptMap.Platforms[platform]

	supportedCommands := make(SupportedCommands, len(p.CommandTranscripts))
	for k, v := range p.CommandTranscripts {
		content, err := readFile(v)
		if err != nil {
			return nil, err
		}
		supportedCommands[k] = content
	}

	// Create our fake device and return it
	myFakeDevice := FakeDevice{
		Vendor:            vendor,
		Platform:          platform,
		Hostname:          p.Hostname,
		DefaultHostname:   p.Hostname,
		Password:          p.Password,
		SupportedCommands: supportedCommands,
		ContextSearch:     p.ContextSearch,
		ContextHierarchy:  p.ContextHierarchy,
	}

	return &myFakeDevice, nil
}
