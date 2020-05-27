package cogs

import (
	"fmt"
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

func help(ctx *commands.Context, command ...string) (err error) {
	isOwner := ctx.Message.Author.ID == utils.GetConfig("ownerID")
	cmd := strings.Join(command, " ")
	if cmd != "" {
		cog, ok := commands.CogMap[strings.Title(cmd)]
		em := commands.NewEmbed()
		if !ok {
			command, ok := commands.CommandMap[cmd]
			if !ok {
				alias, ok := commands.AliasMap[cmd]
				if !ok {
					ctx.Send("No command found.")
				}

				command, _ = commands.CommandMap[alias]
			}

			if command.Usage != "" {
				em.SetTitle(fmt.Sprintf("`Usage: %s%s %s`", ctx.Prefix, command.Name, command.Usage))
			} else {
				em.SetTitle(fmt.Sprintf("`Usage: %s%s`", ctx.Prefix, command.Name))
			}
			em.SetDescription(command.Description)
			em.SetColor(0x2ecc71)
		} else {
			var commandHelp string
			var MaxSig int
			for _, c := range cog.Commands {
				if len(c.Name) > MaxSig {
					MaxSig = len(c.Name)
				}
			}
			MaxSig += 2
			for _, command := range cog.Commands {
				NotAuthorized := false
				for _, c := range command.Checks {
					// Do not expose owner only commands to someone who isn't the owner
					if runtime.FuncForPC(reflect.ValueOf(c).Pointer()).Name() == "OwnerOnly" && !isOwner {
						NotAuthorized = true
					}
				}
				if NotAuthorized == true {
					continue
				}
				commandHelp += "`" + ctx.Prefix + command.Name + strings.Repeat(" ", MaxSig-len(command.Name)) + command.Description + "`\n"
			}

			em.SetTitle(cog.Name)
			em.SetDescription(fmt.Sprintf("*%s*", cog.Description))
			em.SetColor(0x2ecc71)
			em.AddField("Commands", commandHelp)
		}

		_, err := ctx.SendComplex("", em.MessageEmbed)
		return err
	}

	paginator := utils.NewPaginatorForContext(ctx)
	paginator.SetTemplate(func() *commands.Embed {
		return commands.NewEmbed().
			SetColor(0x2ecc71).
			SetFooter(fmt.Sprintf("Type %shelp command for more help on a command.", ctx.Prefix))
	})
	for name, cog := range commands.CogMap {
		var commandHelp string
		var MaxSig int
		for _, c := range cog.Commands {
			if len(c.Name) > MaxSig {
				MaxSig = len(c.Name)
			}
		}
		MaxSig += 2
		for _, command := range cog.Commands {
			NotAuthorized := false
			for _, c := range command.Checks {
				// Do not expose owner only commands to someone who isn't the owner
				if runtime.FuncForPC(reflect.ValueOf(c).Pointer()).Name() == "OwnerOnly" && !isOwner {
					NotAuthorized = true
				}
			}
			if NotAuthorized == true {
				continue
			}
			commandHelp += "`" + ctx.Prefix + command.Name + strings.Repeat(" ", MaxSig-len(command.Name)) + command.Description + "`\n"
		}

		paginator.AddPage(func(em *commands.Embed) *commands.Embed {
			em.SetTitle(name)
			em.SetDescription(fmt.Sprintf("*%s*", cog.Description))
			em.AddField("Commands", commandHelp)
			return em
		})

	}

	paginator.Run()
	return
}

func init() {
	cog := commands.NewCog("General", "General commands")
	cog.AddCommand("ping", "Pong!", "", ping)
	cog.AddCommand("help", "Shows a list of commands", "[command]", help).SetDefaultArg("")
	cog.Load()
}
