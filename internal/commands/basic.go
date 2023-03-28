package commands

import (
	"github.com/bwmarrin/discordgo"
)

type BasicCommand struct {
	command *discordgo.ApplicationCommand
}

func NewBasicCommand() *BasicCommand {
	return &BasicCommand{
		command: &discordgo.ApplicationCommand{
			Name:        "basic-cmd",
			Description: "Comando básico",
		},
	}
}

func (c *BasicCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Olá",
		},
	})
}

func (c *BasicCommand) Command() *discordgo.ApplicationCommand {
	return c.command
}

func (c *BasicCommand) CommandName() string {
	return c.command.Name
}
