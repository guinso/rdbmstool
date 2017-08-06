package query

import (
	"errors"
	"fmt"
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

//SQL generate SQL string for condition statement
func (cond *ConditionDefinition) SQL() (string, error) {
	switch cond.Operator {
	case EQUAL:
		return cond.LeftExpression + " = " + cond.RightExpression, nil
	case NOT_EQUAL:
		return cond.LeftExpression + " <> " + cond.RightExpression, nil
	case GREATER_THAN:
		return cond.LeftExpression + " > " + cond.RightExpression, nil
	case LESS_THAN:
		return cond.LeftExpression + " < " + cond.RightExpression, nil
	case GREATER_THAN_EQUAL:
		return cond.LeftExpression + " >= " + cond.RightExpression, nil
	case LESS_THAN_EQUAL:
		return cond.LeftExpression + " <= " + cond.RightExpression, nil
	case BETWEEN:
		return cond.LeftExpression + " BETWEEN " + cond.RightExpression, nil
	case IN:
		return cond.LeftExpression + " IN " + cond.RightExpression, nil
	default:
		return "", fmt.Errorf("Unsupported condition operator detected: %d, left expr '%s', right expr '%s'",
			cond.Operator, cond.LeftExpression, cond.RightExpression)
	}
}

//ConditionGroupDefinition condition definition with parenthesis support
type ConditionGroupDefinition struct {
	condition *ConditionDefinition
	//OR
	subGroup *ConditionGroupDefinition

	items     []ConditionGroupDefinition
	operators []ConditionGroupOperator
}

//Add concate a simple condition into condition group
//EXAMPLE: (a && b)
func (condGroup *ConditionGroupDefinition) Add(operator ConditionGroupOperator, conditionItem *ConditionDefinition) {
	condGroup.items = append(condGroup.items, ConditionGroupDefinition{condition: conditionItem})
}

//AddGroup concate a sub-condition group into condition group
//EXAMPLE: (a && (new-sub-condition))
func (condGroup *ConditionGroupDefinition) AddGroup(operator ConditionGroupOperator, groupItem *ConditionGroupDefinition) {
	condGroup.items = append(condGroup.items, *groupItem)
}

//SQL generate condition group SQL statement
func (condGroup *ConditionGroupDefinition) SQL() (string, error) {
	result := ""
	if condGroup.condition != nil {
		sql, err := condGroup.condition.SQL()
		if err != nil {
			return "", err
		}
		result = sql

	} else if condGroup.subGroup != nil {
		sql, err := condGroup.subGroup.SQL()
		if err != nil {
			return "", err
		}
		result = sql
	} else {
		return "", errors.New("Both condition field and subGroup field cannot be NULL")
	}

	if len(condGroup.items) > 0 || len(condGroup.operators) > 0 {
		if len(condGroup.items) == len(condGroup.operators) {
			count := len(condGroup.items)

			for i := 0; i < count; i++ {
				sql, err := condGroup.items[i].SQL()
				if err != nil {
					return "", err
				}

				if condGroup.operators[i] == AND {
					result = result + " AND " + sql
				} else if condGroup.operators[i] == OR {
					result = result + " OR " + sql
				} else {
					return "", fmt.Errorf("Unsupported condition group operator %d found at index %d",
						condGroup.operators[i], i)
				}
			}

			return "(" + result + ")", nil
		}

		return "", fmt.Errorf("items length# (%d) is not same as operators length (%d)",
			len(condGroup.items), len(condGroup.operators))

	} else if len(condGroup.items) == len(condGroup.operators) && len(condGroup.items) == 0 {
		return result, nil
	}

	return "", fmt.Errorf("items length* (%d) is not same as operators length (%d)",
		len(condGroup.items), len(condGroup.operators))
}

//NewConditionGroupDefinition create new instance of NewConditionGroupDefinition
func NewConditionGroupDefinition(operator ConditionOperator,
	leftExpression string, rightExpression string) *ConditionGroupDefinition {

	return &ConditionGroupDefinition{
		condition: NewConditionDefinition(operator, leftExpression, rightExpression),
		subGroup:  nil,
		items:     []ConditionGroupDefinition{}}
}

//NewConditionGroupDefinitionSub create new instance of NewConditionGroupDefinition with sub group
func NewConditionGroupDefinitionSub(groupCond *ConditionGroupDefinition) *ConditionGroupDefinition {
	return &ConditionGroupDefinition{
		condition: nil,
		subGroup:  groupCond,
		items:     []ConditionGroupDefinition{},
		operators: []ConditionGroupOperator{}}
}

//NewConditionDefinition create new instance of ConditionDefinition
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
