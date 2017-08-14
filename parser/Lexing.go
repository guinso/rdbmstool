package parser

//Lexical Analysis - convert string into a stack of tokens

import (
	"strconv"
	"strings"
)

//StateFn state function, representation of a fraction of lexing state machine
type StateFn func(item *Lexer) StateFn

const eof = 0

var kw = []TokenItem{
	{TokenSelect, "select", 0, 0},
	{TokenFrom, "from", 0, 0},
	{TokenJoin, "join", 0, 0},
	{TokenHaving, "having", 0, 0},
	{TokenWhere, "where", 0, 0},
	{TokenUnion, "union", 0, 0},
	{TokenLimit, "limit", 0, 0},
	{TokenOffset, "offset", 0, 0},
	{TokenOn, "on", 0, 0},
	{TokenBetween, "between", 0, 0},
	{TokenLike, "like", 0, 0},
	{TokenCreate, "create", 0, 0},
	{TokenTable, "table", 0, 0},
	{TokenView, "view", 0, 0},
	{TokenDrop, "drop", 0, 0},
	{TokenAsc, "asc", 0, 0},
	{TokenDesc, "desc", 0, 0},
}

var symbols = []TokenItem{
	{TokenEqual, "=", 0, 0},
	{TokenNotEqual, "<>", 0, 0},
	{TokenNotEqual, "!=", 0, 0},
	{TokenGreater, ">", 0, 0},
	{TokenGreaterEqual, ">=", 0, 0},
	{TokenLesser, "<", 0, 0},
	{TokenLesserEqual, "<=", 0, 0},
	{TokenQuestionMark, "?", 0, 0},
	{TokenWildcard, "%", 0, 0},
	{TokenAsterisk, "*", 0, 0},
	{TokenAdd, "+", 0, 0},
	{TokenSubtract, "-", 0, 0},
	{TokenDivide, "/", 0, 0},
	{TokenLeftParen, "(", 0, 0},
	{TokenRightParen, ")", 0, 0},
	{TokenDot, ".", 0, 0},
	{TokenColon, ",", 0, 0},
	{TokenSemiColon, ";", 0, 0},
}

func isWhiteSpace(input rune) bool {
	return input == ' ' || input == '\n' || input == '\t'
}

func isAlphaNumeric(input rune) bool {
	return isNumeric(input) || isLetter(input)
}

func isLiteralCharacter(input rune) bool {
	return isAlphaNumeric(input) || input == '_'
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
	nr1 := lex.peekAhead(1)
	nr2 := lex.peekAhead(2)

	//ignore white space
	if isWhiteSpace(nr1) {
		lex.next()
		lex.ignore()
		return lexText
	}

	//looking for keyword if match
	for _, k := range kw {
		if isKeywordMatch(lex, k.Value) {
			return lexKeyword(lex, k.Value, k.Type)
		}
	}

	//handle complex keyword(s)
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

	//looking for SQL parameter
	if nr1 == ':' && isLetter(nr2) {
		return lexParameter(lex)
	}

	//looking for number
	if isNumeric(nr1) || ((nr1 == '+' || nr1 == '-') && isNumeric(nr2)) {
		return lexNumber(lex)
	}

	//looking for quoted string
	if nr1 == '\'' {
		return lexQuoteString(lex)
	}

	//looking for literal (literal can start with backquote or without backquote)
	if nr1 == '`' || isLetter(nr2) {
		return lexLiteral(lex)
	}

	//looking for symbols if match
	for _, s := range symbols {
		if isKeywordMatch(lex, s.Value) {
			return lexKeyword(lex, s.Value, s.Type)
		}
	}

	//go next rune and check is EOF or not
	if nr1 == eof {
		// Correctly reached EOF.
		if lex.pos > lex.start {
			lex.emit(TokenText)
		}

		lex.emit(TokenEOF) //signal end of file token
		return nil         //stop run loop
	}

	//handle unspecified token
	return lex.errorf("Syntax error charactor(%s) not recognize by SQL lexical analysis; line %d, column %d",
		strconv.QuoteRune(nr1), lex.line, lex.pos+1)
}

//lexKeyword
func lexKeyword(lex *Lexer, keyword string, tokenTypee TokenType) StateFn {
	lex.fastForward(len(keyword))
	lex.emit(tokenTypee)

	return lexText
}

func lexJoin(lex *Lexer, tokenTy TokenType, keywordLen int) StateFn {

	//fast forward the keyword
	if err := lex.fastForward(keywordLen); err != nil {
		return lex.errorf(err.Error())
	}

	//skip all white space
	for {
		r := lex.next()
		if isWhiteSpace(r) {
			continue //keep looping
		} else if r == eof {
			//handle syntax error
			return lex.errorf("unexpected end of file reached at column %d", lex.pos)
		} else if lex.matchPrefix("join", "JOIN") {
			break
		}
	}

	//fast forward lexer to end of 'join' keyword
	lex.backup() //move backward one step
	if err := lex.fastForward(len("join")); err != nil {
		return lex.errorf(err.Error())
	}

	if r := lex.peek(); !isWhiteSpace(r) {
		return lex.errorf("syntax error detected for JOIN statement, "+
			"expect white space after 'join' keyword at column %d", lex.pos)
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
			continue
		} else if r == eof {
			//handle syntax error
			return lex.errorf("unexpected end of file reached at column %d", lex.pos)
		} else if lex.matchPrefix("by", "BY") {
			break
		}
	}

	//fast forward lexer to end of 'by' keyword
	if err := lex.fastForward(2); err != nil {
		return lex.errorf(err.Error())
	}

	if r := lex.peek(); !isWhiteSpace(r) {
		return lex.errorf("syntax error detected for Group By statement, "+
			"expect white space after 'by' keyword at column %d", lex.pos)
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

func lexParameter(lex *Lexer) StateFn {
	lex.next() //accept semi colon
	lex.next() //accept letter

	for { //break when reach white space
		r := lex.next()

		if isWhiteSpace(r) || r == eof {
			break
		} else if isLiteralCharacter(r) {
			continue
		} else {
			//handler syntax error
			return lex.errorf(
				"syntax error detected on lexing SQL parameter at line %d, column %d",
				lex.line, lex.pos)
		}
	}

	//push parameter token to feeder
	lex.emit(TokenParameter)

	return lexText
}

func lexQuoteString(lex *Lexer) StateFn {
	lex.next() //consume ' character

	for { //loop till reach ' charactor or EOF
		r := lex.next()

		if r == eof {
			return lex.errorf("Syntax error, quoted string not close before reach end of file")
		} else if r == '\'' {
			lex.emit(TokenString)
			return lexText
		}
	}
}

func lexLiteral(lex *Lexer) StateFn {
	r := lex.next()

	hasBackQuote := r == '`'

	if hasBackQuote {
		for { //loop til reach backquote
			r = lex.next()
			if r == '`' {
				lex.emit(TokenLiteral)
				return lexText
			} else if !isLiteralCharacter(r) {
				return lex.errorf(
					"Quoted literal doesn't close properly at line %d, pos %d",
					lex.line, lex.pos)
			}
		}
	} else {
		for { //loop til reach white space, EOF, or non-alphaNumeric
			r = lex.next()
			if !isLiteralCharacter(r) {
				lex.emit(TokenLiteral)
				return lexText
			}
		}
	}
}
