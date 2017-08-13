package mysql

import (
	"database/sql"
	"fmt"
	"testing"

	//explicitly include GO mysql library
	_ "github.com/go-sql-driver/mysql"
)

func TestGetForeignKey(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8", "root", "", "localhost", 3306, "test"))

	if err != nil {
		t.Error(err.Error())
		return
	}

	fkDefs, fkErr := getForeignKey(db, "test", "hub_tax_invoice_rev0")
	if fkErr != nil {
		t.Error(fkErr)
		return
	}

	if len(fkDefs) == 0 {
		t.Error("test.hub_tax_invoice_rev0 should have 3 foreign keys")
		return
	}

	if len(fkDefs[0].Columns) == 0 {
		t.Error("FK column expect to have one")
	}
}

func TestGetUniqueKey(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8", "root", "", "localhost", 3306, "test"))

	if err != nil {
		t.Error(err.Error())
		return
	}

	fkDefs, fkErr := getUniqueKey(db, "test", "hub_tax_invoice_rev0")
	if fkErr != nil {
		t.Error(fkErr)
		return
	}

	if len(fkDefs) == 0 {
		t.Error("test.hub_tax_invoice_rev0 should have 1 unique key")
		return
	}

	if len(fkDefs[0].ColumnNames) == 0 {
		t.Error("UK column name expect to have one")
	}
}
