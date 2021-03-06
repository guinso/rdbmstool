package parser

//Lexical Analysis - convert string into a stack of tokens

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

//StateFn state function, representation of a fraction of lexing state machine
type StateFn func(item *lexer) StateFn

const eof = 0

var kw = []tokenItem{
	{TokenAs, "as", 0, 0},
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
	{TokenAnd, "and", 0, 0},
	{TokenOr, "or", 0, 0},
	{TokenNot, "not", 0, 0},
	//{tokenDistinct, "distinct", 0, 0},
	{TokenGroupBy, "group by", 0, 0},
}

var symbols = []struct {
	Type  TokenType
	Value string
}{
	{TokenEqual, "="},
	{TokenNotEqual, "<>"},
	{TokenNotEqual, "!="},
	{TokenGreater, ">"},
	{TokenGreaterEqual, ">="},
	{TokenLesser, "<"},
	{TokenLesserEqual, "<="},
	{TokenQuestionMark, "?"},
	{TokenWildcard, "%"},
	{TokenAsterisk, "*"},
	{TokenAdd, "+"},
	{TokenSubtract, "-"},
	{TokenDivide, "/"},
	{TokenLeftParen, "("},
	{TokenRightParen, ")"},
	{TokenDot, "."},
	{TokenColon, ","},
	{TokenSemiColon, ";"},
}

var fns = []tokenItem{
	{TokenSum, "sum", 0, 0},
	{TokenMin, "min", 0, 0},
	{TokenMax, "max", 0, 0},
	{TokenAvg, "avg", 0, 0},
	{TokenCount, "count", 0, 0},
	{TokenGreatest, "greatest", 0, 0},
	//{tokenDistinct, "distinct", 0, 0},
}

func isWhiteSpace(input rune) bool {
	return input == ' ' || input == '\n' || input == '\t'
}

func isSymbol(input rune) bool {
	for _, sym := range symbols {
		r, _ := utf8.DecodeRuneInString(sym.Value)
		if r == input {
			return true
		}
	}

	return false
}

func isAlphaNumeric(input rune) bool {
	return isNumeric(input) || isLetter(input)
}

func isLiteralCharacter(input rune) bool {
	return isAlphaNumeric(input) || input == '_'
}

func isNumeric(input rune) bool {
	return strings.IndexRune(
		"0123456789", input) != -1
}

