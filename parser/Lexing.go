package parser

//Lexical Analysis - convert string into a stack of tokens

import (
	"strings"
)

//StateFn state function, representation of a fraction of lexing state machine
type StateFn func(item *Lexer) StateFn

const eof = 0

var kw = []TokenItem{
	{TokenSelect, "select"},
	{TokenFrom, "from"},
	{TokenJoin, "join"},
	{TokenHaving, "having"},
	{TokenWhere, "where"},
	{TokenUnion, "union"},
	{TokenLimit, "limit"},
	{TokenOffset, "offset"},
	{TokenOn, "on"},
	{TokenBetween, "between"},
	{TokenLike, "like"},
	{TokenCreate, "create"},
	{TokenTable, "table"},
	{TokenView, "view"},
	{TokenDrop, "drop"},
}

func isWhiteSpace(input rune) bool {
	return input == ' ' || input == '\n' || input == '\t'
}

func isAlphaNumeric(input rune) bool {
	return (input >= '0' && input <= '9') ||
		(input >= 'a' && input <= 'z') ||
		(input >= 'A' && input <= 'Z')
}

func isNumeric(input rune) bool {
	return input >= '0' && input <= '9'
}

func isLetter(input rune) bool {
	return (input >= 'a' && input <= 'z') ||
		(input >= 'A' && input <= 'Z')
}

func isKeywordMatch(lex *Lexer, keyword string) bool {
	return lex.matchPrefix(strings.ToUpper(keyword), strings.ToLower(keyword)) &&
		isWhiteSpace(lex.peekAhead(len(keyword)+1))
}

func lexText(lex *Lexer) StateFn {

	//loop infinite until explicitly break
	for {
		//ignore white space
		if isWhiteSpace(lex.peek()) {
			lex.ignore()
		}

		//looking for keyword if match
		for _, k := range kw {
			if isKeywordMatch(lex, k.Value) {
				return lexKeyword(lex, k.Value, k.Type)
			}
		}

		if lex.matchPrefix("group", "GROUP") {
			return lexGroupBy(lex)
		} else if lex.matchPrefix("inner", "INNER") {
			return lexJoin(lex, TokenInnerJoin, 5)
		} else if lex.matchPrefix("outer", "OUTER", "left", "LEFT", "right", "RIGHT") {
			return lexJoin(lex, TokenOuterJoin, 5)
		} else if lex.matchPrefix("left", "LEFT", "right", "RIGHT") {
			return lexJoin(lex, TokenLeftJoin, 4)
		} else if lex.matchPrefix("right", "RIGHT") {
			return lexJoin(lex, TokenRightJoin, 5)
		}

		//looking for parameter

		//looking for number

		//looking for quoted string

		//looking for literal

		//go next rune and check is EOF or not
		if lex.next() == eof {
			break
		}
	}

	lex.emit(TokenEOF) //signal end of file token
	return nil         //stop run loop
}

func lexJoin(lex *Lexer, tokenTy TokenType, keywordLen int) StateFn {
	//fast forward the keyword
	if err := lex.fastForward(keywordLen); err != nil {
		return lex.errorf(err.Error())
	}

	for {
		r := lex.next()
		if isWhiteSpace(r) {
			break
		} else if r == eof {
			//handle syntax error
			return lex.errorf("unexpected end of file reached at column %d", lex.pos)
		} else if lex.matchPrefix("join") {
			break
		}
	}

	//fast forward lexer to end of 'join' keyword
	if err := lex.fastForward(4); err != nil {
		return lex.errorf(err.Error())
	}

	if r := lex.peek(); !isWhiteSpace(r) {
		return lex.errorf("syntax error detected for JOIN statement, "+
			"expect white space after 'join' keyword at column %d", lex.pos+1)
	}

	lex.emit(tokenTy)
	return lexText
}

func lexGroupBy(lex *Lexer) StateFn {
	keywordLen := 5 //group
	tokenTy := TokenGroupBy

	//fast forward the keyword
	if err := lex.fastForward(keywordLen); err != nil {
		return lex.errorf(err.Error())
	}

	for {
		r := lex.next()
		if isWhiteSpace(r) {
			break
		} else if r == eof {
			//handle syntax error
			return lex.errorf("unexpected end of file reached at column %d", lex.pos)
		} else if lex.matchPrefix("join") {
			break
		}
	}

	//fast forward lexer to end of 'by' keyword
	if err := lex.fastForward(2); err != nil {
		return lex.errorf(err.Error())
	}

	if r := lex.peek(); !isWhiteSpace(r) {
		return lex.errorf("syntax error detected for Group By statement, "+
			"expect white space after 'by' keyword at column %d", lex.pos+1)
	}

	lex.emit(tokenTy)
	return lexText
}

func lexNumber(l *Lexer) StateFn {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	// Is it imaginary?
	// Is it imaginary?
	l.accept("i")
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q",
			l.input[l.start:l.pos])
	}
	l.emit(TokenNumber)
	return lexText
}

//lexKeyword
func lexKeyword(lex *Lexer, keyword string, tokenTypee TokenType) StateFn {
	lex.width = len(keyword)
	lex.pos += lex.width

	lex.emit(tokenTypee)
	return lexText
}
