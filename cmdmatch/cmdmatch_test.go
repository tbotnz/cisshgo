package cmdmatch

import "testing"

type inputMatch struct {
	match           bool
	matchedCommand  string
	multipleMatches bool
}

func TestMatch(t *testing.T) {
	mySupportedCommands := map[string]string{
		"show version":    "a version of stuff",
		"show vlan":       "some vlan stuff",
		"show vlan brief": "some more brief vlan stuff",
		"reboot":          "oh noes!",
	}

	inputs := map[string]inputMatch{
		"show version":         {true, "show version", false},    // exact match
		"show ver":             {true, "show version", false},    // prefix match on second word
		"sho vlan":             {true, "show vlan", false},       // prefix match on first word
		"s v":                  {true, "", true},                 // ambiguous: matches show version, show vlan, show vlan brief
		"show version made-up": {false, "", false},               // wrong word count
		"no version":           {false, "", false},               // no prefix match
		"Sho vLan BrIef":       {true, "show vlan brief", false}, // case-insensitive prefix match
		"show vlan!":           {false, "", false},               // no match (! not a prefix of brief)
	}

	for input, expected := range inputs {
		match, matchedCommand, multipleMatches := Match(input, mySupportedCommands)
		if match != expected.match ||
			matchedCommand != expected.matchedCommand ||
			multipleMatches != expected.multipleMatches {
			t.Errorf(
				"Match(%q) = (%t, %q, %t); want (%t, %q, %t)",
				input, match, matchedCommand, multipleMatches,
				expected.match, expected.matchedCommand, expected.multipleMatches,
			)
		}
	}
}

// TestMatch_PrefixNotSubstring verifies HasPrefix semantics:
// "sho" matches "show" but "how" does not (substring, not prefix).
func TestMatch_PrefixNotSubstring(t *testing.T) {
	cmds := map[string]string{"show version": ""}
	match, _, _ := Match("how version", cmds)
	if match {
		t.Error("expected no match: 'how' is a substring but not a prefix of 'show'")
	}
}

// TestMatch_EmptyInput returns no match.
func TestMatch_EmptyInput(t *testing.T) {
	cmds := map[string]string{"show version": ""}
	match, cmd, multi := Match("", cmds)
	if match || cmd != "" || multi {
		t.Errorf("expected (false, '', false) for empty input, got (%t, %q, %t)", match, cmd, multi)
	}
}

// TestMatch_EmptyCommands returns no match.
func TestMatch_EmptyCommands(t *testing.T) {
	match, cmd, multi := Match("show version", map[string]string{})
	if match || cmd != "" || multi {
		t.Errorf("expected (false, '', false) for empty commands, got (%t, %q, %t)", match, cmd, multi)
	}
}
