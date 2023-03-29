package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

type RulesCommand struct {
	log     *zerolog.Logger
	command *discordgo.ApplicationCommand
	Rules   []Rule
}

type Rule struct {
	Text string
}

func NewRulesCommand(log *zerolog.Logger) *RulesCommand {
	return &RulesCommand{
		log: log,
		command: &discordgo.ApplicationCommand{
			Name:        "rules",
			Description: "Comando para gerenciar e consutar as regras do servidor.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "show",
					Description: "Mostra todas regras.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		Rules: []Rule{Rule{Text: "Ado ado ado."}, Rule{Text: "Cado cado cado."}},
	}
}

func (c *RulesCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	switch options[0].Name {
	case "show":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					c.createRulesEmbed(),
				},
			},
		})
	}
}

func (c *RulesCommand) Command() *discordgo.ApplicationCommand {
	return c.command
}

func (c *RulesCommand) CommandName() string {
	return c.command.Name
}

func (c *RulesCommand) createRulesEmbed() *discordgo.MessageEmbed {
	fields := make([]*discordgo.MessageEmbedField, len(c.Rules))

	for i, rule := range c.Rules {
		fields[i] = &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("Regra N°%v", i+1),
			Value: rule.Text,
		}
	}

	return &discordgo.MessageEmbed{
		Title:       "Regras",
		Description: "Regras do servidor.",
		Fields:      fields,
	}
}
