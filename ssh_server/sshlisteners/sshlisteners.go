// Package sshlisteners contains SSH Listeners for cisshgo to utilize when building
// fake devices to emulate network equipment
package sshlisteners

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/gliderlabs/ssh"

	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/ssh_server/handlers"
	"github.com/tbotnz/cisshgo/transcript"
)

// GenericListener starts an SSH server on the given port and blocks until ctx is cancelled.
func GenericListener(
	ctx context.Context,
	fd *fakedevices.FakeDevice,
	port int,
	myHandler handlers.PlatformHandler,
) error {
	return listen(ctx, fd, port, myHandler(fd.Copy()))
}

// ScenarioListener starts an SSH server that plays back a scenario sequence.
func ScenarioListener(
	ctx context.Context,
	fd *fakedevices.FakeDevice,
	sequence []transcript.SequenceStep,
	port int,
) error {
	return listen(ctx, fd, port, handlers.GenericCiscoScenarioHandler(fd.Copy(), sequence))
}

func listen(ctx context.Context, fd *fakedevices.FakeDevice, port int, handler ssh.Handler) error {
	portString := ":" + strconv.Itoa(port)
	if fd.ScenarioName != "" {
		log.Printf("Starting listener on %s [scenario=%s hostname=%s user=%s]",
			portString, fd.ScenarioName, fd.Hostname, fd.Username)
	} else {
		log.Printf("Starting listener on %s [platform=%s hostname=%s user=%s]",
			portString, fd.Platform, fd.Hostname, fd.Username)
	}

	srv := &ssh.Server{
		Addr:    portString,
		Handler: handler,
		PasswordHandler: func(sshCtx ssh.Context, pass string) bool {
			return sshCtx.User() == fd.Username && pass == fd.Password
		},
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	err := srv.ListenAndServe()
	if err == ssh.ErrServerClosed {
		return nil
	}
	return err
}
