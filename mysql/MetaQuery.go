package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/guinso/rdbmstool"
)

//GetTableDefinition get all datatable's column(s) definition
func GetTableDefinition(db rdbmstool.DbHandlerProxy, dbName string,
	tableName string) (*rdbmstool.TableDefinition, error) {

	tableDef := rdbmstool.TableDefinition{
		Name: tableName}

	//get columns definition
	colsDef, colErr := getDataColumnDefinition(db, dbName, tableName)
	if colErr != nil {
		return nil, colErr
	}
	tableDef.Columns = colsDef

	//get primary key
	primaryKeys, priErr := getPrimaryKey(db, dbName, tableName)
	if priErr != nil {
		return nil, priErr
	}
	tableDef.PrimaryKey = primaryKeys

	//get foreign keys
	foreignKeys, fkErr := getForeignKey(db, dbName, tableName)
	if fkErr != nil {
		return nil, fkErr
	}
	tableDef.ForiegnKeys = foreignKeys

	//get unique keys
	uniqueKeys, ukErr := getUniqueKey(db, dbName, tableName)
	if ukErr != nil {
		return nil, ukErr
	}
	tableDef.UniqueKeys = uniqueKeys

	//get index keys
	indeces, indexErr := getIndexDefinitions(db, dbName, tableName)
	if indexErr != nil {
		return nil, indexErr
	}
	tableDef.Indices = indeces

	return &tableDef, nil
}

//GetTableNames get list of datatables' name which start with provided search pattern
//search pattern allow '%' as wild card; example 'hub_%'
func GetTableNames(db rdbmstool.DbHandlerProxy, databaseName string, tableNamePattern string) ([]string, error) {
	rows, err := db.Query("SELECT table_name FROM information_schema.tables"+
		" where table_schema=? AND table_name LIKE '"+tableNamePattern+"'", databaseName)

	if err != nil {
		return nil, err
	}

	var result []string
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)

		if err != nil {
			continue
		}

		result = append(result, tmp)
	}

	rows.Close()

	return result, nil
}

func getDataColumnDefinition(db rdbmstool.DbHandlerProxy, dbName string, tableName string) ([]rdbmstool.ColumnDefinition, error) {
	//TODO: get each column definition
	rows, err := db.Query("SELECT column_name, ordinal_position, column_default, "+
		"is_nullable, data_type, character_maximum_length,  "+
		"character_octet_length, numeric_precision, "+
		"numeric_scale, datetime_precision, "+
		"character_set_name, collation_name, column_key "+
		"FROM information_schema.columns "+
		"WHERE table_schema=? AND table_name=?", dbName, tableName)

	if err != nil {
		return nil, err
	}

	colsDef := []rdbmstool.ColumnDefinition{}

	for rows.Next() {
		var columnName string
		var ordinalPosition int
		var defaultValue sql.NullString
		var isNull string
		var dataType string
		var charMaxLength sql.NullInt64
		var numericLength sql.NullInt64
		var numericPrecision sql.NullInt64
		var numericScale sql.NullInt64
		var datetimePrecision sql.NullInt64
		var charset sql.NullString
		var collation sql.NullString
		var colKey string

		err := rows.Scan(&columnName, &ordinalPosition, &defaultValue,
			&isNull, &dataType, &charMaxLength, &numericLength,
			&numericPrecision, &numericScale, &datetimePrecision, &charset, &collation, &colKey)

		if err != nil {
			rows.Close()
			return nil, err
		}

		colDef := rdbmstool.ColumnDefinition{
			Name:       columnName,
			IsNullable: strings.Compare(isNull, "YES") == 0}

		//check datatype
		switch dataType {
		case "char":
			if !charMaxLength.Valid {
				return nil, fmt.Errorf(
					"%s.%s has null value on charactor length", tableName, columnName)
			}
			colDef.DataType = rdbmstool.CHAR
			colDef.Length = int(charMaxLength.Int64)
			break
		case "int":
			if !numericLength.Valid {
				return nil, fmt.Errorf("%s.%s has null value on numeric length", tableName, columnName)
			}
			colDef.DataType = rdbmstool.INTEGER
			colDef.Length = int(numericLength.Int64)
			break
		case "text":
			colDef.DataType = rdbmstool.TEXT
			break
		case "varchar":
			if !numericLength.Valid {
				return nil, fmt.Errorf("%s.%s has null value on character length", tableName, columnName)
			}
			colDef.DataType = rdbmstool.VARCHAR
			colDef.Length = int(charMaxLength.Int64)
			break
		case "datetime":
			if !datetimePrecision.Valid {
				return nil, fmt.Errorf("%s.%s has null value on datetime length", tableName, columnName)
			}
			colDef.DataType = rdbmstool.DATETIME
			colDef.Length = int(datetimePrecision.Int64)
			break
		case "date":
			colDef.DataType = rdbmstool.DATE
			break
		case "float":
			colDef.DataType = rdbmstool.FLOAT
			break
		case "double":
			colDef.DataType = rdbmstool.DOUBLE
			break
		case "decimal":
			if !numericScale.Valid {
				return nil, fmt.Errorf("%s.%s has null value on decimal precision", tableName, columnName)
			}
			if !numericPrecision.Valid {
				return nil, fmt.Errorf("%s.%s has null value on decimal length", tableName, columnName)
			}
			colDef.DataType = rdbmstool.DECIMAL
			colDef.Length = int(numericPrecision.Int64)
			colDef.DecimalPrecision = int(numericScale.Int64)
			break
		default:
			return nil, fmt.Errorf("%s.%s datatype(%s) is not support by package MySQL.MetaQuery",
				tableName, columnName, dataType)
		}

		colsDef = append(colsDef, colDef)
	}
	rows.Close()

	return colsDef, nil
}

