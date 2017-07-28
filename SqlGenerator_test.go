package rdbmstool

import (
	"strings"
	"testing"
)

func TestGenerateDataTableSQL(t *testing.T) {
	tableDef := TableDefinition{
		Name: "Koko",
		Columns: []ColumnDefinition{
			ColumnDefinition{
				Name:             "Jojo",
				DataType:         CHAR,
				Length:           32,
				IsNullable:       false,
				DecimalPrecision: 0}},
		PrimaryKey:  []string{"Jojo"},
		ForiegnKeys: []ForeignKeyDefinition{},
		UniqueKeys:  []UniqueKeyDefinition{},
		Indices:     []IndexKeyDefinition{}}

	sql, err := GenerateTableSQL(&tableDef)

	if err != nil {
		t.Errorf("Unable to generate SQL: " + err.Error())
		return
	}

	expectedSQL :=
		"CREATE TABLE `Koko`(\n" +
			"`Jojo` char(32) COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
			"PRIMARY KEY(`Jojo`)\n" +
			") ENGINE=innodb DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"

	if strings.Compare(sql, expectedSQL) != 0 {
		t.Errorf(
			"Unexpected SQL statement generated:\n==Expected==\n%s\n\n==Result==\n%s",
			expectedSQL, sql)
	}
}
