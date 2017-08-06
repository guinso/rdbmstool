package mysql

import "github.com/guinso/rdbmstool"

func NewWhereSQLBuilder(operator rdbmstool.WhereOperator,
	leftHandSide string, rightHandSide string) *WhereSQLBuilder {

	group := rdbmstool.GroupWhereDefinition{
		Where: &rdbmstool.WhereDefinition{
			Operator:        operator,
			LeftExpression:  leftHandSide,
			RightExpression: rightHandSide}}

	return &WhereSQLBuilder{groupWhere: []rdbmstool.GroupWhereDefinition{group}}
}

type WhereSQLBuilder struct {
	groupWhere []rdbmstool.GroupWhereDefinition

	//And(WhereOperator, string, string)
	//Or(WhereOperator, string, string)

	AndGroup (*WhereSQLBuilder)
	OrGroup  (*WhereSQLBuilder)
}

func (whereBuilder *WhereSQLBuilder) And(operator rdbmstool.WhereOperator,
	leftExpression string, rightExpression string) *WhereSQLBuilder {

	whereBuilder.groupWhere = append(whereBuilder.groupWhere, 
	})

	return whereBuilder
}

func (whereBuilder *WhereSQLBuilder) Or(operator rdbmstool.WhereOperator,
	leftExpression string, rightExpression string) *WhereSQLBuilder {

	return whereBuilder
}
