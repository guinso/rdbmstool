package rdbmstool

//TableBuilder SQL create table statement builder
type TableBuilder struct {
	tableDefinition *TableDefinition
}

//NewTableBuilder create new SQL table definition builder
func NewTableBuilder() *TableBuilder {
	return &TableBuilder{
		tableDefinition: &TableDefinition{
			Name:        "",
			Columns:     []ColumnDefinition{},
			PrimaryKey:  []string{},
			ForiegnKeys: []ForeignKeyDefinition{},
			UniqueKeys:  []UniqueKeyDefinition{},
			Indices:     []IndexKeyDefinition{}}}
}

//GetTableName get table name
func (builder *TableBuilder) GetTableName() string {
	return builder.tableDefinition.Name
}

//SQL generate table definition SQL statement
func (builder *TableBuilder) SQL() (string, error) {
	return builder.tableDefinition.SQL()
}

//TableName set table name
func (builder *TableBuilder) TableName(tableName string) *TableBuilder {
	builder.tableDefinition.Name = tableName
	return builder
}

//AddPrimaryKey add primary key
func (builder *TableBuilder) AddPrimaryKey(primaryKey string) *TableBuilder {
	builder.tableDefinition.PrimaryKey = append(builder.tableDefinition.PrimaryKey, primaryKey)

	return builder
}

//AddForeignKey add simple foreign key
func (builder *TableBuilder) AddForeignKey(columnName, referenceTableName, referenceColumnName string) *TableBuilder {
	builder.tableDefinition.ForiegnKeys = append(builder.tableDefinition.ForiegnKeys, ForeignKeyDefinition{
		Name:               "",
		ReferenceTableName: referenceTableName,
		Columns: []FKColumnDefinition{
			FKColumnDefinition{
				ColumnName:    columnName,
				RefColumnName: referenceColumnName}}})
	return builder
}

//AddForeignKeyMultiColumn add foreign key with support multi columns reference
func (builder *TableBuilder) AddForeignKeyMultiColumn(referenceTableName string, columns []FKColumnDefinition) *TableBuilder {
	builder.tableDefinition.ForiegnKeys = append(builder.tableDefinition.ForiegnKeys, ForeignKeyDefinition{
		Name:               "",
		ReferenceTableName: referenceTableName,
		Columns:            columns})
	return builder
}

//AddUniqueKey add single column unique key
func (builder *TableBuilder) AddUniqueKey(columnName string) *TableBuilder {
	builder.tableDefinition.UniqueKeys = append(builder.tableDefinition.UniqueKeys, UniqueKeyDefinition{
		ColumnNames: []string{columnName}})

	return builder
}

//AddUniqueKeyMultiColumn add multi column unique key
func (builder *TableBuilder) AddUniqueKeyMultiColumn(columnsName []string) *TableBuilder {
	builder.tableDefinition.UniqueKeys = append(builder.tableDefinition.UniqueKeys, UniqueKeyDefinition{
		ColumnNames: columnsName})

	return builder
}

//AddIndexKey add single column index key
func (builder *TableBuilder) AddIndexKey(columnName string) *TableBuilder {
	builder.tableDefinition.Indices = append(builder.tableDefinition.Indices, IndexKeyDefinition{
		ColumnNames: []string{columnName}})

	return builder
}

//AddIndexKeyMultiColumn add multi column index key
func (builder *TableBuilder) AddIndexKeyMultiColumn(columnsName []string) *TableBuilder {
	builder.tableDefinition.Indices = append(builder.tableDefinition.Indices, IndexKeyDefinition{
		ColumnNames: columnsName})

	return builder
}

//AddColumn add data column definition
func (builder *TableBuilder) AddColumn(columnName string,
	colType ColumnDataType, dataLength int, isNullable bool, decimalPrecision int) *TableBuilder {
	builder.tableDefinition.Columns = append(builder.tableDefinition.Columns, ColumnDefinition{
		Name:             columnName,
		DataType:         colType,
		Length:           dataLength,
		IsNullable:       isNullable,
		DecimalPrecision: decimalPrecision})

	return builder
}

