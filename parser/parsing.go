package parser

import (
	"errors"
	"fmt"
)

func parseParenthesis(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//pattern:
	//expr = ()
	//expr = (<*>)
	//<*> = <*><*>

	startPos := -1
	stopPos := -1

	tmpStartToken := []tokenItem{}

	if source[startIndex].Type != tokenLeftParen {
		return nil, errors.New("start token is not left parenthesis")
	}

	startPos = startIndex
	tmpStartToken = append(tmpStartToken, source[startIndex])

	for i := (startIndex + 1); i < len(source); i++ {
		if source[i].Type == tokenLeftParen {
			tmpStartToken = append(tmpStartToken, source[i])
		} else if source[i].Type == tokenRightParen {

			if len(tmpStartToken) > 1 {
				tmpStartToken = tmpStartToken[:len(tmpStartToken)-1]
			} else {
				tmpStartToken = []tokenItem{}

				stopPos = i
				break
			}
		} else if source[i].Type == tokenEOF {
			break
		}
	}

	if len(tmpStartToken) == 0 {
		return &SyntaxTree{
			childNodes:    []SyntaxTree{},
			StartPosition: startPos,
			EndPosition:   stopPos,
			Source:        source,
			DataType:      "parenthesis"}, nil
	}

	return nil, fmt.Errorf(
		"Syntax error, parenthesis is not complete at line %d, position %d",
		tmpStartToken[len(tmpStartToken)-1].line,
		tmpStartToken[len(tmpStartToken)-1].Pos)
}

func parseField(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//pattern:
	//expr = <expression>
	//expr = (<expression>)
	//expr = <expression> AS <literal>
	//expr = (<expression>) AS <literal>

	//TODO: try parse by expression
	ast, astErr := parseExpresion(source, startIndex)
	if astErr != nil {
		return nil, astErr
	}

	nodes := []SyntaxTree{*ast}

	endIndex := ast.EndPosition

	if len(source) > (ast.EndPosition + 2) {
		if source[ast.EndPosition+1].Type == tokenAs &&
			source[ast.EndPosition+2].Type == tokenLiteral {
			endIndex = ast.EndPosition + 2

			nodes = append(nodes, SyntaxTree{
				childNodes:    []SyntaxTree{},
				StartPosition: ast.StartPosition + 2,
				EndPosition:   ast.StartPosition + 2,
				Source:        source,
				DataType:      "alias",
			})
		}
	}

	return &SyntaxTree{
		childNodes:    nodes,
		StartPosition: startIndex,
		EndPosition:   endIndex,
		Source:        source,
		DataType:      "field",
	}, nil
}

