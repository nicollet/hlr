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

// vim: set sw=2 ts=2 list:
