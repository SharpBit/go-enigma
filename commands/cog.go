package commands

import (
	"fmt"
	"reflect"
)

// CogMap finds the Cog object from the name
var CogMap = make(map[string]*Cog)

/*
Cog structs and functions
*/

// Cog is a similar class to commands.Cog in discord.py
type Cog struct {
	Name        string
	Description string
	Dev         bool
	Commands    []*Command
	Loaded      bool
}

// NewCog creates a new cog instance
func NewCog(name, description string, dev bool) *Cog {
	return &Cog{Name: name, Description: description, Dev: dev, Loaded: false}
}

// AddCommand : Adds a command to the cog
func (cog *Cog) AddCommand(name, description string, usage string, run interface{}) *Command {
	cmd, existing := NewCommand(name, description)
	if existing {
		fmt.Println("error: command " + name + " already exists")
	}
	cmd.Run = reflect.ValueOf(run)
	if cog.Dev == true {
		cmd.Dev = true
	}
	cmd.SetUsage(usage)
	cog.Commands = append(cog.Commands, cmd)

	return cmd
}

// Load : Registers each command in the cog
func (cog *Cog) Load() {
	for _, cmd := range cog.Commands {
		RegisterCommand(cmd)
	}
	CogMap[cog.Name] = cog
	cog.Loaded = true
}

// Unload : Unregisters each command in the cog
func (cog *Cog) Unload() {
	for _, cmd := range cog.Commands {
		UnregisterCommand(cmd)
	}
	delete(CogMap, cog.Name)
	cog.Loaded = false
}