func isLetter(input rune) bool {
	return strings.IndexRune(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", input) != -1
}

func isKeywordMatch(lex *lexer, keyword string) bool {
	aheadRune := lex.peekAhead(len(keyword) + 1)

	return lex.matchPrefix(strings.ToUpper(keyword), strings.ToLower(keyword)) &&
		(isWhiteSpace(aheadRune) || isSymbol(aheadRune) || aheadRune == eof)
}

func lexText(lex *lexer) StateFn {
	nr1 := lex.peekAhead(1)
	nr2 := lex.peekAhead(2)
	subString := lex.input[lex.start:]

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

	//looking for function if match
	for _, fn := range fns {
		if strings.HasPrefix(subString, strings.ToUpper(fn.Value)+"(") ||
			strings.HasPrefix(subString, strings.ToLower(fn.Value)+"(") {
			return lexFunction(lex, fn.Value, fn.Type)
		}
	}

	//handle complex keyword(s)
	if lex.matchPrefix("group", "GROUP") {
		return lexGroupBy(lex)
	} else if lex.matchPrefix("order", "ORDER") {
		return lexOrderBy(lex)
	} else if lex.matchPrefix("inner", "INNER") {
		return lexJoin(lex, TokenInnerJoin, 5)
	} else if lex.matchPrefix("outer", "OUTER") {
		return lexJoin(lex, TokenOuterJoin, 5)
	} else if lex.matchPrefix("left", "LEFT") {
		return lexJoin(lex, TokenLeftJoin, 4)
	} else if lex.matchPrefix("right", "RIGHT") {
		return lexJoin(lex, TokenRightJoin, 5)
	} else if lex.matchPrefix("join", "JOIN") {
		if xErr := lex.fastForward(4); xErr != nil {
			return lex.errorf("fail to tokenize JOIN token at %d", lex.pos)
		}

		lex.emit(TokenJoin)
		return lexText
	}

	//looking for SQL parameter
	if nr1 == ':' && isLetter(nr2) {
		return lexParameter(lex)
	}

	//looking for number
	if isNumeric(nr1) { //|| ((nr1 == '+' || nr1 == '-') && isNumeric(nr2)) {
		return lexNumber(lex)
	}

	//looking for quoted string
	if nr1 == '\'' {
		return lexQuoteString(lex)
	} else if nr1 == '"' {
		return lexQuoteDoubleString(lex)
	}

	//looking for literal (literal can start with backquote or without backquote)
	if (nr1 == '`' && isLetter(nr2)) || isLetter(nr1) {
		return lexLiteral(lex)
	}

	//looking for symbols if match
	for _, s := range symbols {
		if strings.HasPrefix(subString, s.Value) {
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
func lexKeyword(lex *lexer, keyword string, tokenTypee TokenType) StateFn {
	lex.fastForward(len(keyword))
	lex.emit(tokenTypee)

	return lexText
}

//lexFunction
func lexFunction(lex *lexer, keyword string, tokenTypee TokenType) StateFn {
	lex.fastForward(len(keyword))
	lex.emit(tokenTypee)

	return lexText
}

func lexJoin(lex *lexer, tokenTy TokenType, keywordLen int) StateFn {

	//fast forward the keyword
	if err := lex.fastForward(keywordLen); err != nil {
		return lex.errorf(err.Error())
	}

	//skip all white space
	for {
		r := lex.next()
		if lex.matchPrefix("join", "JOIN") {
			if err := lex.fastForward(4); err != nil {
				return lex.errorf(err.Error())
			}
			break
		} else if isWhiteSpace(r) {
			continue //keep looping
		} else if r == eof {
			//handle syntax error
			return lex.errorf("unexpected end of file reached at column %d", lex.pos)
		}

		if r := lex.peek(); !isWhiteSpace(r) {
			return lex.errorf("syntax error detected for JOIN statement, "+
				"expect white space after 'join' keyword at column %d", lex.pos)
		}
	}

	//fast forward lexer to end of 'join' keyword
	//lex.backup() //move backward one step

	lex.emit(tokenTy)
	return lexText
}

func lexGroupBy(lex *lexer) StateFn {
	keywordLen := 5 //group
	tokenTy := TokenGroupBy

	//fast forward the keyword
	if err := lex.fastForward(keywordLen); err != nil {
		return lex.errorf(err.Error())
	}

	for {
		r := lex.next()

		if lex.matchPrefix("by", "BY") {
			//fast forward lexer to end of 'by' keyword
			if err := lex.fastForward(2); err != nil {
				return lex.errorf(err.Error())
			}
			break
		} else if isWhiteSpace(r) {
			continue
		} else if r == eof {
			//handle syntax error
			return lex.errorf("unexpected end of file reached at column %d", lex.pos)
		}

		return lex.errorf("syntax error detected for Group By statement, "+
			"expect white space after 'by' keyword at column %d", lex.pos)
	}

	//fast forward lexer to end of 'by' keyword
	// if err := lex.fastForward(2); err != nil {
	// 	return lex.errorf(err.Error())
	// }

	// if r := lex.peek(); !isWhiteSpace(r) {
	// 	return lex.errorf("syntax error detected for Group By statement, "+
	// 		"expect white space after 'by' keyword at column %d", lex.pos)
	// }

	lex.emit(tokenTy)
	return lexText
}

func lexOrderBy(lex *lexer) StateFn {
	keywordLen := 5 //order
	tokenTy := TokenOrderBy

	//fast forward the keyword
	if err := lex.fastForward(keywordLen); err != nil {
		return lex.errorf(err.Error())
	}

	for {
		r := lex.next()
		if lex.matchPrefix("by", "BY") {
			//fast forward lexer to end of 'by' keyword
			if err := lex.fastForward(2); err != nil {
				return lex.errorf(err.Error())
			}

			break
		} else if isWhiteSpace(r) {
			continue
		} else if r == eof {
			//handle syntax error
			return lex.errorf("unexpected end of file reached at column %d", lex.pos)
		}

		if r := lex.peek(); !isWhiteSpace(r) {
			return lex.errorf("syntax error detected for Order By statement, "+
				"expect white space after 'by' keyword at column %d", lex.pos)
		}
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
	l.emit(TokenNumber)
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
	lex.emit(TokenParameter)

	return lexText
}

func lexQuoteString(lex *lexer) StateFn {
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

func lexQuoteDoubleString(lex *lexer) StateFn {
	lex.next() //consume ' character

	for { //loop till reach ' charactor or EOF
		r := lex.next()

		if r == eof {
			return lex.errorf("Syntax error, quoted string not close before reach end of file")
		} else if r == '"' {
			lex.emit(TokenString)
			return lexText
		}
	}
}

func lexLiteral(lex *lexer) StateFn {
	r := lex.next()

	if r == '`' { //has backquote character
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
				lex.backup() //backward one rune since latest run is not valid literal
				lex.emit(TokenLiteral)
				return lexText
			}
		}
	}
}
