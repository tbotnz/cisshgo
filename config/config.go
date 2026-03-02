// Package config provides CLI argument parsing and inventory loading for cisshgo.
package config

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v2"
)

// CLI defines the command-line interface for cisshgo.
type CLI struct {
	Version       kong.VersionFlag `short:"v" help:"Show version information."`
	Listeners     int              `help:"How many listeners to spawn." default:"50" short:"l" env:"CISSHGO_LISTENERS"`
	StartingPort  int              `help:"Starting port." default:"10000" short:"p" env:"CISSHGO_STARTING_PORT"`
	TranscriptMap string           `help:"Path to transcript map YAML file." default:"transcripts/transcript_map.yaml" short:"t" type:"path" env:"CISSHGO_TRANSCRIPT_MAP"`
	Platform      string           `help:"Platform to use when no inventory is provided." default:"csr1000v" short:"P" env:"CISSHGO_PLATFORM"`
	Inventory     string           `help:"Path to inventory YAML file." optional:"" short:"i" type:"path" env:"CISSHGO_INVENTORY"`
}

// InventoryEntry defines a single platform or scenario and how many listeners to spawn.
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
		if entry.Count < 0 {
			return Inventory{}, fmt.Errorf("inventory entry %d: count must be non-negative, got %d", i, entry.Count)
		}
		if entry.Platform == "" && entry.Scenario == "" {
			return Inventory{}, fmt.Errorf("inventory entry %d: must set either platform or scenario", i)
		}
		if entry.Platform != "" && entry.Scenario != "" {
			return Inventory{}, fmt.Errorf("inventory entry %d: platform and scenario are mutually exclusive", i)
		}
	}
	return inv, nil
}
