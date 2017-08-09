package computed

import (
	"strings"
	"testing"

	"github.com/guinso/rdbmstool/query"
)

func TestViewDefinition_SQL(t *testing.T) {
	builder := query.NewSelectSQLBuilder().
		Select("a.name", "").Select("a.years_old", "age").
		From("student", "a").
		JoinSimple("school", "b", query.INNER_JOIN, "a.school", "b.name", query.EQUAL).
		JoinSimple("family", "c", query.OUTER_JOIN, "a.surname", "c.surname", query.EQUAL).
		WhereAnd(query.EQUAL, "a.l", "4").
		WhereOR(query.GREATER_THAN, "a.age", "4").
		WhereGroup(query.AND, query.NewConditionGroupDefinition(query.EQUAL, "b.name", "'john'").
			And(query.NOT_EQUAL, "b.k", "8")).
		OrderBy("a.name", true).OrderBy("age", true).
		GroupBy("b.name", true).
		Having(query.NewConditionGroupDefinition(query.EQUAL, "a.name", "'john'")).
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
WHERE (a.l = 4 OR a.age > 4 AND (b.name = 'john' AND b.k <> 8))
GROUP BY b.name
HAVING a.name = 'john'
ORDER BY a.name, age
LIMIT 20 OFFSET 5`

	if strings.Compare(expectedSQL, sql) != 0 {
		t.Errorf("Generate SQL not match with expected SQL\n\nExpected:\n%s\n\nActual:\n%s", expectedSQL, sql)
	}
}
