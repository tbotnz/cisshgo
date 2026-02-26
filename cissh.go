// cisshgo is a lightweight SSH server that plays back pre-defined command
// transcripts to emulate network equipment for testing purposes.
//
// Usage:
//
//	cisshgo [-listeners N] [-startingPort N] [-transcriptMap path]
package main

import (
	"log"

	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/ssh_server/handlers"
	"github.com/tbotnz/cisshgo/ssh_server/sshlistners"
	"github.com/tbotnz/cisshgo/utils"
)

func run() error {
	// Parse the command line arguments
	numListeners, startingPortPtr, myTranscriptMap, err := utils.ParseArgs()
	if err != nil {
		return err
	}

	// Initialize our fake device
	myFakeDevice, err := fakedevices.InitGeneric("cisco", "csr1000v", myTranscriptMap)
	if err != nil {
		return err
	}

	// Make a Channel named "done" for handling Goroutines
	done := make(chan bool, 1)

	// Iterate through the server ports and spawn a Goroutine for each
	for portNumber := *startingPortPtr; portNumber < numListeners; portNumber++ { // coverage-ignore // blocks on <-done
		go sshlistners.GenericListener(myFakeDevice, portNumber, handlers.GenericCiscoHandler, done)
	}

	// Wait on channel
	<-done
	return nil
}

func main() { // coverage-ignore
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
