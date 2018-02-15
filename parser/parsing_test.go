package parser

import (
	"fmt"
	"testing"
)

func Test_parseParenthesis(t *testing.T) {
	tokens := tokenize("(b - 222 * (c + 3) - g * (m + c * (b / q)))")
	if _, parenErr := parseParenthesis(tokens, 0); parenErr != nil {
		t.Errorf(parenErr.Error())
	}

	tokens = tokenize("(b - 222 * (c + 3) - g * (m + c * (b / q))) + c")
	if _, parenErr := parseParenthesis(tokens, 0); parenErr != nil {
		t.Errorf(parenErr.Error())
	}

	tokens = tokenize("2 * 4 + (b - 222 * (c + 3) - g * (m + c * (b / q))) + c")
	if _, parenErr := parseParenthesis(tokens, 4); parenErr != nil {
		t.Errorf(parenErr.Error())
	}

	tokens = tokenize("(b - 222 * (c + 3) - g * (m + c * (b / q))")
	if _, parenErr := parseParenthesis(tokens, 0); parenErr == nil {
		t.Errorf("Expect syntax error occur since parenthesis not completely close at end of context")
	}

	tokens = tokenize("(b - 222 * (c + 3) - g * (m + c * (b / q)))")
	if _, tmpErr := parseParenthesis(tokens, 1); tmpErr == nil {
		t.Errorf("expect parse error since start index is not pointing at left parenthesis")
	}

}

func Test_parseField(t *testing.T) {
	token := tokenize("kikilala")
	if _, err := parseField(token, 0); err != nil {
		t.Errorf(err.Error())
	}

	token = tokenize("SELECT a, b")
	if _, err := parseField(token, 0); err == nil {
		t.Errorf("expect expression will not able to detect")
	}

	token = tokenize("SELECT a, b")
	if _, err := parseField(token, 1); err != nil {
		t.Error(err)
	}

	token = tokenize("a.b AS koko")
	if _, err := parseField(token, 0); err != nil {
		t.Errorf(err.Error())
	}

	token = tokenize("(b + c.a) AS frog")
	if _, err := parseField(token, 0); err != nil {
		t.Errorf(err.Error())
	}

	token = tokenize("(b + * c.a) AS frog")
	if _, err := parseField(token, 0); err == nil {
		t.Errorf("expect double operator syntax error will be detected")
	}

	token = tokenize("a.b AS koko, b")
	if _, err := parseField(token, 0); err != nil {
		t.Errorf(err.Error())
	}

	token = tokenize("(a.b AS koko, b")
	if _, err := parseField(token, 0); err == nil {
		t.Errorf("expect imcomplete parenthesis syntax found")
	}
}

