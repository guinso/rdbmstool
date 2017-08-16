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
	lexer := lex("testing", "SELECT * FROM invoice")

	tSelect := lexer.nextItem()
	t.Log(tSelect.String())
	if tSelect.Type != tokenSelect {
		t.Errorf("Expect tokenSelect but get: %s", tSelect.String())
	}

	tStar := lexer.nextItem()
	t.Log(tStar.String())
	if tStar.Type != tokenAsterisk {
		t.Errorf("Expect tokenAsterisk but get: %s", tStar.String())
	}

	tFrom := lexer.nextItem()
	t.Log(tFrom.String())
	if tFrom.Type != tokenFrom {
		t.Errorf("Expect tokenFrom but get: %s", tFrom.String())
	}

	tLiteral := lexer.nextItem()
	t.Log(tLiteral.String())
	if tLiteral.Type != tokenLiteral {
		t.Errorf("Expect tokenLiteral but get: %s", tLiteral.String())
	}

	lexer.drain() //clear all scanning process and exit the goroutine
}
