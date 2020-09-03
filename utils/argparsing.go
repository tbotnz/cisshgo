package utils

import (
	"flag"
	"io/ioutil"
	"log"

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
	Platforms []map[string]TranscriptMapPlatform `yaml:"platforms" json:"platforms"`
}

// ParseArgs parses command line arguments for cisgo-ios
func ParseArgs() (int, *int, TranscriptMap) {
	// Gather command line arguments and parse them
	listnersPtr := flag.Int("listners", 50, "How many listeners do you wish to spawn?")
	startingPortPtr := flag.Int("startingPort", 10000, "What port do you want to start at?")
	transcriptMapPtr := flag.String(
		"transcriptMap",
		"transcripts/transcript_map.yaml",
		"What file contains the map of commands to transcipted output?",
	)
	flag.Parse()

	// How many total listners will we have?
	numListners := *startingPortPtr + *listnersPtr

	// Gather the command transcripts and create a map of vendor/platform/command
	transcriptMapRaw, err := ioutil.ReadFile(*transcriptMapPtr)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// fmt.Printf("Raw Transcript Map from file:\n\n%s\n", transcriptMapRaw)

	myTranscriptMap := TranscriptMap{}
	err = yaml.UnmarshalStrict([]byte(transcriptMapRaw), &myTranscriptMap)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// fmt.Printf("YAML Parsed Transcript Map:\n\n%+v\n", myTranscriptMap)

	return numListners, startingPortPtr, myTranscriptMap
}
