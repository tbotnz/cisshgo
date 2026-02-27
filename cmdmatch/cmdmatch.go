// Package cmdmatch provides fuzzy command matching for cisshgo.
package cmdmatch

import (
	"fmt"
	"strings"
)

// Match searches the provided supportedCommands to find a match for the provided userInput.
// Returns match, matchedCommand, multipleMatches, error.
func Match(userInput string, supportedCommands map[string]string) (bool, string, bool, error) {
	match := false
	matchedCmd := ""
	multipleMatches := false

	possibleMatches := make(map[string][]string)

	userInput = strings.ToLower(userInput)
	userInputFields := strings.Fields(userInput)

	for supportedCommand := range supportedCommands {
		supportedCommand := strings.ToLower(supportedCommand)
		commandFields := strings.Fields(supportedCommand)

		if strings.Contains(commandFields[0], userInputFields[0]) &&
			(len(commandFields) == len(userInputFields)) {
			possibleMatches[supportedCommand] = commandFields
		}
	}

	closestMatch := make(map[string]struct{})

	for possibleMatch := range possibleMatches {
		if userInput == possibleMatch {
			closestMatch[possibleMatch] = struct{}{}
			break
		}

		if strings.Contains(possibleMatch, userInput) {
			closestMatch[possibleMatch] = struct{}{}
			break
		}

		for p, possibleMatchField := range possibleMatches[possibleMatch] {
			if !strings.Contains(possibleMatchField, userInputFields[p]) {
				break
			}
			if p == (len(possibleMatches[possibleMatch]) - 1) {
				closestMatch[possibleMatch] = struct{}{}
			}
		}
	}

	if len(closestMatch) > 1 {
		fmt.Printf("multiple matchedCmds: %s\n", closestMatch)
		match = true
		matchedCmd = ""
		multipleMatches = true
	} else if len(closestMatch) < 1 {
		match = false
		matchedCmd = ""
	} else {
		match = true
		for k := range closestMatch {
			matchedCmd = k
		}
	}

	return match, matchedCmd, multipleMatches, nil
}
