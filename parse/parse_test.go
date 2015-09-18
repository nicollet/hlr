package parse

import (
	"fmt"
	"hlr/file"
	"testing"
)

func TestParse(t *testing.T) {
	rf := file.New()
	text := rf.Read("../includes/test")
	variables := BuildVariables(text)

	for k, v := range *variables {
		fmt.Printf("%v = /%v/\n", k, v)
	}
}

// vim: set sw=2 ts=2 list:
