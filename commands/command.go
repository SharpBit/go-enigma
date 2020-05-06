package commands

import (
	"fmt"
	"reflect"
)

/*
Command Structs and Functions
*/

// CommandMap is a map that gets the user's command input and retrieves its respective function
var CommandMap = make(map[string]*Command)

// AliasMap finds the command of each alias
var AliasMap = make(map[string]string)

// Command is a command object
type Command struct {
	Name           string
	Description    string
	Aliases        []string
	Checks         []func(*Context) (bool, error)
	Usage          string
	Run            reflect.Value
	HasOptionalArg bool
	DefaultArg     reflect.Value
}

func (cmd *Command) Type() reflect.Type {
	return reflect.TypeOf(cmd.Run)
}

func (cmd *Command) SetUsage(usage string) *Command {
	cmd.Usage = usage
	return cmd
}

func (cmd *Command) SetAliases(aliases ...string) *Command {
	for _, alias := range aliases {
		// Check to see if the alias is already registered as a command or alias
		_, existing := AliasMap[alias]
		if existing {
			panic(fmt.Errorf("CogError: alias " + alias + " already exists"))
		}

		_, existing = CommandMap[alias]
		if existing {
			panic(fmt.Errorf("CogError: command " + alias + " already exists"))
		}

		// Not registered yet, add it to the aliases
		cmd.Aliases = append(cmd.Aliases, alias)
	}
	return cmd
}

// SetDefaultArg: Makes the last argument optional and will pass in a default argument if not provided
func (cmd *Command) SetDefaultArg(def interface{}) *Command {
	cmd.DefaultArg = reflect.ValueOf(def)
	cmd.HasOptionalArg = true
	return cmd
}

// AddCheck: Adds a check to the command to be run before command invokation. Must return (bool, err)
func (cmd *Command) AddCheck(check func(*Context) (bool, error)) *Command {
	cmd.Checks = append(cmd.Checks, check)
	return cmd
}

// RemoveChecks: Removes all checks and returns the command and check function
func (cmd *Command) RemoveChecks() (*Command, []func(*Context) (bool, error)) {
	checks := cmd.Checks
	cmd.Checks = []func(*Context) (bool, error){}
	return cmd, checks
}

// NewCommand creates a new command
func NewCommand(name, description string) (cmd *Command, existing bool) {
	_, existing = CommandMap[name]
	if existing {
		return nil, existing
	}

	_, existing = AliasMap[name]
	if existing {
		return nil, existing
	}

	cmd = &Command{Name: name, Description: description}
	return cmd, existing
}

// RegisterCommand adds the command to the CommandMap
func RegisterCommand(cmd *Command) {
	CommandMap[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		AliasMap[alias] = cmd.Name
	}
}

// UnregisterCommand removes the command from the CommandMap
func UnregisterCommand(cmd *Command) {
	delete(CommandMap, cmd.Name)
	for _, alias := range cmd.Aliases {
		delete(AliasMap, alias)
	}
}
