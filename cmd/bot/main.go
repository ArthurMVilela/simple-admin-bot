package main

import (
	"fmt"
	conf "github.com/ardanlabs/conf/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type configuration struct {
	Token string `conf:"help:'Token used to authenticate the bot.',required,noprint"`
}

func main() {
	cfg := new(configuration)
	if help, err := conf.Parse("", cfg); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			os.Exit(0)
		}
		fmt.Printf("Invalid configuration: %v\n", err)
		os.Exit(1)
	}

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
