package query

import (
	"errors"
	"fmt"
	"strings"
)

//JoinType type for JOIN clause: join, inner join, outer join, cross join, etc.
type JoinType uint8

const (
	JOIN       JoinType = iota + 1
	INNER_JOIN JoinType = iota + 1
	OUTER_JOIN JoinType = iota + 1
	LEFT_JOIN  JoinType = iota + 1
	RIGHT_JOIN JoinType = iota + 1
)

type JoinDefinition struct {
	source string
	//OR
	subQuery *SelectDefinition

	Alias string
	Type  JoinType
	Where *ConditionGroupDefinition
}

//SQL generate SQL string for Join link definition
func (join *JoinDefinition) SQL() (string, error) {
	result := ""

	switch join.Type {
	case JOIN:
		result = result + "JOIN "
		break
	case INNER_JOIN:
		result = result + "INNER JOIN "
		break
	case OUTER_JOIN:
		result = result + "OUTER JOIN "
		break
	default:
		return "", fmt.Errorf("Unsupported JOIN type found: %d", join.Type)
	}

	if len(join.source) > 0 {
		result = join.source
	} else if join.subQuery != nil {
		sql, err := join.subQuery.SQL()
		if err != nil {
			return "", err
		}

		result = sql
	} else {
		return "", errors.New("JoinDefinition source field and subQuery field cannot be NULL")
	}

	if strings.Compare(join.Alias, "") == 0 {
		result = result + " AS " + join.Alias
	}

	conditionSQL, sqlErr := join.Where.SQL()
	if sqlErr != nil {
		return "", fmt.Errorf("Unable to generate JOIN condition SQL string: %s", sqlErr.Error())
	}

	return result + " ON " + conditionSQL, nil

}

//NewJoinDefinition create new Join statement definition instance
func NewJoinDefinition(source string, alias string,
	category JoinType, condition *ConditionGroupDefinition) *JoinDefinition {
	return &JoinDefinition{
		source:   source,
		subQuery: nil,
		Alias:    alias,
		Type:     category,
		Where:    condition}
}

//NewJoinDefinitionSubQuery create new Join statement definition instance from sub query as source
func NewJoinDefinitionSubQuery(source *SelectDefinition, alias string,
	category JoinType, condition *ConditionGroupDefinition) *JoinDefinition {
	return &JoinDefinition{
		source:   "",
		subQuery: source,
		Alias:    alias,
		Type:     category,
		Where:    condition}
}
