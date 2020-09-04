// Package handlers contains SSH Handlers for specific device types
// in order to best emulate their actual behavior.
package handlers

import (
	"log"
	"strings"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/tbotnz/cisgo-ios/fakedevices"
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

			// Split user input into fields for use further down the handler
			userInputFields := strings.Fields(userInput)

			// Handle any responses provided at the terminal of the fakeDevice
			if myFakeDevice.SupportedCommands[userInput] != "" {
				// lookup supported commands for the user input
				term.Write(append([]byte(myFakeDevice.SupportedCommands[userInput]), '\n'))

			} else if userInput == "" {
				// return nothing but a newline if nothing is entered
				term.Write([]byte(""))

			} else if myFakeDevice.ContextSearch[userInput] != "" {
				// switch contexts as needed
				term.SetPrompt(string(myFakeDevice.Hostname + myFakeDevice.ContextSearch[userInput]))
				ContextState = myFakeDevice.ContextSearch[userInput]

			} else if userInputFields[0] == "hostname" && ContextState == "(config)#" {
				// Set the hostname to the values after "hostname" in the userInputFields
				myFakeDevice.Hostname = strings.Join(userInputFields[1:], " ")
				log.Printf("Setting hostname to %s\n", myFakeDevice.Hostname)
				term.SetPrompt(myFakeDevice.Hostname + ContextState)

			} else if userInput == "exit" || userInput == "end" {
				// Back out of the lower contexts, i.e. drop from "(config)#" to "#"
				if myFakeDevice.ContextHierarchy[ContextState] == "exit" {
					break
				} else {
					term.SetPrompt(string(myFakeDevice.Hostname + myFakeDevice.ContextHierarchy[ContextState]))
					ContextState = myFakeDevice.ContextHierarchy[ContextState]
				}

			} else {
				// If all else fails, we did not recognize the input!
				term.Write(append([]byte("% Ambiguous command:  \""+userInput+"\""), '\n'))
			}
		}
		log.Println("terminal closed")

	})
}
