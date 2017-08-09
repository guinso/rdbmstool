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

//WhereAnd append where statement with AND operator
//if where statement is NULL, it will init as first condition and omit AND
func (builder *SelectSQLBuilder) WhereAnd(
	operator ConditionOperator, leftExpression string, rightExpression string) *SelectSQLBuilder {

	if builder.selectDefinition.Where == nil {
		builder.selectDefinition.Where = NewConditionGroupDefinition(operator, leftExpression, rightExpression)
	} else {
		builder.selectDefinition.Where.And(operator, leftExpression, rightExpression)
	}

	return builder
}

//WhereOR append where statement with OR operator
//if where statement is NULL, it will init as first condition and omit OR
func (builder *SelectSQLBuilder) WhereOR(
	operator ConditionOperator, leftExpression string, rightExpression string) *SelectSQLBuilder {

	if builder.selectDefinition.Where == nil {
		builder.selectDefinition.Where = NewConditionGroupDefinition(operator, leftExpression, rightExpression)
	} else {
		builder.selectDefinition.Where.Or(operator, leftExpression, rightExpression)
	}

	return builder
}

//WhereGroup append Where statement with group condition
//operator: group condition, example OR, AND
//condition: nested condition, example (a =3 AND (b > 4 OR d <> false))
func (builder *SelectSQLBuilder) WhereGroup(
	operator ConditionGroupOperator, condition *ConditionGroupDefinition) *SelectSQLBuilder {

	if builder.selectDefinition.Where == nil {
		builder.selectDefinition.Where = condition
	} else {
		builder.selectDefinition.Where.AddGroup(operator, condition)
	}

	return builder
}

//GroupBy add group by statment
func (builder *SelectSQLBuilder) GroupBy(expression string, isAscending bool) *SelectSQLBuilder {
	builder.selectDefinition.GroupBy = append(builder.selectDefinition.GroupBy, GroupByDefinition{
		Expression: expression,
		IsAcending: isAscending})

	return builder
}

//Having set having statement
func (builder *SelectSQLBuilder) Having(having *ConditionGroupDefinition) *SelectSQLBuilder {
	builder.selectDefinition.Having = having
	return builder
}

//OrderBy add order by statement
func (builder *SelectSQLBuilder) OrderBy(expression string, isAscending bool) *SelectSQLBuilder {
	builder.selectDefinition.OrderBy = append(builder.selectDefinition.OrderBy, OrderByDefinition{
		Expression:  expression,
		IsAscending: isAscending})

	return builder
}

//Limit set limit statement
func (builder *SelectSQLBuilder) Limit(rowCount int, offset int) *SelectSQLBuilder {
	builder.selectDefinition.Limit = &LimitDefinition{
		RowCount: rowCount,
		Offset:   offset}

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
