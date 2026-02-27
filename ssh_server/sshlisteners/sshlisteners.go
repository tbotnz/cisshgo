// Package sshlisteners contains SSH Listeners for cisshgo to utilize when building
// fake devices to emulate network equipment
package sshlisteners

import (
	"context"
	"log"
	"strconv"

	"github.com/gliderlabs/ssh"

	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/ssh_server/handlers"
)

// GenericListener starts an SSH server on the given port and blocks until ctx is cancelled.
func GenericListener(
	ctx context.Context,
	myFakeDevice *fakedevices.FakeDevice,
	portNumber int,
	myHandler handlers.PlatformHandler,
) error {
	portString := ":" + strconv.Itoa(portNumber)
	log.Printf("Starting cissh.go ssh server on port %s\n", portString)

	srv := &ssh.Server{
		Addr:    portString,
		Handler: myHandler(myFakeDevice.Copy()),
		PasswordHandler: func(sshCtx ssh.Context, pass string) bool {
			return pass == myFakeDevice.Password
		},
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	err := srv.ListenAndServe()
	if err == ssh.ErrServerClosed {
		return nil
	}
	return err
}
