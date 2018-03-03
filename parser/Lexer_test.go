package parser

import (
	"strconv"
	"testing"
)

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
	items := []TokenType{
		TokenSelect,
		TokenLiteral,
		TokenDot,
		TokenAsterisk,
		TokenColon,
		TokenLiteral,
		TokenAs,
		TokenLiteral,
		TokenColon,
		TokenSum,
		TokenLeftParen,
		TokenLiteral,
		TokenRightParen,
		TokenColon,
		TokenString,
		TokenFrom,
		TokenLiteral,
	}

	lexer := lex("testing", "SELECT a.*, name AS koko, SUM(qty), 'qwe' FROM invoice")

	var token tokenItem
	for _, item := range items {
		token = lexer.nextItem()

		if token.Type == TokenEOF {
			break
		} else if token.Type == TokenError {
			t.Error(token.String())
		} else if token.Type != item {
			t.Errorf("Expect %s but get %s", item, token.String())
		}
	}

	lexer.drain() //clear all scanning process and exit the goroutine
}

func TestLexer_tokenize(t *testing.T) {
	expectedTokens := []TokenType{
		TokenSelect,
		TokenLiteral,
		TokenDot,
		TokenAsterisk,
		TokenColon,
		TokenLiteral,
		TokenAs,
		TokenLiteral,
		TokenEOF,
	}

	tokens := tokenize("SELECT a.*, name AS koko")

	if len(expectedTokens) != len(tokens) {
		t.Errorf("tokens quantity not tally, expect %d, actual get %d", len(expectedTokens), len(tokens))
		return
	}

	for i := 0; i < len(expectedTokens); i++ {
		if expectedTokens[i] != tokens[i].Type {
			t.Errorf("Expect token %s at index %d, but get %s",
				expectedTokens[i].String(), i, tokens[i].Type.String())
		}
	}

	/////// test token recognition
	expectedTokens = []TokenType{
		TokenGroupBy,
		TokenJoin,
		TokenOrderBy,
		TokenLeftJoin,
		TokenRightJoin,
		TokenInnerJoin,
		TokenOuterJoin,
		TokenAsc,
		TokenDesc,
		TokenEOF,
	}

	tokens = tokenize("GROUP BY JOIN ORDER BY LEFT JOIN RIGHT JOIN INNER JOIN OUTER JOIN ASC DESC") // JOIN LEFT JOIN RIGHT JOIN INNER JOIN OUTER JOIN")

	if len(expectedTokens) != len(tokens) {
		t.Errorf("tokens quantity not tally, expect %d, actual get %d", len(expectedTokens), len(tokens))
		return
	}

	for i := 0; i < len(expectedTokens); i++ {
		if expectedTokens[i] != tokens[i].Type {
			t.Errorf("Expect token %s at index %d, but get %s",
				expectedTokens[i].String(), i, tokens[i].Type.String())
		}
	}
}
