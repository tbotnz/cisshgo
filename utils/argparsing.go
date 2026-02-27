// Package utils provides CLI argument parsing and transcript map loading
// for cisshgo.
package utils

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

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

// ParseArgs parses command line arguments for cisshgo and returns the configuration.
func ParseArgs() (int, int, TranscriptMap, error) {
	listenersPtr := flag.Int("listeners", 50, "How many listeners do you wish to spawn?")
	startingPortPtr := flag.Int("startingPort", 10000, "What port do you want to start at?")
	transcriptMapPtr := flag.String(
		"transcriptMap",
		"transcripts/transcript_map.yaml",
		"What file contains the map of commands to transcribed output?",
	)
	flag.Parse()

	myTranscriptMap, err := LoadTranscriptMap(*transcriptMapPtr)
	if err != nil {
		return 0, 0, TranscriptMap{}, err
	}

	return *startingPortPtr + *listenersPtr, *startingPortPtr, myTranscriptMap, nil
}
