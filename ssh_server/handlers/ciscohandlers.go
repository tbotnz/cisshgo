// Package handlers contains SSH Handlers for specific device types
// in order to best emulate their actual behavior.
package handlers

import (
	"io"
	"log"
	"strings"

	"github.com/gliderlabs/ssh"
	"golang.org/x/term"

	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/utils"
)

// GenericCiscoHandler function handles generic Cisco style sessions
func GenericCiscoHandler(myFakeDevice *fakedevices.FakeDevice) ssh.Handler {
	return func(s ssh.Session) {

		// Exec mode: client sent a command directly (e.g., ssh host "show version")
		if cmd := s.RawCommand(); cmd != "" {
			log.Printf("exec: %s", cmd)
			match, matchedCommand, multipleMatches, _ := utils.CmdMatch(cmd, myFakeDevice.SupportedCommands)
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

		// Interactive shell mode
		contextState := myFakeDevice.ContextSearch["base"]
		t := term.NewTerminal(s, myFakeDevice.Hostname+contextState)

		for {
			userInput, err := t.ReadLine()
			if err != nil {
				break
			}
			log.Println(userInput)

			done := handleShellInput(t, userInput, myFakeDevice, &contextState)
			if done {
				break
			}
		}
		log.Println("terminal closed")
	}
}

// handleShellInput processes a single line of user input in interactive shell mode.
// Returns true if the session should be terminated.
func handleShellInput(t *term.Terminal, userInput string, fd *fakedevices.FakeDevice, contextState *string) bool {
	if userInput == "" {
		t.Write([]byte(""))
		return false
	}

	// Check for context switching commands
	matchPrompt, matchedPrompt, multiplePromptMatches, err := utils.CmdMatch(userInput, fd.ContextSearch)
	if err != nil {
		log.Println(err) // coverage-ignore // CmdMatch never returns errors
		return true
	}

	if matchPrompt && !multiplePromptMatches {
		t.SetPrompt(fd.Hostname + fd.ContextSearch[matchedPrompt])
		*contextState = fd.ContextSearch[matchedPrompt]
		return false
	}

	if userInput == "exit" || userInput == "end" {
		if fd.ContextHierarchy[*contextState] == "exit" {
			return true
		}
		t.SetPrompt(fd.Hostname + fd.ContextHierarchy[*contextState])
		*contextState = fd.ContextHierarchy[*contextState]
		return false
	}

	if userInput == "reset state" {
		t.Write(append([]byte("Resetting State..."), '\n'))
		*contextState = fd.ContextSearch["base"]
		fd.Hostname = fd.DefaultHostname
		t.SetPrompt(fd.Hostname + *contextState)
		return false
	}

	// Handle hostname changes in config mode
	userInputFields := strings.Fields(userInput)
	if userInputFields[0] == "hostname" && *contextState == "(config)#" {
		fd.Hostname = strings.Join(userInputFields[1:], " ")
		log.Printf("Setting hostname to %s\n", fd.Hostname)
		t.SetPrompt(fd.Hostname + *contextState)
		return false
	}

	// Match against supported commands
	return dispatchCommand(t, userInput, fd)
}

// dispatchCommand matches userInput against supported commands and writes the response.
// Returns true if the session should be terminated.
func dispatchCommand(t *term.Terminal, userInput string, fd *fakedevices.FakeDevice) bool {
	match, matchedCommand, multipleMatches, err := utils.CmdMatch(userInput, fd.SupportedCommands)
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
