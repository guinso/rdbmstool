package query

import (
	"errors"
	"strings"
)

//ColumnDefinition SQL select column definition
type ColumnDefinition struct {
	Expression string
	Alias      string
}

//SQL generate SQL string for select column statement
func (column *ColumnDefinition) SQL() (string, error) {
	if strings.Compare(column.Alias, "") == 0 {
		return column.Expression, nil
	}

	return column.Expression + " AS " + column.Alias, nil
}

//OrderByDefinition SQL Order By statement definition
type OrderByDefinition struct {
	Expression  string
	IsDecending bool
}

//SQL generate SQL string for Order By statement
func (orderBy *OrderByDefinition) SQL() (string, error) {
	if orderBy.IsDecending {
		return orderBy.Expression + " DESC", nil
	}

	return orderBy.Expression, nil
}

//SelectDefinition SQL query definition
type SelectDefinition struct {
	Select  []ColumnDefinition
	From    *FromDefinition
	Join    []JoinDefinition
	Where   *ConditionGroupDefinition
	GroupBy []string
	Having  []ConditionGroupDefinition
	Union   []SelectDefinition
}

//SQL generate SQL string for SELECT statement
func (query *SelectDefinition) SQL() (string, error) {
	return "", errors.New("Not implemented yet")
}
