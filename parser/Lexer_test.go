package parser

import "testing"
import "strconv"

func TestLexer_peekAhead(t *testing.T) {
	lexer := lexer{
		"",
		"abcdefg",
		0,
		0,
		0,
		0,
		make(chan tokenItem)}

	r := lexer.peekAhead(4)

	if r != 'd' {
		t.Errorf("Expect 'd' but get %s", strconv.QuoteRune(r))
	}
}

func TestLexer_nextItem(t *testing.T) {
	items := []tokenType{
		tokenSelect,
		tokenAsterisk,
		tokenColon,
		tokenLiteral,
		tokenColon,
		tokenSum,
		tokenLeftParen,
		tokenLiteral,
		tokenRightParen,
		tokenFrom,
		tokenLiteral,
	}

	lexer := lex("testing", "SELECT *, name, SUM(qty) FROM invoice")

	var token tokenItem
	for _, item := range items {
		token = lexer.nextItem()

		if token.Type == tokenEOF {
			break
		} else if token.Type == tokenError {
			t.Error(token.String())
		} else if token.Type != item {
			t.Errorf("Expect %s but get %s", item, token.String())
		}
	}

	lexer.drain() //clear all scanning process and exit the goroutine
}
