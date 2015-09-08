package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	// include_regexp, _ := regexp.Compile(`^\s*include\s+"(.*?)"`)
	rf := newReadFile()
	lines := rf.readFile("test")
	//for filePos := range getRunes(lines) {
	//fmt.Printf("line: %d, offset: %d, %c \n",
	//	filePos.lineNumber, filePos.lineOffset, filePos.r)
	//	fmt.Printf("%c", filePos.r)
	//}
	chFilePos := getRunes(lines)
	for {
		filePos, err := nextNoSpace(chFilePos)
		if err != nil {
			fmt.Print(err)
		} else {
			fmt.Printf("%c", filePos.r)
		}
	}
}

func nextNoSpace(ch chan FilePos) (FilePos, error) {
	var filePos FilePos
	for {
		filePos = <-ch
		r := filePos.r
		if r != ' ' || r != '\t' || r != '\n' {
			return filePos, nil
		}
		fmt.Printf("%c", filePos.r)
	}
	return filePos, errors.New("Can't find any space")
}

type FilePos struct {
	r          rune // current rune
	lineNumber int
	lineOffset int // offset on current line
}

// generator to advance one rune by one

// generator to get rune by rune a []string
func getRunes(lines []string) chan FilePos {
	ch := make(chan FilePos)
	var filePos FilePos
	go func() {
		for ln, line := range lines {
			filePos.lineNumber = ln
			for offset, r := range line {
				filePos.lineOffset = offset
				filePos.r = r
				ch <- filePos
			}
			filePos.lineOffset += 1
			filePos.r = '\n'
			ch <- filePos
		}
		close(ch)
	}()
	return ch
}

type ReadFile struct {
	doneFiles     map[string]bool
	includeRegexp *regexp.Regexp
}

func newReadFile() ReadFile {
	return ReadFile{
		doneFiles:     make(map[string]bool),
		includeRegexp: regexp.MustCompile(`^\s*include\s+"(.*?)"`),
	}
}

func (rf ReadFile) readFile(file_name string) []string {
	var lines []string
	if file_name[0] == '/' {
		return rf.readOneFile(file_name, lines)
	}

	for _, p := range []string{"hosts", "includes"} {
		name := p + "/" + file_name
		names := []string{name}
		name_strip := stripnb(name)
		if name_strip != name {
			names = append(names, name_strip)
		}
		for _, file_name := range names {
			if _, err := os.Stat(file_name); os.IsNotExist(err) {
				continue
			}
			_, hask := rf.doneFiles[file_name]
			if hask {
				continue
			}
			lines_ := rf.readOneFile(file_name, lines)
			lines = append(lines, lines_...)
		}
	}
	return lines
}

func stripnb(hostname string) string {
	return strings.Trim(hostname, "0123456789")
}

func (rf ReadFile) readOneFile(file_name string, lines []string) []string {
	rf.doneFiles[file_name] = true
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// re := regexp.MustCompile(`^\s*include\s+"(.*?)"`)
	re := rf.includeRegexp
	for scanner.Scan() {
		line := scanner.Text()
		includes := re.FindSubmatch([]byte(line))
		if len(includes) >= 1 {
			// fmt.Printf("! include %s\n", includes[1])
			include_fname := string("includes/" + string(includes[1]))
			lines = append(lines, rf.readFile(include_fname)...)
		} else {
			lines = append(lines, line)
		}
	}
	return lines
}

// Maybe later: easier stuff before
// const (
//	TABLE_DEFAULT = iota
//	TABLE_FILTER
//	TABLE_RAW
//)

//type Table struct {
//	rules map[int][]string // or map[int][]*string ?
//}

//func initVariables(lines []string) map[string]string {
//	variables := map[string]string{
//		"RFC1918": "10.0.0.0/8,192.168.0.0/16,172.16.0.0/12",
//	}
// var newVar string = ""
//}

// vim: set sw=2 ts=2 list:
