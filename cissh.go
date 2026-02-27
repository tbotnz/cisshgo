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

	"github.com/tbotnz/cisshgo/config"
	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/ssh_server/handlers"
	"github.com/tbotnz/cisshgo/ssh_server/sshlisteners"
	"github.com/tbotnz/cisshgo/transcript"
)

type listenerConfig struct {
	fd       *fakedevices.FakeDevice
	port     int
	sequence []transcript.SequenceStep // nil for platform-only listeners
}

func resolveListeners(cli config.CLI, tm transcript.Map, baseDir string) ([]listenerConfig, error) {
	if cli.Inventory == "" {
		fd, err := fakedevices.InitGeneric(cli.Platform, tm, baseDir)
		if err != nil {
			return nil, err
		}
		var configs []listenerConfig
		for port := cli.StartingPort; port < cli.StartingPort+cli.Listeners; port++ {
			configs = append(configs, listenerConfig{fd, port, nil})
		}
		return configs, nil
	}

	inv, err := config.LoadInventory(cli.Inventory)
	if err != nil {
		return nil, err
	}

	var configs []listenerConfig
	port := cli.StartingPort
	for _, entry := range inv.Devices {
		if entry.Scenario != "" {
			fd, seq, err := fakedevices.InitScenario(entry.Scenario, tm, baseDir)
			if err != nil {
				return nil, err
			}
			for i := 0; i < entry.Count; i++ {
				configs = append(configs, listenerConfig{fd, port, seq})
				port++
			}
		} else {
			fd, err := fakedevices.InitGeneric(entry.Platform, tm, baseDir)
			if err != nil {
				return nil, err
			}
			for i := 0; i < entry.Count; i++ {
				configs = append(configs, listenerConfig{fd, port, nil})
				port++
			}
		}
	}
	return configs, nil
}

func run(ctx context.Context, cli config.CLI) error {
	tm, err := transcript.Load(cli.TranscriptMap)
	if err != nil {
		return err
	}

	baseDir := filepath.Dir(cli.TranscriptMap)
	if err := transcript.Validate(tm, baseDir); err != nil {
		return err
	}

	configs, err := resolveListeners(cli, tm, baseDir)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, config := range configs {
		wg.Add(1)
		go func(cfg listenerConfig) {
			defer wg.Done()
			var err error
			if cfg.sequence != nil {
				err = sshlisteners.ScenarioListener(ctx, cfg.fd, cfg.sequence, cfg.port)
			} else {
				err = sshlisteners.GenericListener(ctx, cfg.fd, cfg.port, handlers.GenericCiscoHandler)
			}
			if err != nil {
				log.Printf("listener on port %d: %v", cfg.port, err)
			}
		}(config)
	}

	wg.Wait()
	return nil
}

func main() { // coverage-ignore
	var cli config.CLI
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
