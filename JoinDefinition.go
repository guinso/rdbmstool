package rdbmstool

import (
	"errors"
	"fmt"
)

//JoinType type for JOIN clause: join, inner join, outer join, cross join, etc.
type JoinType uint8

const (
	//Join SQL JOIN type
	Join JoinType = iota + 1
	//InnerJoin SQL INNER JOIN type
	InnerJoin JoinType = iota + 1
	//OuterJoin SQL OUTER JOIN type
	OuterJoin JoinType = iota + 1
	//LeftJoin SQL LEFT JOIN type
	LeftJoin JoinType = iota + 1
	//RightJoin SQL RIGHT JOIN type
	RightJoin JoinType = iota + 1
)

//JoinDefinition SQL Join definition
type JoinDefinition struct {
	source string
	//OR
	subQuery *SelectDefinition

	Alias string
	Type  JoinType
	Where *ConditionDefinition
}

//SQL generate SQL string for Join link definition
func (join *JoinDefinition) SQL() (string, error) {
	result := ""

	switch join.Type {
	case Join:
		result = result + "JOIN"
		break
	case InnerJoin:
		result = result + "INNER JOIN"
		break
	case OuterJoin:
		result = result + "OUTER JOIN"
		break
	case LeftJoin:
		result = result + "LEFT JOIN"
		break
	case RightJoin:

	default:
		return "", fmt.Errorf("Unsupported JOIN type found: %d", join.Type)
	}

	if len(join.source) > 0 {
		result = result + " " + join.source
	} else if join.subQuery != nil {
		sql, err := join.subQuery.SQL()
		if err != nil {
			return "", err
		}

		result = result + " " + sql
	} else {
		return "", errors.New("JoinDefinition source field and subQuery field cannot be NULL")
	}

	if len(join.Alias) > 0 {
		result = result + " AS " + join.Alias
	}

	if join.Where != nil {
		conditionSQL, sqlErr := join.Where.String()
		if sqlErr != nil {
			return "", fmt.Errorf("Unable to generate JOIN condition SQL string: %s", sqlErr.Error())
		}

		result = result + " ON " + conditionSQL
	}

	return result, nil
}

//NewJoinDefinition create new Join statement definition instance
func NewJoinDefinition(source string, alias string,
	category JoinType, condition string) *JoinDefinition {
	return &JoinDefinition{
		source:   source,
		subQuery: nil,
		Alias:    alias,
		Type:     category,
		Where:    NewCondition(condition)}
}

//NewJoinDefinitionComplex create new Join statement definition instance from sub query as source
func NewJoinDefinitionComplex(source *SelectDefinition, alias string,
	category JoinType, condition *ConditionDefinition) *JoinDefinition {
	return &JoinDefinition{
		source:   "",
		subQuery: source,
		Alias:    alias,
		Type:     category,
		Where:    condition}
}
