package utils

import "testing"

func TestCmdMatch(t *testing.T) {

	// Create a fake SupportedCommands map
	mySupportedCommands := map[string]string{
		"show version": "a version of stuff",
		"show vlan":    "some vlan stuff",
	}

	input1 := "show version" // Should match "show version"
	input2 := "sho ver"      // Should match "show version"
	input3 := "sho vlan"     // Should match "show vlan"
	input4 := "s v"          // Should return no match

	match1, matchedCommand1 := CmdMatch(input1, mySupportedCommands)
	if match1 != true {
		t.Errorf(
			"CmdMatch('%s', %v) = (%t, '%s'); want (true, 'a version of stuff')",
			input1, mySupportedCommands, match1, matchedCommand1,
		)
	}

	match2, matchedCommand2 := CmdMatch(input2, mySupportedCommands)
	if match2 != true {
		t.Errorf(
			"CmdMatch('%s', %v) = (%t, '%s'); want (true, 'a version of stuff')",
			input2, mySupportedCommands, match2, matchedCommand2,
		)
	}

	match3, matchedCommand3 := CmdMatch(input3, mySupportedCommands)
	if match3 != true {
		t.Errorf(
			"CmdMatch('%s', %v) = (%t, '%s); want (true, 'some vlan stuff')",
			input3, mySupportedCommands, match3, matchedCommand3,
		)
	}

	match4, matchedCommand4 := CmdMatch(input4, mySupportedCommands)
	if match4 != false {
		t.Errorf(
			"CmdMatch('%s', %v) = (%t, '%s'); want (false, '')",
			input4, mySupportedCommands, match4, matchedCommand4,
		)
	}

}
