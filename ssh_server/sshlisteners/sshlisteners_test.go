package sshlisteners

import (
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

	// Find a free port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	addr := ln.Addr().String()
	ln.Close()

	done := make(chan bool, 1)
	errCh := make(chan error, 1)
	go func() {
		errCh <- GenericListener(fd, port, handlers.GenericCiscoHandler, done)
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

	// Connect via SSH and verify auth works
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
}

func TestGenericListener_PortInUse(t *testing.T) {
	fd := &fakedevices.FakeDevice{
		Password:          "admin",
		SupportedCommands: fakedevices.SupportedCommands{},
		ContextSearch:     map[string]string{"base": ">"},
		ContextHierarchy:  map[string]string{">": "exit"},
	}

	// Occupy a port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	done := make(chan bool, 1)
	// GenericListener should return an error since the port is in use
	err = GenericListener(fd, port, handlers.GenericCiscoHandler, done)
	if err == nil {
		t.Error("expected error when port is already in use")
	}
}
