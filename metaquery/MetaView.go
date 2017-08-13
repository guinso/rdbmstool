package metaquery

import (
	"github.com/guinso/rdbmstool"
	"github.com/guinso/rdbmstool/computed"
)

//MetaView interface to query dataview's meta data
type MetaView interface {
	//DbHandlerProxy : sql.DB or sql.Tx or compatible with it
	//string: view name
	//string: view name patterm (can use % as wildcard; example - 'hub_%')
	GetViewNames(rdbmstool.DbHandlerProxy, string, string) ([]string, error)

	//DbHandlerProxy : sql.DB or sql.Tx or compatible with it
	//string: database name
	//string: view name (example 'tax_invoice')
	GetViewDefinition(rdbmstool.DbHandlerProxy, string, string) (*computed.ViewDefinition, error)
}
