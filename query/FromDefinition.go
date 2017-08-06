package query

import (
	"errors"
	"strings"
)

//FromDefinition SQL FROM statement definition
type FromDefinition struct {
	expression string
	//OR
	queryBuilder *SelectDefinition

	alias string
}

//SQL generate SQL string for FROM statement
func (from *FromDefinition) SQL() (string, error) {
	result := "FROM "

	if from.queryBuilder != nil {
		sql, err := from.queryBuilder.SQL()
		if err != nil {
			return "", err
		}
		result = result + "(" + sql + ")"

	} else if len(from.expression) > 0 {
		result = result + from.expression

	} else {
		return "", errors.New("FromDefinition cannot have empty value on " +
			"both expression and queryBuilder field")
	}

	if strings.Compare(from.alias, "") == 0 {
		return result, nil
	}

	return result + " AS " + from.alias, nil
}

//NewFromDefinition create new FROM statement definition
func NewFromDefinition(sourceExpression string, aliasName string) *FromDefinition {
	return &FromDefinition{
		expression:   sourceExpression,
		queryBuilder: nil,
		alias:        aliasName}
}

//NewFromDefinitionSubQuery create new FROM statement definition from sub-query
func NewFromDefinitionSubQuery(subQuery *SelectDefinition, aliasName string) *FromDefinition {
	return &FromDefinition{
		expression:   "",
		queryBuilder: subQuery,
		alias:        aliasName}
}
