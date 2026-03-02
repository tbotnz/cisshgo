// Package handlers contains SSH Handlers for specific device types
// in order to best emulate their actual behavior.
package handlers

import (
	"io"
	"log"
	"strings"

	"github.com/gliderlabs/ssh"
	"golang.org/x/term"

	"github.com/tbotnz/cisshgo/cmdmatch"
	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/transcript"
)

// GenericCiscoHandler function handles generic Cisco style sessions
func GenericCiscoHandler(myFakeDevice *fakedevices.FakeDevice) ssh.Handler {
	return genericCiscoSession(myFakeDevice, nil)
}

// GenericCiscoScenarioHandler returns an ssh.Handler that plays back a scenario sequence.
func GenericCiscoScenarioHandler(myFakeDevice *fakedevices.FakeDevice, sequence []transcript.SequenceStep) ssh.Handler {
	return genericCiscoSession(myFakeDevice, sequence)
}

func genericCiscoSession(myFakeDevice *fakedevices.FakeDevice, sequence []transcript.SequenceStep) ssh.Handler {
	return func(s ssh.Session) {

		// Exec mode: client sent a command directly (e.g., ssh host "show version")
		if cmd := s.RawCommand(); cmd != "" {
			log.Printf("exec: %s", cmd)
			match, matchedCommand, multipleMatches := cmdmatch.Match(cmd, myFakeDevice.SupportedCommands)
			if match && !multipleMatches {
				output, err := fakedevices.TranscriptReader(
					myFakeDevice.SupportedCommands[matchedCommand], myFakeDevice,
				)
				if err == nil {
					io.WriteString(s, output)
				}
			}
			s.Exit(0)
			return
		}

		// Interactive shell mode — sequence pointer resets per session
		seqIdx := 0
		contextState := myFakeDevice.ContextSearch["base"]
		t := term.NewTerminal(s, devicePrompt(myFakeDevice, contextState))

		for {
			userInput, err := t.ReadLine()
			if err != nil {
				break
			}
			log.Println(userInput)

			done := handleShellInput(t, userInput, myFakeDevice, &contextState, sequence, &seqIdx)
			if done {
				break
			}
		}
		log.Println("terminal closed")
	}
}

// handleShellInput processes a single line of user input in interactive shell mode.
// Returns true if the session should be terminated.
func handleShellInput(t *term.Terminal, userInput string, fd *fakedevices.FakeDevice, contextState *string, sequence []transcript.SequenceStep, seqIdx *int) bool {
	if userInput == "" {
		t.Write([]byte(""))
		return false
	}

	// Check sequence step FIRST — if it matches, advance the pointer and write the
	// transcript. Then fall through to also apply any context switch side effect.
	sequenceHandled, terminate := handleSequenceStep(t, userInput, fd, sequence, seqIdx)
	if terminate {
		return true
	}

	// Apply context switch if the input matches a context_search key.
	// Uses starts-with-N-words semantics so "interface Gi0/0/2" matches key "interface".
	// In scenario mode (active sequence), context switches only fire when the sequence
	// step was just handled — enforcing strict command ordering.
	inScenario := seqIdx != nil && *seqIdx < len(sequence)
	if !inScenario || sequenceHandled {
		if matchedCtx, ok := matchContextKey(userInput, fd.ContextSearch); ok {
			t.SetPrompt(devicePrompt(fd, fd.ContextSearch[matchedCtx]))
			*contextState = fd.ContextSearch[matchedCtx]
			return false
		}
	}

	if sequenceHandled {
		return false
	}

	if userInput == "exit" {
		return handleExitEnd(t, userInput, fd, contextState)
	}

	// In scenario mode, "end" is blocked unless it was the current sequence step
	if userInput == "end" && !inScenario {
		return handleExitEnd(t, userInput, fd, contextState)
	}

	if userInput == "reset state" {
		t.Write(append([]byte("Resetting State..."), '\n'))
		*contextState = fd.ContextSearch["base"]
		fd.Hostname = fd.DefaultHostname
		t.SetPrompt(devicePrompt(fd, *contextState))
		return false
	}

	// Handle hostname changes in config mode
	userInputFields := strings.Fields(userInput)
	if userInputFields[0] == "hostname" && *contextState == "(config)#" {
		fd.Hostname = strings.Join(userInputFields[1:], " ")
		log.Printf("Setting hostname to %s\n", fd.Hostname)
		t.SetPrompt(devicePrompt(fd, *contextState))
		return false
	}

	return dispatchCommand(t, userInput, fd)
}

// handleSequenceStep checks if userInput matches the current sequence step.
// Returns (handled, terminate).
func handleSequenceStep(t *term.Terminal, userInput string, fd *fakedevices.FakeDevice, sequence []transcript.SequenceStep, seqIdx *int) (bool, bool) {
	if seqIdx == nil || *seqIdx >= len(sequence) {
		return false, false
	}
	step := sequence[*seqIdx]
	if !matchSequenceStep(userInput, step.Command) {
		return false, false
	}
	output, err := fakedevices.TranscriptReader(step.Transcript, fd)
	if err != nil {
		log.Println(err)
		return false, true
	}
	t.Write(append([]byte(output), '\n'))
	*seqIdx++
	return true, false
}

