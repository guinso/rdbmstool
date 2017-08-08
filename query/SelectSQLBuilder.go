package query

//SelectSQLBuilder SQl Select statement builder interface
type SelectSQLBuilder struct {
	selectDefinition *SelectDefinition
}

//NewSelectSQLBuilder create new Select SQL string builder
func NewSelectSQLBuilder() *SelectSQLBuilder {
	return &SelectSQLBuilder{
		selectDefinition: &SelectDefinition{
			Select:  []ColumnDefinition{},
			From:    nil,
			Join:    []JoinDefinition{},
			Where:   nil,
			GroupBy: []GroupByDefinition{},
			Having:  nil,
			OrderBy: []OrderByDefinition{},
			Limit:   nil,
			Union:   []SelectDefinition{}}}
}

//Select add select column
func (builder *SelectSQLBuilder) Select(expression string, alias string) *SelectSQLBuilder {
	builder.selectDefinition.Select = append(builder.selectDefinition.Select,
		ColumnDefinition{
			Expression: expression,
			Alias:      alias})
	return builder
}

//From set from statement
func (builder *SelectSQLBuilder) From(expression string, alias string) *SelectSQLBuilder {
	builder.selectDefinition.From = NewFromDefinition(expression, alias)
	return builder
}

//Join  add join statement
func (builder *SelectSQLBuilder) Join(join *JoinDefinition) *SelectSQLBuilder {
	builder.selectDefinition.Join = append(builder.selectDefinition.Join, *join)
	return builder
}

//JoinSimple add simple Join statement
func (builder *SelectSQLBuilder) JoinSimple(source string, alias string, category JoinType,
	leftCond string, rightCond string, condOpr ConditionOperator) *SelectSQLBuilder {

	builder.selectDefinition.Join = append(
		builder.selectDefinition.Join,
		*NewJoinDefinition(source, alias, category,
			NewConditionGroupDefinition(condOpr, leftCond, rightCond)))

	return builder
}

//Where set Where statement
func (builder *SelectSQLBuilder) Where(where *ConditionGroupDefinition) *SelectSQLBuilder {
	builder.selectDefinition.Where = where
	return builder
}

//GroupBy add group by statment
func (builder *SelectSQLBuilder) GroupBy(groupBy *GroupByDefinition) *SelectSQLBuilder {
	builder.selectDefinition.GroupBy = append(builder.selectDefinition.GroupBy, *groupBy)
	return builder
}

//Having set having statement
func (builder *SelectSQLBuilder) Having(having *ConditionGroupDefinition) *SelectSQLBuilder {
	builder.selectDefinition.Having = having
	return builder
}

//OrderBy add order by statement
func (builder *SelectSQLBuilder) OrderBy(orderBy *OrderByDefinition) *SelectSQLBuilder {
	builder.selectDefinition.OrderBy = append(builder.selectDefinition.OrderBy, *orderBy)
	return builder
}

//Limit set limit statement
func (builder *SelectSQLBuilder) Limit(limit *LimitDefinition) *SelectSQLBuilder {
	builder.selectDefinition.Limit = limit
	return builder
}

//Union add union statement
func (builder *SelectSQLBuilder) Union(union *SelectDefinition) *SelectSQLBuilder {
	builder.selectDefinition.Union = append(builder.selectDefinition.Union, *union)
	return builder
}

//SQL generate SQL string
func (builder *SelectSQLBuilder) SQL() (string, error) {
	return builder.selectDefinition.SQL()
}
