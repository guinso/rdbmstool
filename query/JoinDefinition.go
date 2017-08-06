package query

//JoinCondition category for JOIN clause: join, inner join, outer join, cross join, etc.
type JoinCondition uint8

const (
	INNER_JOIN JoinCondition = iota + 1
	OUTER_JOIN JoinCondition = iota + 1
	JOIN       JoinCondition = iota + 1
)

type JoinDefinition struct {
	Source string
	//OR
	QueryBuilder *SelectSQLBuilder

	Alias     string
	Condition JoinCondition
	Where     *ConditionGroupDefinition
}
