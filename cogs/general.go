package cogs

import "strings"

func ping(ctx *Context) (err error) {
	_, err = ctx.Send("Pong!")
	return
}

func help(ctx *Context) (err error) {
	var msg string
	owner := ctx.Message.Author.ID == "281821029490229251"

	var MaxSig int
	for name := range CommandMap {
		if len(name) > MaxSig {
			MaxSig = len(name)
		}
	}
	MaxSig += 2

	for name, cog := range CogMap {
		if cog.Dev && !owner {
			continue
		}
		msg += "= " + name + " =\n"
		for _, command := range cog.Commands {
			if command.Dev && !owner {
				continue
			}
			msg += command.Name + strings.Repeat(" ", MaxSig-len(command.Name)) + ":: " + command.Description + "\n"
		}
		msg += "\n"
	}

	_, err = ctx.Send(ctx.CodeBlock(msg, "asciidoc"))
	return err
}

func init() {
	cog := NewCog("General", "", false)
	cog.AddCommand("ping", "Pong!", "", ping)
	cog.AddCommand("help", "Shows this message", "", help)
	cog.Load()
}
