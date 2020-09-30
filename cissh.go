package main

import (
	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/ssh_server/handlers"
	"github.com/tbotnz/cisshgo/ssh_server/sshlistners"
	"github.com/tbotnz/cisshgo/utils"
)

func main() {

	// Parse the command line arguments
	numListeners, startingPortPtr, myTranscriptMap := utils.ParseArgs()

	// Initialize our fake device
	myFakeDevice := fakedevices.InitGeneric(
		"cisco",         // Vendor
		"csr1000v",      // Platform
		myTranscriptMap, // Transcript map with locations of command output to play back
	)

	// Make a Channel named "done" for handling Goroutines, which expects a bool as return value
	done := make(chan bool, 1)

	// Iterate through the server ports and spawn a Goroutine for each
	for portNumber := *startingPortPtr; portNumber < numListeners; portNumber++ {
		// Today this is just spawning a generic listener.
		// In the future, this is where we could split out listeners/handlers by device type.
		go sshlistners.GenericListener(myFakeDevice, portNumber, handlers.GenericCiscoHandler, done)
	}

	// Receive all the values from the channel (essentially wait on it to be empty)
	<-done
}
