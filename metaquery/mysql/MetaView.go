package mysql

import (
	"errors"

	"github.com/guinso/rdbmstool"
	"github.com/guinso/rdbmstool/computed"
)

//GetViewNames show all database's view table
//db : sql.DB or sql.Tx or compatible with it
//dbName: database schema name
//searchPattern: view name patterm (can use % as wildcard; example - 'view_%')
func (meta *MetaQuery) GetViewNames(db rdbmstool.DbHandlerProxy, dbName string, searchPattern string) ([]string, error) {
	rows, queryErr := db.Query("SELECT table_name FROM information_schema.views"+
		" WHERE table_schema = ? AND table_name LIKE ?", dbName, searchPattern)

	if queryErr != nil {
		return nil, queryErr
	}

	result := make([]string, 5)
	for rows.Next() {
		tmp := ""
		rows.Scan(&tmp)
		result = append(result, tmp)
	}

	return result, nil
}

//GetViewDefinition get data view definition
//db : sql.DB or sql.Tx or compatible with it
//dbName: database name
//viewName: view name (example 'tax_invoice')
func (meta *MetaQuery) GetViewDefinition(db rdbmstool.DbHandlerProxy, dbName string, viewName string) (
	*computed.ViewDefinition, error) {
	//TODO: implement SQL parser... -_-|||

	return nil, errors.New("Not implemented yet")
}
