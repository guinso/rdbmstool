package query

//SelectSQLBuilder SQl Select statement builder interface
type SelectSQLBuilder interface {
	Select(*ColumnDefinition) *SelectSQLBuilder
	From(*FromDefinition) *SelectSQLBuilder
	Join(*JoinDefinition) *SelectSQLBuilder
	Where(*ConditionGroupDefinition) *SelectSQLBuilder
	GroupBy(*GroupByDefinition) *SelectSQLBuilder
	Having(*ConditionGroupDefinition) *SelectSQLBuilder
	OrderBy(*OrderByDefinition) *SelectSQLBuilder
	Limit(*LimitDefinition) *SelectSQLBuilder
	Union(*SelectSQLBuilder) *SelectSQLBuilder

	SQL() (string, error)
	//GenerateDefinition() *SelectDefinition
}
