// cisshgo is a lightweight SSH server that plays back pre-defined command
// transcripts to emulate network equipment for testing purposes.
//
// Usage:
//
//	cisshgo [-listeners N] [-startingPort N] [-transcriptMap path]
package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/ssh_server/handlers"
	"github.com/tbotnz/cisshgo/ssh_server/sshlisteners"
	"github.com/tbotnz/cisshgo/utils"
)

func run(ctx context.Context) error {
	numListeners, startingPort, myTranscriptMap, err := utils.ParseArgs()
	if err != nil {
		return err
	}

	myFakeDevice, err := fakedevices.InitGeneric("cisco", "csr1000v", myTranscriptMap)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for portNumber := startingPort; portNumber < numListeners; portNumber++ { // coverage-ignore
		wg.Add(1)           // coverage-ignore
		go func(port int) { // coverage-ignore
			defer wg.Done()                                                                                             // coverage-ignore
			if err := sshlisteners.GenericListener(ctx, myFakeDevice, port, handlers.GenericCiscoHandler); err != nil { // coverage-ignore
				log.Printf("listener on port %d: %v", port, err) // coverage-ignore
			} // coverage-ignore
		}(portNumber) // coverage-ignore
	} // coverage-ignore

	wg.Wait()  // coverage-ignore
	return nil // coverage-ignore
}

func main() { // coverage-ignore
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}
