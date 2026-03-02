package transcript

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestMapParsing(t *testing.T) {
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
	var tm Map
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

func TestMapParsing_MultiplePlatforms(t *testing.T) {
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
	var tm Map
	if err := yaml.UnmarshalStrict([]byte(raw), &tm); err != nil {
		t.Fatalf("UnmarshalStrict error: %v", err)
	}
	if len(tm.Platforms) != 2 {
		t.Fatalf("Platforms len = %d, want 2", len(tm.Platforms))
	}
}

func TestMapParsing_InvalidYAML(t *testing.T) {
	var tm Map
	if err := yaml.UnmarshalStrict([]byte(`not: valid: yaml: [[[`), &tm); err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func writeTestMap(t *testing.T) string {
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

func TestLoad(t *testing.T) {
	tm, err := Load(writeTestMap(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(tm.Platforms) != 1 {
		t.Fatalf("Platforms len = %d, want 1", len(tm.Platforms))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/file.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "bad.yaml")
	os.WriteFile(tmpFile, []byte(`not: valid: yaml: [[[`), 0644)
	_, err := Load(tmpFile)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestValidate(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "show_version.txt")
	if err := os.WriteFile(f, []byte("output"), 0644); err != nil {
		t.Fatal(err)
	}
	tm := Map{
		Platforms: map[string]Platform{
			"csr1000v": {Username: "admin", CommandTranscripts: map[string]string{"show version": "show_version.txt"}},
		},
	}
	if err := Validate(tm, dir); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_MissingUsername(t *testing.T) {
	tm := Map{
		Platforms: map[string]Platform{
			"csr1000v": {CommandTranscripts: map[string]string{}},
		},
	}
	if err := Validate(tm, t.TempDir()); err == nil {
		t.Error("expected error for platform with empty username")
	}
}

func TestValidate_MissingFile(t *testing.T) {
	tm := Map{
		Platforms: map[string]Platform{
			"csr1000v": {Username: "admin", CommandTranscripts: map[string]string{"show version": "nonexistent.txt"}},
		},
	}
	if err := Validate(tm, t.TempDir()); err == nil {
		t.Error("expected error for missing transcript file")
	}
}

func TestValidate_UnknownScenarioPlatform(t *testing.T) {
	tm := Map{
		Platforms: map[string]Platform{},
		Scenarios: map[string]Scenario{"test": {Platform: "nonexistent"}},
	}
	if err := Validate(tm, "."); err == nil {
		t.Error("expected error for scenario referencing unknown platform")
	}
}

func TestValidate_ScenarioMissingFile(t *testing.T) {
	tm := Map{
		Platforms: map[string]Platform{"csr1000v": {Username: "admin"}},
		Scenarios: map[string]Scenario{
			"test": {
				Platform: "csr1000v",
				Sequence: []SequenceStep{{Command: "show version", Transcript: "nonexistent.txt"}},
			},
		},
	}
	if err := Validate(tm, t.TempDir()); err == nil {
		t.Error("expected error for missing scenario transcript file")
	}
}
