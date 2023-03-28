package main

import (
	"fmt"
	conf "github.com/ardanlabs/conf/v3"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"time"
)

type configuration struct {
	Token string `conf:"help:'Token used to authenticate the bot.',required,noprint"`
}

func (c *configuration) getToken() string {
	return fmt.Sprintf("Bot %s", c.Token)
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

	if err := run(&log, cfg); err != nil {
		log.Fatal().Err(err).Msg("Fatal error.")
	}
}

func run(log *zerolog.Logger, c *configuration) error {
	log.Info().Msg("Starting up bot.")

	log.Info().Msg("Creating session.")
	session, err := discordgo.New(c.getToken())
	if err != nil {
		return errors.Wrap(err, "Unable to create session.")
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Info().Str("username", s.State.User.Username).Str("discriminator", s.State.User.Discriminator).Msg("Session logged in.")
	})

	log.Info().Msg("Starting session.")
	err = session.Open()
	if err != nil {
		return errors.Wrap(err, "Unable to open session.")
	}

	defer func() {
		log.Info().Msg("Closing session.")
		session.Close()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	log.Info().Msg("Press ctrl+c to stop.")
	<-stop

	log.Info().Msg("Shutting down bot.")
	return nil
}
