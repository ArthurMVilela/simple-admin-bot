package commands

import (
	"github.com/bwmarrin/discordgo"
)

type BasicCommand struct {
	Command *discordgo.ApplicationCommand
}

func (c *BasicCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Olá",
		},
	})
}

func NewBasicCommand() *BasicCommand {
	return &BasicCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "basic-cmd",
			Description: "Comando básico",
		},
	}
}
