package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Command() *discordgo.ApplicationCommand
	CommandName() string
	Handle(*discordgo.Session, *discordgo.InteractionCreate)
}

type CommandList map[string]Command

func NewCommandList() CommandList {
	return make(CommandList)
}

func (list CommandList) Append(cmd Command) {
	list[cmd.CommandName()] = cmd
}

func (list CommandList) GetCommand(cmdName string) (Command, bool) {
	cmd, ok := list[cmdName]
	return cmd, ok
}