func parseExpresion(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//pattern:
	//expr = <column>
	//expr = <number>
	//expr = <string>
	//expr = (<expr>)
	//expr = <operator><expr>
	//expr = <expr><operator><expr>
	//expr = <fn>

	expectOperator := false
	endIndex := -1
	unarySign := 0

	nodes := []SyntaxTree{}

	//complete at EOF or next token not match pattern
	for i := startIndex; i < len(source); i++ {

		if expectOperator {
			if isOperatorToken(source[i]) {
				expectOperator = false
				unarySign = 0
				nodes = append(nodes, SyntaxTree{
					childNodes:    []SyntaxTree{},
					StartPosition: i,
					EndPosition:   i,
					Source:        source,
					DataType:      "operator",
				})
				continue
			}

			endIndex = i - 1
			break
			//return nil, fmt.Errorf(
			//	"expect to get operator token at position %d, but it is not", source[i].Pos)
		}

		currentLoop := i

		//check is operand or not
		if isOperandToken(source[i]) {

			if source[i].Type == tokenLiteral &&
				(i+2) < len(source) &&
				source[i+1].Type == tokenDot &&
				(source[i+2].Type == tokenLiteral || source[i+2].Type == tokenAsterisk) {
				i += 2
			}

			expectOperator = true
			if unarySign == 1 {
				nodes = append(nodes, SyntaxTree{
					childNodes:    []SyntaxTree{},
					StartPosition: currentLoop - 1,
					EndPosition:   currentLoop,
					Source:        source,
					DataType:      "unary-operator",
				})

				unarySign = 0
			}
			nodes = append(nodes, SyntaxTree{
				childNodes:    []SyntaxTree{},
				StartPosition: currentLoop,
				EndPosition:   i,
				Source:        source,
				DataType:      "operand",
			})
			continue
		} else if source[i].Type == tokenAdd || source[i].Type == tokenSubtract {
			if unarySign > 0 {
				return nil, fmt.Errorf(
					"invalid syntax found at position %d (%s)",
					source[i].Pos, source[i].String())
			}

			unarySign++
			continue
		} else if isFunctionToken(source[i]) {
			funcc, funccErr := parseFunction(source, i)
			if funccErr != nil {
				return nil, funccErr
			}

			i = funcc.EndPosition
			expectOperator = true
			if unarySign == 1 {
				nodes = append(nodes, SyntaxTree{
					childNodes:    []SyntaxTree{},
					StartPosition: currentLoop - 1,
					EndPosition:   currentLoop,
					Source:        source,
					DataType:      "unary-operator",
				})

				unarySign = 0
			}
			nodes = append(nodes, *funcc)
			expectOperator = true
			continue
		} else if source[i].Type == tokenLeftParen {
			//test parse parenthesis
			tmpBracket, bracketErr := parseParenthesis(source, i)
			if bracketErr != nil {
				return nil, fmt.Errorf(
					"syntax error, incomplete parenthesis found at position %d (%s)", source[i].Pos, source[i].Value)
			}

			//try parse inner expression
			expr2, exprErr := parseExpresion(source, i+1)
			if exprErr != nil {
				return nil, exprErr
			}
			if expr2.EndPosition+1 != tmpBracket.EndPosition {
				return nil, fmt.Errorf(
					"expect expression ended at position %d, but get at %d instead",
					tmpBracket.EndPosition-1,
					expr2.EndPosition)
			}

			i = tmpBracket.EndPosition
			expectOperator = true
			if unarySign == 1 {
				nodes = append(nodes, SyntaxTree{
					childNodes:    []SyntaxTree{},
					StartPosition: currentLoop - 1,
					EndPosition:   currentLoop,
					Source:        source,
					DataType:      "unary-operator",
				})

				unarySign = 0
			}
			nodes = append(nodes, *expr2)
			continue
		}

		return nil, fmt.Errorf(
			"expect to get operand token at position %d, but it is not (%s)", source[i].Pos, source[i].String())
	}

	if endIndex == -1 {
		return nil, fmt.Errorf("no expression syntax found from position %d", startIndex)
	}

	return &SyntaxTree{
		childNodes:    nodes,
		StartPosition: startIndex,
		EndPosition:   endIndex,
		Source:        source,
		DataType:      "expression",
	}, nil
}

func isOperatorToken(item tokenItem) bool {
	return item.Type == tokenAdd || item.Type == tokenBetween ||
		item.Type == tokenDivide || item.Type == tokenEqual ||
		item.Type == tokenGreater || item.Type == tokenGreaterEqual ||
		item.Type == tokenLesser || item.Type == tokenLesserEqual ||
		item.Type == tokenLike || item.Type == tokenNot ||
		item.Type == tokenNotEqual || item.Type == tokenSubtract ||
		item.Type == tokenAsterisk
}

func isOperandToken(item tokenItem) bool {
	return item.Type == tokenLiteral ||
		item.Type == tokenNumber ||
		item.Type == tokenString ||
		item.Type == tokenAsterisk
}

func isFunctionToken(item tokenItem) bool {
	return item.Type == tokenAvg ||
		item.Type == tokenCount ||
		//item.Type == tokenDistinct ||
		item.Type == tokenMax ||
		item.Type == tokenMin ||
		item.Type == tokenSum
}

