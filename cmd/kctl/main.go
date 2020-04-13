package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/milosgajdos/kraph/cmd/kctl/app"
)

func main() {
	sigChan := setupSigHandler()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		signal.Stop(sigChan)
		cancel()
	}()

	go func() {
		select {
		case <-sigChan: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-sigChan // second signal, hard exit
		os.Exit(1)
	}()

	app := app.New()

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// setupSigHandler makes signal handler for catching os.Interrupt and returns it
func setupSigHandler() chan os.Signal {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	return signalChan
}
