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
	tokenOrderBy                // ORDER BY keyword
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
	tokenAnd                    // AND keyword
	tokenOr                     // OR keyword
	tokenNot                    // not keyword
	tokenDistinct               // distinct keyword
	tokenAs                     // AS keyword
	tokenMin                    // MIN()
	tokenMax                    //MAX()
	tokenGreatest               //GREATEST()
	tokenCount                  //COUNT()
	tokenAvg                    //AVG()
	tokenSum                    //SUM()

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

func (item tokenType) String() string {
	switch item {
	case tokenAdd:
		return "add"
	case tokenAnd:
		return "and"
	case tokenAs:
		return "as"
	case tokenAsc:
		return "asc"
	case tokenAsterisk:
		return "*"
	case tokenAvg:
		return "avg"
	case tokenBetween:
		return "between"
	case tokenColon:
		return ","
	case tokenCount:
		return "count"
	case tokenCreate:
		return "create"
	case tokenDesc:
		return "desc"
	case tokenDistinct:
		return "distinct"
	case tokenDivide:
		return "/"
	case tokenDot:
		return "."
	case tokenDrop:
		return "drop"
	case tokenEOF:
		return "EOF"
	case tokenEqual:
		return "="
	case tokenError:
		return "ERROR"
	case tokenFrom:
		return "from"
	case tokenGreater:
		return ">"
	case tokenGreaterEqual:
		return ">="
	case tokenGreatest:
		return "greatest"
	case tokenGroupBy:
		return "group-by"
	case tokenOrderBy:
		return "order-by"
	case tokenHaving:
		return "having"
	case tokenIn:
		return "in"
	case tokenInnerJoin:
		return "inner-join"
	case tokenJoin:
		return "join"
	case tokenLeftJoin:
		return "left-join"
	case tokenLeftParen:
		return "("
	case tokenLesser:
		return "<"
	case tokenLesserEqual:
		return "<="
	case tokenLike:
		return "like"
	case tokenLimit:
		return "limit"
	case tokenLiteral:
		return "literal"
	case tokenMax:
		return "max"
	case tokenMin:
		return "min"
	case tokenNot:
		return "not"
	case tokenNotEqual:
		return "<>"
	case tokenNumber:
		return "numeric"
	case tokenOffset:
		return "offset"
	case tokenOn:
		return "on"
	case tokenOr:
		return "or"
	case tokenOuterJoin:
		return "outer-join"
	case tokenParameter:
		return "parameter"
	case tokenQuestionMark:
		return "?"
	case tokenRightJoin:
		return "right-join"
	case tokenRightParen:
		return ")"
	case tokenSelect:
		return "select"
	case tokenSemiColon:
		return ";"
	case tokenString:
		return "string"
	case tokenSubtract:
		return "-"
	case tokenSum:
		return "sum"
	case tokenTable:
		return "table"
	case tokenText:
		return "text"
	case tokenUnion:
		return "union"
	case tokenView:
		return "view"
	case tokenWhere:
		return "where"
	case tokenWildcard:
		return "%"
	default:
		return "undefined"
	}
}
