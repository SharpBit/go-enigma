package handlers

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/SharpBit/go-enigma/commands"
)

// HandleCommandError: Handles and recovers from errors without panicking
func HandleCommandError(ctx *commands.Context, err error) {
	// No command found
	if strings.HasPrefix(err.Error(), "InvokeCommandError") {
		return
	}

	// Guild or channel could not be retrieved
	if ctx == nil {
		panic(err)
	}

	// Something is wrong with the user inputted arguments
	if strings.HasPrefix(err.Error(), "ArgumentError") {
		ctx.SendError(err, true)
		return
	}

	// Command check failed
	if strings.HasPrefix(err.Error(), "CommandCheckError") {
		ctx.SendError(err, false)
		return
	}

	// Insufficient permissions
	if strings.HasPrefix(err.Error(), "HTTP 403 Forbidden") {
		ctx.Send("I do not have the permissions to perform this command.")
		return
	}

	// Log any other errors without panicking
	fmt.Println(err)
	debug.PrintStack()
}
