// Package fakedevices provides types and initialization for simulated
// network devices used by cisshgo.
package fakedevices

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tbotnz/cisshgo/transcript"
)

// SupportedCommands is a map of the commands a FakeDevice supports and it's corresponding output
type SupportedCommands map[string]string

// FakeDevice Struct for the device we will be simulating
type FakeDevice struct {
	Vendor             string              // Vendor of this fake device
	Platform           string              // Platform of this fake device
	ScenarioName       string              // Scenario name if this device was initialized from a scenario (empty otherwise)
	Hostname           string              // Hostname of the fake device
	DefaultHostname    string              // Default Hostname of the fake device (for resetting)
	Username           string              // Expected SSH username (empty = any username accepted)
	Password           string              // Password of the fake device
	PromptFormat       string              // Optional prompt format string (e.g. "{username}@{hostname}{context}")
	SupportedCommands  SupportedCommands   // What commands this fake device supports
	ContextSearch      map[string]string   // The available CLI prompt/contexts on this fake device
	ContextHierarchy   map[string]string   // The hierarchy of the available contexts
	ContextPrefixLines map[string]string   // Optional prefix lines above the prompt, keyed by context value
	ContextCommands    map[string][]string // Optional per-context command whitelist; nil means all commands allowed
	EndContext         string              // If set, "end" jumps directly to this context (e.g. "#") instead of traversing hierarchy
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
	c.ContextPrefixLines = make(map[string]string, len(fd.ContextPrefixLines))
	for k, v := range fd.ContextPrefixLines {
		c.ContextPrefixLines[k] = v
	}
	if fd.ContextCommands != nil {
		c.ContextCommands = make(map[string][]string, len(fd.ContextCommands))
		for k, v := range fd.ContextCommands {
			c.ContextCommands[k] = append([]string(nil), v...)
		}
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

// InitScenario builds a FakeDevice for a named scenario, loading the base platform
// and returning the pre-loaded sequence steps alongside the device.
func InitScenario(name string, tm transcript.Map, baseDir string) (*FakeDevice, []transcript.SequenceStep, error) {
	s, ok := tm.Scenarios[name]
	if !ok {
		return nil, nil, fmt.Errorf("scenario %q not found in transcript map", name)
	}
	fd, err := InitGeneric(s.Platform, tm, baseDir)
	if err != nil {
		return nil, nil, err
	}
	fd.ScenarioName = name
	steps := make([]transcript.SequenceStep, len(s.Sequence))
	for i, step := range s.Sequence {
		path := step.Transcript
		if !filepath.IsAbs(path) {
			path = filepath.Join(baseDir, path)
		}
		content, err := readFile(path)
		if err != nil {
			return nil, nil, err
		}
		steps[i] = transcript.SequenceStep{Command: step.Command, Transcript: content}
	}
	return fd, steps, nil
}

// InitGeneric builds a FakeDevice struct for use with cisshgo.
// baseDir is the directory from which transcript paths are resolved (typically
// the directory containing the transcript map file).
func InitGeneric(platform string, tm transcript.Map, baseDir string) (*FakeDevice, error) {
	p, ok := tm.Platforms[platform]
	if !ok {
		return nil, fmt.Errorf("platform %q not found in transcript map", platform)
	}

	supportedCommands := make(SupportedCommands, len(p.CommandTranscripts))
	for k, v := range p.CommandTranscripts {
		if !filepath.IsAbs(v) {
			v = filepath.Join(baseDir, v)
		}
		content, err := readFile(v)
		if err != nil {
			return nil, err
		}
		supportedCommands[k] = content
	}

	return &FakeDevice{
		Vendor:             p.Vendor,
		Platform:           platform,
		Hostname:           p.Hostname,
		DefaultHostname:    p.Hostname,
		Username:           p.Username,
		Password:           p.Password,
		PromptFormat:       p.PromptFormat,
		SupportedCommands:  supportedCommands,
		ContextSearch:      p.ContextSearch,
		ContextHierarchy:   p.ContextHierarchy,
		ContextPrefixLines: p.ContextPrefixLines,
		ContextCommands:    p.ContextCommands,
		EndContext:         p.EndContext,
	}, nil
}
