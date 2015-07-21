package main

import (
	// "fmt"
	"regexp"
	"testing"
)

func BenchmarkReadFile(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var rf ReadFile = ReadFile{
			doneFiles:     make(map[string]bool),
			includeRegexp: regexp.MustCompile(`^\s*include\s+"(.*?)"`),
		}
		// for _, line := range rf.readFile("test") {
		//	fmt.Println(line)
		// }
		rf.readFile("test")

	}
}
