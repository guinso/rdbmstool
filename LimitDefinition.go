package rdbmstool

import (
	"fmt"
)

//LimitDefinition SQL LIMIT statement definition
type LimitDefinition struct {
	RowCount int
	Offset   int
}

//SQL generate LIMIT statement SQL string
func (limit *LimitDefinition) SQL() (string, error) {
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit.RowCount, limit.Offset), nil
}

//NewLimitDefinition create new LIMIT SQL statement definition
func NewLimitDefinition(rowCount int, offset int) *LimitDefinition {
	return &LimitDefinition{
		RowCount: rowCount,
		Offset:   offset}
}
