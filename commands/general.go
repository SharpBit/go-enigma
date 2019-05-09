package commands

func ping(ctx *Context) {
	ctx.Send("Pong!")
}

func init() {
	cog := NewCog("General", "", false)
	cog.AddCommand("ping", "Pong!", []string{}, ping)
	cog.Load()
}
