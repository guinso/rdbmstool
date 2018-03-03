package parser

import (
	"fmt"
)

//TokenType token type for SQL string
type TokenType uint8

//TokenItem single token definition
type tokenItem struct {
	Type  TokenType //token type
	Value string    //extracted string value
	Pos   int       //column index of the line
	line  int       //line number of the string
}

//Token type constants
const (
	TokenError        TokenType = iota //error token identifier
	TokenEOF                           //end of file
	TokenEqual                         // =
	TokenNotEqual                      //<>, or !=
	TokenGreater                       // >
	TokenGreaterEqual                  // >=
	TokenLesser                        // <
	TokenLesserEqual                   // <=
	TokenLike                          // LIKE
	TokenBetween                       // BETWEEN
	TokenIn                            // IN
	TokenQuestionMark                  // ?
	TokenWildcard                      // %
	TokenAsterisk                      // *
	TokenAdd                           // +
	TokenSubtract                      // -
	TokenDivide                        // /
	TokenDot                           // .
	TokenLeftParen                     // (
	TokenRightParen                    // )
	TokenColon                         // ,
	TokenSemiColon                     // ;
	TokenString                        // quoted string; 'sample'
	TokenNumber                        // number; -1.23 or 34.6
	TokenParameter                     // parameter; :param1
	TokenLiteral                       // `asd` or asd
	TokenSelect                        // SELECT keyword
	TokenFrom                          // FROM keyword
	TokenCreate                        // CREATE keyword
	TokenTable                         // TABLE keyword
	TokenView                          // VIEW keyword
	TokenWhere                         // WHERE keyword
	TokenGroupBy                       // GROUP BY keyword
	TokenOrderBy                       // ORDER BY keyword
	TokenHaving                        // HAVING keyword
	TokenUnion                         // UNION keyword
	TokenJoin                          // JOIN keyword
	TokenInnerJoin                     // INNER JOIN keyword
	TokenOuterJoin                     // OUTER JOIN keyword
	TokenLeftJoin                      // LEFT JOIN keyword
	TokenRightJoin                     // RIGHT JOIN keyword
	TokenOn                            // ON keyword
	TokenLimit                         // LIMIT keyword
	TokenOffset                        // OFFSET keyword
	TokenAsc                           // ASC keyword
	TokenDesc                          // DESC keyword
	TokenDrop                          // DROP keyword
	TokenText                          // anonymous text token to take care EOF case (please refer Lexing.go -> LexText())
	TokenAnd                           // AND keyword
	TokenOr                            // OR keyword
	TokenNot                           // not keyword
	TokenAs                            // AS keyword
	TokenMin                           // MIN()
	TokenMax                           //MAX()
	TokenGreatest                      //GREATEST()
	TokenCount                         //COUNT()
	TokenAvg                           //AVG()
	TokenSum                           //SUM()
	//tokenDistinct               // distinct keyword

)

func (item tokenItem) String() string {
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

func (item TokenType) String() string {
	switch item {
	case TokenAdd:
		return "add"
	case TokenAnd:
		return "and"
	case TokenAs:
		return "as"
	case TokenAsc:
		return "asc"
	case TokenAsterisk:
		return "*"
	case TokenAvg:
		return "avg"
	case TokenBetween:
		return "between"
	case TokenColon:
		return ","
	case TokenCount:
		return "count"
	case TokenCreate:
		return "create"
	case TokenDesc:
		return "desc"
	// case tokenDistinct:
	// 	return "distinct"
	case TokenDivide:
		return "/"
	case TokenDot:
		return "."
	case TokenDrop:
		return "drop"
	case TokenEOF:
		return "EOF"
	case TokenEqual:
		return "="
	case TokenError:
		return "ERROR"
	case TokenFrom:
		return "from"
	case TokenGreater:
		return ">"
	case TokenGreaterEqual:
		return ">="
	case TokenGreatest:
		return "greatest"
	case TokenGroupBy:
		return "group-by"
	case TokenOrderBy:
		return "order-by"
	case TokenHaving:
		return "having"
	case TokenIn:
		return "in"
	case TokenInnerJoin:
		return "inner-join"
	case TokenJoin:
		return "join"
	case TokenLeftJoin:
		return "left-join"
	case TokenLeftParen:
		return "("
	case TokenLesser:
		return "<"
	case TokenLesserEqual:
		return "<="
	case TokenLike:
		return "like"
	case TokenLimit:
		return "limit"
	case TokenLiteral:
		return "literal"
	case TokenMax:
		return "max"
	case TokenMin:
		return "min"
	case TokenNot:
		return "not"
	case TokenNotEqual:
		return "<>"
	case TokenNumber:
		return "numeric"
	case TokenOffset:
		return "offset"
	case TokenOn:
		return "on"
	case TokenOr:
		return "or"
	case TokenOuterJoin:
		return "outer-join"
	case TokenParameter:
		return "parameter"
	case TokenQuestionMark:
		return "?"
	case TokenRightJoin:
		return "right-join"
	case TokenRightParen:
		return ")"
	case TokenSelect:
		return "select"
	case TokenSemiColon:
		return ";"
	case TokenString:
		return "string"
	case TokenSubtract:
		return "-"
	case TokenSum:
		return "sum"
	case TokenTable:
		return "table"
	case TokenText:
		return "text"
	case TokenUnion:
		return "union"
	case TokenView:
		return "view"
	case TokenWhere:
		return "where"
	case TokenWildcard:
		return "%"
	default:
		return "undefined"
	}
}
