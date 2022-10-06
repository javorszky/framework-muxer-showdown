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

	a := app.New(appLogger, errchan)

	go func() {
		errchan <- a.Start()
	}()

	select {
	case sig := <-shutdownChan:
		appLogger.Error().Msg("doing the shutdown because shutdownchan did a thing")
		a.Stop(sig.String())
		os.Exit(0)

	case err := <-errchan:
		appLogger.Error().Msg("doing the shutdown because errchan did a thing")
		a.Stop(err.Error())
		os.Exit(1)
	}

}
