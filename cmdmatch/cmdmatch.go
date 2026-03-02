// Package cmdmatch provides fuzzy command matching for cisshgo.
package cmdmatch

import (
	"log"
	"strings"
)

// Match searches supportedCommands for a prefix match against userInput.
// Each word in userInput must be a prefix of the corresponding word in the command.
// Returns (matched, matchedCommand, multipleMatches).
func Match(userInput string, supportedCommands map[string]string) (bool, string, bool) {
	userInput = strings.ToLower(strings.TrimSpace(userInput))
	if userInput == "" {
		return false, "", false
	}
	userFields := strings.Fields(userInput)

	var matches []string
	for cmd := range supportedCommands {
		cmdFields := strings.Fields(strings.ToLower(cmd))
		if len(cmdFields) != len(userFields) {
			continue
		}
		if prefixMatch(userFields, cmdFields) {
			matches = append(matches, cmd)
		}
	}

	switch len(matches) {
	case 0:
		return false, "", false
	case 1:
		return true, matches[0], false
	default:
		log.Printf("ambiguous command %q matches: %v", userInput, matches)
		return true, "", true
	}
}

// prefixMatch returns true if every userField is a prefix of the corresponding cmdField.
func prefixMatch(userFields, cmdFields []string) bool {
	for i, f := range userFields {
		if !strings.HasPrefix(cmdFields[i], f) {
			return false
		}
	}
	return true
}
