package parser

import (
	"fmt"
)

//TokenType token type for SQL string
type TokenType uint8

//TokenItem single token definition
type TokenItem struct {
	Type  TokenType //token type
	Value string    //extracted string value
	Pos   int       //column index of the line
	line  int       //line number of the string
}

//Token type constants
const (
	TokenError        TokenType = iota
	TokenEOF                    //end of file
	TokenEqual                  // =
	TokenNotEqual               //<>, or !=
	TokenGreater                // >
	TokenGreaterEqual           // >=
	TokenLesser                 // <
	TokenLesserEqual            // <=
	TokenLike                   // LIKE
	TokenBetween                // BETWEEN
	TokenIn                     // IN
	TokenQuestionMark           // ?
	TokenWildcard               // %
	TokenAsterisk               // *
	TokenAdd                    // +
	TokenSubtract               // -
	TokenDivide                 // /
	TokenDot                    // .
	TokenLeftParen              // (
	TokenRightParen             // )
	TokenColon                  // ,
	TokenSemiColon              // ;
	TokenString                 // quoted string; 'sample'
	TokenNumber                 // number; -1.23 or 34.6
	TokenParameter              // parameter; :param1
	TokenLiteral                // `asd` or asd
	TokenSelect                 // SELECT keyword
	TokenFrom                   // FROM keyword
	TokenCreate                 // CREATE keyword
	TokenTable                  // TABLE keyword
	TokenView                   // VIEW keyword
	TokenWhere                  // WHERE keyword
	TokenGroupBy                // GROUP BY keyword
	TokenHaving                 // HAVING keyword
	TokenUnion                  // UNION keyword
	TokenJoin                   // JOIN keyword
	TokenInnerJoin              // INNER JOIN keyword
	TokenOuterJoin              // OUTER JOIN keyword
	TokenLeftJoin               // LEFT JOIN keyword
	TokenRightJoin              // RIGHT JOIN keyword
	TokenOn                     // ON keyword
	TokenLimit                  // LIMIT keyword
	TokenOffset                 // OFFSET keyword
	TokenAsc                    // ASC keyword
	TokenDesc                   // DESC keyword
	TokenDrop                   // DROP keyword
	TokenText                   // anonymous text token to take care EOF case (please refer Lexing.go -> LexText())

	//Agregate token --> SUM(), AVR(), MAX(),
)

func (item TokenItem) String() string {
	switch item.Type {
	case TokenError:
		return item.Value
	case TokenEOF:
		return "EOF"
	default:
		if len(item.Value) > 10 {
			return fmt.Sprintf("%.10q...", item.Value)
		}

		return fmt.Sprintf("%q", item.Value)
	}
}
