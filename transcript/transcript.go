// Package transcript provides types and loading for cisshgo transcript maps and scenarios.
package transcript

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

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

// Platform defines a simulated device platform and its command transcripts.
type Platform struct {
	Vendor             string            `yaml:"vendor" json:"vendor"`
	Hostname           string            `yaml:"hostname" json:"hostname"`
	Username           string            `yaml:"username" json:"username"`
	Password           string            `yaml:"password" json:"password"`
	CommandTranscripts map[string]string `yaml:"command_transcripts" json:"command_transcripts"`
	ContextSearch      map[string]string `yaml:"context_search" json:"context_search"`
	ContextHierarchy   map[string]string `yaml:"context_hierarchy" json:"context_hierarchy"`
}

// Map holds all platforms and scenarios defined in a transcript map YAML file.
type Map struct {
	Platforms map[string]Platform `yaml:"platforms" json:"platforms"`
	Scenarios map[string]Scenario `yaml:"scenarios" json:"scenarios"`
}

// Load reads and parses a transcript map YAML file.
func Load(path string) (Map, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Map{}, fmt.Errorf("reading transcript map: %w", err)
	}
	var tm Map
	if err = yaml.UnmarshalStrict(raw, &tm); err != nil {
		return Map{}, fmt.Errorf("parsing transcript map: %w", err)
	}
	return tm, nil
}

// Validate checks that all transcript file paths in the map exist on disk.
// baseDir is prepended to relative paths.
// Returns an error listing all missing files and invalid references.
func Validate(tm Map, baseDir string) error {
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
