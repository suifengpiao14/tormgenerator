package ddlparser

import (
	"fmt"
	"strings"

	executor "github.com/bytewatch/ddl-executor"
	"github.com/suifengpiao14/generaterepository/pkg"
)

//ParseDDL 解析sql ddl
func ParseDDL(ddl string, dbConfig DatabaseConfig) (tables []*Table, err error) {
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
			DatabaseConfig: dbConfig,
			TableName:      tableName,
			Columns:        make([]*Column, 0),
			EnumsConst:     _Enums{},
			Comment:        tableDef.Comment,
			DeleteColumn:   dbConfig.DeletedAtColumn,
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

			columnName := columnDef.Name
			if dbConfig.ColumnPrefix != "" {
				columnName = strings.TrimPrefix(columnName, dbConfig.ColumnPrefix)
			}

			columnPt := &Column{
				ColumnName:    columnDef.Name, // 这个地方记录数据库原始字段，包含前缀
				Name:          columnName,
				CamelName:     pkg.ToCamel(columnName),
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
