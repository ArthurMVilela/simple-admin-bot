package commands

import (
	"github.com/bwmarrin/discordgo"
)

func OptionMap(
	options []*discordgo.ApplicationCommandInteractionDataOption,
) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optMap[opt.Name] = opt
	}

	return optMap
}
