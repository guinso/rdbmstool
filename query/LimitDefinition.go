package query

import (
	"fmt"
)

//LimitDefinition SQL LIMIT statement definition
type LimitDefinition struct {
	RowCount     int
	Offset       int
	EnableOffset bool
}

//SQL generate LIMIT statement SQL string
func (limit *LimitDefinition) SQL() (string, error) {
	if limit.EnableOffset {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit.RowCount, limit.Offset), nil
	}

	return fmt.Sprintf("LIMIT %d", limit.RowCount), nil
}

//NewLimitDefinition create new LIMIT SQL statement definition
func NewLimitDefinition(rowCount int, offset int) *LimitDefinition {
	return &LimitDefinition{
		RowCount:     rowCount,
		Offset:       offset,
		EnableOffset: true}
}

//NewLimitDefinitionNoOffset create new LIMIT SQL statement definition without offset enabled
func NewLimitDefinitionNoOffset(rowCount int) *LimitDefinition {
	return &LimitDefinition{
		RowCount:     rowCount,
		Offset:       0,
		EnableOffset: false}
}
