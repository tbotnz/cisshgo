package fakedevices

import (
	"os"
	"testing"

	"github.com/tbotnz/cisshgo/utils"
)

func testTranscriptMap() utils.TranscriptMap {
	return utils.TranscriptMap{
		Platforms: map[string]utils.TranscriptMapPlatform{
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
	tm := utils.TranscriptMap{
		Platforms: map[string]utils.TranscriptMapPlatform{
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
	tm := utils.TranscriptMap{
		Platforms: map[string]utils.TranscriptMapPlatform{
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
