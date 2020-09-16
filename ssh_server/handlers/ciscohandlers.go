// Package handlers contains SSH Handlers for specific device types
// in order to best emulate their actual behavior.
package handlers

import (
	"log"
	"strings"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/tbotnz/cisgo-ios/fakedevices"
	"github.com/tbotnz/cisgo-ios/utils"
)

// GenericCiscoHandler function handles generic Cisco style sessions
func GenericCiscoHandler(myFakeDevice *fakedevices.FakeDevice) {

	// Prepare the "ssh.DefaultHandler", this houses our device specific functionality
	ssh.Handle(func(s ssh.Session) {

		// Setup our initial "context" or prompt
		ContextState := myFakeDevice.ContextSearch["base"]

		// Setup a terminal with the hostname + initial context state as a prompt
		term := terminal.NewTerminal(s, myFakeDevice.Hostname+ContextState)

		// Iterate over any user input that is provided at the terminal
		for {
			userInput, err := term.ReadLine()
			if err != nil {
				break
			}
			log.Println(userInput)

			// Handle any empty input (assumed to just be a carriage return)
			if userInput == "" {
				// return nothing but a newline if nothing is entered
				term.Write([]byte(""))
				continue
			}

			// Handle any context switching
			if myFakeDevice.ContextSearch[userInput] != "" {
				// switch contexts as needed
				term.SetPrompt(string(myFakeDevice.Hostname + myFakeDevice.ContextSearch[userInput]))
				ContextState = myFakeDevice.ContextSearch[userInput]
				continue
			} else if userInput == "exit" || userInput == "end" {
				// Back out of the lower contexts, i.e. drop from "(config)#" to "#"
				if myFakeDevice.ContextHierarchy[ContextState] == "exit" {
					break
				} else {
					term.SetPrompt(string(myFakeDevice.Hostname + myFakeDevice.ContextHierarchy[ContextState]))
					ContextState = myFakeDevice.ContextHierarchy[ContextState]
					continue
				}
			}

			// Split user input into fields
			userInputFields := strings.Fields(userInput)

			// Handle hostname changes
			if userInputFields[0] == "hostname" && ContextState == "(config)#" {
				// Set the hostname to the values after "hostname" in the userInputFields
				myFakeDevice.Hostname = strings.Join(userInputFields[1:], " ")
				log.Printf("Setting hostname to %s\n", myFakeDevice.Hostname)
				term.SetPrompt(myFakeDevice.Hostname + ContextState)
				continue
			}

			// Run userInput through the command matcher
			match, matchedCommand, multipleMatches, err := utils.CmdMatch(userInput, myFakeDevice.SupportedCommands)
			if err != nil {
				log.Println(err)
				break
			}

			if match && !multipleMatches {
				// Render the matched command output
				output, err := fakedevices.TranscriptReader(
					myFakeDevice.SupportedCommands[matchedCommand],
					myFakeDevice,
				)
				if err != nil {
					log.Fatal(err)
				}

				// Write the output of our matched command
				term.Write(append([]byte(output), '\n'))
				continue
			} else if multipleMatches {
				// Multiple commands were matched, throw ambigious command
				term.Write(append([]byte("% Ambiguous command:  \""+userInput+"\""), '\n'))
				continue
			} else {
				// If all else fails, we did not recognize the input!
				term.Write(append([]byte("% Ambiguous command:  \""+userInput+"\""), '\n'))
				continue
			}
		}
		log.Println("terminal closed")

	})
}
