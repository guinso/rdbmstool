package rdbmstool

import (
	"errors"
	"fmt"
	"strings"
)

// ColumnDataType is enum for data column's data type
type ColumnDataType uint8

// Data table column's data type definition
const (
	CHAR     ColumnDataType = iota + 1
	INTEGER  ColumnDataType = iota + 1
	DECIMAL  ColumnDataType = iota + 1
	FLOAT    ColumnDataType = iota + 1
	TEXT     ColumnDataType = iota + 1
	DATE     ColumnDataType = iota + 1
	DATETIME ColumnDataType = iota + 1
	BOOLEAN  ColumnDataType = iota + 1
	VARCHAR  ColumnDataType = iota + 1
	DOUBLE   ColumnDataType = iota + 1
)

//String return database's data column type in string format
func (colType ColumnDataType) String() string {
	if colType == CHAR {
		return "CHAR"
	} else if colType == INTEGER {
		return "INTEGER"
	} else if colType == DECIMAL {
		return "DECIMAL"
	} else if colType == FLOAT {
		return "FLOAT"
	} else if colType == TEXT {
		return "TEXT"
	} else if colType == DATE {
		return "DATE"
	} else if colType == DATETIME {
		return "DATETIME"
	} else if colType == BOOLEAN {
		return "BOOLEAN"
	} else if colType == VARCHAR {
		return "VARCHAR"
	} else if colType == DOUBLE {
		return "DOUBLE"
	} else {
		return "unknown"
	}
}

// TableDefinition is information to create a data table
type TableDefinition struct {
	Name        string
	Columns     []ColumnDefinition
	PrimaryKey  []string //PK can form by more than one column
	ForiegnKeys []ForeignKeyDefinition
	UniqueKeys  []UniqueKeyDefinition
	Indices     []IndexKeyDefinition
	//do I need to include encoding as well?
}

// ColumnDefinition is information to defined a data table column
type ColumnDefinition struct {
	Name             string
	DataType         ColumnDataType
	Length           int
	IsNullable       bool
	DecimalPrecision int
}

// ForeignKeyDefinition is information to create a RDBMS FK
type ForeignKeyDefinition struct {
	Name               string //foreign key name display on database (optional)
	ReferenceTableName string
	//ColumnName         string
	//ReferenceSchemaName string
	//ReferenceColumnName string
	Columns []FKColumnDefinition
}

//FKColumnDefinition information for FK column reference
type FKColumnDefinition struct {
	ColumnName    string
	RefColumnName string
}

// UniqueKeyDefinition is information to create an unique key
//NOTE: a single unique key can made up from multiple columns
type UniqueKeyDefinition struct {
	ColumnNames []string
}

//IndexKeyDefinition information to hold a single index key
//NOTE: a single index key can made up from multiple columns
type IndexKeyDefinition struct {
	ColumnNames []string
}