func parseFunction(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//TODO: expand this function to verify various functions
	//e.g. MAX(), MIN(), COUNT(), AVG(), etc.

	//pattern:
	//expr = min(<literal>)
	//expr = max(<literal>)
	//expr = count(<literal>)
	//expr = avg(<literal>)
	//expr = sum(<literal>)
	//TODO: expr = if(<expr>, <expr>, <expr>)
	if !isFunctionToken(source[startIndex]) {
		return nil, fmt.Errorf("no function syntax found at position %d (%s)",
			startIndex,
			source[startIndex].String())
	}

	if len(source) < (startIndex + 1) {
		return nil, fmt.Errorf("incomplete function syntax found at position %d", startIndex)
	}

	paren, parenErr := parseParenthesis(source, startIndex+1)
	if parenErr != nil {
		return nil, fmt.Errorf("no complete parenthesis found at position %d", startIndex+1)
	}

	expr, exprErr := parseExpresion(source, paren.StartPosition+1)
	if exprErr != nil {
		return nil, exprErr
	}

	if paren.EndPosition != (expr.EndPosition + 1) {
		return nil, fmt.Errorf(
			"expect after expression (position %d) will follow by close parenthesis (position %d)",
			expr.EndPosition,
			paren.EndPosition)
	}

	return &SyntaxTree{
		childNodes:    []SyntaxTree{*expr},
		StartPosition: startIndex,
		EndPosition:   paren.EndPosition,
		Source:        source,
		DataType:      "function",
	}, nil
}

func parseJoin(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//patern:
	//expr = <join> <src>
	//expr = <join> <src> ON <condition>
	//join = JOIN
	//join = LEFT JOIN
	//join = RIGHT JOIN
	//join = INNER JOIN
	//src = <literal>
	//src = <literal> <literal>
	//src = <literal> AS <literal>
	//src = <selectExpr> <literal>
	//src = <selectExpr> AS <literal>

	sourceLen := len(source)
	endPos := -1
	nodes := []SyntaxTree{}

	tmp := source[startIndex]
	if !isJoinToken(tmp) {
		return nil, fmt.Errorf("Expect position %d is Join token but it is not", tmp.Pos)
	}

	//TODO: support select expression

	if sourceLen <= (startIndex + 1) {
		return nil, fmt.Errorf("incomplete join syntax at position %d (%s)",
			source[startIndex].Pos,
			source[startIndex].String())
	}

	//source
	tmp1 := source[startIndex+1]
	if tmp1.Type != tokenLiteral {
		return nil, fmt.Errorf(
			"syntax error found next to JOIN token at position %d (%s)",
			tmp1.Pos,
			tmp1.String())
	}
	endPos = startIndex + 1

	if sourceLen > (endPos+2) &&
		source[endPos+1].Type == tokenDot &&
		source[endPos+2].Type == tokenLiteral {
		endPos += 2
	}

	nodes = append(nodes, SyntaxTree{
		childNodes:    []SyntaxTree{},
		StartPosition: startIndex + 1,
		EndPosition:   endPos,
		Source:        source,
		DataType:      "source",
	})

	//AS
	if sourceLen > (endPos+2) &&
		source[endPos+1].Type == tokenAs &&
		source[endPos+2].Type == tokenLiteral {
		endPos += 2

		nodes = append(nodes, SyntaxTree{
			childNodes:    []SyntaxTree{},
			StartPosition: endPos,
			EndPosition:   endPos,
			Source:        source,
			DataType:      "alias",
		})
	}

	//ON
	if sourceLen > (endPos+1) &&
		source[endPos+1].Type == tokenOn {
		// if sourceLen < (endPos + 2) {
		// 	return nil, fmt.Errorf("Syntax error after ON token at position %d (%s)",
		// 		source[endPos+1].Pos,
		// 		source[endPos+1].String())
		// }

		cond, condErr := parseCondition(source, endPos+2)
		if condErr != nil {
			return nil, condErr
		}

		endPos = cond.EndPosition
		nodes = append(nodes, *cond)
	}

	return &SyntaxTree{
		childNodes:    nodes,
		StartPosition: startIndex,
		EndPosition:   endPos,
		Source:        source,
		DataType:      "join",
	}, nil
}

func isJoinToken(token tokenItem) bool {
	return token.Type == tokenJoin || token.Type == tokenLeftJoin ||
		token.Type == tokenInnerJoin || token.Type == tokenOuterJoin ||
		token.Type == tokenRightJoin
}

func parseCondition(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//pattern:
	//expr = <expression>
	//expr = (<expr>)
	//expr = <expr> AND <expr>
	//expr = <expr> OR <expr>

	checkCondSymbol := false
	endIndex := -1
	nodes := []SyntaxTree{}
	currentLoop := -1

	for i := startIndex; i < len(source); i++ {
		currentLoop = i

		//check have AND / OR token present or not
		if checkCondSymbol {
			if source[i].Type == tokenAnd || source[i].Type == tokenOr {
				checkCondSymbol = false
				endIndex = i
				nodes = append(nodes, SyntaxTree{
					childNodes:    []SyntaxTree{},
					StartPosition: currentLoop,
					EndPosition:   i,
					Source:        source,
					DataType:      source[i].Type.String(),
				})
				continue
			}

			//since no matching token found, can terminate loop
			break
		}

		if source[i].Type == tokenLeftParen { //check is parenthesis or not
			paren, parenErr := parseParenthesis(source, i)
			if parenErr != nil {
				return nil, parenErr
			}

			subCond, condErr := parseCondition(source, i+1)
			if condErr != nil {
				return nil, condErr
			}

			if paren.EndPosition != subCond.EndPosition+1 {
				return nil, fmt.Errorf(
					"sub condition expression not ended at %d, but %d (%s)",
					paren.EndPosition-1,
					subCond.EndPosition,
					source[subCond.EndPosition].String())
			}

			checkCondSymbol = true
			i = paren.EndPosition
			endIndex = paren.EndPosition
			nodes = append(nodes, *subCond)
			continue
		}

		expr, exprErr := parseExpresion(source, i)
		if exprErr == nil { //check is expression or not
			checkCondSymbol = true
			i = expr.EndPosition
			endIndex = expr.EndPosition
			nodes = append(nodes, *expr)
			continue
		}

		//syntax error, return error
		return nil, fmt.Errorf(
			"Expect expression syntax at position %d (%s)",
			source[i].Pos,
			source[i].String())
	}

	return &SyntaxTree{
		childNodes:    nodes,
		StartPosition: startIndex,
		EndPosition:   endIndex,
		Source:        source,
		DataType:      "condition",
	}, nil
}

func parseSelect(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//pattern:
	//expr = SELECT <cols>
	//cols = <expression>
	//cols = <expression>, <cols>
	if source[startIndex].Type != tokenSelect {
		return nil, fmt.Errorf(
			"Expect token SELECT but get %s instead at position %d",
			source[startIndex].String(),
			source[startIndex].Pos)
	}

	index := startIndex + 1
	nodes := []SyntaxTree{}

	checkColon := false
	for i := index; i < len(source); i++ {
		if checkColon == true {
			if source[i].Type == tokenColon {
				checkColon = false
				index = i
				continue
			}

			//no more selectable columns
			break
		}

		col, colErr := parseExpresion(source, i)
		if colErr != nil {
			return nil, colErr
		}

		checkColon = true
		index = col.EndPosition
		i = col.EndPosition

		colNodes := []SyntaxTree{*col}

		//check has alias or not
		if len(source) > (i + 1) {
			//log.Printf("i = %d, (%s)", i+1, source[i+1].String())
			if source[i+1].Type == tokenLiteral {
				i++
				colNodes = append(colNodes, SyntaxTree{
					childNodes:    []SyntaxTree{},
					StartPosition: i,
					EndPosition:   i,
					Source:        source,
					DataType:      "alias",
				})
				nodes = append(nodes, SyntaxTree{
					childNodes:    colNodes,
					StartPosition: col.StartPosition,
					EndPosition:   i,
					Source:        source,
					DataType:      "column",
				})
				index = i
				continue
			} else if source[i+1].Type == tokenAs {
				if len(source) > (i + 2) {

					if source[i+2].Type == tokenLiteral {
						i += 2
						colNodes = append(colNodes, SyntaxTree{
							childNodes:    []SyntaxTree{},
							StartPosition: i,
							EndPosition:   i,
							Source:        source,
							DataType:      "alias",
						})
						nodes = append(nodes, SyntaxTree{
							childNodes:    colNodes,
							StartPosition: col.StartPosition,
							EndPosition:   i,
							Source:        source,
							DataType:      "column",
						})
						index = i
						continue
					}

					return nil, fmt.Errorf(
						"Expect alias name after token AS but found %s; position %d",
						source[i+2].String(),
						source[i+2].Pos)
				}

				return nil, fmt.Errorf("syntax error; unexpected content ended with token AS")
			}
		}

		nodes = append(nodes, SyntaxTree{
			childNodes:    colNodes,
			StartPosition: col.StartPosition,
			EndPosition:   col.EndPosition,
			Source:        source,
			DataType:      "column",
		})
		i = col.EndPosition
	}

	return &SyntaxTree{
		childNodes:    nodes,
		StartPosition: startIndex,
		EndPosition:   index,
		Source:        source,
		DataType:      "select",
	}, nil
}

