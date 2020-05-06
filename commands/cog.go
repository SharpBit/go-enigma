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
	Commands    []*Command
	Loaded      bool
}

// NewCog creates a new cog instance
func NewCog(name, description string) *Cog {
	return &Cog{Name: name, Description: description, Loaded: false}
}

// AddCommand : Adds a command to the cog
func (cog *Cog) AddCommand(name, description string, usage string, run interface{}) *Command {
	cmd, existing := NewCommand(name, description)
	if existing {
		panic(fmt.Errorf("CogError: command/alias " + name + " already exists"))
	}
	cmd.Run = reflect.ValueOf(run)
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
