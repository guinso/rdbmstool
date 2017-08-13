package metaquery

import "github.com/guinso/rdbmstool"

//MetaTable interface to query datatable's meta data
type MetaTable interface {
	//DbHandlerProxy : sql.DB or sql.Tx or compatible with it
	//string: database name
	//string: datatable name patterm (can use % as wildcard; example - 'hub_%')
	GetTableNames(rdbmstool.DbHandlerProxy, string, string) ([]string, error)

	//DbHandlerProxy : sql.DB or sql.Tx or compatible with it
	//string: database name
	//string: datatable name (example 'tax_invoice')
	GetTableDefinition(rdbmstool.DbHandlerProxy, string, string) (*rdbmstool.TableDefinition, error)
}
