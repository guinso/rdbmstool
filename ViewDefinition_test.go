package rdbmstool

import (
	"strings"
	"testing"
)

func TestViewDefinition_SQL(t *testing.T) {
	builder := NewQueryBuilder().
		Select("a.name", "").
		Select("a.years_old", "age").
		From("student", "a").
		Join("school", "b", InnerJoin, "a.school = b.name").
		JoinAdd("family", "c", OuterJoin, "a.surname = c.surname").
		WhereAddAnd("a.l = 4").
		WhereAddOr("a.age > 4").
		WhereAddComplex(And, NewCondition("b.name = 'john'").
			AddAnd("b.k <> 8")).
		OrderBy("a.name", true).
		OrderByAdd("age", true).
		GroupBy("b.name", true).
		Having("a.name = 'john'").
		Limit(20, 5)

	viewDef := ViewDefinition{
		Name:  "student",
		Query: builder}

	sql, err := viewDef.SQL()
	if err != nil {
		t.Error(err.Error())
	}

	expectedSQL := `CREATE VIEW student AS 
SELECT a.name, a.years_old AS age
FROM student AS a
INNER JOIN school AS b ON a.school = b.name
OUTER JOIN family AS c ON a.surname = c.surname
WHERE a.l = 4 OR a.age > 4 AND (b.name = 'john' AND b.k <> 8)
GROUP BY b.name
HAVING a.name = 'john'
ORDER BY a.name, age
LIMIT 20 OFFSET 5`

	if strings.Compare(expectedSQL, sql) != 0 {
		t.Errorf("Generated SQL not match with expected SQL\n\nExpected:\n%s\n\nActual:\n%s", expectedSQL, sql)
	}
}
