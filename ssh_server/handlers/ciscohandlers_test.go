package handlers

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	gossh "golang.org/x/crypto/ssh"

	"github.com/gliderlabs/ssh"
	"github.com/tbotnz/cisshgo/fakedevices"
)

// newTestDevice creates a FakeDevice for testing without reading files from disk.
func newTestDevice() *fakedevices.FakeDevice {
	return &fakedevices.FakeDevice{
		Vendor:   "cisco",
		Platform: "csr1000v",
		Hostname: "testhost",
		DefaultHostname: "testhost",
		Password: "admin",
		SupportedCommands: fakedevices.SupportedCommands{
			"show version":            "FakeOS version 1.0\n{{.Hostname}} uptime is 1 hour\n",
			"show ip interface brief":  "Interface  IP-Address  OK?\n",
			"terminal length 0":        "",
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
	}
}

// startTestServer starts an SSH server on a random port and returns the address and a cleanup func.
func startTestServer(t *testing.T, fd *fakedevices.FakeDevice) (string, func()) {
	t.Helper()

	// Register the handler
	GenericCiscoHandler(fd)

	// Find a free port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close()

	srv := &ssh.Server{
		Addr: addr,
		Handler: ssh.DefaultHandler,
		PasswordHandler: func(ctx ssh.Context, pass string) bool {
			return pass == fd.Password
		},
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != ssh.ErrServerClosed {
			// server stopped
		}
	}()

	// Wait for server to be ready
	for i := 0; i < 20; i++ {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	return addr, func() { srv.Close() }
}

// sshSession connects to the test server and returns a function to send commands and read output.
func sshSession(t *testing.T, addr string) (*gossh.Session, func()) {
	t.Helper()
	config := &gossh.ClientConfig{
		User:            "admin",
		Auth:            []gossh.AuthMethod{gossh.Password("admin")},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout:         2 * time.Second,
	}
	client, err := gossh.Dial("tcp", addr, config)
	if err != nil {
		t.Fatalf("ssh dial: %v", err)
	}
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		t.Fatalf("new session: %v", err)
	}

	// Request a PTY so the server gives us a terminal
	if err := session.RequestPty("xterm", 80, 200, gossh.TerminalModes{}); err != nil {
		session.Close()
		client.Close()
		t.Fatalf("request pty: %v", err)
	}

	return session, func() {
		session.Close()
		client.Close()
	}
}

// interact sends commands over stdin/stdout pipes and collects output.
func interact(t *testing.T, addr string, commands []string) string {
	t.Helper()
	session, cleanup := sshSession(t, addr)
	defer cleanup()

	stdin, err := session.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	if err := session.Shell(); err != nil {
		t.Fatalf("shell: %v", err)
	}

	// Small delay to let prompt appear
	time.Sleep(200 * time.Millisecond)

	for _, cmd := range commands {
		fmt.Fprintf(stdin, "%s\n", cmd)
		time.Sleep(200 * time.Millisecond)
	}

	// Send exit to close
	fmt.Fprintf(stdin, "exit\n")
	time.Sleep(200 * time.Millisecond)

	buf := make([]byte, 64*1024)
	n, _ := stdout.Read(buf)
	return string(buf[:n])
}

func TestHandler_ShowVersion(t *testing.T) {
	fd := newTestDevice()
	addr, cleanup := startTestServer(t, fd)
	defer cleanup()

	out := interact(t, addr, []string{"show version"})
	if !strings.Contains(out, "FakeOS version 1.0") {
		t.Errorf("expected 'FakeOS version 1.0' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "testhost uptime is 1 hour") {
		t.Errorf("expected template-rendered hostname in output, got:\n%s", out)
	}
}

func TestHandler_AbbreviatedCommand(t *testing.T) {
	fd := newTestDevice()
	addr, cleanup := startTestServer(t, fd)
	defer cleanup()

	out := interact(t, addr, []string{"sho ver"})
	if !strings.Contains(out, "FakeOS version 1.0") {
		t.Errorf("expected abbreviated 'sho ver' to match show version, got:\n%s", out)
	}
}

func TestHandler_UnknownCommand(t *testing.T) {
	fd := newTestDevice()
	addr, cleanup := startTestServer(t, fd)
	defer cleanup()

	out := interact(t, addr, []string{"do something weird"})
	if !strings.Contains(out, "Unknown command") {
		t.Errorf("expected 'Unknown command' in output, got:\n%s", out)
	}
}

func TestHandler_ContextSwitching(t *testing.T) {
	fd := newTestDevice()
	addr, cleanup := startTestServer(t, fd)
	defer cleanup()

	out := interact(t, addr, []string{"enable", "configure terminal", "end"})
	if !strings.Contains(out, "#") {
		t.Errorf("expected '#' prompt after enable, got:\n%s", out)
	}
}

func TestHandler_HostnameChange(t *testing.T) {
	fd := newTestDevice()
	addr, cleanup := startTestServer(t, fd)
	defer cleanup()

	out := interact(t, addr, []string{"enable", "configure terminal", "hostname newname", "end"})
	if !strings.Contains(out, "newname") {
		t.Errorf("expected 'newname' in prompt after hostname change, got:\n%s", out)
	}
}

func TestHandler_ResetState(t *testing.T) {
	fd := newTestDevice()
	addr, cleanup := startTestServer(t, fd)
	defer cleanup()

	out := interact(t, addr, []string{"enable", "configure terminal", "hostname changed", "reset state"})
	if !strings.Contains(out, "Resetting State") {
		t.Errorf("expected 'Resetting State' in output, got:\n%s", out)
	}
}

func TestHandler_EmptyInput(t *testing.T) {
	fd := newTestDevice()
	addr, cleanup := startTestServer(t, fd)
	defer cleanup()

	// Just send empty line then exit — should not crash
	out := interact(t, addr, []string{""})
	if !strings.Contains(out, "testhost") {
		t.Errorf("expected prompt with hostname, got:\n%s", out)
	}
}

func TestHandler_AmbiguousCommand(t *testing.T) {
	fd := newTestDevice()
	// Add commands that will be ambiguous with "s v"
	fd.SupportedCommands["show vlan"] = "vlan info\n"
	addr, cleanup := startTestServer(t, fd)
	defer cleanup()

	out := interact(t, addr, []string{"s v"})
	if !strings.Contains(out, "Ambiguous command") {
		t.Errorf("expected 'Ambiguous command' in output, got:\n%s", out)
	}
}

func TestHandler_TranscriptReaderError(t *testing.T) {
	fd := newTestDevice()
	// Set a command with an invalid template to trigger TranscriptReader error
	fd.SupportedCommands["show bad"] = "{{.Bad"
	addr, cleanup := startTestServer(t, fd)
	defer cleanup()

	// Should not crash, just close the session
	_ = interact(t, addr, []string{"show bad"})
}
