package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

func TestRun_BadTranscriptMap(t *testing.T) {
	oldArgs := os.Args
	oldFlags := flag.CommandLine
	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldFlags
	})

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cisshgo", "-transcriptMap", "/nonexistent/file.yaml"}

	err := run()
	if err == nil {
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

	oldArgs := os.Args
	oldFlags := flag.CommandLine
	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldFlags
	})

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cisshgo", "-transcriptMap", tmpFile}

	err := run()
	if err == nil {
		t.Error("expected error for bad transcript file reference")
	}
}
