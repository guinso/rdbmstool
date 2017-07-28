package rdbmstool

import "database/sql"

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

//DbHandlerProxy interface to accept sql.Db or sql.Tx in order to allow execute SQL without
//knowning presence of transaction or not
type DbHandlerProxy interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
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
	//ColumnName         string
	ReferenceTableName string
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
