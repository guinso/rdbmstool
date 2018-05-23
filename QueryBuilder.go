package rdbmstool

import (
	"strings"
)

//QueryBuilder SQl Select statement builder
type QueryBuilder struct {
	selectDefinition *SelectDefinition
}

//NewQueryBuilder create new Select SQL string builder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		selectDefinition: &SelectDefinition{
			Select:  nil,
			From:    nil,
			Join:    nil,
			Where:   nil,
			GroupBy: nil,
			Having:  nil,
			OrderBy: nil,
			Limit:   nil,
			Union:   nil,
		}}
}

//Select add select column
func (builder *QueryBuilder) Select(expression string, alias string) *QueryBuilder {
	builder.selectDefinition.Select = append(builder.selectDefinition.Select,
		SelectColumnDefinition{
			Expression: expression,
			Alias:      alias})
	return builder
}

//SelectClear clear all select columns
func (builder *QueryBuilder) SelectClear() *QueryBuilder {
	builder.selectDefinition.Select = nil
	return builder
}

//From set from statement
func (builder *QueryBuilder) From(expression string, alias string) *QueryBuilder {
	builder.selectDefinition.From = NewFromDefinition(expression, alias)
	return builder
}

//JoinComplex  set join statement
func (builder *QueryBuilder) JoinComplex(join *JoinDefinition) *QueryBuilder {
	builder.selectDefinition.Join = []JoinDefinition{*join}
	return builder
}

//JoinComplexAdd append join statement
func (builder *QueryBuilder) JoinComplexAdd(join *JoinDefinition) *QueryBuilder {
	builder.selectDefinition.Join = append(builder.selectDefinition.Join, *join)
	return builder
}

//Join set simple Join statement
func (builder *QueryBuilder) Join(source string, alias string, joinType JoinType,
	condition string) *QueryBuilder {

	builder.selectDefinition.Join = []JoinDefinition{
		*NewJoinDefinition(
			source,
			alias,
			joinType,
			condition)}

	return builder
}

//JoinAdd append simple Join statement
func (builder *QueryBuilder) JoinAdd(source string, alias string, joinType JoinType,
	condition string) *QueryBuilder {

	builder.selectDefinition.Join = append(
		builder.selectDefinition.Join,
		*NewJoinDefinition(
			source,
			alias,
			joinType,
			condition))

	return builder
}

//Where set Where condition with simple expression string
func (builder *QueryBuilder) Where(condition string) *QueryBuilder {
	if strings.Compare(condition, "") == 0 {
		builder.selectDefinition.Where = nil
	} else if builder.selectDefinition.Where == nil {
		builder.selectDefinition.Where = NewCondition(condition)
	} else {
		builder.selectDefinition.Where.SetCondition(condition)
	}

	return builder
}

//WhereAddAnd append AND Where condition with simple expression string
func (builder *QueryBuilder) WhereAddAnd(condition string) *QueryBuilder {
	if builder.selectDefinition.Where == nil {
		builder.selectDefinition.Where = NewCondition(condition)
	} else {
		builder.selectDefinition.Where.AddAnd(condition)
	}

	return builder
}

//WhereAddOr append OR Where condition with simple expression string
func (builder *QueryBuilder) WhereAddOr(condition string) *QueryBuilder {
	if builder.selectDefinition.Where == nil {
		builder.selectDefinition.Where = NewCondition(condition)
	} else {
		builder.selectDefinition.Where.AddOr(condition)
	}

	return builder
}

//WhereComplex set where condition with ConditionDefinition
func (builder *QueryBuilder) WhereComplex(conditionDef *ConditionDefinition) *QueryBuilder {
	builder.selectDefinition.Where = conditionDef

	return builder
}

//WhereAddComplex append where condition with ConditionDefinition
func (builder *QueryBuilder) WhereAddComplex(operator ConditionOperator,
	conditionDef *ConditionDefinition) *QueryBuilder {

	if builder.selectDefinition.Where == nil {
		builder.selectDefinition.Where = conditionDef
	} else {
		builder.selectDefinition.Where.AddComplex(operator, conditionDef)
	}

	return builder
}

//WhereClear clear WHERE statement
func (builder *QueryBuilder) WhereClear() *QueryBuilder {
	builder.selectDefinition.Where = nil
	return builder
}

//GroupBy set group by statment
func (builder *QueryBuilder) GroupBy(expression string, isAscending bool) *QueryBuilder {
	builder.selectDefinition.GroupBy = []GroupByDefinition{GroupByDefinition{
		Expression: expression,
		IsAcending: isAscending}}

	return builder
}

//GroupByAdd append group by statement
func (builder *QueryBuilder) GroupByAdd(expression string, isAscending bool) *QueryBuilder {
	builder.selectDefinition.GroupBy = append(builder.selectDefinition.GroupBy, GroupByDefinition{
		Expression: expression,
		IsAcending: isAscending})

	return builder
}

//GroupByClear clear GROUP BY statement
func (builder *QueryBuilder) GroupByClear() *QueryBuilder {
	builder.selectDefinition.GroupBy = nil
	return builder
}

//Having set Having statement with expression string
func (builder *QueryBuilder) Having(condition string) *QueryBuilder {
	if strings.Compare(condition, "") == 0 {
		builder.selectDefinition.Having = nil
	} else if builder.selectDefinition.Having == nil {
		builder.selectDefinition.Having = NewCondition(condition)
	} else {
		builder.selectDefinition.Having.SetCondition(condition)
	}

	return builder
}

//HavingComplex set having statement with ConditionDefinition
func (builder *QueryBuilder) HavingComplex(having *ConditionDefinition) *QueryBuilder {
	builder.selectDefinition.Having = having
	return builder
}

//OrderBy set order by statement
func (builder *QueryBuilder) OrderBy(expression string, isAscending bool) *QueryBuilder {
	builder.selectDefinition.OrderBy = []OrderByDefinition{OrderByDefinition{
		Expression:  expression,
		IsAscending: isAscending}}

	return builder
}

//OrderByAdd append ORDER BY statement
func (builder *QueryBuilder) OrderByAdd(expression string, isAscending bool) *QueryBuilder {
	builder.selectDefinition.OrderBy = append(builder.selectDefinition.OrderBy, OrderByDefinition{
		Expression:  expression,
		IsAscending: isAscending})

	return builder
}

//OrderByClear clear ORDER BY statement
func (builder *QueryBuilder) OrderByClear() *QueryBuilder {
	builder.selectDefinition.OrderBy = nil
	return builder
}

//Limit set limit statement
func (builder *QueryBuilder) Limit(rowCount int, offset int) *QueryBuilder {
	builder.selectDefinition.Limit = &LimitDefinition{
		RowCount: rowCount,
		Offset:   offset}

	return builder
}

//Union append union statement
func (builder *QueryBuilder) Union(union *SelectDefinition) *QueryBuilder {
	builder.selectDefinition.Union = append(builder.selectDefinition.Union, *union)
	return builder
}

//UnionClear clear UNION statement
func (builder *QueryBuilder) UnionClear() *QueryBuilder {
	builder.selectDefinition.Union = nil
	return builder
}

//SQL generate SQL string
func (builder *QueryBuilder) SQL() (string, error) {
	return builder.selectDefinition.SQL()
}
