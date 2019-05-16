package cogs

import "strings"

func ping(ctx *Context) {
	ctx.Send("Pong!")
}

func help(ctx *Context) {
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

	ctx.Send(ctx.CodeBlock(msg, "asciidoc"))
}

func init() {
	cog := NewCog("General", "", false)
	cog.AddCommand("ping", "Pong!", []string{}, ping)
	cog.AddCommand("help", "Shows this message", []string{}, help)
	cog.Load()
}
