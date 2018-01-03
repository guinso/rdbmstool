package rdbmstool

import (
	"database/sql"
)

const (
	collate = "utf8mb4_unicode_ci"
)

//SQLGenerator generate SQL statement interface
type SQLGenerator interface {
	SQL() (string, error)
	Validate() error
}

//DbHandlerProxy interface to accept sql.Db or sql.Tx in order to allow execute SQL without
//knowning presence of transaction or not
type DbHandlerProxy interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