func Test_parseExpression(t *testing.T) {
	token := tokenize("5 * g")
	expr, err := parseExpresion(token, 0)
	if err != nil {
		t.Error(err)
	} else if expr.EndPosition != 2 {
		t.Errorf("expect expression detect at position 2 but get %d", expr.EndPosition)
	}

	token = tokenize("-12 + 5 * g")
	if _, err := parseExpresion(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("-12 + (5 * -g)")
	if _, err := parseExpresion(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("-12 + (5 * -+g)")
	if _, err := parseExpresion(token, 0); err == nil {
		t.Errorf("expect syntax error found at position 7")
	}

	token = tokenize("-12 + (5 * - g)")
	if _, err := parseExpresion(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("-12 + (5 * (a.b - g)) / bahamut.dark_flare")
	if _, err := parseExpresion(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("-b + (5 * -g)")
	if _, err := parseExpresion(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("-12 + (5 * -g) / 1R2D2")
	if _, err := parseExpresion(token, 0); err == nil {
		t.Error("expect error since fail to tokenize")
	}

	token = tokenize("-12 + SUM(5 * g)")
	if _, err := parseExpresion(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("-12 + SUM(5 * g)")
	expr, err = parseExpresion(token, 0)
	if err != nil {
		t.Error(err)
		return
	}
	if expr.EndPosition != 8 {
		t.Errorf("expect expression detect at position 8 but ended at %d", expr.EndPosition)
	}

	token = tokenize("-12 + (5 * (a.b - g)) / bahamut.dark_flare")
	expr, err = parseExpresion(token, 0)
	if err != nil {
		t.Error(err)
		return
	}
	if expr.EndPosition != 17 {
		t.Errorf("expect expression detect at position 17 but ended at %d", expr.EndPosition)
	}
}

func Test_parseFunction(t *testing.T) {
	token := tokenize("MIN(b.g)")
	if _, err := parseFunction(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("COUNT(student)")
	if _, err := parseFunction(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("COUNT(3 + 2)")
	if _, err := parseFunction(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("COUNT(3 + a.g)")
	if _, err := parseFunction(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("COUNT(3 + (b.c * 4)))")
	if _, err := parseFunction(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("3 + (b.c * 4))")
	if _, err := parseFunction(token, 0); err == nil {
		t.Errorf("expect syntax error since no function token found")
	}

	token = tokenize("count a.g")
	if _, err := parseFunction(token, 0); err == nil {
		t.Errorf("expect syntax error since incomplete function syntax")
	}

	token = tokenize("count(b.c * 4")
	if _, err := parseFunction(token, 0); err == nil {
		t.Errorf("expect syntax error since incomplete function's parenthesis")
	}

	token = tokenize("SUM(a.b, gg)")
	if _, err := parseFunction(token, 0); err == nil {
		t.Errorf("expect syntax error since expression need to ended side to right parenthesis")
	}

	token = tokenize("3 + MAX")
	if _, err := parseFunction(token, 2); err == nil {
		t.Errorf("expect syntax error since it is incomplete syntax")
	}

	token = tokenize("MAX(a.b + SELECT)")
	if _, err := parseFunction(token, 2); err == nil {
		t.Errorf("expect syntax error since it is incomplete syntax")
	}
}

func Test_parseCondition(t *testing.T) {

	token := tokenize("a.b > 45 AND 5 != 3")
	if _, err := parseCondition(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("a.b > 45 OR 5 != 3")
	if _, err := parseCondition(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("a.b > 45 AND (5 != 3)")
	if _, err := parseCondition(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("a.b > 45 AND ((5 != 3 OR b = c) OR v.bobo <> jojo)")
	if _, err := parseCondition(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("a.b > 45")
	if _, err := parseCondition(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("a.b > 45 AND ")
	if cond, err := parseCondition(token, 0); err == nil {
		t.Errorf("expect syntax error since ended with AND token but parse complete and ended at %d (%s)",
			cond.EndPosition, token[cond.EndPosition].String())
	}

	token = tokenize("a.b > 45 AND (a.b = 5")
	if _, err := parseCondition(token, 0); err == nil {
		t.Errorf("expect syntax error since parenthesis is not completely closed")
	}

	token = tokenize("a.b > 45 AND (a.b = 5 AND)")
	if _, err := parseCondition(token, 0); err == nil {
		t.Errorf("expect syntax error since sub condition expression has syntax error")
	}

	token = tokenize("a.b > 45 AND (a.b = 5 AND b > 4, k.f)")
	if _, err := parseCondition(token, 0); err == nil {
		t.Errorf("expect syntax error since inner condition not ended next to right parenthesis")
	}

	token = tokenize("WHERE a.g")
	if _, err := parseCondition(token, 0); err == nil {
		t.Errorf("expect syntax error since it is not a valid condition expression")
	}
}

func Test_parseJoin(t *testing.T) {
	token := tokenize("JOIN student AS stu ON a.name = stu.name AND a.age = stu.age")
	if _, err := parseJoin(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("JOIN student ON a.name = stu.name AND a.age = stu.age")
	if _, err := parseJoin(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("JOIN student AS stu")
	if _, err := parseJoin(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("LEFT JOIN student AS stu")
	if _, err := parseJoin(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("RIGHT JOIN student AS stu ON a.name = stu.name AND a.age = stu.age")
	if _, err := parseJoin(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("INNER JOIN student AS stu ON a.name = stu.name AND a.age = stu.age")
	if _, err := parseJoin(token, 0); err != nil {
		t.Error(err)
	}
}
func Test_parseSelect(t *testing.T) {
	token := tokenize("SELECT a, a.b, MAX(c), SUM(go), a + b, k.hoho AS valueA, a.b koko, gg")
	ast, err := parseSelect(token, 0)
	if err != nil {
		t.Error(err)
	}
	if len(ast.childNodes) != 8 {
		t.Errorf("expect parsing SELECT return 8 columns but get %d instead", len(ast.childNodes))
	}
	tmpLog := ""
	for i := 0; i < len(ast.childNodes); i++ {
		tmpLog = ""
		for j := ast.childNodes[i].StartPosition; j <= ast.childNodes[i].EndPosition; j++ {
			tmpLog = tmpLog + " " + ast.childNodes[i].Source[j].Value
		}
		fmt.Println(tmpLog)
	}

	token = tokenize("JOIN student AS stu ON a.name = stu.name AND a.age = stu.age")
	if _, err := parseSelect(token, 0); err == nil {
		t.Errorf("expect synxtax error")
	}
}

func Test_parseFrom(t *testing.T) {
	token := tokenize("FROM a")
	if _, err := parseFrom(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("FROM bangbang.student")
	if _, err := parseFrom(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("FROM (bangbang.student)")
	if _, err := parseFrom(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("FROM (bangbang.student 123 asd)")
	if _, err := parseFrom(token, 0); err == nil {
		t.Errorf("expect syntax error inside parenthesis")
	}

	token = tokenize("FROM (bangbang.student")
	if _, err := parseFrom(token, 0); err == nil {
		t.Errorf("expect syntax error since parenthesis is incomplete")
	}

	token = tokenize("FROM 123 + a.b")
	if _, err := parseFrom(token, 0); err == nil {
		t.Errorf("expect syntax error")
	}

	token = tokenize("WHERE bangbang.student > 123")
	if _, err := parseFrom(token, 0); err == nil {
		t.Errorf("expect syntax error since it is not WHERE clause syntax")
	}

	token = tokenize("FROM bangbang.student b")
	ast, err := parseFrom(token, 0)
	if err != nil {
		t.Error(err)
	}
	if ast.EndPosition != 4 {
		t.Errorf("Expect ended at index 4 but %d (%s) instead",
			ast.EndPosition,
			token[ast.EndPosition].String())
	}
	if len(ast.childNodes) != 2 {
		t.Errorf("expect to have source alias but it is not")
	}

	token = tokenize("FROM bangbang.student AS b")
	ast, err = parseFrom(token, 0)
	if err != nil {
		t.Error(err)
	}
	if ast.EndPosition != 5 {
		t.Errorf("Expect ended at index 5 but %d (%s) instead",
			ast.EndPosition,
			token[ast.EndPosition].String())
	}
	if len(ast.childNodes) != 2 {
		t.Errorf("expect to have source alias but it is not")
	}

	// token = tokenize("FROM (SELECT a.b, a.c FROM student a)")
	// if _, err := parseFrom(token, 0); err != nil {
	// 	t.Error(err)
	// }
}

func Test_parseWhere(t *testing.T) {
	token := tokenize("WHERE a.b = c.b")
	if _, err := parseWhere(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("WHERE (a.b = c.b)")
	if _, err := parseWhere(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("WHERE (a.b = c.b AND 23 > a.b)")
	if _, err := parseWhere(token, 0); err != nil {
		t.Error(err)
	}
}

func Test_parseGroupBy(t *testing.T) {
	token := tokenize("GROUP BY a.b, c.b")
	if _, err := parseGroupBy(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("GROUP BY a.b ASC, c.b")
	if _, err := parseGroupBy(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("GROUP BY a.b DESC, c.b")
	if _, err := parseGroupBy(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("GROUP BY a.b ASC, c")
	if _, err := parseGroupBy(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("GROUP BY a.b DESC, c ASC")
	gb, err := parseGroupBy(token, 0)
	if err != nil {
		t.Error(err)
	}
	if len(gb.childNodes) != 2 {
		t.Errorf("expect parse result gives 2 columns but get %d instead",
			len(gb.childNodes))
	} else {
		if len(gb.childNodes[0].childNodes) != 2 {
			t.Errorf("expect column 1 gives 2 nodes (col, order)")
		}

		if len(gb.childNodes[1].childNodes) != 2 {
			t.Errorf("expect column 2 gives 2 nodes (col, order)")
		}
	}
}

func Test_parseOrderBy(t *testing.T) {
	token := tokenize("ORDER BY a.b, c.b")
	if _, err := parseOrderBy(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("ORDER BY a.b ASC, c.b")
	if _, err := parseOrderBy(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("ORDER BY a.b DESC, c.b")
	if _, err := parseOrderBy(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("ORDER BY a.b ASC, c")
	if _, err := parseOrderBy(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("ORDER BY a.b DESC, c ASC")
	gb, err := parseOrderBy(token, 0)
	if err != nil {
		t.Error(err)
	}
	if len(gb.childNodes) != 2 {
		t.Errorf("expect parse result gives 2 columns but get %d instead",
			len(gb.childNodes))
	} else {
		if len(gb.childNodes[0].childNodes) != 2 {
			t.Errorf("expect column 1 gives 2 nodes (col, order)")
		}

		if len(gb.childNodes[1].childNodes) != 2 {
			t.Errorf("expect column 2 gives 2 nodes (col, order)")
		}
	}
}

func Test_parseHaving(t *testing.T) {
	token := tokenize("HAVING COUNT(student) > 3")
	if _, err := parseHaving(token, 0); err != nil {
		t.Error(err)
	}
}

func Test_parseLimit(t *testing.T) {
	token := tokenize("LIMIT 2")
	if _, err := parseLimit(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("LIMIT 2, 3")
	if _, err := parseLimit(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("LIMIT 2 OFFSET 3")
	if _, err := parseLimit(token, 0); err != nil {
		t.Error(err)
	}
}

func Test_parseQuerySelect(t *testing.T) {
	token := tokenize("SELECT name, age FROM student WHERE age > 16 AND age < 55")
	if _, err := parseQuerySelect(token, 0); err != nil {
		t.Error(err)
	}

	token = tokenize("SELECT name, age " +
		"FROM student a " +
		"LEFT JOIN classs b ON a.name = b.name " +
		"WHERE age > 16 AND age < 55 " +
		"GROUP BY a.name, b.leader " +
		"HAVING a.age > 10 " +
		"ORDER BY a.name " +
		"LIMIT 10 OFFSET 2")
	query, err := parseQuerySelect(token, 0)
	if err != nil {
		t.Error(err)
	} else if len(query.childNodes) != 8 {
		t.Errorf("expect has 8 nodes but get %d instead", len(query.childNodes))
		for _, node := range query.childNodes {
			t.Logf("%s %d-%d", node.DataType, node.StartPosition, node.EndPosition)
		}
	}
}
