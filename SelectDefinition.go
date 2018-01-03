package rdbmstool

import (
	"errors"
	"fmt"
	"strings"
)

//SelectColumnDefinition SQL select column definition
type SelectColumnDefinition struct {
	Expression string
	Alias      string
}

//SQL generate SQL string for select column statement
func (column *SelectColumnDefinition) SQL() (string, error) {
	if strings.Compare(column.Alias, "") == 0 {
		return column.Expression, nil
	}

	return column.Expression + " AS " + column.Alias, nil
}

//SelectDefinition SQL query definition
type SelectDefinition struct {
	Select  []SelectColumnDefinition
	From    *FromDefinition
	Join    []JoinDefinition
	Where   *ConditionGroupDefinition
	GroupBy []GroupByDefinition
	Having  *ConditionGroupDefinition
	OrderBy []OrderByDefinition
	Limit   *LimitDefinition
	Union   []SelectDefinition
}

//SQL generate SQL string for SELECT statement
func (query *SelectDefinition) SQL() (string, error) {
	result := ""

	//Column
	if len(query.Select) == 0 {
		return "", errors.New("Select column must atlest have one item to select")
	}
	for index, col := range query.Select {
		sql, err := col.SQL()
		if err != nil {
			return "", fmt.Errorf("Failed to generate SELECT column (index %d) SQL string: %s", index, err.Error())
		}

		if index == 0 {
			result = "SELECT " + sql
		} else {
			result = result + ", " + sql
		}
	}

	//From
	fromSQL, fromErr := query.From.SQL()
	if fromErr != nil {
		return "", errors.New("Failed to generate FROM SQL string: " + fromErr.Error())
	}
	result = result + "\n" + fromSQL

	//Join
	for index, join := range query.Join {
		joinSQL, joinErr := join.SQL()
		if joinErr != nil {
			return "", fmt.Errorf("Failed to generate JOIN (index %d) SQL string: %s", index, joinErr.Error())
		}
		result = result + "\n" + joinSQL
	}

	//Where
	if query.Where != nil {
		whereSQL, whrErr := query.Where.SQL()
		if whrErr != nil {
			return "", errors.New("Failed to generate WHERE SQL string: " + whrErr.Error())
		}
		result = result + "\nWHERE " + whereSQL
	}

	//Group By
	if len(query.GroupBy) > 0 {
		for index, groupBy := range query.GroupBy {
			groupBySQL, groupErr := groupBy.SQL()
			if groupErr != nil {
				return "", fmt.Errorf("Failed to generate GROUP BY (index %d) SQL string: %s", index, groupErr.Error())
			}

			if index == 0 {
				result = result + "\nGROUP BY " + groupBySQL
			} else {
				result = result + ", " + groupBySQL
			}
		}
	}

	//Having
	if query.Having != nil {
		havingSQL, haveErr := query.Having.SQL()
		if haveErr != nil {
			return "", fmt.Errorf("Failed to generate HAVING SQL string: %s", haveErr.Error())
		}
		result = result + "\nHAVING " + havingSQL
	}

	//Order By
	if len(query.OrderBy) > 0 {
		for index, orderBy := range query.OrderBy {
			orderBySQL, orderErr := orderBy.SQL()
			if orderErr != nil {
				return "", fmt.Errorf("Failed to generate ORDER BY (index %d) SQL string: %s", index, orderErr.Error())
			}

			if index == 0 {
				result = result + "\nORDER BY " + orderBySQL
			} else {
				result = result + ", " + orderBySQL
			}
		}
	}

	//Limit
	if query.Limit != nil {
		limitSQL, limitErr := query.Limit.SQL()
		if limitErr != nil {
			return "", fmt.Errorf("Failed to generate LIMIT SQL string: %s", limitErr.Error())
		}
		result = result + "\n" + limitSQL
	}

	//Union
	if len(query.Union) > 0 {
		for index, q := range query.Union {
			qSQL, qErr := q.SQL()
			if qErr != nil {
				return "", fmt.Errorf("Failed to generate UNION (index %d) SQL string: %s", index, qErr.Error())
			}
			result = result + "\n" + qSQL
		}
	}

	return result, nil
}
