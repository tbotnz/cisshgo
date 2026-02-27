// cisshgo is a lightweight SSH server that plays back pre-defined command
// transcripts to emulate network equipment for testing purposes.
//
// Usage:
//
//	cisshgo [--inventory path] [--platform name] [--listeners N] [--starting-port N] [--transcript-map path]
package main

import (
	"context"
	"log"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/alecthomas/kong"

	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/ssh_server/handlers"
	"github.com/tbotnz/cisshgo/ssh_server/sshlisteners"
	"github.com/tbotnz/cisshgo/utils"
)

type listenerSpec struct {
	fd       *fakedevices.FakeDevice
	port     int
	sequence []utils.SequenceStep // nil for platform-only listeners
}

func run(ctx context.Context, cli utils.CLI) error {
	myTranscriptMap, err := utils.LoadTranscriptMap(cli.TranscriptMap)
	if err != nil {
		return err
	}

	baseDir := filepath.Dir(cli.TranscriptMap)
	if err := utils.ValidateTranscriptMap(myTranscriptMap, baseDir); err != nil {
		return err
	}

	var specs []listenerSpec

	if cli.Inventory != "" {
		inv, err := utils.LoadInventory(cli.Inventory)
		if err != nil {
			return err
		}
		port := cli.StartingPort
		for _, entry := range inv.Devices {
			if entry.Scenario != "" {
				fd, seq, err := fakedevices.InitScenario(entry.Scenario, myTranscriptMap, baseDir)
				if err != nil {
					return err
				}
				for i := 0; i < entry.Count; i++ {
					specs = append(specs, listenerSpec{fd, port, seq})
					port++
				}
			} else {
				fd, err := fakedevices.InitGeneric(entry.Platform, myTranscriptMap, baseDir)
				if err != nil {
					return err
				}
				for i := 0; i < entry.Count; i++ {
					specs = append(specs, listenerSpec{fd, port, nil})
					port++
				}
			}
		}
	} else {
		fd, err := fakedevices.InitGeneric(cli.Platform, myTranscriptMap, baseDir)
		if err != nil {
			return err
		}
		for port := cli.StartingPort; port < cli.StartingPort+cli.Listeners; port++ {
			specs = append(specs, listenerSpec{fd, port, nil})
		}
	}

	var wg sync.WaitGroup
	for _, spec := range specs { // coverage-ignore
		wg.Add(1)                 // coverage-ignore
		go func(s listenerSpec) { // coverage-ignore
			defer wg.Done()        // coverage-ignore
			var err error          // coverage-ignore
			if s.sequence != nil { // coverage-ignore
				err = sshlisteners.ScenarioListener(ctx, s.fd, s.sequence, s.port) // coverage-ignore
			} else { // coverage-ignore
				err = sshlisteners.GenericListener(ctx, s.fd, s.port, handlers.GenericCiscoHandler) // coverage-ignore
			} // coverage-ignore
			if err != nil { // coverage-ignore
				log.Printf("listener on port %d: %v", s.port, err) // coverage-ignore
			} // coverage-ignore
		}(spec) // coverage-ignore
	} // coverage-ignore

	wg.Wait()  // coverage-ignore
	return nil // coverage-ignore
}

func main() { // coverage-ignore
	var cli utils.CLI
	kong.Parse(&cli,
		kong.Name("cisshgo"),
		kong.Description("Lightweight SSH server that emulates network equipment for testing."),
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, cli); err != nil {
		log.Fatal(err)
	}
}
