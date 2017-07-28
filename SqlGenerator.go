package rdbmstool

import (
	"errors"
	"fmt"
)

const (
	collate = "utf8mb4_unicode_ci"
)

//GenerateTableSQL to generate "create table" SQL statement
func GenerateTableSQL(tableDef *TableDefinition) (string, error) {
	if tableDef == nil {
		return "", errors.New("input parameter is null")
	}

	//validate tableDef integrity
	tableDefValidErr := ValidateTableDefinition(tableDef)
	if tableDefValidErr != nil {
		return "", tableDefValidErr
	}

	//generate based on tableDef variable
	var colSQL string
	var tmpSQL string
	var tmpErr error

	//generate column SQL statement
	for index, col := range tableDef.Columns {
		tmpSQL, tmpErr = generateColumnSQL(&col)

		if tmpErr != nil {
			return "", tmpErr
		}

		if index == 0 {
			colSQL = tmpSQL
		} else {
			colSQL = colSQL + ",\n" + tmpSQL
		}
	}

	//generate PK SQL statement
	pkSQL, pkErr := generatePrimaryKeySQL(tableDef)
	if pkErr != nil {
		return "", pkErr
	}

	if pkSQL != "" {
		colSQL = colSQL + ",\n" + pkSQL
	}

	//generate Unique key SQL statement
	ukSQL, ukErr := generateUniqueKeySQL(tableDef)
	if ukErr != nil {
		return "", ukErr
	}

	if ukSQL != "" {
		colSQL = colSQL + ",\n" + ukSQL
	}

	//generate index key SQL statement
	ikSQL, ikErr := generateIndexSQL(tableDef)
	if ikErr != nil {
		return "", ikErr
	}

	if ikSQL != "" {
		colSQL = colSQL + ",\n" + ikSQL
	}

	//generate FK SQL statement
	fkSQL, fkErr := generateForeignKeySQL(tableDef)
	if fkErr != nil {
		return "", fkErr
	}

	if fkSQL != "" {
		colSQL = colSQL + ",\n" + fkSQL
	}

	sqlStatement := fmt.Sprintf(
		"CREATE TABLE `%s`(\n%s\n) ENGINE=innodb DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;",
		tableDef.Name, colSQL)

	return sqlStatement, nil
}

func generateColumnSQL(colDef *ColumnDefinition) (string, error) {
	switch colDef.DataType {
	case CHAR:
		return fmt.Sprintf("`%s` char(%d) COLLATE %s %s",
			colDef.Name, colDef.Length, collate, generateIsNullSQL(colDef.IsNullable)), nil
	case INTEGER:
		return fmt.Sprintf("`%s` int(%d) %s",
			colDef.Name, colDef.Length, generateIsNullSQL(colDef.IsNullable)), nil
	case DECIMAL:
		return fmt.Sprintf("`%s` decimal(%d,%d) %s",
			colDef.Name, colDef.Length, colDef.DecimalPrecision, generateIsNullSQL(colDef.IsNullable)), nil
	case FLOAT:
		return fmt.Sprintf("`%s` float %s",
			colDef.Name, generateIsNullSQL(colDef.IsNullable)), nil
	case TEXT:
		return fmt.Sprintf("`%s` text COLLATE %s %s",
			colDef.Name, collate, generateIsNullSQL(colDef.IsNullable)), nil
	case DATE:
		return fmt.Sprintf("`%s` date %s",
			colDef.Name, generateIsNullSQL(colDef.IsNullable)), nil
	case DATETIME:
		return fmt.Sprintf("`%s` datetime %s",
			colDef.Name, generateIsNullSQL(colDef.IsNullable)), nil
	case BOOLEAN:
		return fmt.Sprintf("`%s` tinyint(1) %s", colDef.Name, generateIsNullSQL(colDef.IsNullable)), nil
	default:
		return "", fmt.Errorf(
			"unknown data column (%s) type: %d", colDef.Name, colDef.DataType)
	}
}

func generateIsNullSQL(isNullable bool) string {
	if isNullable {
		return "NULL"
	}

	return "NOT NULL"

}

