package rules

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"reflect"
	"time"
)

type RulesCommand struct {
	log     *zerolog.Logger
	command *discordgo.ApplicationCommand
	Rules   []Rule
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
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "n",
							Description: "Número da regra a mostrar.",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    false,
						},
						{
							Name:        "ping",
							Description: "Ping",
							Type:        discordgo.ApplicationCommandOptionMentionable,
							Required:    false,
						},
					},
				},
				{
					Name:        "add",
					Description: "Adiciona um nova regra.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "text",
							Description: "Texto da regra.",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "remove",
					Description: "Remove uma regra.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "n",
							Description: "Número da regra a retirar.",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
					},
				},
				{
					Name:        "move",
					Description: "Move uma regra.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "n",
							Description: "Número da regra a ser movida.",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
						{
							Name:        "to",
							Description: "Para que posição mover a regra.",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
					},
				},
				{
					Name:        "edit",
					Description: "Edita um comando.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "n",
							Description: "Número da regra a ser movida.",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
						{
							Name:        "text",
							Description: "Texto da regra.",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
		Rules: []Rule{Rule{Text: "Regra 1"}, Rule{Text: "Regra 2"}, Rule{Text: "Regra 3"}},
	}
}

func (c *RulesCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	c.log.Debug().Msgf("%v", options)

	switch options[0].Name {
	case "show":
		c.handleShow(s, i)
		return
	case "add":
		c.handleAdd(s, i)
		return
	case "remove":
		c.handleRemove(s, i)
		return
	case "move":
		c.handleMove(s, i)
		return
	case "edit":
		c.handleEdit(s, i)
		return
	}
}

func (c *RulesCommand) Command() *discordgo.ApplicationCommand {
	return c.command
}

func (c *RulesCommand) CommandName() string {
	return c.command.Name
}

func (c *RulesCommand) handleShow(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	subOptions := options[0].Options

	optionsMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(subOptions))
	for _, opt := range subOptions {
		optionsMap[opt.Name] = opt
	}

	content := ""
	if pingOpt, ok := optionsMap["ping"]; ok {
		c.log.Debug().Msgf("%s", pingOpt.UserValue(nil).ID)
		c.log.Debug().Msgf("%v", reflect.TypeOf(pingOpt))

		if user := pingOpt.UserValue(nil); user != nil {
			_, err := s.User(user.ID)

			if role := pingOpt.RoleValue(nil, i.GuildID); err != nil && role != nil {
				content = fmt.Sprintf("<@&%s>", role.ID)
			} else {
				content = fmt.Sprintf("<@!%s>", user.ID)
			}
		}
	}

	if nOpt, ok := optionsMap["n"]; ok {
		n := int(nOpt.IntValue())
		index := 0

		if n == -1 {
			index = len(c.Rules) - 1
		} else {
			index = n - 1
		}

		if index < 0 || index > (len(c.Rules)-1) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Número de regra inválida.",
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
				Embeds: []*discordgo.MessageEmbed{
					c.createRuleEmbed(index, s),
				},
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Embeds: []*discordgo.MessageEmbed{
				c.createRulesEmbed(),
			},
		},
	})
	return
}

func (c *RulesCommand) handleAdd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	text := options[0].Options[0].StringValue()

	userId := ""
	if i.Member != nil {
		if i.Member.User != nil {
			userId = i.Member.User.ID
		}
	}

	rule := NewRule(text, userId, time.Now())

	c.Rules = append(c.Rules, rule)

	c.log.Debug().Msgf("Added new rule: %v", rule)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				c.createRuleEmbed((len(c.Rules)), s),
			},
		},
	})
}

func (c *RulesCommand) handleRemove(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	index := options[0].Options[0].IntValue()

	if index < 1 || index > int64(len(c.Rules)) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Número de regra inválida.",
			},
		})
		return
	}

	c.Rules = append(c.Rules[:(index-1)], c.Rules[(index):]...)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				c.createRulesEmbed(),
			},
		},
	})
	return
}

func (c *RulesCommand) handleMove(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options

	optionsMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionsMap[opt.Name] = opt
	}

	index := optionsMap["n"].IntValue()
	dest := optionsMap["to"].IntValue()

	if dest == -1 {
		dest = int64(len(c.Rules))
	}

	if index == dest {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					c.createRulesEmbed(),
				},
			},
		})
	}

	if dest < 1 || dest > int64(len(c.Rules)) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Destino da regra inválida.",
			},
		})
		return
	}

	if index < 1 || index > int64(len(c.Rules)) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Número de regra inválida.",
			},
		})
		return
	}

	rule := c.Rules[index-1]

	c.log.Debug().Msgf("%v", c.Rules)
	c.Rules = append(c.Rules[:(index-1)], c.Rules[(index):]...)
	c.log.Debug().Msgf("%v", c.Rules)

	rulesBefore := c.Rules[:dest-1]
	rulesAfter := c.Rules[dest-1:]

	c.Rules = append(rulesBefore, rule)
	c.log.Debug().Msgf("%v", c.Rules)
	c.Rules = append(c.Rules, rulesAfter...)
	c.log.Debug().Msgf("%v", c.Rules)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				c.createRulesEmbed(),
			},
		},
	})
	return
}

func (c *RulesCommand) handleEdit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options

	optionsMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionsMap[opt.Name] = opt
	}

	index := optionsMap["n"].IntValue()
	text := optionsMap["text"].StringValue()

	if index < 1 || index > int64(len(c.Rules)) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Número de regra inválida.",
			},
		})
		return
	}

	userId := ""
	if i.Member != nil {
		if i.Member.User != nil {
			userId = i.Member.User.ID
		}
	}

	c.Rules[index-1].Edit(text, userId, time.Now())

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				c.createRulesEmbed(),
			},
		},
	})
	return
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

func (c *RulesCommand) createRuleEmbed(index int, s *discordgo.Session) *discordgo.MessageEmbed {
	rule := c.Rules[index]

	footer := &discordgo.MessageEmbedFooter{}

	author, err := s.User(rule.creator)
	if err != nil {
		c.log.Debug().Err(err).Msg("Unable to get rule creator.")
	}
	lastuser, err := s.User(rule.lastChangeBy)
	if err != nil {
		c.log.Debug().Err(err).Msg("Unable to get rule last modifier user.")
	}

	if author != nil {
		footer.Text = fmt.Sprintf("Criado por: %v.\nCriado em: %v.\n", author.Username, rule.createdAt)
	} else {
		footer.Text = fmt.Sprintf("Criado em: %v.\n", rule.createdAt)
	}

	if lastuser != nil {
		footer.Text = fmt.Sprintf(
			"%sÚltima modificação por: %v.\nÚltima modificação em: %v.\n",
			footer.Text,
			lastuser.Username, rule.lastChangeAt,
		)
	} else {
		footer.Text = fmt.Sprintf(
			"%sÚltima modificação em: %v.\n",
			footer.Text, rule.lastChangeAt,
		)
	}

	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Regra N°%v", index+1),
		Description: rule.Text,
		Footer:      footer,
	}
}
