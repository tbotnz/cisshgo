// Package sshlistners contains SSH Listners for cisshgo to utilize when building
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
) {

	// Prepare an SSH Handler for our fake device.
	// This will allow for per-device-type handling/features
	myHandler(myFakeDevice)

	portString := ":" + strconv.Itoa(portNumber)
	log.Printf("Starting cissh.go ssh server on port %s\n", portString)

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
