package main

import (
	"bufio"
	//"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	// include_regexp, _ := regexp.Compile(`^\s*include\s+"(.*?)"`)
	rf := newReadFile()
	rf.readFile("test")
	//for _, line := range rf.readFile("test") {
	//	fmt.Println(line)
	//}
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

type itemType int
const (
	itemEqual = iota,
	itemVariable, // a variable
	itemString, // value of a variable (with quotes)
	itemEOF,
	item.Error, // an error
)

type item struct {
	typ itemType
	val string
}

// just a handy print function
func (i item) String() string {
	switch i.typ {
	case item.EOF:
		return "EOF"
	case item.Error:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// stateFn represents the state of the scanner
// as a function that returns the next state
type stateFn func(*lexer) stateFn

// notice the channel, ignore the rest
type lexer struct {
	name string // used only for error report
	input string // the string being scanned
	start int // start position of this item
	pos int // current position in the input
	width int // width of last rune read from input
	items chan item // channel of scanned items
}

// start a lexer
func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name: name,
		input: input,
		items: make(chan item),
	}
	go l.run() // Concurently run state machine
	return l, l.items
}

func (l *lexer) run() {
	for state := lexText; state != nil {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered
}

// emit passes an item back to the client
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func initVariables(lines []string) map[string]string {
	variables := map[string]string{
		"RFC1918": "10.0.0.0/8,192.168.0.0/16,172.16.0.0/12",
	}
}

func (remote Remote) generate() {
	tables := make(map[string]string) // retour: type table

	rf := newReadFile()
	lines := rf.readFile(remote.configFile)
	if len(s) == 0 {
		panic(fmt.Sprintf("%s can't be read. Aborting."))
	}

	// Here build and show variables

}

// vim: set sw=2 ts=2 list:
