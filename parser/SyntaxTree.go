package parser

//SyntaxTree data structure to keep AST
type SyntaxTree struct {
	childNodes    []SyntaxTree
	StartPosition int
	EndPosition   int
	Source        []tokenItem
	DataType      string
}
