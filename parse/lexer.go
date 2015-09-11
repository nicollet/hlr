package parse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const eof = -1

type itemType int

type item struct {
	typ itemType
	val string
}

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

func lex(text string) *lexer {
	l := &lexer{
		input: text,
		items: make(chan item),
	}
	go l.run()
	return l
}

func (l *lexer) run() {
	for state := l.getUpWord(); state != nil; {
		state = state()
	}
	close(l.items)
}

// helper functions
func isSpace(r rune) bool {
	return unicode.IsSpace(r)
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

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
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
		if r == eof {
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
		if r == eof {
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

// State Functions: for our state machine

func (l *lexer) getEqual() stateFn {
	r := l.next()
	if r == eof {
		return nil
	}
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
	if r == eof {
		l.emit(itemEOF)
		return nil
	}
	if !(isUpper(r) || r == '_') {
		l.nextLine()
		return l.getUpWord
	}

	for {
		r = l.next()
		if r == eof {
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

// return an error item
func (l *lexer) errorF(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, fmt.Sprintf(format, args...)}
	return nil
}

// general print function for an item
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

// vim: set sw=2 ts=2 list:
