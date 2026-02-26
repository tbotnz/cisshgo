package utils

import (
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
	p, ok := tm.Platforms[0]["csr1000v"]
	if !ok {
		t.Fatal("missing csr1000v platform")
	}
	if p.Hostname != "testhost" {
		t.Errorf("Hostname = %q, want %q", p.Hostname, "testhost")
	}
	if p.Password != "admin" {
		t.Errorf("Password = %q, want %q", p.Password, "admin")
	}
	if p.Vendor != "cisco" {
		t.Errorf("Vendor = %q, want %q", p.Vendor, "cisco")
	}
	if p.CommandTranscripts["show version"] != "transcripts/cisco/csr1000v/show_version.txt" {
		t.Errorf("unexpected command transcript path")
	}
	if p.ContextSearch["base"] != ">" {
		t.Errorf("ContextSearch[base] = %q, want %q", p.ContextSearch["base"], ">")
	}
	if p.ContextHierarchy["#"] != ">" {
		t.Errorf("ContextHierarchy[#] = %q, want %q", p.ContextHierarchy["#"], ">")
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
	if _, ok := tm.Platforms[1]["asa"]; !ok {
		t.Error("missing asa platform")
	}
}

func TestTranscriptMapParsing_InvalidYAML(t *testing.T) {
	raw := `not: valid: yaml: [[[`
	var tm TranscriptMap
	if err := yaml.UnmarshalStrict([]byte(raw), &tm); err == nil {
		t.Error("expected error for invalid YAML")
	}
}
