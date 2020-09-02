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

// FakeDevice Struct for the device we will be simulating
type FakeDevice struct {
	vendor            string            // Vendor of this fake device
	platform          string            // Platform of this fake device
	hostname          string            // Hostname of the fake device
	password          string            // Password of the fake device
	supportedCommands map[string]string // What commands this fake device supports
	contextSearch     map[string]string // The available CLI prompt/contexts on this fake device
	contextHierarchy  map[string]string // The heiarchy of the available contexts
}

// TranscriptMapPlatform struct for use inside of a TranscriptMap struct
type TranscriptMapPlatform struct {
	Vendor             string            `yaml:"vendor" json:"vendor"`
	Hostname           string            `yaml:"hostname" json:"hostname"`
	Password           string            `yaml:"password" json:"password"`
	CommandTranscripts map[string]string `yaml:"command_transcripts" json:"command_transcripts"`
	ContextSearch      map[string]string `yaml:"context_search" json:"context_search"`
	ContextHierarchy   map[string]string `yaml:"context_hierarchy" json:"context_hierarchy"`
}

// TranscriptMap Struct for modeling the TranscriptMap YAML
type TranscriptMap struct {
	Platforms []map[string]TranscriptMapPlatform `yaml:"platforms" json:"platforms"`
}

func readFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func sshFakeDeviceInit(
	vendor string,
	platform string,
	myTranscriptMap TranscriptMap,
) *FakeDevice {

	supportedCommands := make(map[string]string)
	contextSearch := make(map[string]string)
	contextHierarchy := make(map[string]string)
	commandTranscriptFiles := make(map[string]string)

	// Find the hostname, password, and other info in the data for this device
	var deviceHostname string
	var devicePassword string
	for _, fakeDevicePlatform := range myTranscriptMap.Platforms {
		// fmt.Printf("\nPlatform Map:\n%+v\n", fakeDevicePlatform)
		for k, v := range fakeDevicePlatform {
			if k == platform {
				// fmt.Printf("\nKey: %+v\n", k)
				// fmt.Printf("Value: %+v\n", v)
				deviceHostname = v.Hostname
				devicePassword = v.Password
				contextSearch = v.ContextSearch
				contextHierarchy = v.ContextHierarchy
				commandTranscriptFiles = v.CommandTranscripts
			}
		}
	}

	// Iterate through the command transcripts and read their contents into our fake device
	for k, v := range commandTranscriptFiles {
		supportedCommands[k] = readFile(v)
	}

	// Create our fake device and return it
	myFakeDevice := FakeDevice{
		vendor:            vendor,
		platform:          platform,
		hostname:          deviceHostname,
		password:          devicePassword,
		supportedCommands: supportedCommands,
		contextSearch:     contextSearch,
		contextHierarchy:  contextHierarchy,
	}

	//fmt.Printf("\n%+v\n", myFakeDevice)
	return &myFakeDevice
}

// ssh listner function that creates a fake device and terminal session
func sshListener(myFakeDevice *FakeDevice, portNumber int, done chan bool) {

	ssh.Handle(func(s ssh.Session) {

		// io.WriteString(s, fmt.Sprintf(SHOW_VERSION_PAGING_DISABLED))
		term := terminal.NewTerminal(s, myFakeDevice.hostname+myFakeDevice.contextSearch["base"])
		contextState := myFakeDevice.contextSearch["base"]
		for {
			userInput, err := term.ReadLine()
			if err != nil {
				break
			}
			log.Println(userInput)

			// Handle any responses provided at the terminal of the fakeDevice
			if myFakeDevice.supportedCommands[userInput] != "" {
				// lookup supported commands for the user input
				term.Write(append([]byte(myFakeDevice.supportedCommands[userInput]), '\n'))

			} else if userInput == "" {
				// return if nothing is entered
				term.Write(append([]byte(userInput)))

			} else if myFakeDevice.contextSearch[userInput] != "" {
				// switch contexts as needed
				term.SetPrompt(string(myFakeDevice.hostname + myFakeDevice.contextSearch[userInput]))
				contextState = myFakeDevice.contextSearch[userInput]

			} else if userInput == "exit" || userInput == "end" {
				// drop down configs if required
				if myFakeDevice.contextHierarchy[contextState] == "exit" {
					break
				} else {
					term.SetPrompt(string(myFakeDevice.hostname + myFakeDevice.contextHierarchy[contextState]))
					contextState = myFakeDevice.contextHierarchy[contextState]
				}

			} else {
				term.Write(append([]byte("% Ambiguous command:  \""+userInput+"\""), '\n'))
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
	// fmt.Printf("Raw Transcript Map from file:\n\n%s\n", transcriptMapRaw)

	myTranscriptMap := TranscriptMap{}
	err = yaml.UnmarshalStrict([]byte(transcriptMapRaw), &myTranscriptMap)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// fmt.Printf("YAML Parsed Transcript Map:\n\n%+v\n", myTranscriptMap)

	// Init our fake device
	myFakeDevice := sshFakeDeviceInit(
		"cisco",
		"csr1000v",
		myTranscriptMap,
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
