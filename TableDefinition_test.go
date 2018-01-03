package rdbmstool

import (
	"strings"
	"testing"
)

func TestGenerateDataTableSQL(t *testing.T) {
	def := TableDefinition{
		Name:        "account_role",
		Columns:     []ColumnDefinition{},
		PrimaryKey:  []string{"id"},
		ForiegnKeys: []ForeignKeyDefinition{},
		UniqueKeys: []UniqueKeyDefinition{
			UniqueKeyDefinition{
				ColumnNames: []string{"account_id", "role_id"},
			},
		},
		Indices: []IndexKeyDefinition{},
	}

	def.Columns = append(def.Columns, ColumnDefinition{
		Name:             "id",
		DataType:         CHAR,
		Length:           100,
		IsNullable:       false,
		DecimalPrecision: 0,
	})
	def.Columns = append(def.Columns, ColumnDefinition{
		Name:             "account_id",
		DataType:         CHAR,
		Length:           100,
		IsNullable:       false,
		DecimalPrecision: 0,
	})
	def.Columns = append(def.Columns, ColumnDefinition{
		Name:             "role_id",
		DataType:         CHAR,
		Length:           100,
		IsNullable:       false,
		DecimalPrecision: 0,
	})

	def.Indices = append(def.Indices, IndexKeyDefinition{
		ColumnNames: []string{"account_id"},
	})
	def.Indices = append(def.Indices, IndexKeyDefinition{
		ColumnNames: []string{"role_id"},
	})
	def.Indices = append(def.Indices, IndexKeyDefinition{
		ColumnNames: []string{"account_id", "role_id"},
	})

	def.ForiegnKeys = append(def.ForiegnKeys, ForeignKeyDefinition{
		ReferenceTableName: "account",
		Columns: []FKColumnDefinition{
			FKColumnDefinition{
				ColumnName:    "account_id",
				RefColumnName: "id",
			},
		},
	})
	def.ForiegnKeys = append(def.ForiegnKeys, ForeignKeyDefinition{
		ReferenceTableName: "role",
		Columns: []FKColumnDefinition{
			FKColumnDefinition{
				ColumnName:    "role_id",
				RefColumnName: "id",
			},
		},
	})

	sql, err := def.GenerateTableSQL()

	if err != nil {
		t.Errorf("Unable to generate SQL: " + err.Error())
		return
	}

	expectedSQL :=
		"CREATE TABLE `account_role`(\n" +
			"`id` char(100) COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
			"`account_id` char(100) COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
			"`role_id` char(100) COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
			"PRIMARY KEY(`id`),\n" +
			"UNIQUE KEY `account_id_role_id` (`account_id`,`role_id`),\n" +
			"KEY `account_id` (`account_id`),\n" +
			"KEY `role_id` (`role_id`),\n" +
			"KEY `account_id_role_id` (`account_id`,`role_id`),\n" +
			"CONSTRAINT `account_role_ibfk_1` FOREIGN KEY (`account_id`) REFERENCES `account` (`id`),\n" +
			"CONSTRAINT `account_role_ibfk_2` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`)\n" +
			") ENGINE=innodb DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"

	if strings.Compare(sql, expectedSQL) != 0 {
		t.Errorf(
			"Unexpected SQL statement generated:\n==Expected==\n%s\n\n==Result==\n%s",
			expectedSQL, sql)
	}
}