func generateIndexSQL(tableDef *TableDefinition) (string, error) {
	if tableDef == nil || tableDef.Indices == nil {
		return "", errors.New("Index SQL generator: Cannot pass null parameter")
	}

	var sql string
	var tmpSQL string
	length := len(tableDef.Indices)
	if length > 0 {
		for index, ik := range tableDef.Indices {

			var uName string
			var uCols string
			for colIndex, colName := range ik.ColumnNames {
				if colIndex == 0 {
					uName = colName
					uCols = colName
				} else {
					uName = uName + "_" + colName
					uCols = uCols + "," + colName
				}
			}
			tmpSQL = fmt.Sprintf("KEY `%s` (`%s`)", uName, uCols)

			if index == 0 {
				sql = tmpSQL
			} else {
				sql = sql + ",\n" + tmpSQL
			}
		}
	} else {
		sql = ""
	}

	return sql, nil
}

func generateUniqueKeySQL(tableDef *TableDefinition) (string, error) {
	if tableDef == nil || tableDef.UniqueKeys == nil {
		return "", errors.New("Unique Key SQL generator: Cannot pass null parameter")
	}

	var sql string
	var tmpSQL string
	length := len(tableDef.UniqueKeys)
	if length > 0 {
		for index, uk := range tableDef.UniqueKeys {

			var uName string
			var uCols string
			for colIndex, colName := range uk.ColumnNames {
				if colIndex == 0 {
					uName = colName
					uCols = colName
				} else {
					uName = uName + "_" + colName
					uCols = uCols + "," + colName
				}
			}
			tmpSQL = fmt.Sprintf("UNIQUE KEY `%s` (`%s`)", uName, uCols)

			if index == 0 {
				sql = tmpSQL
			} else {
				sql = sql + ",\n" + tmpSQL
			}
		}
	} else {
		sql = ""
	}

	return sql, nil
}

func generateForeignKeySQL(tableDef *TableDefinition) (string, error) {
	if tableDef == nil || tableDef.ForiegnKeys == nil {
		return "", errors.New("FK SQL generator: Cannot pass null parameter")
	}

	var sql string
	var tmpSQL string
	length := len(tableDef.ForiegnKeys)
	if length > 0 {
		for index, fk := range tableDef.ForiegnKeys {
			//build columns strings
			baseCols := ""
			refCols := ""
			for colIndex, coll := range fk.Columns {
				if colIndex == 0 {
					baseCols = "`" + coll.ColumnName + "`"
					refCols = "`" + coll.RefColumnName + "`"
				}

				baseCols += "," + "`" + coll.ColumnName + "`"
				refCols += "," + "`" + coll.RefColumnName + "`"
			}

			tmpSQL = fmt.Sprintf("CONSTRAINT `%s_ibfk_%d` FOREIGN KEY (%s) REFERENCES `%s` (%s)",
				tableDef.Name, index+1, baseCols,
				fk.ReferenceTableName, refCols)

			if index == 0 {
				sql = tmpSQL
			} else {
				sql = sql + ",\n" + tmpSQL
			}
		}
	} else {
		sql = ""
	}

	return sql, nil
}

func generatePrimaryKeySQL(tableDef *TableDefinition) (string, error) {
	if tableDef == nil || tableDef.PrimaryKey == nil {
		return "", errors.New("PK SQL generator: Cannot pass null parameter")
	}

	var sql string
	length := len(tableDef.PrimaryKey)
	if length > 0 {
		for index, pk := range tableDef.PrimaryKey {
			if index == 0 {
				sql = pk
			} else {
				sql = sql + "," + pk
			}
		}

		sql = fmt.Sprintf("PRIMARY KEY(`%s`)", sql)
	} else {
		sql = ""
	}

	return sql, nil
}

func ValidateTableDefinition(tableDef *TableDefinition) error {
	//TODO: implement validation
	return nil
}

func ValidateColumnDefinition(colDef *ColumnDefinition) error {
	//TODO: implement validation
	return nil
}
