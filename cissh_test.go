package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	gossh "golang.org/x/crypto/ssh"

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

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	return port
}

func waitReady(t *testing.T, addr string) {
	t.Helper()
	for i := 0; i < 20; i++ {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("server at %s never became ready", addr)
}

func sshDial(t *testing.T, addr string) {
	t.Helper()
	cfg := &gossh.ClientConfig{
		User:            "admin",
		Auth:            []gossh.AuthMethod{gossh.Password("admin")},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout:         2 * time.Second,
	}
	client, err := gossh.Dial("tcp", addr, cfg)
	if err != nil {
		t.Fatalf("ssh dial: %v", err)
	}
	client.Close()
}

func TestRun_PlatformListener(t *testing.T) {
	port := freePort(t)
	addr := net.JoinHostPort("127.0.0.1", fmt.Sprintf("%d", port))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli := utils.CLI{Listeners: 1, StartingPort: port, Platform: "csr1000v", TranscriptMap: validTranscriptMap(t)}
	done := make(chan error, 1)
	go func() { done <- run(ctx, cli) }()

	waitReady(t, addr)
	sshDial(t, addr)

	cancel()
	if err := <-done; err != nil {
		t.Errorf("run() returned error: %v", err)
	}
}

func TestRun_ScenarioListener(t *testing.T) {
	dir := t.TempDir()

	// Write transcript files
	transcriptFile := filepath.Join(dir, "show_version.txt")
	os.WriteFile(transcriptFile, []byte("version output\n"), 0644)
	seqFile := filepath.Join(dir, "seq_step.txt")
	os.WriteFile(seqFile, []byte("seq output\n"), 0644)

	tmContent := `---
platforms:
  csr1000v:
    vendor: "cisco"
    hostname: "testhost"
    password: "admin"
    command_transcripts:
      "show version": "show_version.txt"
    context_search:
      base: ">"
    context_hierarchy:
      ">": "exit"
scenarios:
  test-scenario:
    platform: csr1000v
    sequence:
      - command: "show running-config"
        transcript: "seq_step.txt"
`
	tmFile := filepath.Join(dir, "transcript_map.yaml")
	os.WriteFile(tmFile, []byte(tmContent), 0644)

	invContent := `---
devices:
  - scenario: test-scenario
    count: 1
`
	invFile := filepath.Join(dir, "inventory.yaml")
	os.WriteFile(invFile, []byte(invContent), 0644)

	port := freePort(t)
	addr := net.JoinHostPort("127.0.0.1", fmt.Sprintf("%d", port))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli := utils.CLI{StartingPort: port, Platform: "csr1000v", TranscriptMap: tmFile, Inventory: invFile}
	done := make(chan error, 1)
	go func() { done <- run(ctx, cli) }()

	waitReady(t, addr)
	sshDial(t, addr)

	cancel()
	if err := <-done; err != nil {
		t.Errorf("run() returned error: %v", err)
	}
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

func TestRun_InventoryZeroCounts(t *testing.T) {
	inv := `---
devices:
  - platform: csr1000v
    count: 0
`
	invFile := filepath.Join(t.TempDir(), "inventory.yaml")
	if err := os.WriteFile(invFile, []byte(inv), 0644); err != nil {
		t.Fatal(err)
	}
	cli := utils.CLI{Platform: "csr1000v", TranscriptMap: validTranscriptMap(t), Inventory: invFile}
	if err := run(context.Background(), cli); err != nil {
		t.Errorf("unexpected error: %v", err)
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
