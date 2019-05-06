package utils

import "../commands"

// CommandMap is a map that gets the user's command input and retrieves its respective function
var CommandMap = map[string]interface{}{
	"ping": commands.Ping,
}