func parseFrom(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//pattern:
	//expr = FROM <src>
	//expr = FROM <src> <literal>
	//src = <literal>
	//src = <literal>.<literal>
	//src = <selectExpr>
	if source[startIndex].Type != tokenFrom {
		return nil, fmt.Errorf(
			"Syntax error expect token SELECT at position %d", source[startIndex].Pos)
	}

	nodes := []SyntaxTree{}
	index := startIndex + 1
	var bracket *SyntaxTree
	var fromSource *SyntaxTree

	if len(source) > index && source[index].Type == tokenLeftParen {
		tmp, bracketErr := parseParenthesis(source, index)
		if bracketErr != nil {
			return nil, bracketErr
		}

		bracket = tmp
		index++
	}

	if expr, exprErr := parseQuerySelect(source, index); exprErr == nil {
		fromSource = &SyntaxTree{
			childNodes:    []SyntaxTree{*expr},
			StartPosition: bracket.StartPosition,
			EndPosition:   bracket.EndPosition,
			Source:        source,
			DataType:      "source",
		}
		index = expr.EndPosition
	} else if len(source) > (index+3) &&
		source[index].Type == tokenLiteral &&
		source[index+1].Type == tokenDot &&
		source[index+2].Type == tokenLiteral {
		index += 2
		fromSource = &SyntaxTree{
			childNodes:    []SyntaxTree{},
			StartPosition: startIndex + 1,
			EndPosition:   index,
			Source:        source,
			DataType:      "source",
		}
	} else if len(source) > (index) && source[index].Type == tokenLiteral {
		fromSource = &SyntaxTree{
			childNodes:    []SyntaxTree{},
			StartPosition: index,
			EndPosition:   index,
			Source:        source,
			DataType:      "source",
		}
	} else {
		return nil, fmt.Errorf(
			"syntax error; no source found for FROM token at position %d", source[startIndex+1].Pos)
	}

	if bracket != nil && fromSource.EndPosition+1 != bracket.EndPosition {
		return nil, fmt.Errorf("expect inner query ended next to close bracket at position %d (%s)",
			source[fromSource.EndPosition].Pos,
			source[fromSource.EndPosition].String())
	}

	nodes = append(nodes, *fromSource)
	if bracket != nil {
		index++
	}

	//check for alias
	if len(source) > index+1 && source[index+1].Type == tokenLiteral {
		index++
		nodes = append(nodes, SyntaxTree{
			childNodes:    []SyntaxTree{},
			StartPosition: index - 1,
			EndPosition:   index,
			Source:        source,
			DataType:      "alias",
		})
	} else if len(source) > index+2 &&
		source[index+1].Type == tokenAs &&
		source[index+2].Type == tokenLiteral {
		index += 2
		nodes = append(nodes, SyntaxTree{
			childNodes:    []SyntaxTree{},
			StartPosition: index - 2,
			EndPosition:   index,
			Source:        source,
			DataType:      "alias",
		})
	}

	return &SyntaxTree{
		childNodes:    nodes,
		StartPosition: startIndex,
		EndPosition:   index,
		Source:        source,
		DataType:      "from",
	}, nil
}

