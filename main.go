package main

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/suborbital/framework-muxer-showdown/app"
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	appLogger := log.With().Str("component", "app").Logger()

	a := app.New(appLogger)

	a.Start()

	time.Sleep(3 * time.Second)
	a.Stop()
}
