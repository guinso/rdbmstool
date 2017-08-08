package query

//GroupByDefinition SQL Group By statement definition
type GroupByDefinition struct {
	Expression string
	IsAcending bool
}

//SQL generate SQL string for Group By statement
func (groupBy *GroupByDefinition) SQL() (string, error) {
	if !groupBy.IsAcending {
		return groupBy.Expression + " DESC", nil
	}

	return groupBy.Expression, nil
}
