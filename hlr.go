package main

import (
	"fmt"
	"hlr/file"
)

func main() {
	rf := file.New()
	text := rf.Read("test")

	fmt.Print(text)
}