func parseWhere(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//pattern
	//expr = WHERE <condition>

	if source[startIndex].Type != tokenWhere {
		return nil, fmt.Errorf(
			"Expect token WHERE found at position %d", source[startIndex].Pos)
	}

	index := startIndex + 1

	condition, condErr := parseCondition(source, index)
	if condErr != nil {
		return nil, condErr
	}

	return &SyntaxTree{
		childNodes:    []SyntaxTree{*condition},
		StartPosition: startIndex,
		EndPosition:   condition.EndPosition,
		Source:        source,
		DataType:      "where",
	}, nil
}

func parseGroupBy(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//pattern
	//expr = GROUP BY <cols>
	//cols = <col>, <cols>
	//cols = <col> <order>, <cols>
	//cols = <col>
	//cols = <col> <order>
	//order = ASC
	//order = DESC

	if source[startIndex].Type != tokenGroupBy {
		return nil, fmt.Errorf("Expect token GROUP BY at position %d", source[startIndex].Pos)
	}

	srcLen := len(source)
	nodes := []SyntaxTree{}
	index := startIndex + 1

	var col, order *SyntaxTree

	for index < srcLen {
		tmpStartIndex := index
		col = nil
		order = nil

		if srcLen > index+2 &&
			source[index].Type == tokenLiteral &&
			source[index+1].Type == tokenDot &&
			source[index+2].Type == tokenLiteral {
			col = &SyntaxTree{
				childNodes:    []SyntaxTree{},
				StartPosition: index,
				EndPosition:   index + 2,
				Source:        source,
				DataType:      "colname",
			}
			index += 2
		} else if srcLen > index && source[index].Type == tokenLiteral {
			col = &SyntaxTree{
				childNodes:    []SyntaxTree{},
				StartPosition: index,
				EndPosition:   index,
				Source:        source,
				DataType:      "colname",
			}
		} else {
			return nil, fmt.Errorf(
				"Syntax error, no matching column clause found at %d (%s)",
				source[index].Pos,
				source[index].String())
		}

		//check for order token
		if srcLen > (index+1) &&
			(source[index+1].Type == tokenAsc || source[index+1].Type == tokenDesc) {
			order = &SyntaxTree{
				childNodes:    []SyntaxTree{},
				StartPosition: index + 1,
				EndPosition:   index + 1,
				Source:        source,
				DataType:      "order",
			}
			index++
		}

		colAst := []SyntaxTree{*col}
		if order != nil {
			colAst = append(colAst, *order)
		}

		nodes = append(nodes, SyntaxTree{
			childNodes:    colAst,
			StartPosition: tmpStartIndex,
			EndPosition:   index,
			Source:        source,
			DataType:      "column",
		})

		//check for colon token
		if srcLen > (index + 1) {
			if source[index+1].Type == tokenColon {
				index += 2
				continue
			}

			//no comma, no reason to continue
			break
		}

		//escape loop if no colon token found
		break
	}

	return &SyntaxTree{
		childNodes:    nodes,
		StartPosition: startIndex,
		EndPosition:   index,
		Source:        source,
		DataType:      "groupby",
	}, nil
}

