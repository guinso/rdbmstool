package rdbmstool

//QueryBuilder SQl Select statement builder
type QueryBuilder struct {
	selectDefinition *SelectDefinition
}

//NewQueryBuilder create new Select SQL string builder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		selectDefinition: &SelectDefinition{
			Select:  []SelectColumnDefinition{},
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
func (builder *QueryBuilder) Select(expression string, alias string) *QueryBuilder {
	builder.selectDefinition.Select = append(builder.selectDefinition.Select,
		SelectColumnDefinition{
			Expression: expression,
			Alias:      alias})
	return builder
}

//From set from statement
func (builder *QueryBuilder) From(expression string, alias string) *QueryBuilder {
	builder.selectDefinition.From = NewFromDefinition(expression, alias)
	return builder
}

//Join  add join statement
func (builder *QueryBuilder) Join(join *JoinDefinition) *QueryBuilder {
	builder.selectDefinition.Join = append(builder.selectDefinition.Join, *join)
	return builder
}

//JoinSimple add simple Join statement
func (builder *QueryBuilder) JoinSimple(source string, alias string, category JoinType,
	leftCond string, rightCond string, condOpr ConditionOperator) *QueryBuilder {

	builder.selectDefinition.Join = append(
		builder.selectDefinition.Join,
		*NewJoinDefinition(source, alias, category,
			NewConditionGroupDefinition(condOpr, leftCond, rightCond)))

	return builder
}

//Where set Where statement
func (builder *QueryBuilder) Where(where *ConditionGroupDefinition) *QueryBuilder {
	builder.selectDefinition.Where = where
	return builder
}

//WhereAnd append where statement with AND operator
//if where statement is NULL, it will init as first condition and omit AND
func (builder *QueryBuilder) WhereAnd(
	operator ConditionOperator, leftExpression string, rightExpression string) *QueryBuilder {

	if builder.selectDefinition.Where == nil {
		builder.selectDefinition.Where = NewConditionGroupDefinition(operator, leftExpression, rightExpression)
	} else {
		builder.selectDefinition.Where.And(operator, leftExpression, rightExpression)
	}

	return builder
}

//WhereOR append where statement with OR operator
//if where statement is NULL, it will init as first condition and omit OR
func (builder *QueryBuilder) WhereOR(
	operator ConditionOperator, leftExpression string, rightExpression string) *QueryBuilder {

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
func (builder *QueryBuilder) WhereGroup(
	operator ConditionGroupOperator, condition *ConditionGroupDefinition) *QueryBuilder {

	if builder.selectDefinition.Where == nil {
		builder.selectDefinition.Where = condition
	} else {
		builder.selectDefinition.Where.AddGroup(operator, condition)
	}

	return builder
}

//GroupBy add group by statment
func (builder *QueryBuilder) GroupBy(expression string, isAscending bool) *QueryBuilder {
	builder.selectDefinition.GroupBy = append(builder.selectDefinition.GroupBy, GroupByDefinition{
		Expression: expression,
		IsAcending: isAscending})

	return builder
}

//Having set having statement
func (builder *QueryBuilder) Having(having *ConditionGroupDefinition) *QueryBuilder {
	builder.selectDefinition.Having = having
	return builder
}

//OrderBy add order by statement
func (builder *QueryBuilder) OrderBy(expression string, isAscending bool) *QueryBuilder {
	builder.selectDefinition.OrderBy = append(builder.selectDefinition.OrderBy, OrderByDefinition{
		Expression:  expression,
		IsAscending: isAscending})

	return builder
}

//Limit set limit statement
func (builder *QueryBuilder) Limit(rowCount int, offset int) *QueryBuilder {
	builder.selectDefinition.Limit = &LimitDefinition{
		RowCount: rowCount,
		Offset:   offset}

	return builder
}

//Union add union statement
func (builder *QueryBuilder) Union(union *SelectDefinition) *QueryBuilder {
	builder.selectDefinition.Union = append(builder.selectDefinition.Union, *union)
	return builder
}

//SQL generate SQL string
func (builder *QueryBuilder) SQL() (string, error) {
	return builder.selectDefinition.SQL()
}
