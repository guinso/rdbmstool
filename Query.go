package rdbmstool

type WhereOperator uint8
type ConditionOperator uint8
type JoinCondition uint8

const (
	AND ConditionOperator = iota + 1
	OR  ConditionOperator = iota + 1
)

const (
	INNER_JOIN JoinCondition = iota + 1
	OUTER_JOIN JoinCondition = iota + 1
	JOIN       JoinCondition = iota + 1
)

const (
	EQUAL              WhereOperator = iota + 1
	NOT_EQUAL          WhereOperator = iota + 1
	GREATER_THAN       WhereOperator = iota + 1
	LESS_THAN          WhereOperator = iota + 1
	GREATER_THAN_EQUAL WhereOperator = iota + 1
	LESS_THAN_EQUAL    WhereOperator = iota + 1
	BETWEEN            WhereOperator = iota + 1
	LIKE               WhereOperator = iota + 1
	IN                 WhereOperator = iota + 1
)

//SelectSQLBuilder SQl Select statement builder interface
type QuerySQLBuilder interface {
	Select(*SelectDefinition) *QuerySQLBuilder
	From(*FromDefinition) *QuerySQLBuilder
	Join(*JoinDefinition) *QuerySQLBuilder
	Where(*WhereSQLBuilder) *QuerySQLBuilder
	GroupBy(string) *QuerySQLBuilder
	Having(string) *QuerySQLBuilder
	OrderBy(string) *QuerySQLBuilder
	Limit(int, int) *QuerySQLBuilder
	Union(*QuerySQLBuilder) *QuerySQLBuilder

	SQL() (string, error)
	GenerateDefinition() *QueryDefinition
}

type WhereSQLBuilder interface {
	And(WhereOperator, string, string)
	Or(WhereOperator, string, string)

	AndGroup(*WhereSQLBuilder)
	OrGroup(*WhereSQLBuilder)
}

type QueryDefinition struct {
	Select []SelectDefinition
	From   *FromDefinition
	Join   []JoinDefinition
	Where  *WhereSQLBuilder
	Union  []QuerySQLBuilder
}

type SelectDefinition struct {
	Expression string
	Alias      string
}

type FromDefinition struct {
	Expression string
	//OR
	QueryBuilder *QuerySQLBuilder

	Alias string
}

type JoinDefinition struct {
	Source string
	//OR
	QueryBuilder *QuerySQLBuilder

	Alias     string
	Condition JoinCondition
	Where     *GroupWhereDefinition
}

type GroupWhereDefinition struct {
	Where *WhereDefinition
	//OR
	Group *GroupWhereDefinition

	ConcateWhere []ConcateWhereDefinition
}

type ConcateWhereDefinition struct {
	Where *WhereDefinition
	//OR
	Group *GroupWhereDefinition

	Operator *ConditionOperator
}

type WhereDefinition struct {
	LeftExpression  string
	RightExpression string
	Operator        WhereOperator
}
