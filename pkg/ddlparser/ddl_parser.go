package ddlparser

import (
	"fmt"
	"strings"

	executor "github.com/bytewatch/ddl-executor"
	"github.com/suifengpiao14/tormrepository/pkg"
)

//ParseDDLWithConfig 解析sql ddl,并格式化
func ParseDDLWithConfig(ddl string, dbCfg DatabaseConfig) (tables []*Table, err error) {
	tables, err = ParseDDL(ddl)
	if err != nil {
		return nil, err
	}
	tables = FormatWithConfig(tables, dbCfg)
	return tables, nil
}

//ParseDDL 解析sql ddl
func ParseDDL(ddl string) (tables []*Table, err error) {
	tables = make([]*Table, 0)
	conf := executor.NewDefaultConfig()
	inst := executor.NewExecutor(conf)
	databaseName := "test"
	ddl = fmt.Sprintf("create database `%s`;use `%s`;%s", databaseName, databaseName, ddl)
	err = inst.Exec(ddl)
	if err != nil {
		return nil, err
	}

	tableList, err := inst.GetTables(databaseName)
	if err != nil {
		return
	}

	for i := 0; i < len(tableList); i++ {
		tableName := tableList[i]
		tableDef, err := inst.GetTableDef(databaseName, tableName)
		if err != nil {
			return nil, err
		}

		table := &Table{
			TableName:  tableName,
			Columns:    make([]*Column, 0),
			EnumsConst: _Enums{},
			Comment:    tableDef.Comment,
		}
		for _, indice := range tableDef.Indices {
			if indice.Name == "PRIMARY" {
				table.PrimaryKey = indice.Columns[0] // 暂时取第一个为主键，不支持多字段主键
			}
		}
		for _, columnDef := range tableDef.Columns {

			goType, err := mysql2GoType(columnDef.Type, true)
			if err != nil {
				return nil, err
			}
			if isTimeMysqlType(columnDef.Type) && strings.Contains(columnDef.Name, "deleted_at") {
				table.DeleteColumn = columnDef.Name
			}

			columnPt := &Column{
				ColumnName:    columnDef.Name, // 这个地方记录数据库原始字段，包含前缀
				Type:          goType,
				Comment:       columnDef.Comment,
				Nullable:      columnDef.Nullable,
				Enums:         columnDef.Elems,
				AutoIncrement: columnDef.AutoIncrement,
				DefaultValue:  columnDef.DefaultValue,
				OnUpdate:      columnDef.OnUpdate,
			}
			if len(columnPt.Enums) > 0 {
				subEnumConst := enumsConst(table.TableNameTrimPrefix(), columnPt)
				table.EnumsConst = append(table.EnumsConst, subEnumConst...)
			}
			columnPt.OnCreate = columnPt.IsDefaultValueCurrentTimestamp() && !columnPt.OnUpdate    // 自动填充时间，但是更新时不变，认为是创建时间列
			if table.DeleteColumn == columnPt.Name || columnPt.Name == DEFAULT_COLUMN_DELETED_AT { // 删除记录列，通过配置指定，或者列名称为 DEFAULT_COLUMN_DELETED_AT 的值
				columnPt.OnDelete = true
			}

			table.Columns = append(table.Columns, columnPt)
		}
		tables = append(tables, table)
	}
	return
}

//FormatWithConfig 用DB配置，将ddl解析出的原始Table对象，格式化
func FormatWithConfig(tables []*Table, dbConfig DatabaseConfig) (formatedTables []*Table) {
	formatedTables = make([]*Table, len(tables))
	for index, table := range tables {
		table.DatabaseConfig = dbConfig
		if table.DeleteColumn == "" {
			table.DeleteColumn = dbConfig.DeletedAtColumn
		}

		for i, column := range table.Columns {
			columnName := column.ColumnName
			if dbConfig.ColumnPrefix != "" {
				columnName = strings.TrimPrefix(columnName, dbConfig.ColumnPrefix)
			}
			jsonName := pkg.ToLowerCamel(columnName)
			column.Name = columnName
			column.CamelName = pkg.ToCamel(columnName)
			column.Tag = fmt.Sprintf("`json:\"%s\" gorm:\"column:%s\"`", jsonName, column.Name)
			table.Columns[i] = column
		}
		formatedTables[index] = table
	}
	return formatedTables
}

