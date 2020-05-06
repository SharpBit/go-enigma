package cogs

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/SharpBit/go-enigma/commands"
	"github.com/SharpBit/go-enigma/utils"
)

func ping(ctx *commands.Context) (err error) {
	_, err = ctx.Send("Pong!")
	return
}

func help(ctx *commands.Context) (err error) {
	var msg string
	owner := ctx.Message.Author.ID == utils.GetConfig("ownerID")

	var MaxSig int
	for name := range commands.CommandMap {
		if len(name) > MaxSig {
			MaxSig = len(name)
		}
	}
	MaxSig += 2

	for name, cog := range commands.CogMap {
		msg += "= " + name + " =\n"
		for _, command := range cog.Commands {
			for _, c := range command.Checks {
				// Do not expose owner only commands to someone who isn't the owner
				if runtime.FuncForPC(reflect.ValueOf(c).Pointer()).Name() == "OwnerOnly" && !owner {
					continue
				}
			}
			msg += command.Name + strings.Repeat(" ", MaxSig-len(command.Name)) + ":: " + command.Description + "\n"
		}
		msg += "\n"
	}

	_, err = ctx.Send(ctx.CodeBlock(msg, "asciidoc"))
	return err
}

func init() {
	cog := commands.NewCog("General", "General commands")
	cog.AddCommand("ping", "Pong!", "", ping)
	cog.AddCommand("help", "Shows this message", "", help)
	cog.Load()
}
