package main

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

func main() {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		With().Timestamp().Logger()

	if err := run(&log); err != nil {
		log.Fatal().Err(err).Msg("Fatal error.")
	}
}

func run(log *zerolog.Logger) error {
	log.Info().Msg("Starting up bot.")
	log.Info().Msg("Shutting down bot.")
	return nil
}
