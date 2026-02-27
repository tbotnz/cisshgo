package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/tbotnz/cisshgo/utils"
)

func validTranscriptMap(t *testing.T) string {
	t.Helper()
	content := `---
platforms:
  csr1000v:
    vendor: "cisco"
    hostname: "testhost"
    password: "admin"
    command_transcripts: {}
    context_search:
      base: ">"
    context_hierarchy:
      ">": "exit"
`
	tmpFile := filepath.Join(t.TempDir(), "transcript_map.yaml")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return tmpFile
}

func TestRun_ZeroListeners(t *testing.T) {
	cli := utils.CLI{Listeners: 0, StartingPort: 10000, Platform: "csr1000v", TranscriptMap: validTranscriptMap(t)}
	if err := run(context.Background(), cli); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_BadTranscriptMap(t *testing.T) {
	cli := utils.CLI{Listeners: 0, StartingPort: 10000, Platform: "csr1000v", TranscriptMap: "/nonexistent/file.yaml"}
	if err := run(context.Background(), cli); err == nil {
		t.Error("expected error for missing transcript map")
	}
}

func TestRun_BadTranscriptContent(t *testing.T) {
	content := `---
platforms:
  csr1000v:
    vendor: "cisco"
    hostname: "testhost"
    password: "admin"
    command_transcripts:
      "show version": "/nonexistent/transcript.txt"
    context_search: {}
    context_hierarchy: {}
`
	tmpFile := filepath.Join(t.TempDir(), "transcript_map.yaml")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	cli := utils.CLI{Listeners: 1, StartingPort: 10000, Platform: "csr1000v", TranscriptMap: tmpFile}
	if err := run(context.Background(), cli); err == nil {
		t.Error("expected error for bad transcript file reference")
	}
}

func TestRun_BadInventory(t *testing.T) {
	cli := utils.CLI{Platform: "csr1000v", TranscriptMap: validTranscriptMap(t), Inventory: "/nonexistent/inventory.yaml"}
	if err := run(context.Background(), cli); err == nil {
		t.Error("expected error for missing inventory file")
	}
}

func TestRun_InventoryBadPlatform(t *testing.T) {
	inv := `---
devices:
  - platform: nonexistent
    count: 1
`
	invFile := filepath.Join(t.TempDir(), "inventory.yaml")
	if err := os.WriteFile(invFile, []byte(inv), 0644); err != nil {
		t.Fatal(err)
	}
	cli := utils.CLI{Platform: "csr1000v", TranscriptMap: validTranscriptMap(t), Inventory: invFile}
	if err := run(context.Background(), cli); err == nil {
		t.Error("expected error for nonexistent platform in inventory")
	}
}
