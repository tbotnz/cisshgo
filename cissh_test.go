package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

func validTranscriptMap(t *testing.T) string {
	t.Helper()
	content := `---
platforms:
  - csr1000v:
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

func resetFlags(t *testing.T) {
	t.Helper()
	oldArgs := os.Args
	oldFlags := flag.CommandLine
	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldFlags
	})
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func TestRun_ZeroListeners(t *testing.T) {
	resetFlags(t)
	os.Args = []string{"cisshgo", "-listeners", "0", "-transcriptMap", validTranscriptMap(t)}
	if err := run(context.Background()); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_BadTranscriptMap(t *testing.T) {
	resetFlags(t)
	os.Args = []string{"cisshgo", "-transcriptMap", "/nonexistent/file.yaml"}
	if err := run(context.Background()); err == nil {
		t.Error("expected error for missing transcript map")
	}
}

func TestRun_BadTranscriptContent(t *testing.T) {
	content := `---
platforms:
  - csr1000v:
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

	resetFlags(t)
	os.Args = []string{"cisshgo", "-transcriptMap", tmpFile}
	if err := run(context.Background()); err == nil {
		t.Error("expected error for bad transcript file reference")
	}
}
