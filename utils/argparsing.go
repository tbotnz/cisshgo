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

// InventoryEntry defines a single platform and how many listeners to spawn for it.
type InventoryEntry struct {
	Platform string `yaml:"platform"`
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
	return inv, nil
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
// Returns an error listing all missing files, not just the first.
func ValidateTranscriptMap(tm TranscriptMap, baseDir string) error {
	var missing []string
	for platform, p := range tm.Platforms {
		for cmd, path := range p.CommandTranscripts {
			resolved := path
			if !filepath.IsAbs(path) {
				resolved = filepath.Join(baseDir, path)
			}
			if _, err := os.Stat(resolved); err != nil {
				missing = append(missing, fmt.Sprintf("  platform %q command %q: %s", platform, cmd, resolved))
			}
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("transcript map validation failed:\n%s", strings.Join(missing, "\n"))
	}
	return nil
}
