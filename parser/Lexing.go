package parser

//Lexical Analysis - convert string into a stack of tokens

import (
	"strconv"
	"strings"
)

//StateFn state function, representation of a fraction of lexing state machine
type StateFn func(item *lexer) StateFn

const eof = 0

var kw = []tokenItem{
	{tokenSelect, "select", 0, 0},
	{tokenFrom, "from", 0, 0},
	{tokenJoin, "join", 0, 0},
	{tokenHaving, "having", 0, 0},
	{tokenWhere, "where", 0, 0},
	{tokenUnion, "union", 0, 0},
	{tokenLimit, "limit", 0, 0},
	{tokenOffset, "offset", 0, 0},
	{tokenOn, "on", 0, 0},
	{tokenBetween, "between", 0, 0},
	{tokenLike, "like", 0, 0},
	{tokenCreate, "create", 0, 0},
	{tokenTable, "table", 0, 0},
	{tokenView, "view", 0, 0},
	{tokenDrop, "drop", 0, 0},
	{tokenAsc, "asc", 0, 0},
	{tokenDesc, "desc", 0, 0},
}

var symbols = []tokenItem{
	{tokenEqual, "=", 0, 0},
	{tokenNotEqual, "<>", 0, 0},
	{tokenNotEqual, "!=", 0, 0},
	{tokenGreater, ">", 0, 0},
	{tokenGreaterEqual, ">=", 0, 0},
	{tokenLesser, "<", 0, 0},
	{tokenLesserEqual, "<=", 0, 0},
	{tokenQuestionMark, "?", 0, 0},
	{tokenWildcard, "%", 0, 0},
	{tokenAsterisk, "*", 0, 0},
	{tokenAdd, "+", 0, 0},
	{tokenSubtract, "-", 0, 0},
	{tokenDivide, "/", 0, 0},
	{tokenLeftParen, "(", 0, 0},
	{tokenRightParen, ")", 0, 0},
	{tokenDot, ".", 0, 0},
	{tokenColon, ",", 0, 0},
	{tokenSemiColon, ";", 0, 0},
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

func isKeywordMatch(lex *lexer, keyword string) bool {
	return lex.matchPrefix(strings.ToUpper(keyword), strings.ToLower(keyword)) &&
		isWhiteSpace(lex.peekAhead(len(keyword)+1))
}

func lexText(lex *lexer) StateFn {
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
		return lexJoin(lex, tokenInnerJoin, 5)
	} else if lex.matchPrefix("outer", "OUTER", "left", "LEFT", "right", "RIGHT") {
		return lexJoin(lex, tokenOuterJoin, 5)
	} else if lex.matchPrefix("left", "LEFT", "right", "RIGHT") {
		return lexJoin(lex, tokenLeftJoin, 4)
	} else if lex.matchPrefix("right", "RIGHT") {
		return lexJoin(lex, tokenRightJoin, 5)
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
			lex.emit(tokenText)
		}

		lex.emit(tokenEOF) //signal end of file token
		return nil         //stop run loop
	}

	//handle unspecified token
	return lex.errorf("Syntax error charactor(%s) not recognize by SQL lexical analysis; line %d, column %d",
		strconv.QuoteRune(nr1), lex.line, lex.pos+1)
}

//lexKeyword
func lexKeyword(lex *lexer, keyword string, tokenTypee tokenType) StateFn {
	lex.fastForward(len(keyword))
	lex.emit(tokenTypee)

	return lexText
}

func lexJoin(lex *lexer, tokenTy tokenType, keywordLen int) StateFn {

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

func lexGroupBy(lex *lexer) StateFn {
	keywordLen := 5 //group
	tokenTy := tokenGroupBy

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

func lexNumber(l *lexer) StateFn {
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
	l.emit(tokenNumber)
	return lexText
}

func lexParameter(lex *lexer) StateFn {
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
	lex.emit(tokenParameter)

	return lexText
}

func lexQuoteString(lex *lexer) StateFn {
	lex.next() //consume ' character

	for { //loop till reach ' charactor or EOF
		r := lex.next()

		if r == eof {
			return lex.errorf("Syntax error, quoted string not close before reach end of file")
		} else if r == '\'' {
			lex.emit(tokenString)
			return lexText
		}
	}
}

func lexLiteral(lex *lexer) StateFn {
	r := lex.next()

	hasBackQuote := r == '`'

	if hasBackQuote {
		for { //loop til reach backquote
			r = lex.next()
			if r == '`' {
				lex.emit(tokenLiteral)
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
				lex.emit(tokenLiteral)
				return lexText
			}
		}
	}
}
