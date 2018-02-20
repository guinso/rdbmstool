package parser

//BuildAST build SQL abstract syntax tree
func BuildAST(sql string) {
	lexer := lex("build-AST", sql)

	tokens := []tokenItem{}

	tmpToken := lexer.nextItem()
	for tmpToken.Type != tokenEOF {
		tokens = append(tokens, tmpToken)

		tmpToken = lexer.nextItem()
	}

	convertAST(tokens)
}

func convertAST(tokens []tokenItem) {
	tmpToken := tokens[0]

	if tmpToken.Type == tokenCreate {

	} else if tmpToken.Type == tokenSelect {

	} else if tmpToken.Type == tokenUnion {

	} else if tmpToken.Type == tokenSemiColon {

	} else {
		//TODO: handle syntax error
	}
}
