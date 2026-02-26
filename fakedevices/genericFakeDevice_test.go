package fakedevices

import (
	"os"
	"testing"

	"github.com/tbotnz/cisshgo/utils"
)

func testTranscriptMap() utils.TranscriptMap {
	return utils.TranscriptMap{
		Platforms: []map[string]utils.TranscriptMapPlatform{
			{
				"csr1000v": utils.TranscriptMapPlatform{
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
		},
	}
}

func TestInitGeneric(t *testing.T) {
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir("fakedevices") })

	fd, err := InitGeneric("cisco", "csr1000v", testTranscriptMap())
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
	tm := utils.TranscriptMap{
		Platforms: []map[string]utils.TranscriptMapPlatform{
			{"other": utils.TranscriptMapPlatform{Hostname: "other"}},
		},
	}
	fd, err := InitGeneric("cisco", "csr1000v", tm)
	if err != nil {
		t.Fatal(err)
	}
	if fd.Hostname != "" {
		t.Errorf("Hostname = %q, want empty for unknown platform", fd.Hostname)
	}
	if len(fd.SupportedCommands) != 0 {
		t.Errorf("SupportedCommands should be empty for unknown platform")
	}
}

func TestInitGeneric_BadTranscriptFile(t *testing.T) {
	tm := utils.TranscriptMap{
		Platforms: []map[string]utils.TranscriptMapPlatform{
			{
				"csr1000v": utils.TranscriptMapPlatform{
					CommandTranscripts: map[string]string{
						"show version": "/nonexistent/file.txt",
					},
				},
			},
		},
	}
	_, err := InitGeneric("cisco", "csr1000v", tm)
	if err == nil {
		t.Error("expected error for missing transcript file")
	}
}