//getUniqueKey get unique key for targeted
func getUniqueKey(db rdbmstool.DbHandlerProxy,
	dbName string, tableName string) ([]rdbmstool.UniqueKeyDefinition, error) {

	rows, err := db.Query(
		"SELECT a.CONSTRAINT_NAME, a.COLUMN_NAME "+
			"FROM `INFORMATION_SCHEMA`.`KEY_COLUMN_USAGE` a "+
			"INNER JOIN `INFORMATION_SCHEMA`.`TABLE_CONSTRAINTS` b "+
			"ON a.`CONSTRAINT_NAME` =   b.`CONSTRAINT_NAME`  AND "+
			"a.`TABLE_SCHEMA` =   b.`TABLE_SCHEMA` AND "+
			"a.`TABLE_NAME` =   b.`TABLE_NAME` "+
			"WHERE a.`TABLE_SCHEMA` = ? AND "+
			"a.`TABLE_NAME` = ? AND "+
			"b.`CONSTRAINT_TYPE` = 'UNIQUE' "+
			"ORDER BY a.`CONSTRAINT_NAME`", dbName, tableName)
	if err != nil {
		return nil, err
	}

	result := []rdbmstool.UniqueKeyDefinition{}
	currentConstraintName := ""
	for rows.Next() {
		var constraintName, colName string

		scanErr := rows.Scan(
			&constraintName, &colName)
		if scanErr != nil {
			rows.Close()
			return nil, scanErr
		}

		var uKey rdbmstool.UniqueKeyDefinition
		if strings.Compare(currentConstraintName, constraintName) != 0 {
			//since order by constraint_name column, same group UK shall keep in same row
			//new row created if new UK key discovered
			result = append(result, rdbmstool.UniqueKeyDefinition{
				ColumnNames: []string{}})
			currentConstraintName = constraintName
		}
		uKey = result[len(result)-1]
		result[len(result)-1].ColumnNames = append(uKey.ColumnNames, colName)
	}
	rows.Close()

	return result, nil
}

