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

type fakeDevice struct {
	vendor            string            // Vendor of this fake device
	platform          string            // Platform of this fake device
	hostname          string            // Hostname of the fake device
	password          string            // Password of the fake device
	supportedCommands map[string]string // What commands this fake device supports
	contextSearch     map[string]string // The available CLI prompt/contexts on this fake device
	contextHierarchy  map[string]string // The heiarchy of the available contexts
}

func sshFakeDeviceInit(
	vendor string,
	platform string,
	transcriptMapYAML map[interface{}]interface{},
) *fakeDevice {

	myFakeDevice := fakeDevice{
		vendor:   vendor,
		platform: platform,
		hostname: transcriptMapYAML[vendor][platform]["hostname"],
		password: transcriptMapYAML["cisco"]["csr1000v"]["password"],
	}

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

	supportedCommands["show version"] = ``
	supportedCommands["show ip interface brief"] = ``
	supportedCommands["show running-config"] = ``

	return &myFakeDevice
}

// ssh listernet
func sshListener(myFakeDevice *fakeDevice, portNumber int, done chan bool) {

	ssh.Handle(func(s ssh.Session) {

		// io.WriteString(s, fmt.Sprintf(SHOW_VERSION_PAGING_DISABLED))
		term := terminal.NewTerminal(s, myFakeDevice.hostname+myFakeDevice.contextSearch["base"])
		contextState := ">"
		for {
			line, err := term.ReadLine()
			if err != nil {
				break
			}
			response := line
			log.Println(line)
			if myFakeDevice.supportedCommands[response] != "" {
				// lookup supported commands for response
				term.Write(append([]byte(myFakeDevice.supportedCommands[response]), '\n'))
			} else if response == "" {
				// return if nothing is entered
				term.Write(append([]byte(response)))
			} else if myFakeDevice.contextSearch[response] != "" {
				// switch contexts as needed
				term.SetPrompt(string(myFakeDevice.hostname + myFakeDevice.contextSearch[response]))
				contextState = myFakeDevice.contextSearch[response]
			} else if response == "exit" {
				// drop down configs if required
				if myFakeDevice.contextHierarchy[contextState] == "exit" {
					break
				} else {
					term.SetPrompt(string(myFakeDevice.hostname + myFakeDevice.contextHierarchy[contextState]))
					contextState = myFakeDevice.contextHierarchy[contextState]
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
					return pass == myFakeDevice.password
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

	// Init our fake device
	myFakeDevice := sshFakeDeviceInit(
		"cisco",
		"csr100v",
		transcriptMapYAML,
	)

	// Make a Channel for handling Goroutines, name of `done` expects a bool as return value
	done := make(chan bool, 1)

	// Iterate through the server ports and spawn a Goroutine for each
	for portNumber := *startingPortPtr; portNumber < listners; portNumber++ {
		go sshListener(myFakeDevice, portNumber, done)
	}

	// Recieve all the values from the channel (essentially wait on it to be empty)
	<-done
}
