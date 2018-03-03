package rdbmstool

import (
	"fmt"
	"strings"
)

//ConditionOperator logic operator for WHERE clause: =,<>,>,>=,<,<=, etc.
type ConditionOperator uint8

//Condition operator constants
const (
	EQUAL ConditionOperator = iota
	NOT_EQUAL
	GREATER_THAN
	LESS_THAN
	GREATER_THAN_EQUAL
	LESS_THAN_EQUAL
	BETWEEN
	LIKE
	IN
	UNARY_OPERATOR
)

func (operator ConditionOperator) String() string {
	switch operator {
	case EQUAL:
		return "="
	case NOT_EQUAL:
		return "<>"
	case GREATER_THAN:
		return ">"
	case LESS_THAN:
		return "<"
	case GREATER_THAN_EQUAL:
		return "<="
	case LESS_THAN_EQUAL:
		return "<="
	case BETWEEN:
		return "BETWEEN"
	case LIKE:
		return "LIKE"
	case IN:
		return "IN"
	default:
		return ""
	}
}

//ConditionDefinition simplest definition for a where statement
type ConditionDefinition struct {
	LeftExpression  string
	RightExpression string            //set to empty string if condition is unary operation
	Operator        ConditionOperator //set to UNARY_OPERATOR if condition is unary operation

	LeftComplex  *ConditionDefinition //use leftComplex for 'left hand side' expression if the expression is nested
	RightComplex *ConditionDefinition //use rightComplex for 'right hand side' expression if the expression is nested
}

func (cond *ConditionDefinition) leftString() (string, error) {
	if cond.LeftComplex == nil {
		return cond.LeftExpression, nil
	}

	leftStr, err := cond.LeftComplex.SQL()
	if err != nil {
		return "", err
	}

	return "(" + leftStr + ")", nil
}

func (cond *ConditionDefinition) rightString() (string, error) {
	if cond.Operator == UNARY_OPERATOR {
		return "", nil
	} else if cond.RightComplex == nil {
		return cond.RightExpression, nil
	} else {
		rightStr, err := cond.RightComplex.SQL()

		if err != nil {
			return "", err
		}

		return "(" + rightStr + ")", nil
	}
}

//SQL generate SQL string for condition statement
func (cond *ConditionDefinition) SQL() (string, error) {
	if cond.Operator == UNARY_OPERATOR {
		return cond.leftString()
	}

	if strings.Compare(cond.Operator.String(), "") == 0 {
		return "", fmt.Errorf("unsupported operator detected %T", cond.Operator)
	}

	leftStr, leftErr := cond.leftString()
	if leftErr != nil {
		return "", leftErr
	}

	rightStr, rightErr := cond.rightString()
	if rightErr != nil {
		return "", rightErr
	}

	return leftStr + " " + cond.Operator.String() + " " + rightStr, nil
}

//NewConditionDefinition create new condition definition instance
func NewConditionDefinition(lhs string, operator ConditionOperator,
	rhs string) *ConditionDefinition {
	return &ConditionDefinition{
		LeftExpression:  lhs,
		RightExpression: rhs,
		Operator:        operator,
		LeftComplex:     nil,
		RightComplex:    nil,
	}
}

//SetAsUnaryOperation set condition as unary operation
func (cond *ConditionDefinition) SetAsUnaryOperation(expression string) {
	cond.LeftComplex = nil
	cond.RightComplex = nil
	cond.Operator = UNARY_OPERATOR
	cond.RightExpression = ""
	cond.LeftExpression = expression
}
