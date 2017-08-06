package query

//SelectSQLBuilder SQl Select statement builder interface
type SelectSQLBuilder interface {
	Select(*ColumnDefinition) *SelectSQLBuilder
	From(*FromDefinition) *SelectSQLBuilder
	Join(*JoinDefinition) *SelectSQLBuilder
	Where(*ConditionGroupDefinition) *SelectSQLBuilder
	GroupBy(string) *SelectSQLBuilder
	Having(string) *SelectSQLBuilder
	OrderBy(string) *SelectSQLBuilder
	Limit(int, int) *SelectSQLBuilder
	Union(*SelectSQLBuilder) *SelectSQLBuilder

	SQL() (string, error)
	GenerateDefinition() *SelectDefinition
}
