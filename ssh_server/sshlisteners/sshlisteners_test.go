package sshlisteners

import (
	"context"
	"net"
	"testing"
	"time"

	gossh "golang.org/x/crypto/ssh"

	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/ssh_server/handlers"
)

func TestGenericListener(t *testing.T) {
	fd := &fakedevices.FakeDevice{
		Vendor:          "cisco",
		Platform:        "csr1000v",
		Hostname:        "testhost",
		DefaultHostname: "testhost",
		Password:        "admin",
		SupportedCommands: fakedevices.SupportedCommands{
			"show version": "version 1.0\n",
		},
		ContextSearch:    map[string]string{"base": ">"},
		ContextHierarchy: map[string]string{">": "exit"},
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	addr := ln.Addr().String()
	ln.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		GenericListener(ctx, fd, port, handlers.GenericCiscoHandler)
	}()

	// Wait for server to be ready
	for i := 0; i < 20; i++ {
		conn, dialErr := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if dialErr == nil {
			conn.Close()
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

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
	client.Close()

	// Verify wrong password is rejected
	badConfig := &gossh.ClientConfig{
		User:            "admin",
		Auth:            []gossh.AuthMethod{gossh.Password("wrong")},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout:         2 * time.Second,
	}
	_, err = gossh.Dial("tcp", addr, badConfig)
	if err == nil {
		t.Error("expected auth failure with wrong password")
	}

	// Verify graceful shutdown via context cancellation
	cancel()
	time.Sleep(100 * time.Millisecond)
	_, err = gossh.Dial("tcp", addr, config)
	if err == nil {
		t.Error("expected server to be shut down after context cancel")
	}
}

func TestGenericListener_PortInUse(t *testing.T) {
	fd := &fakedevices.FakeDevice{
		Password:          "admin",
		SupportedCommands: fakedevices.SupportedCommands{},
		ContextSearch:     map[string]string{"base": ">"},
		ContextHierarchy:  map[string]string{">": "exit"},
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	err = GenericListener(context.Background(), fd, port, handlers.GenericCiscoHandler)
	if err == nil {
		t.Error("expected error when port is already in use")
	}
}

func TestGenericListener_UsernameEnforcement(t *testing.T) {
	fd := &fakedevices.FakeDevice{
		Hostname:          "testhost",
		DefaultHostname:   "testhost",
		Username:          "admin",
		Password:          "admin",
		SupportedCommands: fakedevices.SupportedCommands{},
		ContextSearch:     map[string]string{"base": ">"},
		ContextHierarchy:  map[string]string{">": "exit"},
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	addr := ln.Addr().String()
	ln.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { GenericListener(ctx, fd, port, handlers.GenericCiscoHandler) }()

	for i := 0; i < 20; i++ {
		conn, dialErr := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if dialErr == nil {
			conn.Close()
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	cfg := &gossh.ClientConfig{
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout:         2 * time.Second,
	}

	// Correct username + password — should succeed
	cfg.User = "admin"
	cfg.Auth = []gossh.AuthMethod{gossh.Password("admin")}
	client, err := gossh.Dial("tcp", addr, cfg)
	if err != nil {
		t.Fatalf("expected success with correct credentials: %v", err)
	}
	client.Close()

	// Wrong username — should fail
	cfg.User = "wronguser"
	_, err = gossh.Dial("tcp", addr, cfg)
	if err == nil {
		t.Error("expected auth failure with wrong username")
	}
}
