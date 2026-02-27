// Package utils provides CLI argument parsing and transcript map loading
// for cisshgo.
package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// CLI defines the command-line interface for cisshgo.
type CLI struct {
	Listeners     int    `help:"How many listeners to spawn." default:"50" short:"l"`
	StartingPort  int    `help:"Starting port." default:"10000" short:"p"`
	TranscriptMap string `help:"Path to transcript map YAML file." default:"transcripts/transcript_map.yaml" short:"t" type:"path"`
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
