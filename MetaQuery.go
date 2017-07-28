package rdbmstool

//MetaQuery interface to query datatable's meta data
type MetaQuery interface {
	//DbHandlerProxy : sql.DB or sql.Tx or compatible with it
	//string: database name
	//string: datatable name patterm (can use % as wildcard; example - 'hub_%')
	GetTableNames(DbHandlerProxy, string, string) ([]string, error)

	//DbHandlerProxy : sql.DB or sql.Tx or compatible with it
	//string: database name
	//string: datatable name (example 'tax_invoice')
	GetTableDefinition(DbHandlerProxy, string, string) (*TableDefinition, error)
}
