package query

import "testing"

func TestSelectSQLBuilder(t *testing.T) {
	builder := NewSelectSQLBuilder()

	builder.Select("a.name", "").Select("a.years_old", "age").
		From("student", "a").
		JoinSimple("school", "b", INNER_JOIN, "a.school", "b.name", EQUAL).
		JoinSimple("family", "c", OUTER_JOIN, "a.surname", "c.surname", EQUAL).
		OrderBy("a.name", true).OrderBy("age", true).
		GroupBy("b.name", true).
		Limit(20, 5)

	sql, sqlErr := builder.SQL()
	if sqlErr != nil {
		t.Error(sqlErr.Error())
	}

	t.Error(sql)
}