//AddColumnInt add integer column definition
func (builder *TableBuilder) AddColumnInt(columnName string, dataLength int, isNullable bool) *TableBuilder {
	builder.tableDefinition.Columns = append(builder.tableDefinition.Columns, ColumnDefinition{
		Name:             columnName,
		DataType:         INTEGER,
		Length:           dataLength,
		IsNullable:       isNullable,
		DecimalPrecision: 0})

	return builder
}

//AddColumnChar add char column definition
func (builder *TableBuilder) AddColumnChar(columnName string, dataLength int, isNullable bool) *TableBuilder {
	builder.tableDefinition.Columns = append(builder.tableDefinition.Columns, ColumnDefinition{
		Name:             columnName,
		DataType:         CHAR,
		Length:           dataLength,
		IsNullable:       isNullable,
		DecimalPrecision: 0})

	return builder
}

//AddColumnVarchar add varchar column definition
func (builder *TableBuilder) AddColumnVarchar(columnName string, dataLength int, isNullable bool) *TableBuilder {
	builder.tableDefinition.Columns = append(builder.tableDefinition.Columns, ColumnDefinition{
		Name:             columnName,
		DataType:         VARCHAR,
		Length:           dataLength,
		IsNullable:       isNullable,
		DecimalPrecision: 0})

	return builder
}

//AddColumnDecimal add decimal column definition
func (builder *TableBuilder) AddColumnDecimal(columnName string, dataLength int, decimalPrecision int, isNullable bool) *TableBuilder {
	builder.tableDefinition.Columns = append(builder.tableDefinition.Columns, ColumnDefinition{
		Name:             columnName,
		DataType:         DECIMAL,
		Length:           dataLength,
		IsNullable:       isNullable,
		DecimalPrecision: decimalPrecision,
	})

	return builder
}

//AddColumnFloat add float column definition
func (builder *TableBuilder) AddColumnFloat(columnName string, isNullable bool) *TableBuilder {
	builder.tableDefinition.Columns = append(builder.tableDefinition.Columns, ColumnDefinition{
		Name:             columnName,
		DataType:         FLOAT,
		Length:           0,
		IsNullable:       isNullable,
		DecimalPrecision: 0,
	})

	return builder
}

//AddColumnDate add Date column definition (day, month, and year)
func (builder *TableBuilder) AddColumnDate(columnName string, isNullable bool) *TableBuilder {
	builder.tableDefinition.Columns = append(builder.tableDefinition.Columns, ColumnDefinition{
		Name:             columnName,
		DataType:         DATE,
		Length:           0,
		IsNullable:       isNullable,
		DecimalPrecision: 0})

	return builder
}

//AddColumnDateTime add Datetime column definition (day, month, and year, hour, minute, and second)
func (builder *TableBuilder) AddColumnDateTime(columnName string, isNullable bool) *TableBuilder {
	builder.tableDefinition.Columns = append(builder.tableDefinition.Columns, ColumnDefinition{
		Name:             columnName,
		DataType:         DATETIME,
		Length:           0,
		IsNullable:       isNullable,
		DecimalPrecision: 0})

	return builder
}

//AddColumnBoolean add boolean column definition
func (builder *TableBuilder) AddColumnBoolean(columnName string, isNullable bool) *TableBuilder {
	builder.tableDefinition.Columns = append(builder.tableDefinition.Columns, ColumnDefinition{
		Name:             columnName,
		DataType:         BOOLEAN,
		Length:           0,
		IsNullable:       isNullable,
		DecimalPrecision: 0})

	return builder
}

//AddColumnText add text string column definition
func (builder *TableBuilder) AddColumnText(columnName string, isNullable bool) *TableBuilder {
	builder.tableDefinition.Columns = append(builder.tableDefinition.Columns, ColumnDefinition{
		Name:             columnName,
		DataType:         TEXT,
		Length:           0,
		IsNullable:       isNullable,
		DecimalPrecision: 0})

	return builder
}
