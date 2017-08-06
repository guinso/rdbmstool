package query

//JoinCondition category for JOIN clause: join, inner join, outer join, cross join, etc.
type JoinCondition uint8

const (
	INNER_JOIN JoinCondition = iota + 1
	OUTER_JOIN JoinCondition = iota + 1
	JOIN       JoinCondition = iota + 1
)

//QuerySQLBuilder SQl Select statement builder interface
type QuerySQLBuilder interface {
	Select(*SelectDefinition) *QuerySQLBuilder
	From(*FromDefinition) *QuerySQLBuilder
	Join(*JoinDefinition) *QuerySQLBuilder
	Where(*ConditionGroupDefinition) *QuerySQLBuilder
	GroupBy(string) *QuerySQLBuilder
	Having(string) *QuerySQLBuilder
	OrderBy(string) *QuerySQLBuilder
	Limit(int, int) *QuerySQLBuilder
	Union(*QuerySQLBuilder) *QuerySQLBuilder

	SQL() (string, error)
	GenerateDefinition() *QueryDefinition
}

type QueryDefinition struct {
	Select []SelectDefinition
	From   *FromDefinition
	Join   []JoinDefinition
	Where  *ConditionGroupDefinition
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
	Where     *ConditionGroupDefinition
}