// handleExitEnd processes exit and end commands.
// Returns true if the session should be terminated.
func handleExitEnd(t *term.Terminal, userInput string, fd *fakedevices.FakeDevice, contextState *string) bool {
	target := fd.ContextHierarchy[*contextState]
	if userInput == "end" && fd.EndContext != "" {
		target = fd.EndContext
	}
	if target == "exit" {
		return true
	}
	t.SetPrompt(devicePrompt(fd, target))
	*contextState = target
	return false
}

// matchContextKey returns the context_search key that the input starts with (word-prefix match).
// The key's words must be a prefix of the input's words; extra input words are allowed.
// Returns the matched key and true, or "" and false if no match or multiple matches.
// Skips the "base" key since it is not a real command.
func matchContextKey(userInput string, contextSearch map[string]string) (string, bool) {
	inputFields := strings.Fields(strings.ToLower(userInput))
	var matches []string
	for key := range contextSearch {
		if key == "base" {
			continue
		}
		keyFields := strings.Fields(strings.ToLower(key))
		if len(inputFields) < len(keyFields) {
			continue
		}
		matched := true
		for i, kf := range keyFields {
			if !strings.HasPrefix(kf, inputFields[i]) {
				matched = false
				break
			}
		}
		if matched {
			matches = append(matches, key)
		}
	}
	if len(matches) == 1 {
		return matches[0], true
	}
	return "", false
}

// dispatchCommand matches userInput against supported commands.
// Returns true if the session should be terminated.
func dispatchCommand(t *term.Terminal, userInput string, fd *fakedevices.FakeDevice) bool {
	match, matchedCommand, multipleMatches := cmdmatch.Match(userInput, fd.SupportedCommands)

	if multipleMatches {
		t.Write(append([]byte("% Ambiguous command:  \""+userInput+"\""), '\n'))
		return false
	}

	if !match {
		t.Write(append([]byte("% Unknown command:  \""+userInput+"\""), '\n'))
		return false
	}

	output, err := fakedevices.TranscriptReader(fd.SupportedCommands[matchedCommand], fd)
	if err != nil {
		log.Println(err)
		return true
	}
	t.Write(append([]byte(output), '\n'))
	return false
}

// matchSequenceStep returns true if userInput matches the sequence step command.
// Uses standard prefix matching for regular words, and for tokens containing a
// letter/digit boundary (e.g. "g0/0/2"), matches the alpha prefix against the
// step's alpha prefix and requires an exact suffix match.
// This allows "int g0/0/2" to match "interface GigabitEthernet0/0/2" without
// a hardcoded abbreviation table.
func matchSequenceStep(userInput, stepCmd string) bool {
	userFields := strings.Fields(strings.ToLower(userInput))
	stepFields := strings.Fields(strings.ToLower(stepCmd))
	if len(userFields) != len(stepFields) {
		return false
	}
	for i, uf := range userFields {
		sf := stepFields[i]
		uAlpha, uSuffix := splitIfaceToken(uf)
		sAlpha, sSuffix := splitIfaceToken(sf)
		if uSuffix != "" || sSuffix != "" {
			// Interface-style token: alpha prefix must match, suffix must be equal
			if !strings.HasPrefix(sAlpha, uAlpha) || uSuffix != sSuffix {
				return false
			}
		} else {
			// Regular word: step word must start with user word (abbreviation)
			if !strings.HasPrefix(sf, uf) {
				return false
			}
		}
	}
	return true
}

// splitIfaceToken splits a token like "gigabitethernet0/0/2" into ("gigabitethernet", "0/0/2").
// Only splits when the suffix starts with a digit (interface-style tokens).
// Returns (token, "") if there is no letter/digit boundary.
func splitIfaceToken(token string) (alpha, suffix string) {
	i := 0
	for i < len(token) && token[i] >= 'a' && token[i] <= 'z' {
		i++
	}
	if i == 0 || i == len(token) || token[i] < '0' || token[i] > '9' {
		return token, ""
	}
	return token[:i], token[i:]
}

// buildPrompt constructs the terminal prompt string.
// If format is empty, falls back to hostname+context (default Cisco style).
// If prefixLine is non-empty, it is prepended above the prompt on its own line.
func buildPrompt(format, hostname, username, context, prefixLine string) string {
	prompt := hostname + context
	if format != "" {
		prompt = strings.NewReplacer(
			"{hostname}", hostname,
			"{username}", username,
			"{context}", context,
		).Replace(format)
	}
	if prefixLine != "" {
		return prefixLine + "\n" + prompt
	}
	return prompt
}

// devicePrompt builds the prompt for a FakeDevice at the given context state.
func devicePrompt(fd *fakedevices.FakeDevice, context string) string {
	return buildPrompt(fd.PromptFormat, fd.Hostname, fd.Username, context, fd.ContextPrefixLines[context])
}
