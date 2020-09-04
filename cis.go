package main

import (
<<<<<<< HEAD
	"github.com/tbotnz/cisgo-ios/fakedevices"
	"github.com/tbotnz/cisgo-ios/ssh_server/handlers"
	"github.com/tbotnz/cisgo-ios/ssh_server/sshlistners"
	"github.com/tbotnz/cisgo-ios/utils"
)

func main() {
=======
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	defaultHostname     = "cisgo1000v"
	defaultContextState = ">"
	password            = "admin"
)

func readFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

type commandGroup struct {
	basic    map[string]string
	hostname map[string]string
	mode     map[string]string
}

func newCommandGroup() *commandGroup {
	cmds := new(commandGroup)
	cmds.basic = make(map[string]string)
	cmds.basic["terminal length 0"] = " "
	cmds.basic["terminal width 511"] = " "
	cmds.basic["show ip interface brief"] = readFile("config/show_ip_int_bri.txt")

	cmds.hostname = make(map[string]string)
	cmds.hostname["show version"] = readFile("config/show_version.txt")
	cmds.hostname["show running-config"] = readFile("config/show_running-config.txt")

	cmds.mode = make(map[string]string)
	cmds.mode["conf t"] = "(config)#"
	cmds.mode["configure terminal"] = "(config)#"
	cmds.mode["configure t"] = "(config)#"
	cmds.mode["enable"] = "#"
	cmds.mode["en"] = "#"
	cmds.mode["base"] = ">"
	return cmds
}

type internalState struct {
	hostname    string
	currentMode string // >, #, or (config)#
	prompt      string
}

func (s *internalState) setMode(mode string) {
	s.currentMode = mode
	s.prompt = s.hostname + s.currentMode
}

func (s *internalState) setHostname(hostname string) {
	s.hostname = hostname
	s.prompt = s.hostname + s.currentMode
}

func (s *internalState) exit() bool {
	switch s.currentMode {
	case ">":
		return false
	case "#":
		s.setMode(">")
	case "(config)#":
		s.setMode("#")
	}
	return true
}

func newState() *internalState {
	// log.Println("created new internalState")
	return &internalState{defaultHostname, defaultContextState, defaultHostname + defaultContextState}
}

// ssh listener
func sshListener(portNumber int, done chan bool) {

	commandGroup := newCommandGroup()

	contextHierarchy := make(map[string]string)

	contextHierarchy["(config)#"] = "#"
	contextHierarchy["#"] = ">"
	contextHierarchy[">"] = "exit"

	thisState := newState()

	ssh.Handle(func(s ssh.Session) {

		term := terminal.NewTerminal(s, thisState.prompt)
		for {
			response, err := term.ReadLine()
			if err != nil {
				break
			}

			log.Println(response)
			if response == "reset state" {
				log.Println("resetting internal state")
				thisState = newState()
				term.SetPrompt(thisState.prompt)

			} else if response == "" {
				// return if nothing is entered
				term.Write(append([]byte(response)))

			} else if commandGroup.basic[response] != "" {
				// lookup supported commands for response
				term.Write(append([]byte(commandGroup.basic[response]), '\n'))

			} else if commandGroup.mode[response] != "" {
				// switch contexts as needed
				thisState.setMode(commandGroup.mode[response])
				term.SetPrompt(thisState.prompt)

			} else if response == "exit" || response == "end" {
				// drop down configs if required
				if thisState.exit() { // "true" means we're still active, "false" means we're done
					term.SetPrompt(thisState.prompt)
				} else {
					break
				}

			} else if commandGroup.hostname[response] != "" {
				term.Write([]byte(fmt.Sprintf(commandGroup.hostname[response], thisState.hostname)))

			} else if thisState.currentMode != ">" { // we're in config mode
				fields := strings.Fields(response)
				command := fields[0]
				if command == "hostname" {
					thisState.setHostname(strings.Join(fields[1:], " "))
					log.Printf("Setting hostname to %s\n", thisState.hostname)
					term.SetPrompt(thisState.prompt)

				} else {
					term.Write([]byte("% Ambiguous command:  \"" + response + "\"\n"))
				}

			} else {
				term.Write([]byte("% Ambiguous command:  \"" + response + "\"\n"))
			}

		}
		log.Println("terminal closed")

	})

	portString := ":" + strconv.Itoa(portNumber)
	//prt :=  portString
	log.Printf("starting cis.go ssh server on port %s\n", portString)
>>>>>>> 30e8bb2558b88bf8f310984529947f7d42335afb

	// Parse the command line arguments
	numListners, startingPortPtr, myTranscriptMap := utils.ParseArgs()

	// Initialize our fake device
	myFakeDevice := fakedevices.InitGenric(
		"cisco",         // Vendor
		"csr1000v",      // Platform
		myTranscriptMap, // Transcript map with locations of command output to play back
	)

	// Make a Channel named "done" for handling Goroutines, which expects a bool as return value
	done := make(chan bool, 1)

	// Iterate through the server ports and spawn a Goroutine for each
	for portNumber := *startingPortPtr; portNumber < numListners; portNumber++ {
		// Today this is just spawning a generic listner.
		// In the future, this is where we could split out listners/handlers by device type.
		go sshlistners.GenericListner(myFakeDevice, portNumber, handlers.GenericCiscoHandler, done)
	}

	// Recieve all the values from the channel (essentially wait on it to be empty)
	<-done
}
