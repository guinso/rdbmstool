package parser

import (
	"fmt"
)

//tokenType token type for SQL string
type tokenType uint8

//tokenItem single token definition
type tokenItem struct {
	Type  tokenType //token type
	Value string    //extracted string value
	Pos   int       //column index of the line
	line  int       //line number of the string
}

//Token type constants
const (
	tokenError        tokenType = iota
	tokenEOF                    //end of file
	tokenEqual                  // =
	tokenNotEqual               //<>, or !=
	tokenGreater                // >
	tokenGreaterEqual           // >=
	tokenLesser                 // <
	tokenLesserEqual            // <=
	tokenLike                   // LIKE
	tokenBetween                // BETWEEN
	tokenIn                     // IN
	tokenQuestionMark           // ?
	tokenWildcard               // %
	tokenAsterisk               // *
	tokenAdd                    // +
	tokenSubtract               // -
	tokenDivide                 // /
	tokenDot                    // .
	tokenLeftParen              // (
	tokenRightParen             // )
	tokenColon                  // ,
	tokenSemiColon              // ;
	tokenString                 // quoted string; 'sample'
	tokenNumber                 // number; -1.23 or 34.6
	tokenParameter              // parameter; :param1
	tokenLiteral                // `asd` or asd
	tokenSelect                 // SELECT keyword
	tokenFrom                   // FROM keyword
	tokenCreate                 // CREATE keyword
	tokenTable                  // TABLE keyword
	tokenView                   // VIEW keyword
	tokenWhere                  // WHERE keyword
	tokenGroupBy                // GROUP BY keyword
	tokenHaving                 // HAVING keyword
	tokenUnion                  // UNION keyword
	tokenJoin                   // JOIN keyword
	tokenInnerJoin              // INNER JOIN keyword
	tokenOuterJoin              // OUTER JOIN keyword
	tokenLeftJoin               // LEFT JOIN keyword
	tokenRightJoin              // RIGHT JOIN keyword
	tokenOn                     // ON keyword
	tokenLimit                  // LIMIT keyword
	tokenOffset                 // OFFSET keyword
	tokenAsc                    // ASC keyword
	tokenDesc                   // DESC keyword
	tokenDrop                   // DROP keyword
	tokenText                   // anonymous text token to take care EOF case (please refer Lexing.go -> LexText())

	//Agregate token --> SUM(), AVR(), MAX(),
)

func (item tokenItem) String() string {
	switch item.Type {
	case tokenError:
		return item.Value
	case tokenEOF:
		return "EOF"
	default:
		if len(item.Value) > 10 {
			return fmt.Sprintf("%.10q...", item.Value)
		}

		return fmt.Sprintf("%q", item.Value)
	}
}
