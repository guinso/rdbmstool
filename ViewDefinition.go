package rdbmstool

import "fmt"

//ViewDefinition Data Views definition
type ViewDefinition struct {
	Name   string
	Query  *QueryBuilder
	Select *SelectDefinition
}

//SQL generate View SQL string
func (viewDef *ViewDefinition) SQL() (string, error) {

	query, err := viewDef.Query.SQL()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("CREATE VIEW %s AS \n%s", viewDef.Name, query), nil
}
