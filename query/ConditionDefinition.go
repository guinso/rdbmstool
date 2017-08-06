package query

import (
	"errors"
)

//ConditionOperator logic operator for WHERE clause: =,<>,>,>=,<,<=, etc.
type ConditionOperator uint8

const (
	EQUAL              ConditionOperator = iota + 1
	NOT_EQUAL          ConditionOperator = iota + 1
	GREATER_THAN       ConditionOperator = iota + 1
	LESS_THAN          ConditionOperator = iota + 1
	GREATER_THAN_EQUAL ConditionOperator = iota + 1
	LESS_THAN_EQUAL    ConditionOperator = iota + 1
	BETWEEN            ConditionOperator = iota + 1
	LIKE               ConditionOperator = iota + 1
	IN                 ConditionOperator = iota + 1
)

//ConditionGroupOperator logic operator for WHERE and JOIN clause: AND & OR
type ConditionGroupOperator uint8

const (
	AND ConditionGroupOperator = iota + 1
	OR  ConditionGroupOperator = iota + 1
)

//ConditionDefinition simplest definition for a where statement
type ConditionDefinition struct {
	LeftExpression  string
	RightExpression string
	Operator        ConditionOperator
}

type ConditionGroupDefinition struct {
	condition *ConditionDefinition
	//OR
	subGroup *ConditionGroupDefinition

	items []ConditionGroupDefinition
}

func (condGroup *ConditionGroupDefinition) Add(operator ConditionGroupOperator, conditionItem *ConditionDefinition) {
	condGroup.items = append(condGroup.items, ConditionGroupDefinition{condition: conditionItem})
}

func (condGroup *ConditionGroupDefinition) AddGroup(operator ConditionGroupOperator, groupItem *ConditionGroupDefinition) {
	condGroup.items = append(condGroup.items, *groupItem)
}

//SQL generate condition group SQL statement
func (condGroup *ConditionGroupDefinition) SQL() (string, error) {
	if len(condGroup.items) == 0 {
		//TODO: check wheather condition or subGroup field is empty...
		return "", errors.New("Not implemented yet")
	}

	return "()", errors.New("Not implemented yet")
}

//NewConditionGroupDefinition create new instance for NewConditionGroupDefinition
func NewConditionGroupDefinition(operator ConditionOperator,
	leftExpression string, rightExpression string) *ConditionGroupDefinition {

	return ConditionGroupDefinition{
		condition: NewConditionDefinition(operator, leftExpression, rightExpression),
		subGroup:  nil,
		items:     []ConditionGroupDefinition{}}
}

//NewConditionGroupDefinitionSub create new instance for NewConditionGroupDefinition with sub group
func NewConditionGroupDefinitionSub(groupCond *ConditionGroupDefinition) *ConditionGroupDefinition {
	return ConditionGroupDefinition{
		condition: nil,
		subGroup:  groupCond,
		items:     []ConditionGroupDefinition{}}
}

//NewConditionDefinition create new instance for ConditionDefinition
//SYNTAX:   <leftExpression><operator><rightExpression>
//EXAMPLES: tableA.columnB=tableB.columnC
//          tableA.columnB=:valueX
//
//operator: where operator e.g. =, <>, >, >=, <, <=
//leftExpression: left hand side expression e.g. tableA.columnB, SUM(tableA.columnB)
//rightExpression: right hand side expression e.e. tableA.columnB, SUM(tableA.columnB)
func NewConditionDefinition(operator ConditionOperator,
	leftExpression string, rightExpression string) *ConditionDefinition {

	return &ConditionDefinition{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
		Operator:        operator,
	}
}
