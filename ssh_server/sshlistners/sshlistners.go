// Package sshlistners SSH related functions and helpers for cisgo-ios
package sshlistners

import (
	"log"
	"strconv"

	"github.com/gliderlabs/ssh"

	"github.com/tbotnz/cisgo-ios/fakedevices"
	"github.com/tbotnz/cisgo-ios/ssh_server/handlers"
)

// GenericListner function that creates a fake device and terminal session
func GenericListner(
	myFakeDevice *fakedevices.FakeDevice,
	portNumber int,
	myHandler handlers.PlatformHandler,
	done chan bool,
) {

	// Prepare an SSH Handler for our fake device.
	// This will allow for per-device-type handling/features
	myHandler(myFakeDevice)

	portString := ":" + strconv.Itoa(portNumber)
	log.Printf("Starting cis.go ssh server on port %s\n", portString)

	log.Fatal(
		// Actually kick off the SSH server and listen on the given port
		ssh.ListenAndServe(
			portString, // Address string in the form of "ip:port"
			nil,        // ssh.Handler (we're using the DefaultHandler assigned above)
			ssh.PasswordAuth(
				func(ctx ssh.Context, pass string) bool {
					return pass == myFakeDevice.Password
				}, // Handle SSH authentication with the provided password
			), // Additional ssh.Options, in this case Password handling
		),
	)

	done <- true
}
