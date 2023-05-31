package commands

import (
	"os/exec"
	"strings"
)

const SpaceWithinQuote = "{SPACE_WITHIN_QUOTE}"
const DoubleQuote = "{DOUBLE_QUOTE}"

// Create creates a command from a single string
// This allows you to pass parameters which include spaces as commands.
// You just need to add "double-quotes" around the parameter, and it will be treated as one parameter and not be split by whitespace.
func Create(plainCommand string) *exec.Cmd {
	/*
		Example
		     plainCommand   : az role assignment create --role "Network Contributor" --assignee ABC --scope abc
		     escapedCommand : az role assignment create --role "Network{SPACE_WITHIN_QUOTE}Contributor" --assignee ABC --scope abc
		     commandParts   : ["az", "role", "assignment", "create", "--role", "Network Contributor", "--assignee", "ABC", "--scope", "abc"]
	*/
	escapedCommand := markQuotedSpaces(plainCommand)
	commandParts := handleSpacesInQuotes(strings.Fields(escapedCommand))

	var cmd *exec.Cmd

	if len(commandParts) > 1 {
		cmd = exec.Command(commandParts[0], commandParts[1:]...)
	} else {
		cmd = exec.Command(commandParts[0])
	}

	return cmd
}

func markQuotedSpaces(plainCommand string) string {
	escapeMode := false
	var escapedCommand strings.Builder
	enteredEscapeModeInThisIteration := false

	// using previous with a variable instead of index, since in case of index 0 we would run out of bounds, and code is cleaner in this was
	previousCharacter := ""

	for _, char := range plainCommand {
		currentCharacter := string(char)

		// Enter Escape Mode when " occurs, but only if preceded with space, which signals to us that it is an argument which
		// requires this special handling. Example:
		// correct quotes to escape: command -a "some arg"
		// incorrect quotes to escape: command -a object["property"] <- in this case the whole object["property"] should be left unchanged!
		if !escapeMode && currentCharacter == "\"" && previousCharacter == " " {
			escapeMode = true
			enteredEscapeModeInThisIteration = true
		}

		// Handle spaces when in Escape Mode
		if escapeMode && currentCharacter == " " {
			escapedCommand.WriteString(SpaceWithinQuote)
		} else if escapeMode && currentCharacter == "\"" {
			escapedCommand.WriteString(DoubleQuote)
		} else {
			escapedCommand.WriteRune(char)
		}

		// Exit Escape Mode when " occurs
		if !enteredEscapeModeInThisIteration && escapeMode && currentCharacter == "\"" {
			escapeMode = false
		}

		enteredEscapeModeInThisIteration = false
		previousCharacter = currentCharacter
	}

	return escapedCommand.String()
}

func handleSpacesInQuotes(parts []string) []string {
	handledParts := make([]string, len(parts))

	for pos, part := range parts {
		newPartWithSpaces := strings.Replace(part, SpaceWithinQuote, " ", -1)
		newPartWithoutQuotes := strings.Replace(newPartWithSpaces, DoubleQuote, "", -1)
		handledParts[pos] = newPartWithoutQuotes
	}

	return handledParts
}