//getIndexDefinitions get indeces for targeted data table
func getIndexDefinitions(db rdbmstool.DbHandlerProxy,
	dbName string, tableName string) ([]rdbmstool.IndexKeyDefinition, error) {

	rows, err := db.Query("SELECT index_name, column_name "+
		"FROM information_schema.statistics WHERE table_schema = ? AND table_name = ? AND "+
		"non_unique = 1 ORDER BY index_name", dbName, tableName)

	if err != nil {
		return nil, err
	}

	result := []rdbmstool.IndexKeyDefinition{}
	currentIndex := ""
	for rows.Next() {
		var indexName, columnName string
		scanErr := rows.Scan(&indexName, &columnName)
		if scanErr != nil {
			rows.Close()
			return nil, scanErr
		}

		if strings.Compare(currentIndex, indexName) != 0 {
			currentIndex = indexName
			result = append(result, rdbmstool.IndexKeyDefinition{
				ColumnNames: []string{}})
		}
		indexDef := result[len(result)-1]
		indexDef.ColumnNames = append(indexDef.ColumnNames, columnName)
	}

	rows.Close()

	return result, nil
}

func getPrimaryKey(db rdbmstool.DbHandlerProxy, dbName string, tableName string) ([]string, error) {
	rows, err := db.Query("SELECT column_name FROM information_schema.key_column_usage "+
		"WHERE table_schema = ? AND table_name = ? AND "+
		"constraint_name = 'PRIMARY'", dbName, tableName)
	if err != nil {
		rows.Close()
		return nil, err
	}

	result := []string{}
	for rows.Next() {
		var tmp string
		scanErr := rows.Scan(&tmp)

		if scanErr != nil {
			return nil, scanErr
		}

		result = append(result, tmp)
	}
	rows.Close()

	return result, nil
}

func getForeignKey(db rdbmstool.DbHandlerProxy,
	dbName string, tableName string) ([]rdbmstool.ForeignKeyDefinition, error) {

	rows, err := db.Query(
		"SELECT a.CONSTRAINT_NAME, a.COLUMN_NAME, a.REFERENCED_TABLE_SCHEMA, "+
			"a.REFERENCED_TABLE_NAME, a.REFERENCED_COLUMN_NAME "+
			"FROM `information_schema`.`KEY_COLUMN_USAGE` a "+
			"INNER JOIN `information_schema`.`TABLE_CONSTRAINTS` b "+
			"ON a.`CONSTRAINT_NAME` =   b.`CONSTRAINT_NAME`  AND "+
			"a.`TABLE_SCHEMA` =   b.`TABLE_SCHEMA` AND "+
			"a.`TABLE_NAME` =   b.`TABLE_NAME` "+
			"WHERE a.`TABLE_SCHEMA` = ? AND "+
			"a.`TABLE_NAME` = ? AND "+
			"b.`CONSTRAINT_TYPE` = 'FOREIGN KEY' "+
			"ORDER BY a.`CONSTRAINT_NAME`", dbName, tableName)
	if err != nil {
		return nil, err
	}

	result := []rdbmstool.ForeignKeyDefinition{}
	currentConstraintName := ""
	for rows.Next() {
		var constraintName, colName, refDbName, refTableName, refColName string

		scanErr := rows.Scan(
			&constraintName, &colName, &refDbName, &refTableName, &refColName)
		if scanErr != nil {
			rows.Close()
			return nil, scanErr
		}

		var fkKey rdbmstool.ForeignKeyDefinition
		if strings.Compare(currentConstraintName, constraintName) != 0 {
			//since order by constraint_name column, same group FK shall keep in same row
			//new row created if new FK key discovered
			result = append(result, rdbmstool.ForeignKeyDefinition{
				Name:               currentConstraintName,
				ReferenceTableName: refTableName,
				Columns:            []rdbmstool.FKColumnDefinition{}})
			currentConstraintName = constraintName
		}
		fkKey = result[len(result)-1]
		tmpCol := rdbmstool.FKColumnDefinition{
			ColumnName:    colName,
			RefColumnName: refColName}
		result[len(result)-1].Columns = append(fkKey.Columns, tmpCol)
	}
	rows.Close()

	return result, nil
}