//GenerateTableSQL to generate "create table" SQL statement
func (tableDef *TableDefinition) GenerateTableSQL() (string, error) {
	if tableDef == nil {
		return "", errors.New("input parameter is null")
	}

	//validate tableDef integrity
	tableDefValidErr := tableDef.ValidateTableDefinition()
	if tableDefValidErr != nil {
		return "", tableDefValidErr
	}

	//generate based on tableDef variable
	var colSQL string
	var tmpSQL string
	var tmpErr error

	//generate column SQL statement
	for index, col := range tableDef.Columns {
		tmpSQL, tmpErr = tableDef.generateColumnSQL(&col)

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
	pkSQL, pkErr := tableDef.generatePrimaryKeySQL()
	if pkErr != nil {
		return "", pkErr
	}

	if pkSQL != "" {
		colSQL = colSQL + ",\n" + pkSQL
	}

	//generate Unique key SQL statement
	ukSQL, ukErr := tableDef.generateUniqueKeySQL()
	if ukErr != nil {
		return "", ukErr
	}

	if ukSQL != "" {
		colSQL = colSQL + ",\n" + ukSQL
	}

	//generate index key SQL statement
	ikSQL, ikErr := tableDef.generateIndexSQL()
	if ikErr != nil {
		return "", ikErr
	}

	if ikSQL != "" {
		colSQL = colSQL + ",\n" + ikSQL
	}

	//generate FK SQL statement
	fkSQL, fkErr := tableDef.generateForeignKeySQL()
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

func (tableDef *TableDefinition) generateColumnSQL(colDef *ColumnDefinition) (string, error) {
	switch colDef.DataType {
	case CHAR:
		return fmt.Sprintf("`%s` char(%d) COLLATE %s %s",
			colDef.Name, colDef.Length, collate, tableDef.generateIsNullSQL(colDef.IsNullable)), nil
	case INTEGER:
		return fmt.Sprintf("`%s` int(%d) %s",
			colDef.Name, colDef.Length, tableDef.generateIsNullSQL(colDef.IsNullable)), nil
	case DECIMAL:
		return fmt.Sprintf("`%s` decimal(%d,%d) %s",
			colDef.Name, colDef.Length, colDef.DecimalPrecision, tableDef.generateIsNullSQL(colDef.IsNullable)), nil
	case FLOAT:
		return fmt.Sprintf("`%s` float %s",
			colDef.Name, tableDef.generateIsNullSQL(colDef.IsNullable)), nil
	case TEXT:
		return fmt.Sprintf("`%s` text COLLATE %s %s",
			colDef.Name, collate, tableDef.generateIsNullSQL(colDef.IsNullable)), nil
	case DATE:
		return fmt.Sprintf("`%s` date %s",
			colDef.Name, tableDef.generateIsNullSQL(colDef.IsNullable)), nil
	case DATETIME:
		return fmt.Sprintf("`%s` datetime %s",
			colDef.Name, tableDef.generateIsNullSQL(colDef.IsNullable)), nil
	case BOOLEAN:
		return fmt.Sprintf("`%s` tinyint(1) %s", colDef.Name, tableDef.generateIsNullSQL(colDef.IsNullable)), nil
	default:
		return "", fmt.Errorf(
			"unknown data column (%s) type: %d", colDef.Name, colDef.DataType)
	}
}

func (def *TableDefinition) generateIsNullSQL(isNullable bool) string {
	if isNullable {
		return "NULL"
	}

	return "NOT NULL"

}

func (tableDef *TableDefinition) generateIndexSQL() (string, error) {
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
					uCols = "`" + colName + "`"
				} else {
					uName = uName + "_" + colName
					uCols = uCols + ",`" + colName + "`"
				}
			}
			tmpSQL = fmt.Sprintf("KEY `%s` (%s)", uName, uCols)

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

func (tableDef *TableDefinition) generateUniqueKeySQL() (string, error) {
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
					uCols = "`" + colName + "`"
				} else {
					uName = uName + "_" + colName
					uCols = uCols + ",`" + colName + "`"
				}
			}
			tmpSQL = fmt.Sprintf("UNIQUE KEY `%s` (%s)", uName, uCols)

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

func (tableDef *TableDefinition) generateForeignKeySQL() (string, error) {
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
				} else {
					baseCols += "," + "`" + coll.ColumnName + "`"
					refCols += "," + "`" + coll.RefColumnName + "`"
				}
			}

			if strings.Compare(fk.Name, "") == 0 {
				tmpSQL = fmt.Sprintf("CONSTRAINT `%s_ibfk_%d` FOREIGN KEY (%s) REFERENCES `%s` (%s)",
					tableDef.Name, index+1, baseCols,
					fk.ReferenceTableName, refCols)
			} else {
				tmpSQL = fk.Name
			}

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

func (tableDef *TableDefinition) generatePrimaryKeySQL() (string, error) {
	if tableDef == nil || tableDef.PrimaryKey == nil {
		return "", errors.New("PK SQL generator: Cannot pass null parameter")
	}

	var sql string
	length := len(tableDef.PrimaryKey)
	if length > 0 {
		for index, pk := range tableDef.PrimaryKey {
			if index == 0 {
				sql = "`" + pk + "`"
			} else {
				sql = sql + ",`" + pk + "`"
			}
		}

		sql = fmt.Sprintf("PRIMARY KEY(%s)", sql)
	} else {
		sql = ""
	}

	return sql, nil
}

func (tableDef *TableDefinition) ValidateTableDefinition() error {
	//TODO: implement validation
	return nil
}

func (tableDef *TableDefinition) ValidateColumnDefinition() error {
	//TODO: implement validation
	return nil
}