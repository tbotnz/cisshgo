package utils

import (
	"fmt"
	"strings"
)

// CmdMatch searches the provided supportedCommands to find a match for the provided userInput
func CmdMatch(userInput string, supportedCommands map[string]string) (bool, string) {

	// Setup our return variables
	match := false
	matchedCmd := ""

	// Setup a Slice to hold any possibleMatches
	possibleMatches := make([]string, len(supportedCommands))

	// Turn our input string into fields
	inputFields := strings.Fields(userInput)

	// Iterate through all the keys in the SupportedCommands map
	for k := range supportedCommands {
		commandFields := strings.Fields(k)

		// Match field by field to discount any non matches
		for _, inputField := range inputFields {
			for _, commandField := range commandFields {
				if strings.Contains(commandField, inputField) {
					fmt.Printf("supportedCommand: %s\n", k)
					possibleMatches = append(possibleMatches, k)
					match = true
				}
			}
		}
	}

	// Iterate through all possibleMatches to find the best match
	for _, possibleMatch := range possibleMatches {
		fmt.Printf("possibleMatch: %s\n", possibleMatch)
	}

	return match, matchedCmd
}
