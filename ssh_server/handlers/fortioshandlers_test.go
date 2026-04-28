package handlers

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/transcript"
	gossh "golang.org/x/crypto/ssh"
)

func newFortiTestDevice() *fakedevices.FakeDevice {
	return &fakedevices.FakeDevice{
		Vendor:          "fortinet",
		Platform:        "fortios",
		Hostname:        "FGT",
		DefaultHostname: "FGT",
		Username:        "admin",
		Password:        "admin",
		PromptFormat:    "{hostname} {context}",
		SupportedCommands: fakedevices.SupportedCommands{
			"get system status":             "Hostname: {{.Hostname}}\nCurrent virtual domain: root\n",
			"get system interface":          "== [ port1 ]\nname: port1   mode: static\n",
			"get router info bgp neighbors": "Neighbor        V         AS MsgRcvd MsgSent   TblVer  InQ OutQ Up/Down  State/PfxRcd\n203.0.113.2     4      64513    1542    1479       18    0    0 2d04h12m           12\n",
			"show firewall policy":          "config firewall policy\n    edit 1\n    next\nend\n",
		},
		ContextSearch: map[string]string{
			"base": "#",
		},
		ContextHierarchy: map[string]string{
			"#": "exit",
		},
	}
}

func startFortiTestServer(t *testing.T, fd *fakedevices.FakeDevice) (string, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close()

	srv := &ssh.Server{
		Addr:    addr,
		Handler: FortiOSHandler(fd),
		PasswordHandler: func(ctx ssh.Context, pass string) bool {
			return pass == fd.Password
		},
	}
	go func() { _ = srv.ListenAndServe() }()

	for i := 0; i < 20; i++ {
		conn, dialErr := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if dialErr == nil {
			conn.Close()
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	return addr, func() { srv.Close() }
}

func TestResolvePlatformHandlers_FortiOS(t *testing.T) {
	platformHandler, scenarioHandler := ResolvePlatformHandlers("fortios")
	if platformHandler == nil || scenarioHandler == nil {
		t.Fatal("expected handlers for fortios")
	}
}

func TestResolvePlatformHandlers_DefaultsToCisco(t *testing.T) {
	platformHandler, scenarioHandler := ResolvePlatformHandlers("csr1000v")
	if platformHandler == nil || scenarioHandler == nil {
		t.Fatal("expected default handlers")
	}
}

func TestFortiOSPromptAndContextTransitions(t *testing.T) {
	fd := newFortiTestDevice()
	addr, cleanup := startFortiTestServer(t, fd)
	defer cleanup()

	out := interact(t, addr, []string{"config system interface", "edit port1", "next", "end"})
	for _, expected := range []string{"FGT #", "FGT (interface) #", "FGT (port1) #"} {
		if !strings.Contains(out, expected) {
			t.Fatalf("expected %q in output, got:\n%s", expected, out)
		}
	}
}

func TestFortiOSCommandTranscriptMatch(t *testing.T) {
	fd := newFortiTestDevice()
	addr, cleanup := startFortiTestServer(t, fd)
	defer cleanup()

	out := interact(t, addr, []string{"get system status", "get system interface", "get router info bgp neighbors", "show firewall policy"})
	if !strings.Contains(out, "Current virtual domain: root") {
		t.Fatalf("expected system status transcript output, got:\n%s", out)
	}
	if !strings.Contains(out, "config firewall policy") {
		t.Fatalf("expected firewall policy transcript output, got:\n%s", out)
	}
	if !strings.Contains(out, "== [ port1 ]") {
		t.Fatalf("expected system interface transcript output, got:\n%s", out)
	}
	if !strings.Contains(out, "State/PfxRcd") {
		t.Fatalf("expected BGP neighbor transcript output, got:\n%s", out)
	}
}

func TestFortiOSUnknownCommand(t *testing.T) {
	fd := newFortiTestDevice()
	addr, cleanup := startFortiTestServer(t, fd)
	defer cleanup()

	out := interact(t, addr, []string{"diagnose sniffer packet any"})
	if !strings.Contains(out, "Command fail. Return code -61") {
		t.Fatalf("expected FortiOS-like unknown command error, got:\n%s", out)
	}
	if !strings.Contains(out, "FGT #") {
		t.Fatalf("expected prompt after unknown command, got:\n%s", out)
	}
}

func TestFortiOSCIStyleSession(t *testing.T) {
	fd := newFortiTestDevice()
	addr, cleanup := startFortiTestServer(t, fd)
	defer cleanup()

	out := interact(t, addr, []string{"get system status", "config system interface", "edit port1", "next", "end", "show firewall policy"})
	if !strings.Contains(out, "Hostname: FGT") {
		t.Fatalf("expected rendered hostname, got:\n%s", out)
	}
	if !strings.Contains(out, "FGT (port1) #") {
		t.Fatalf("expected edit context prompt, got:\n%s", out)
	}
	if !strings.Contains(out, "config firewall policy") {
		t.Fatalf("expected final command output, got:\n%s", out)
	}
}

func TestFortiOSExecModeKnownAndUnknown(t *testing.T) {
	fd := newFortiTestDevice()
	addr, cleanup := startFortiTestServer(t, fd)
	defer cleanup()

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
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		t.Fatalf("new session: %v", err)
	}
	out, err := session.Output("get system status")
	session.Close()
	if err != nil {
		t.Fatalf("exec known command: %v", err)
	}
	if !strings.Contains(string(out), "Current virtual domain: root") {
		t.Fatalf("expected known exec output, got:\n%s", string(out))
	}

	session2, err := client.NewSession()
	if err != nil {
		t.Fatalf("new session: %v", err)
	}
	out2, err := session2.Output("get not real")
	session2.Close()
	if err != nil {
		t.Fatalf("exec unknown command should still exit 0: %v", err)
	}
	if len(out2) != 0 {
		t.Fatalf("expected empty output for unknown exec command, got:\n%s", string(out2))
	}
}

func TestFortiOSScenarioHandler(t *testing.T) {
	fd := newFortiTestDevice()
	sequence := []transcript.SequenceStep{
		{Command: "show firewall policy", Transcript: "scenario first\n"},
		{Command: "show firewall policy", Transcript: "scenario second\n"},
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close()

	srv := &ssh.Server{
		Addr:    addr,
		Handler: FortiOSScenarioHandler(fd, sequence),
		PasswordHandler: func(ctx ssh.Context, pass string) bool {
			return pass == fd.Password
		},
	}
	go func() { _ = srv.ListenAndServe() }()
	defer srv.Close()

	for i := 0; i < 20; i++ {
		conn, dialErr := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if dialErr == nil {
			conn.Close()
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	out := interact(t, addr, []string{"show firewall policy", "show firewall policy"})
	if !strings.Contains(out, "scenario first") || !strings.Contains(out, "scenario second") {
		t.Fatalf("expected scenario transcripts in output, got:\n%s", out)
	}
}
