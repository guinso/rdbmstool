package query

//OrderByDefinition SQL Order By statement definition
type OrderByDefinition struct {
	Expression  string
	IsAscending bool
}

//SQL generate SQL string for Order By statement
func (orderBy *OrderByDefinition) SQL() (string, error) {
	if !orderBy.IsAscending {
		return orderBy.Expression + " DESC", nil
	}

	return orderBy.Expression, nil
}
