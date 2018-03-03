//Package parser lexer.go is mainly referenced from github.com/golang/go/src/text/template/parse/lex.go
package parser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

//lex create a new scanner for the input string
func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan tokenItem),
		line:  1}

	go l.run()
	return l
}

//lexer hold the state of the scanner
type lexer struct {
	name  string         //used only for error report
	input string         //the string being scanned
	start int            //start position of this scan
	pos   int            //position of current scan
	line  int            //line number of input
	width int            //length of last rune read from in
	items chan tokenItem //channel of scanned items
}

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) StateFn {
	l.items <- tokenItem{
		TokenError,
		fmt.Sprintf(format, args...),
		l.start,
		l.line}

	return nil
}

//Tokenize convert string context into array of tokens
//token array is use for next process: transform into abstract syntax tree
func tokenize(inputStr string) []tokenItem {
	result := []tokenItem{}

	lexer := lex("tokenize", inputStr)

	for true {
		tmpToken := lexer.nextItem()

		if tmpToken.Type == TokenError {
			break
		}

		result = append(result, tmpToken)

		if tmpToken.Type == eof {
			break
		}
	}

	return result
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() tokenItem {
	item := <-l.items
	//x l.lastPos = item.pos
	return item
}

// drain drains the output so the lexing goroutine will exit.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) drain() {
	for range l.items {
	}
}

//emit pass a Token item back to client
func (l *lexer) emit(t TokenType) {
	l.items <- tokenItem{t, l.input[l.start:l.pos], l.start, l.line}

	// Some items contain text internally. If so, count their newlines.
	switch t {
	case TokenText:
		l.line += strings.Count(l.input[l.start:l.pos], "\n")
	}

	l.start = l.pos
}

// run lexes the input by executing state functions until
// the state is nil.
func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered.
}

func (l *lexer) matchPrefix(searchPattern ...string) bool {
	for _, pattern := range searchPattern {
		if strings.HasPrefix(l.input[l.pos:], pattern) {
			return true
		}
	}

	return false
}

func (l *lexer) matchSuffix(searchPattern ...string) bool {
	for _, pattern := range searchPattern {
		if strings.HasSuffix(l.input[l.pos:], pattern) {
			return true
		}
	}

	return false
}

// peek returns but does not consume
// the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()

	return r
}

// peekAhead read rune ahead current without alter state machine
func (l *lexer) peekAhead(step int) rune {
	var r rune

	tmpW := 0
	pos := l.pos
	for i := 0; i < step; i++ {
		if pos >= len(l.input) {
			return eof
		}

		r, tmpW = utf8.DecodeRuneInString(l.input[pos:])
		pos += tmpW
	}

	return r
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width

	// Correct newline count.
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	if r == '\n' {
		l.line++
	}

	return r
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// move position forward
func (l *lexer) fastForward(step int) error {
	backupStart := l.start
	backupPos := l.pos
	backupLine := l.line
	backupWidth := l.width

	for i := 0; i < step; {
		r := l.next()

		if r == eof {
			//restore back state
			l.start = backupStart
			l.pos = backupPos
			l.width = backupWidth
			l.line = backupLine

			return fmt.Errorf("unable to fast forward as it reach end of file(EOF)")
		}
		i += l.width
	}

	return nil
}

// accept consumes the next rune
// if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}
