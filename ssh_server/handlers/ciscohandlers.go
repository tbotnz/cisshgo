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
			match, matchedCommand, multipleMatches, _ := cmdmatch.Match(cmd, myFakeDevice.SupportedCommands)
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

	// Check for context switching commands
	matchPrompt, matchedPrompt, multiplePromptMatches, err := cmdmatch.Match(userInput, fd.ContextSearch)
	if err != nil {
		log.Println(err) // coverage-ignore // CmdMatch never returns errors
		return true
	}

	if matchPrompt && !multiplePromptMatches {
		t.SetPrompt(devicePrompt(fd, fd.ContextSearch[matchedPrompt]))
		*contextState = fd.ContextSearch[matchedPrompt]
		return false
	}

	if userInput == "exit" || userInput == "end" {
		if fd.ContextHierarchy[*contextState] == "exit" {
			return true
		}
		t.SetPrompt(devicePrompt(fd, fd.ContextHierarchy[*contextState]))
		*contextState = fd.ContextHierarchy[*contextState]
		return false
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

	// Match against supported commands
	return dispatchCommand(t, userInput, fd, sequence, seqIdx)
}

// dispatchCommand matches userInput against the active sequence step first, then supported commands.
// Returns true if the session should be terminated.
func dispatchCommand(t *term.Terminal, userInput string, fd *fakedevices.FakeDevice, sequence []transcript.SequenceStep, seqIdx *int) bool {
	// Check if the next sequence step matches
	if seqIdx != nil && *seqIdx < len(sequence) {
		step := sequence[*seqIdx]
		match, _, multipleMatches, _ := cmdmatch.Match(userInput, map[string]string{step.Command: ""})
		if match && !multipleMatches {
			output, err := fakedevices.TranscriptReader(step.Transcript, fd)
			if err != nil {
				log.Println(err)
				return true
			}
			t.Write(append([]byte(output), '\n'))
			*seqIdx++
			return false
		}
	}

	match, matchedCommand, multipleMatches, err := cmdmatch.Match(userInput, fd.SupportedCommands)
	if err != nil {
		log.Println(err) // coverage-ignore // CmdMatch never returns errors
		return true
	}

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