func parseOrderBy(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//pattern
	//expr = ORDER BY <cols>
	//cols = <expression>, <cols>
	//cols = <expression> <order>, <cols>
	//cols = <expression>
	//cols = <expression> <order>
	//order = ASC
	//order = DESC
	if source[startIndex].Type != tokenOrderBy {
		return nil, fmt.Errorf("Expect token ORDER BY at position %d", source[startIndex].Pos)
	}

	srcLen := len(source)
	nodes := []SyntaxTree{}
	index := startIndex + 1

	var col, order *SyntaxTree

	for index < srcLen {
		tmpStartIndex := index
		col = nil
		order = nil

		if srcLen > index+2 &&
			source[index].Type == tokenLiteral &&
			source[index+1].Type == tokenDot &&
			source[index+2].Type == tokenLiteral {
			col = &SyntaxTree{
				childNodes:    []SyntaxTree{},
				StartPosition: index,
				EndPosition:   index + 2,
				Source:        source,
				DataType:      "colname",
			}
			index += 2
		} else if srcLen > index && source[index].Type == tokenLiteral {
			col = &SyntaxTree{
				childNodes:    []SyntaxTree{},
				StartPosition: index,
				EndPosition:   index,
				Source:        source,
				DataType:      "colname",
			}
		} else {
			return nil, fmt.Errorf(
				"Syntax error, no matching column clause found at %d (%s)",
				source[index].Pos,
				source[index].String())
		}

		//check for order token
		if srcLen > (index+1) &&
			(source[index+1].Type == tokenAsc || source[index+1].Type == tokenDesc) {
			order = &SyntaxTree{
				childNodes:    []SyntaxTree{},
				StartPosition: index + 1,
				EndPosition:   index + 1,
				Source:        source,
				DataType:      "order",
			}
			index++
		}

		colAst := []SyntaxTree{*col}
		if order != nil {
			colAst = append(colAst, *order)
		}

		nodes = append(nodes, SyntaxTree{
			childNodes:    colAst,
			StartPosition: tmpStartIndex,
			EndPosition:   index,
			Source:        source,
			DataType:      "column",
		})

		//check for colon token
		if srcLen > (index + 1) {
			if source[index+1].Type == tokenColon {
				index += 2
				continue
			}

			//no comma, no reason to continue
			break
		}

		//escape loop if no colon token found
		break
	}

	return &SyntaxTree{
		childNodes:    nodes,
		StartPosition: startIndex,
		EndPosition:   index,
		Source:        source,
		DataType:      "orderby",
	}, nil
}

func parseHaving(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//pattern
	//expr = HAVING <condition>

	if source[startIndex].Type != tokenHaving {
		return nil, fmt.Errorf("expect token HAVING at position %d", source[startIndex].Pos)
	}

	cond, condErr := parseCondition(source, startIndex+1)
	if condErr != nil {
		return nil, condErr
	}

	return &SyntaxTree{
		childNodes:    []SyntaxTree{*cond},
		StartPosition: startIndex,
		EndPosition:   cond.EndPosition,
		Source:        source,
		DataType:      "having",
	}, nil
}

func parseLimit(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//pattern
	//expr = LIMIT <integer>
	//expr = LIMIT <integer>,<integer>
	//expr = LIMIT <integer> OFFSET <integer>

	if source[startIndex].Type != tokenLimit {
		return nil, fmt.Errorf("Expect token LIMIT at position %d", source[startIndex].Pos)
	}

	sourceLen := len(source)

	//LIMIT <integer> OFFSET <integer>
	if sourceLen > (startIndex+3) &&
		source[startIndex+1].Type == tokenNumber &&
		source[startIndex+2].Type == tokenOffset &&
		source[startIndex+3].Type == tokenNumber {
		return &SyntaxTree{
			childNodes:    []SyntaxTree{},
			StartPosition: startIndex,
			EndPosition:   startIndex + 3,
			Source:        source,
			DataType:      "limit",
		}, nil
		//LIMIT <integer>,<integer>
	} else if sourceLen > (startIndex+3) &&
		source[startIndex+1].Type == tokenNumber &&
		source[startIndex+2].Type == tokenColon &&
		source[startIndex+3].Type == tokenNumber {
		return &SyntaxTree{
			childNodes:    []SyntaxTree{},
			StartPosition: startIndex,
			EndPosition:   startIndex + 3,
			Source:        source,
			DataType:      "limit",
		}, nil
		//LIMIT <integer>
	} else if sourceLen > (startIndex+1) &&
		source[startIndex+1].Type == tokenNumber {
		return &SyntaxTree{
			childNodes:    []SyntaxTree{},
			StartPosition: startIndex,
			EndPosition:   startIndex + 1,
			Source:        source,
			DataType:      "limit",
		}, nil
	}

	return nil, fmt.Errorf(
		"Synxtax error for token LIMIT at position %d", source[startIndex].Pos)
}

