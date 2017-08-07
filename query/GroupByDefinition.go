package query

//GroupByDefinition SQL Group By statement definition
type GroupByDefinition struct {
	Expression  string
	IsDecending bool
}

//SQL generate SQL string for Group By statement
func (groupBy *GroupByDefinition) SQL() (string, error) {
	if groupBy.IsDecending {
		return groupBy.Expression + " DESC", nil
	}

	return groupBy.Expression, nil
}
