// Package sshlistners contains SSH Listeners for cisshgo to utilize when building
// fake devices to emulate network equipment
package sshlistners

import (
	"log"
	"strconv"

	"github.com/gliderlabs/ssh"

	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/ssh_server/handlers"
)

// GenericListener function that creates a fake device and terminal session
func GenericListener(
	myFakeDevice *fakedevices.FakeDevice,
	portNumber int,
	myHandler handlers.PlatformHandler,
	done chan bool,
) error {

	// Prepare an SSH Handler for our fake device.
	myHandler(myFakeDevice)

	portString := ":" + strconv.Itoa(portNumber)
	log.Printf("Starting cissh.go ssh server on port %s\n", portString)

	err := ssh.ListenAndServe(
		portString,
		nil,
		ssh.PasswordAuth(
			func(ctx ssh.Context, pass string) bool {
				return pass == myFakeDevice.Password
			},
		),
	)

	done <- true
	return err
}
