package utils

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestTranscriptMapParsing(t *testing.T) {
	raw := `---
platforms:
  csr1000v:
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
	p := tm.Platforms["csr1000v"]
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
  csr1000v:
    vendor: "cisco"
    hostname: "host1"
    password: "pass1"
    command_transcripts: {}
    context_search: {}
    context_hierarchy: {}
  asa:
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
  csr1000v:
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

func TestLoadInventory(t *testing.T) {
	content := `---
devices:
  - platform: csr1000v
    count: 10
  - platform: iosxr
    count: 5
`
	tmpFile := filepath.Join(t.TempDir(), "inventory.yaml")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	inv, err := LoadInventory(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(inv.Devices) != 2 {
		t.Fatalf("Devices len = %d, want 2", len(inv.Devices))
	}
	if inv.Devices[0].Platform != "csr1000v" || inv.Devices[0].Count != 10 {
		t.Errorf("Devices[0] = %+v, want {csr1000v 10}", inv.Devices[0])
	}
}

func TestLoadInventory_MissingFile(t *testing.T) {
	_, err := LoadInventory("/nonexistent/inventory.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadInventory_InvalidYAML(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "bad.yaml")
	os.WriteFile(tmpFile, []byte(`not: valid: yaml: [[[`), 0644)
	_, err := LoadInventory(tmpFile)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}