func parseQuerySelect(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//pattern:
	//expr = <selectExpr> <fromExpr> <opt>
	//opt = (<joinExpr>,<whereExpr>,<orderbyExpr>,<havingExpr>,<groupbyExpr>,<limitExpr>)
	index := startIndex
	nodes := []SyntaxTree{}

	selectSynxtax, selectSyntaxErr := parseSelect(source, index)
	if selectSyntaxErr != nil {
		return nil, selectSyntaxErr
	}
	nodes = append(nodes, *selectSynxtax)
	index = selectSynxtax.EndPosition

	//try parse FROM statement
	if len(source) <= (index + 1) {
		return nil, fmt.Errorf(
			"incomplete SELECT statement found at position %d",
			index)
	}
	fromSyntax, fromSyntaxErr := parseFrom(source, index+1)
	if fromSyntaxErr != nil {
		return nil, fromSyntaxErr
	}
	nodes = append(nodes, *fromSyntax)
	index = fromSyntax.EndPosition

	//****** optional statements
	//parse JOIN statement
	if len(source) > index+1 && (source[index+1].Type == tokenJoin ||
		source[index+1].Type == tokenLeftJoin ||
		source[index+1].Type == tokenInnerJoin ||
		source[index+1].Type == tokenOuterJoin ||
		source[index+1].Type == tokenRightJoin) {
		joinAST, joinErr := parseJoin(source, index+1)
		if joinErr == nil {
			nodes = append(nodes, *joinAST)
			index = joinAST.EndPosition
		}
	}

	//parse WHERE statement
	if len(source) > (index+1) && source[index+1].Type == tokenWhere {
		whereSyntax, whereSyntaxErr := parseWhere(source, index+1)
		if whereSyntaxErr == nil {
			nodes = append(nodes, *whereSyntax)
			index = whereSyntax.EndPosition
		}
	}

	//group by
	if len(source) > (index+1) && source[index+1].Type == tokenGroupBy {
		groupbyAST, groupbyASTErr := parseGroupBy(source, index+1)
		if groupbyASTErr == nil {
			nodes = append(nodes, *groupbyAST)
			index = groupbyAST.EndPosition
		}
	}

	//having
	if len(source) > (index+1) && source[index+1].Type == tokenHaving {
		havingAST, havingErr := parseHaving(source, index+1)
		if havingErr == nil {
			nodes = append(nodes, *havingAST)
			index = havingAST.EndPosition
		}
	}

	//order by
	if len(source) > (index+1) && source[index+1].Type == tokenOrderBy {
		orderbyAST, orderbyASTErr := parseOrderBy(source, index+1)
		if orderbyASTErr == nil {
			nodes = append(nodes, *orderbyAST)
			index = orderbyAST.EndPosition
		}
	}

	//limit
	if len(source) > (index+1) && source[index+1].Type == tokenLimit {
		limitAST, limitASTErr := parseLimit(source, index+1)
		if limitASTErr == nil {
			nodes = append(nodes, *limitAST)
			index = limitAST.EndPosition
		}
	}

	return &SyntaxTree{
		childNodes:    nodes,
		StartPosition: startIndex,
		EndPosition:   index,
		Source:        source,
		DataType:      "queryselect",
	}, nil
}

func parseQuery(source []tokenItem, startIndex int) (*SyntaxTree, error) {
	//pattern:
	//expr = <selectQuery> UNION expr
	//expr = <selectQuery>
	index := startIndex

	tmp, tmpErr := parseQuerySelect(source, index)
	if tmpErr != nil {
		return nil, tmpErr
	}
	index = tmp.EndPosition

	for i := (index + 1); i < len(source); i++ {
		if source[i].Type == tokenUnion {

			if len(source) > (i + 1) {
				tmp1, tmp1Err := parseQuerySelect(source, i+1)
				if tmp1Err != nil {
					return nil, tmp1Err
				}

				i = tmp1.EndPosition
				index = i
			}

			return nil, fmt.Errorf(
				"incomplete syntax found at position %d", source[i].Pos)
		}

		break
	}

	return &SyntaxTree{
		childNodes:    []SyntaxTree{},
		StartPosition: startIndex,
		EndPosition:   index,
		Source:        source,
		DataType:      "query",
	}, nil
}
