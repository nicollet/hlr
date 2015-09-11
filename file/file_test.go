package file

import (
	// "fmt"
	"testing"
)

const testFile = "../includes/test"

func BenchmarkFile(b *testing.B) {
	rf := New()
	// text := rf.Read(testFile)
	// fmt.Printf("result: %q", text)
	for n := 0; n < b.N; n++ {
		rf.Read(testFile)
	}
}
