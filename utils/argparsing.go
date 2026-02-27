// Package utils provides CLI argument parsing and transcript map loading
// for cisshgo.
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

// CLI defines the command-line interface for cisshgo.
type CLI struct {
	Listeners     int    `help:"How many listeners to spawn." default:"50" short:"l" env:"CISSHGO_LISTENERS"`
	StartingPort  int    `help:"Starting port." default:"10000" short:"p" env:"CISSHGO_STARTING_PORT"`
	TranscriptMap string `help:"Path to transcript map YAML file." default:"transcripts/transcript_map.yaml" short:"t" type:"path" env:"CISSHGO_TRANSCRIPT_MAP"`
	Platform      string `help:"Platform to use when no inventory is provided." default:"csr1000v" short:"P" env:"CISSHGO_PLATFORM"`
	Inventory     string `help:"Path to inventory YAML file." optional:"" short:"i" type:"path" env:"CISSHGO_INVENTORY"`
}

// InventoryEntry defines a single platform or scenario and how many listeners to spawn for it.
// Exactly one of Platform or Scenario must be set.
type InventoryEntry struct {
	Platform string `yaml:"platform"`
	Scenario string `yaml:"scenario"`
	Count    int    `yaml:"count"`
}

// Inventory defines a set of devices to spawn.
type Inventory struct {
	Devices []InventoryEntry `yaml:"devices"`
}

// LoadInventory reads and parses an inventory YAML file.
func LoadInventory(path string) (Inventory, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Inventory{}, fmt.Errorf("reading inventory: %w", err)
	}
	var inv Inventory
	if err = yaml.UnmarshalStrict(raw, &inv); err != nil {
		return Inventory{}, fmt.Errorf("parsing inventory: %w", err)
	}
	for i, entry := range inv.Devices {
		if entry.Platform == "" && entry.Scenario == "" {
			return Inventory{}, fmt.Errorf("inventory entry %d: must set either platform or scenario", i)
		}
		if entry.Platform != "" && entry.Scenario != "" {
			return Inventory{}, fmt.Errorf("inventory entry %d: platform and scenario are mutually exclusive", i)
		}
	}
	return inv, nil
}

// SequenceStep defines a single expected command and its response transcript path.
type SequenceStep struct {
	Command    string `yaml:"command"`
	Transcript string `yaml:"transcript"`
}

// Scenario defines an ordered sequence of command/response pairs layered on top of a platform.
type Scenario struct {
	Platform string         `yaml:"platform"`
	Sequence []SequenceStep `yaml:"sequence"`
}

// TranscriptMapPlatform struct for use inside of a TranscriptMap struct
type TranscriptMapPlatform struct {
	Vendor             string            `yaml:"vendor" json:"vendor"`
	Hostname           string            `yaml:"hostname" json:"hostname"`
	Password           string            `yaml:"password" json:"password"`
	CommandTranscripts map[string]string `yaml:"command_transcripts" json:"command_transcripts"`
	ContextSearch      map[string]string `yaml:"context_search" json:"context_search"`
	ContextHierarchy   map[string]string `yaml:"context_hierarchy" json:"context_hierarchy"`
}

// TranscriptMap Struct for modeling the TranscriptMap YAML
type TranscriptMap struct {
	Platforms map[string]TranscriptMapPlatform `yaml:"platforms" json:"platforms"`
	Scenarios map[string]Scenario              `yaml:"scenarios" json:"scenarios"`
}

// LoadTranscriptMap reads and parses a transcript map YAML file.
func LoadTranscriptMap(path string) (TranscriptMap, error) {
	transcriptMapRaw, err := os.ReadFile(path)
	if err != nil {
		return TranscriptMap{}, fmt.Errorf("reading transcript map: %w", err)
	}

	myTranscriptMap := TranscriptMap{}
	if err = yaml.UnmarshalStrict(transcriptMapRaw, &myTranscriptMap); err != nil {
		return TranscriptMap{}, fmt.Errorf("parsing transcript map: %w", err)
	}

	return myTranscriptMap, nil
}

// ValidateTranscriptMap checks that all transcript file paths in the map exist on disk.
// baseDir is prepended to relative paths (typically filepath.Dir of the transcript map file).
// Returns an error listing all missing files and invalid references, not just the first.
func ValidateTranscriptMap(tm TranscriptMap, baseDir string) error {
	var errs []string

	checkPath := func(label, path string) {
		resolved := path
		if !filepath.IsAbs(path) {
			resolved = filepath.Join(baseDir, path)
		}
		if _, err := os.Stat(resolved); err != nil {
			errs = append(errs, fmt.Sprintf("  %s: %s", label, resolved))
		}
	}

	for platform, p := range tm.Platforms {
		for cmd, path := range p.CommandTranscripts {
			checkPath(fmt.Sprintf("platform %q command %q", platform, cmd), path)
		}
	}

	for name, s := range tm.Scenarios {
		if _, ok := tm.Platforms[s.Platform]; !ok {
			errs = append(errs, fmt.Sprintf("  scenario %q: unknown platform %q", name, s.Platform))
			continue
		}
		for i, step := range s.Sequence {
			checkPath(fmt.Sprintf("scenario %q step %d (%q)", name, i, step.Command), step.Transcript)
		}
	}

	if len(errs) > 0 {
		sort.Strings(errs)
		return fmt.Errorf("transcript map validation failed:\n%s", strings.Join(errs, "\n"))
	}
	return nil
}
