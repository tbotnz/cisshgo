package fakedevices

import "testing"

func TestTranscriptReader_PlainText(t *testing.T) {
	fd := &FakeDevice{Hostname: "router1"}
	out, err := TranscriptReader("plain text no templates", fd)
	if err != nil {
		t.Fatal(err)
	}
	if out != "plain text no templates" {
		t.Errorf("got %q, want %q", out, "plain text no templates")
	}
}

func TestTranscriptReader_WithTemplate(t *testing.T) {
	fd := &FakeDevice{
		Hostname: "myrouter",
		Vendor:   "cisco",
		Platform: "csr1000v",
	}
	out, err := TranscriptReader("{{.Hostname}} uptime is 4 hours", fd)
	if err != nil {
		t.Fatal(err)
	}
	want := "myrouter uptime is 4 hours"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestTranscriptReader_MultipleFields(t *testing.T) {
	fd := &FakeDevice{
		Hostname: "sw1",
		Vendor:   "cisco",
		Platform: "csr1000v",
	}
	out, err := TranscriptReader("{{.Vendor}} {{.Platform}} {{.Hostname}}", fd)
	if err != nil {
		t.Fatal(err)
	}
	want := "cisco csr1000v sw1"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}
