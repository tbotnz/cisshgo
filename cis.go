package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
)

// ssh listernet
func sshListener(portNumber int, done chan bool) {

	supportedCommands := make(map[string]string)
	contextSearch := make(map[string]string)
	contextHierarchy := make(map[string]string)

	contextSearch["conf t"] = "(config)#"
	contextSearch["configure terminal"] = "(config)#"
	contextSearch["configure t"] = "(config)#"
	contextSearch["enable"] = "#"
	contextSearch["en"] = "#"
	contextSearch["base"] = ">"

	contextHierarchy["(config)#"] = "#"
	contextHierarchy["#"] = ">"
	contextHierarchy[">"] = "exit"

	hostname := "cisgo1000v"
	password := "admin"

	supportedCommands["show version"] = ``

	supportedCommands["show ip interface brief"] = ``

	supportedCommands["show running-config"] = ``

	ssh.Handle(func(s ssh.Session) {

		// io.WriteString(s, fmt.Sprintf(SHOW_VERSION_PAGING_DISABLED))
		term := terminal.NewTerminal(s, hostname+contextSearch["base"])
		contextState := ">"
		for {
			line, err := term.ReadLine()
			if err != nil {
				break
			}
			response := line
			log.Println(line)
			if supportedCommands[response] != "" {
				// lookup supported commands for response
				term.Write(append([]byte(supportedCommands[response]), '\n'))
			} else if response == "" {
				// return if nothing is entered
				term.Write(append([]byte(response)))
			} else if contextSearch[response] != "" {
				// switch contexts as needed
				term.SetPrompt(string(hostname + contextSearch[response]))
				contextState = contextSearch[response]
			} else if response == "exit" {
				// drop down configs if required
				if contextHierarchy[contextState] == "exit" {
					break
				} else {
					term.SetPrompt(string(hostname + contextHierarchy[contextState]))
					contextState = contextHierarchy[contextState]
				}
			} else {
				term.Write(append([]byte("% Ambiguous command:  \""+response+"\""), '\n'))
			}
		}
		log.Println("terminal closed")
	})

	portString := ":" + strconv.Itoa(portNumber)
	//prt :=  portString
	log.Printf("starting cis.go ssh server on port %s\n", portString)

	log.Fatal(
		ssh.ListenAndServe(
			portString,
			nil,
			ssh.PasswordAuth(
				func(ctx ssh.Context, pass string) bool {
					return pass == password
				},
			),
		),
	)

	done <- true
}

func main() {

	// Gather command line arguments and parse them
	listnersPtr := flag.Int("listners", 50, "How many listeners do you wish to spawn?")
	startingPortPtr := flag.Int("startingPort", 10000, "What port do you want to start at?")
	transcriptMapPtr := flag.String(
		"transcriptMap",
		"transcripts/transcript_map.yaml",
		"What file contains the map of commands to transcipted output?",
	)
	flag.Parse()

	// How many total listners will we have?
	listners := *startingPortPtr + *listnersPtr

	// Gather the command transcripts and create a map of vendor/platform/command
	transcriptMapRaw, err := ioutil.ReadFile(*transcriptMapPtr)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	transcriptMapYAML := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(transcriptMapRaw), &transcriptMapYAML)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// fmt.Printf("%s", transcriptMapYAML)

	// Make a Channel for handling Goroutines, name of `done` expects a bool as return value
	done := make(chan bool, 1)

	// Iterate through the server ports and spawn a Goroutine for each
	for portNumber := *startingPortPtr; portNumber < listners; portNumber++ {
		go sshListener(portNumber, done)
	}

	// Recieve all the values from the channel (essentially wait on it to be empty)
	<-done
}
