package commands

import "fmt"

func ping(ctx *Context) {
	ctx.Send("Pong!")
}

func init() {
	cmd, existing := NewCommand("ping", "Pong!")
	if existing {
		fmt.Println("error: command ping already exists")
		return
	}
	cmd.Run = ping
	cmd.Aliases = []string{}
	RegisterCommand(cmd)
}
