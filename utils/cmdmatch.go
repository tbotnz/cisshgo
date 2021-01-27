package utils

import (
	"fmt"
	"strings"
)

// CmdMatch searches the provided supportedCommands to find a match for the provided userInput
// Returns:
//	match: bool
// 	matchedCommand: string
//  multipleMatches: bool
//	error
func CmdMatch(userInput string, supportedCommands map[string]string) (bool, string, bool, error) {

	// Setup our return variables
	match := false
	matchedCmd := ""
	multipleMatches := false

	// Setup a Map to hold any possibleMatches as keys, and the string.Fields as values
	possibleMatches := make(map[string][]string)

	// Turn our input string into fields
	// fmt.Printf("userInput: %s\n", userInput)
	userInput = strings.ToLower(userInput) // Lowercase the user input
	userInputFields := strings.Fields(userInput)

	// Iterate through all the commands in the supportedCommands map
	for supportedCommand := range supportedCommands {
		supportedCommand := strings.ToLower(supportedCommand) // Lowercase our supported command
		commandFields := strings.Fields(supportedCommand)

		// Match against the 1st field in each command,
		// and that the number of fields is the same,
		// to find any possibleMatches.
		if strings.Contains(commandFields[0], userInputFields[0]) &&
			(len(commandFields) == len(userInputFields)) {
			// fmt.Printf("supportedCommand: %s\n", k)
			possibleMatches[supportedCommand] = commandFields
		}
	}

	// Setup a map to hold our closestMatch(es)
	closestMatch := make(map[string]struct{})

	// Iterate through all possibleMatches to find the best match
	// fmt.Printf("possibleMatches: %+v\n", possibleMatches)
	for possibleMatch := range possibleMatches {

		// First evaluate if we have an exact string match and break/return that
		if userInput == possibleMatch {
			closestMatch[possibleMatch] = struct{}{}
			break
		}

		// Next, test if the entire input is contained within one of our commands
		if strings.Contains(possibleMatch, userInput) {
			closestMatch[possibleMatch] = struct{}{}
			break
		}

		// Next delve into the fields and find best match
		for p, possibleMatchField := range possibleMatches[possibleMatch] {
			// fmt.Printf("possibleMatchField: %s\n", possibleMatchField)
			if !strings.Contains(possibleMatchField, userInputFields[p]) {
				// We did not get a match on this field, break
				break
			}
			// fmt.Printf("%d\n", p)
			// fmt.Printf("length of possibleMatch fields: %d\n", len(possibleMatches[possibleMatch]))
			if p == (len(possibleMatches[possibleMatch]) - 1) {
				closestMatch[possibleMatch] = struct{}{}
			}
		}
	}

	// Evaluate our closestMatch(es)
	if len(closestMatch) > 1 {
		// We had more than two matches to all conditions, return no match!
		fmt.Printf("multiple matchedCmds: %s\n", closestMatch)
		match = true
		matchedCmd = ""
		multipleMatches = true
	} else if len(closestMatch) < 1 {
		// We had _NO_ matches to any conditions, return no match!
		// fmt.Printf("no matchedCmds\n")
		match = false
		matchedCmd = ""
	} else {
		match = true
		for k := range closestMatch {
			matchedCmd = k
		}

	}

	// fmt.Printf("matchedCmd: %s\n\n", matchedCmd)
	return match, matchedCmd, multipleMatches, nil
}
