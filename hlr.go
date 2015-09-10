package main

import (
	"bufio"
	// "errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

func main() {
	rf := newReadFile()
	text := rf.readFile("test")

	l := lexer{
		input: text,
		items: make(chan item),
	}

	go l.run()
	for i := 0; i < 1000; i++ {
		fmt.Printf("item: %v\n", <-l.items)
	}

	//	for line := 0; line < 10; line++ {
	//		l.nextNoSpace()
	//		if l.getUpWord() {
	//			varName := <-l.items
	//			fmt.Printf("item: %v\n", varName)
	//
	//			// looking for =
	//			l.nextNoSpace()
	//			if l.getEqual() {
	//				equal := <-l.items
	//				fmt.Printf("item: %v\n", equal)
	//
	//				l.nextNoSpace()
	//				if l.getVarValue() {
	//					value := <-l.items
	//					fmt.Printf("item: %v\n", value)
	//				} else {
	//					fmt.Printf("No variable")
	//					l.nextLine()
	//				}
	//			} else {
	//				l.nextLine()
	//			}
	//		} else {
	//			l.nextLine()
	//		}
	//	}
}

func isSpace(r rune) bool {
	if r == ' ' || r == '\t' || r == '\n' {
		return true
	}
	return false
}

func isUpper(r rune) bool {
	if r >= 'A' && r <= 'Z' {
		return true
	}
	return false
}

func isUpperOrDigit(r rune) bool {
	if isUpper(r) || r >= '0' && r <= '9' {
		return true
	}
	return false
}

type item struct {
	typ itemType
	val string
}

type itemType int

const (
	itemError itemType = iota
	itemVarName
	itemEqual
	itemVarValue
	itemEOF
)

type stateFn func() stateFn

type lexer struct {
	name      string // to debug
	input     string // text
	itemStart int
	pos       int // current position
	width     int // width of current rune
	items     chan item
}

func (l *lexer) run() {
	for state := l.getUpWord(); state != nil; {
		state = state()
	}
	close(l.items)
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return 'Z'
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) ignore() {
	l.itemStart = l.pos
}

// return the number of spaces encountered
func (l *lexer) nextNoSpace() (ret int) {
	for {
		r := l.next()
		if r == 'Z' {
			break
		}
		// place the cursor at the beginning of next word
		if !isSpace(r) {
			l.backup()
			break
		}
		ret++
	}
	l.ignore()
	return ret
}

func (l *lexer) nextLine() {
	for {
		r := l.next()
		if r == 'Z' {
			break
		}
		if r == '\n' {
			l.ignore() // cursor at the beginning of next line
			break
		}
	}
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.itemStart:l.pos]}
	l.itemStart = l.pos
}

func (l *lexer) emitString(t itemType, val string) {
	l.items <- item{t, val}
	l.itemStart = l.pos
}

func (l *lexer) getEqual() stateFn {
	r := l.next()
	if r == '=' {
		l.emit(itemEqual)
		l.nextNoSpace()
		return l.getVarValue
	}
	l.nextLine()
	return l.getUpWord
}

// Have we found an upcase word
func (l *lexer) getUpWord() stateFn {
	l.nextNoSpace()
	r := l.next()
	if r == 'Z' {
		return nil
	}
	if !(isUpper(r) || r == '_') {
		l.nextLine()
		return l.getUpWord
	}

	for {
		r = l.next()
		if r == 'Z' {
			return nil
		}

		if !(isUpperOrDigit(r) || r == '_') {
			break
		}
	}
	if isSpace(r) || r == '=' {
		l.backup()
		l.emit(itemVarName)
		l.nextNoSpace()
		return l.getEqual
	}
	l.nextLine()
	return l.getUpWord
}

// we are already at the first non space character
func (l *lexer) getVarValue() stateFn {
	r := l.next()
	if r != '"' {
		fmt.Printf("not first double-quotes: %c\n", r)
		return nil
	}
	l.nextNoSpace()

	var val string
	successive_spaces := 0
	for {
		r = l.next()
		if r == '"' {
			// remove last comma if needed
			if strings.HasSuffix(val, ",") {
				val = val[:len(val)-1]
			}
			l.emitString(itemVarValue, val)
			l.nextNoSpace()
			return l.getUpWord
		}
		if isSpace(r) {
			successive_spaces++
			continue
		}
		if successive_spaces > 0 {
			val += ","
			successive_spaces = 0
		}
		val += string(r)
	}
	return nil // never get there: errors to be handled
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemVarName:
		return fmt.Sprintf("variable: %s", i.val)
	case itemEqual:
		return "="
	case itemVarValue:
		return fmt.Sprintf("varValue: %s", i.val)
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
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

func (rf ReadFile) readFile(file_name string) string {
	var text string
	if file_name[0] == '/' {
		return rf.readOneFile(file_name, text)
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
			new_text := rf.readOneFile(file_name, text)
			text += new_text
		}
	}
	return text
}

func stripnb(hostname string) string {
	return strings.Trim(hostname, "0123456789")
}

func (rf ReadFile) readOneFile(file_name string, text string) string {
	rf.doneFiles[file_name] = true
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	re := rf.includeRegexp
	for scanner.Scan() {
		line := scanner.Text()
		includes := re.FindSubmatch([]byte(line))
		if len(includes) >= 1 {
			include_fname := string("includes/" + string(includes[1]))
			text += rf.readFile(include_fname)
		} else {
			text += line + string('\n')
		}
	}
	return text
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
