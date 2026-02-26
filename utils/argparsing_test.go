package utils

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestTranscriptMapParsing(t *testing.T) {
	raw := `---
platforms:
  - csr1000v:
      vendor: "cisco"
      hostname: "testhost"
      password: "admin"
      command_transcripts:
        "show version": "transcripts/cisco/csr1000v/show_version.txt"
      context_search:
        "enable": "#"
        "base": ">"
      context_hierarchy:
        "#": ">"
        ">": "exit"
`
	var tm TranscriptMap
	if err := yaml.UnmarshalStrict([]byte(raw), &tm); err != nil {
		t.Fatalf("UnmarshalStrict error: %v", err)
	}
	if len(tm.Platforms) != 1 {
		t.Fatalf("Platforms len = %d, want 1", len(tm.Platforms))
	}
	p := tm.Platforms[0]["csr1000v"]
	if p.Hostname != "testhost" {
		t.Errorf("Hostname = %q, want %q", p.Hostname, "testhost")
	}
	if p.Password != "admin" {
		t.Errorf("Password = %q, want %q", p.Password, "admin")
	}
	if p.Vendor != "cisco" {
		t.Errorf("Vendor = %q, want %q", p.Vendor, "cisco")
	}
}

func TestTranscriptMapParsing_MultiplePlatforms(t *testing.T) {
	raw := `---
platforms:
  - csr1000v:
      vendor: "cisco"
      hostname: "host1"
      password: "pass1"
      command_transcripts: {}
      context_search: {}
      context_hierarchy: {}
  - asa:
      vendor: "cisco"
      hostname: "host2"
      password: "pass2"
      command_transcripts: {}
      context_search: {}
      context_hierarchy: {}
`
	var tm TranscriptMap
	if err := yaml.UnmarshalStrict([]byte(raw), &tm); err != nil {
		t.Fatalf("UnmarshalStrict error: %v", err)
	}
	if len(tm.Platforms) != 2 {
		t.Fatalf("Platforms len = %d, want 2", len(tm.Platforms))
	}
}

func TestTranscriptMapParsing_InvalidYAML(t *testing.T) {
	raw := `not: valid: yaml: [[[`
	var tm TranscriptMap
	if err := yaml.UnmarshalStrict([]byte(raw), &tm); err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func writeTestTranscriptMap(t *testing.T) string {
	t.Helper()
	content := `---
platforms:
  - csr1000v:
      vendor: "cisco"
      hostname: "testhost"
      password: "admin"
      command_transcripts: {}
      context_search: {}
      context_hierarchy: {}
`
	tmpFile := filepath.Join(t.TempDir(), "transcript_map.yaml")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return tmpFile
}

func TestLoadTranscriptMap(t *testing.T) {
	tm, err := LoadTranscriptMap(writeTestTranscriptMap(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(tm.Platforms) != 1 {
		t.Fatalf("Platforms len = %d, want 1", len(tm.Platforms))
	}
}

func TestLoadTranscriptMap_MissingFile(t *testing.T) {
	_, err := LoadTranscriptMap("/nonexistent/file.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadTranscriptMap_InvalidYAML(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "bad.yaml")
	os.WriteFile(tmpFile, []byte(`not: valid: yaml: [[[`), 0644)
	_, err := LoadTranscriptMap(tmpFile)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestParseArgs(t *testing.T) {
	tmpFile := writeTestTranscriptMap(t)

	oldArgs := os.Args
	oldFlags := flag.CommandLine
	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldFlags
	})

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cisshgo", "-listeners", "2", "-startingPort", "20000", "-transcriptMap", tmpFile}

	numListeners, startingPort, tm, err := ParseArgs()
	if err != nil {
		t.Fatal(err)
	}
	if numListeners != 20002 {
		t.Errorf("numListeners = %d, want 20002", numListeners)
	}
	if startingPort != 20000 {
		t.Errorf("startingPort = %d, want 20000", startingPort)
	}
	if len(tm.Platforms) != 1 {
		t.Errorf("Platforms len = %d, want 1", len(tm.Platforms))
	}
}

func TestParseArgs_BadFile(t *testing.T) {
	oldArgs := os.Args
	oldFlags := flag.CommandLine
	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldFlags
	})

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cisshgo", "-transcriptMap", "/nonexistent/file.yaml"}

	_, _, _, err := ParseArgs()
	if err == nil {
		t.Error("expected error for missing transcript map")
	}
}
