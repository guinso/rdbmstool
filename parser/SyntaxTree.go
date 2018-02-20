package parser

import (
	"fmt"
)

//SyntaxTree data structure to keep AST
type SyntaxTree struct {
	ChildNodes    []SyntaxTree
	StartPosition int
	EndPosition   int
	Source        []tokenItem
	DataType      NodeType
}

//RawString generate input text based on start position and end position
func (ast *SyntaxTree) RawString() string {
	if ast.Source != nil &&
		len(ast.Source) > ast.StartPosition &&
		len(ast.Source) > ast.EndPosition &&
		ast.StartPosition <= ast.EndPosition {
		result := ""
		for i := ast.StartPosition; i <= ast.EndPosition; i++ {
			result = result + ast.Source[i].String() + " "
		}

		return result
	}

	return ""
}

// NodeType identifies the type of a parse tree node.
type NodeType int

const (
	//NodeBool a boolean constant
	NodeBool NodeType = iota
	//NodeField a db field such as table name, column name
	NodeField
	//NodeIdentifier a function name
	NodeIdentifier
	//NodeNumber a numerical constant
	NodeNumber
	//NodeString a string constant
	NodeString
	//NodeParam an SQL parameter
	NodeParam
	//NodeList a list of SQL expression
	NodeList
	//NodeCondition a logical comparison expression; e.g. x > y
	NodeCondition
	//NodeSelect SQL select statement
	NodeSelect
	//NodeFrom SQL from statement
	NodeFrom
	//NodeJoin SQL join statement
	NodeJoin
	//NodeWhere SQL where statement
	NodeWhere
	//NodeHaving SQL having statement
	NodeHaving
	//NodeGroupBy SQL group by statement
	NodeGroupBy
	//NodeOrderBy SQL order by statement
	NodeOrderBy
	//NodeLimit SQl limit statement
	NodeLimit
	//NodeUnion SQL union statement
	NodeUnion
	//NodeQuery SQL query statement with union
	NodeQuery
	//NodeQuerySelect SQL SELECT query statement
	NodeQuerySelect
	//NodeSource source selector
	NodeSource
	//NodeOperator oprator token
	NodeOperator
	//NodeUnaryOperator unary operator token
	NodeUnaryOperator
	//NodeOperand operand statement
	NodeOperand
	//NodeParenthesis parenthesis statement
	NodeParenthesis
	//NodeAlias alias for source selector
	NodeAlias
	//NodeExpression SQL expression statement
	NodeExpression
	//NodeFunction SQL function statement
	NodeFunction
	//NodeColName column name selector
	NodeColName
	//NodeColumn column source selector
	NodeColumn
	//NodeOrder order token (acending/descending)
	NodeOrder
)

//ParseSQL parse SQL string input into abstract syntax tree
//currently only support query syntax
func ParseSQL(inputText string) (*SyntaxTree, error) {
	tokens := tokenize(inputText)

	if tokens == nil || len(tokens) == 0 {
		return nil, fmt.Errorf("input text has no matching token to tokenize")
	}

	return parseSelect(tokens, 0)
}
