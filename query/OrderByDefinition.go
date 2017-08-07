package query

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
