package parser

import "testing"
import "strconv"

func TestLexer_peekAhead(t *testing.T) {
	lexer := Lexer{
		"",
		"abcdefg",
		0,
		0,
		0,
		make(chan TokenItem)}

	r := lexer.peekAhead(4)

	if r != 'd' {
		t.Errorf("Expect 'd' but get %s", strconv.QuoteRune(r))
	}
}
