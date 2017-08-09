package query

import "testing"
import "strings"

func TestSelectSQLBuilder(t *testing.T) {
	builder := NewSelectSQLBuilder()

	builder.Select("a.name", "").Select("a.years_old", "age").
		From("student", "a").
		JoinSimple("school", "b", INNER_JOIN, "a.school", "b.name", EQUAL).
		JoinSimple("family", "c", OUTER_JOIN, "a.surname", "c.surname", EQUAL).
		WhereAnd(EQUAL, "a.l", "4").
		WhereOR(GREATER_THAN, "a.age", "4").
		WhereGroup(AND,
			NewConditionGroupDefinition(EQUAL, "b.name", "'john'").
				And(NOT_EQUAL, "b.k", "8")).
		OrderBy("a.name", true).OrderBy("age", true).
		GroupBy("b.name", true).
		Having(NewConditionGroupDefinition(EQUAL, "a.name", "'john'")).
		Limit(20, 5)

	sql, sqlErr := builder.SQL()
	if sqlErr != nil {
		t.Error(sqlErr.Error())
	}

	expectedSQL := `SELECT a.name, a.years_old AS age
FROM student AS a
INNER JOIN school AS b ON a.school = b.name
OUTER JOIN family AS c ON a.surname = c.surname
WHERE (a.l = 4 OR a.age > 4 AND (b.name = 'john' AND b.k <> 8))
GROUP BY b.name
HAVING a.name = 'john'
ORDER BY a.name, age
LIMIT 20 OFFSET 5`

	if strings.Compare(expectedSQL, sql) != 0 {
		t.Errorf("Generate SQL not match with expected SQL\n\nExpected:\n%s\n\nActual:\n%s", expectedSQL, sql)
	}
}
