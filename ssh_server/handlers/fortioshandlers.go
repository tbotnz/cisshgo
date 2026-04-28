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

const fortiRootContext = "#"

// FortiOSHandler handles FortiOS style sessions.
func FortiOSHandler(fd *fakedevices.FakeDevice) ssh.Handler {
	return fortiGateSession(fd, nil)
}

// FortiOSScenarioHandler returns an ssh.Handler that plays back a scenario sequence.
func FortiOSScenarioHandler(fd *fakedevices.FakeDevice, sequence []transcript.SequenceStep) ssh.Handler {
	return fortiGateSession(fd, sequence)
}

func fortiGateSession(fd *fakedevices.FakeDevice, sequence []transcript.SequenceStep) ssh.Handler {
	return func(s ssh.Session) {
		if cmd := s.RawCommand(); cmd != "" {
			match, matchedCommand, multipleMatches := cmdmatch.Match(cmd, fd.SupportedCommands)
			if match && !multipleMatches {
				output, err := fakedevices.TranscriptReader(fd.SupportedCommands[matchedCommand], fd)
				if err == nil {
					io.WriteString(s, output)
				}
			}
			s.Exit(0)
			return
		}

		seqIdx := 0
		contextStack := []string{fortiRootContext}
		t := term.NewTerminal(s, devicePrompt(fd, contextStack[len(contextStack)-1]))

		for {
			userInput, err := t.ReadLine()
			if err != nil {
				break
			}
			if strings.TrimSpace(userInput) == "" {
				t.Write([]byte(""))
				continue
			}

			sequenceHandled, terminate := handleSequenceStep(t, userInput, fd, sequence, &seqIdx)
			if terminate {
				break
			}

			if fortiHandleStatefulCommand(t, strings.TrimSpace(userInput), fd, &contextStack) {
				break
			}

			if sequenceHandled {
				continue
			}

			match, matchedCommand, multipleMatches := cmdmatch.Match(userInput, fd.SupportedCommands)
			if multipleMatches {
				t.Write([]byte("Command parse error before '" + userInput + "'\n"))
				continue
			}
			if !match {
				t.Write([]byte("Command fail. Return code -61\n"))
				continue
			}
			output, err := fakedevices.TranscriptReader(fd.SupportedCommands[matchedCommand], fd)
			if err != nil {
				log.Println(err)
				break
			}
			t.Write(append([]byte(output), '\n'))
		}
	}
}

func fortiHandleStatefulCommand(t *term.Terminal, userInput string, fd *fakedevices.FakeDevice, contextStack *[]string) bool {
	fields := strings.Fields(userInput)
	if len(fields) == 0 {
		return false
	}

	switch fields[0] {
	case "config":
		if len(fields) < 2 {
			t.Write([]byte("Command fail. Return code -61\n"))
			return false
		}
		component := fields[len(fields)-1]
		*contextStack = append(*contextStack, "("+component+") #")
		t.SetPrompt(devicePrompt(fd, (*contextStack)[len(*contextStack)-1]))
		return false
	case "edit":
		if len(fields) < 2 {
			t.Write([]byte("Command fail. Return code -61\n"))
			return false
		}
		name := strings.Join(fields[1:], " ")
		*contextStack = append(*contextStack, "("+name+") #")
		t.SetPrompt(devicePrompt(fd, (*contextStack)[len(*contextStack)-1]))
		return false
	case "next":
		if len(*contextStack) > 1 {
			*contextStack = (*contextStack)[:len(*contextStack)-1]
		}
		t.SetPrompt(devicePrompt(fd, (*contextStack)[len(*contextStack)-1]))
		return false
	case "end":
		*contextStack = []string{fortiRootContext}
		t.SetPrompt(devicePrompt(fd, fortiRootContext))
		return false
	case "exit":
		return true
	}

	return false
}
