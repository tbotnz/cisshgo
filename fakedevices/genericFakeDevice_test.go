package fakedevices

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/tbotnz/cisshgo/transcript"
)

func testTranscriptMap() transcript.Map {
	return transcript.Map{
		Platforms: map[string]transcript.Platform{
			"csr1000v": {
				Vendor:   "cisco",
				Hostname: "testhost",
				Password: "secret",
				CommandTranscripts: map[string]string{
					"show version": "transcripts/cisco/csr1000v/show_version.txt",
				},
				ContextSearch: map[string]string{
					"base":               ">",
					"enable":             "#",
					"configure terminal": "(config)#",
				},
				ContextHierarchy: map[string]string{
					">":         "exit",
					"#":         ">",
					"(config)#": "#",
				},
			},
		},
	}
}

func TestInitGeneric(t *testing.T) {
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir("fakedevices") })

	fd, err := InitGeneric("csr1000v", testTranscriptMap(), ".")
	if err != nil {
		t.Fatal(err)
	}
	if fd.Vendor != "cisco" {
		t.Errorf("Vendor = %q, want %q", fd.Vendor, "cisco")
	}
	if fd.Platform != "csr1000v" {
		t.Errorf("Platform = %q, want %q", fd.Platform, "csr1000v")
	}
	if fd.Hostname != "testhost" {
		t.Errorf("Hostname = %q, want %q", fd.Hostname, "testhost")
	}
	if fd.DefaultHostname != "testhost" {
		t.Errorf("DefaultHostname = %q, want %q", fd.DefaultHostname, "testhost")
	}
	if fd.Password != "secret" {
		t.Errorf("Password = %q, want %q", fd.Password, "secret")
	}
	if _, ok := fd.SupportedCommands["show version"]; !ok {
		t.Error("SupportedCommands missing 'show version'")
	}
	if fd.ContextSearch["base"] != ">" {
		t.Errorf("ContextSearch[base] = %q, want %q", fd.ContextSearch["base"], ">")
	}
	if fd.ContextHierarchy["(config)#"] != "#" {
		t.Errorf("ContextHierarchy[(config)#] = %q, want %q", fd.ContextHierarchy["(config)#"], "#")
	}
}

func TestInitGeneric_UnknownPlatform(t *testing.T) {
	tm := transcript.Map{
		Platforms: map[string]transcript.Platform{
			"other": {Hostname: "other"},
		},
	}
	_, err := InitGeneric("csr1000v", tm, ".")
	if err == nil {
		t.Error("expected error for unknown platform")
	}
}

func TestFakeDevice_Copy(t *testing.T) {
	fd := &FakeDevice{
		Vendor:            "cisco",
		Platform:          "csr1000v",
		Hostname:          "original",
		DefaultHostname:   "original",
		Password:          "secret",
		SupportedCommands: SupportedCommands{"show version": "output"},
		ContextSearch:     map[string]string{"base": ">"},
		ContextHierarchy:  map[string]string{">": "exit"},
	}

	c := fd.Copy()

	if c.Hostname != fd.Hostname {
		t.Errorf("Hostname = %q, want %q", c.Hostname, fd.Hostname)
	}

	// Mutate copy and verify original is unchanged
	c.Hostname = "modified"
	c.SupportedCommands["show version"] = "changed"
	c.ContextSearch["base"] = "changed"
	c.ContextHierarchy[">"] = "changed"

	if fd.Hostname == "modified" {
		t.Error("Copy shares Hostname with original")
	}
	if fd.SupportedCommands["show version"] == "changed" {
		t.Error("Copy shares SupportedCommands with original")
	}
	if fd.ContextSearch["base"] == "changed" {
		t.Error("Copy shares ContextSearch with original")
	}
	if fd.ContextHierarchy[">"] == "changed" {
		t.Error("Copy shares ContextHierarchy with original")
	}
}

func TestInitGeneric_BadTranscriptFile(t *testing.T) {
	tm := transcript.Map{
		Platforms: map[string]transcript.Platform{
			"csr1000v": {
				CommandTranscripts: map[string]string{
					"show version": "/nonexistent/file.txt",
				},
			},
		},
	}
	_, err := InitGeneric("csr1000v", tm, ".")
	if err == nil {
		t.Error("expected error for missing transcript file")
	}
}

func TestTranscriptMapIntegrity(t *testing.T) {
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir("fakedevices") })

	tm, err := transcript.Load("transcripts/transcript_map.yaml")
	if err != nil {
		t.Fatalf("loading transcript map: %v", err)
	}

	// Track all referenced paths for orphan detection
	// Paths in the map are relative to the map file's directory (transcripts/)
	const mapDir = "transcripts"
	referenced := map[string]bool{}

	for platform, p := range tm.Platforms {
		for cmd, path := range p.CommandTranscripts {
			resolved := filepath.Join(mapDir, path)
			referenced[resolved] = true
			if _, err := os.Stat(resolved); err != nil {
				t.Errorf("platform %q command %q: file not found: %s", platform, cmd, resolved)
			}
		}
	}
	for name, s := range tm.Scenarios {
		for i, step := range s.Sequence {
			resolved := filepath.Join(mapDir, step.Transcript)
			referenced[resolved] = true
			if _, err := os.Stat(resolved); err != nil {
				t.Errorf("scenario %q step %d (%q): file not found: %s", name, i, step.Command, resolved)
			}
		}
	}

	// Walk transcripts/ and flag unreferenced .txt files
	err = filepath.WalkDir("transcripts", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".txt" {
			return err
		}
		// generic_empty_return.txt is intentionally shared — skip
		if filepath.Base(path) == "generic_empty_return.txt" {
			return nil
		}
		if !referenced[path] {
			t.Errorf("orphan transcript not referenced in transcript_map.yaml: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking transcripts dir: %v", err)
	}
}

func TestInitScenario(t *testing.T) {
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir("fakedevices") })

	tm, err := transcript.Load("transcripts/transcript_map.yaml")
	if err != nil {
		t.Fatal(err)
	}

	fd, steps, err := InitScenario("csr1000v-add-interface", tm, "transcripts")
	if err != nil {
		t.Fatalf("InitScenario: %v", err)
	}
	if fd.Platform != "csr1000v" {
		t.Errorf("Platform = %q, want csr1000v", fd.Platform)
	}
	if len(steps) != 9 {
		t.Errorf("steps len = %d, want 9", len(steps))
	}
	if steps[0].Command != "enable" {
		t.Errorf("steps[0].Command = %q, want 'enable'", steps[0].Command)
	}
}

func TestInitScenario_UnknownScenario(t *testing.T) {
	tm := transcript.Map{Platforms: map[string]transcript.Platform{}}
	_, _, err := InitScenario("nonexistent", tm, ".")
	if err == nil {
		t.Error("expected error for unknown scenario")
	}
}
