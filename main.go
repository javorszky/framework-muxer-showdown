package main

import (
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/suborbital/framework-muxer-showdown/app"
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	appLogger := zerolog.New(os.Stderr).With().Str("component", "app").Logger()

	shutdownChan := make(chan os.Signal, 1)
	errchan := make(chan error)

	signal.Notify(shutdownChan, os.Kill, os.Interrupt)

	a := app.New(appLogger)

	go func() {
		errchan <- a.Start()
	}()

	select {
	case sig := <-shutdownChan:
		a.Stop(sig.String())

	case err := <-errchan:
		a.Stop(err.Error())
	}
}
