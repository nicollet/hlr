package main

import (
	// "fmt"
	"testing"
)

func BenchmarkReadFile(b *testing.B) {
	rf := newReadFile()
	for n := 0; n < b.N; n++ {
		rf.readFile("test")
	}
}
