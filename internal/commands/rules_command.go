package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

type RulesCommand struct {
	log     *zerolog.Logger
	command *discordgo.ApplicationCommand
}

func NewRulesCommand(log *zerolog.Logger) *RulesCommand {
	return &RulesCommand{
		log: log,
		command: &discordgo.ApplicationCommand{
			Name:        "rules",
			Description: "Comando para gerenciar e consutar as regras do servidor.",
		},
	}
}

func (c *RulesCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {

}

func (c *RulesCommand) Command() *discordgo.ApplicationCommand {
	return c.command
}

func (c *RulesCommand) CommandName() string {
	return c.command.Name
}
