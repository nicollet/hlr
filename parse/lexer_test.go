package parse

import (
	"fmt"
	"hlr/file"
	"testing"
)

func TestLex(t *testing.T) {
	rf := file.New()
	text := rf.Read("../includes/test")

	// fmt.Print(text)

	l := lex(text)
	for {
		item, more := <-l.items
		if !more {
			break
		}
		fmt.Printf("item: %v\n", item)
	}
}

type testVarValue struct {
	output   string
	attended string
}

var varValueTests = []testVarValue{
	{"TITI=\"un,deux,trois\"",
		"|TITI|=|un|deux|trois|"},
	{"TITI   =  \" un  deux \n, trois \" ",
		"|TITI|=|un|deux|trois|"},
	{"TITI = \" un deux\" TOTO = \"\n trois \"",
		"|TITI|=|un|deux|TOTO|=|trois|"},
}

func testOneVarValue(t *testing.T, input string, attended string) {
	l := lex(input)
	var output string
	for {
		item, more := <-l.items
		if !more {
			break
		}
		output += "|" + item.val
	}
	// note: got a last | for the EOF event
	if output != attended {
		t.Errorf("got: %v, not %v", output, attended)
	}
}

func TestVarValue(t *testing.T) {
	for _, tVarValue := range varValueTests {
		testOneVarValue(t, tVarValue.output, tVarValue.attended)
	}
}

// vim: set sw=2 ts=2 list:
