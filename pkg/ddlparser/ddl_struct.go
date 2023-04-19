package ddlparser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/tormrepository/pkg"

	executor "github.com/bytewatch/ddl-executor"
)

type DatabaseConfig struct {
	DatabaseName    string `mapstructure:"databaseName"`
	TablePrefix     string `mapstructure:"tablePrefix"`
	ColumnPrefix    string `mapstructure:"columnPrefix"`
	DeletedAtColumn string `mapstructure:"deletedAtColumn"`
}

// map for converting mysql type to golang types
var typeForMysqlToGo = map[string]string{
	"int":                "int",
	"integer":            "int",
	"tinyint":            "int",
	"smallint":           "int",
	"mediumint":          "int",
	"bigint":             "int",
	"int unsigned":       "int",
	"integer unsigned":   "int",
	"tinyint unsigned":   "int",
	"smallint unsigned":  "int",
	"mediumint unsigned": "int",
	"bigint unsigned":    "int",
	"bit":                "int",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "time.Time", // time.Time or string
	"datetime":           "time.Time", // time.Time or string
	"timestamp":          "time.Time", // time.Time or string
	"time":               "time.Time", // time.Time or string
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
}

func isTimeMysqlType(mysqlType string) bool {
	timeMap := map[string]string{
		"date":      "date",
		"datetime":  "datetime",
		"timestamp": "timestamp",
		"time":      "time",
	}
	_, ok := timeMap[mysqlType]
	return ok
}
func mysql2GoType(mysqlType string, time2str bool) (goType string, err error) {
	if time2str {
		typeForMysqlToGo["date"] = "string"
		typeForMysqlToGo["datetime"] = "string"
		typeForMysqlToGo["timestamp"] = "string"
		typeForMysqlToGo["time"] = "string"
	}
	index := strings.Index(mysqlType, "(")
	if index > -1 {
		mysqlType = mysqlType[:index]
	}
	goType, ok := typeForMysqlToGo[mysqlType]
	if !ok {
		err = errors.Errorf("mysql2GoType: not found mysql type %s to go type", mysqlType)
	}
	return

}

const (
	DEFAULT_VALUE_CURRENT_TIMESTAMP = "current_timestamp"
	DEFAULT_COLUMN_DELETED_AT       = "deleted_at" // 默认删除列名称
)

type Column struct {
	Prefix        string
	CamelName     string
	ColumnName    string // 数据库名称
	Name          string
	Type          string
	Comment       string
	Tag           string
	Nullable      bool
	Enums         []string
	AutoIncrement bool
	DefaultValue  string
	OnCreate      bool // 根据数据表ddl及默认 值为current_timestap 判断
	OnUpdate      bool // 根据数据表ddl 配置
	OnDelete      bool // 手动设置
}

// IsDefaultValueCurrentTimestamp 判断默认值是否为自动填充时间
func (c *Column) IsDefaultValueCurrentTimestamp() bool {
	return strings.Contains(strings.ToLower(c.DefaultValue), DEFAULT_VALUE_CURRENT_TIMESTAMP) // 测试发现有 current_timestamp() 情况
}

type _enum struct {
	ConstKey        string // 枚举类型定义 常量 名称
	ConstValue      string // 枚举类型定义值
	Title           string // 枚举类型 标题（中文）
	ColumnNameCamel string //枚举类型分组（字段名，每个枚举字段有多个值，全表通用，需要分组）
	Type            string // 类型 int-整型，string-字符串，默认string
}

type _Enums []*_enum

func (e _Enums) Len() int { // 重写 Len() 方法
	return len(e)
}
func (e _Enums) Swap(i, j int) { // 重写 Swap() 方法
	e[i], e[j] = e[j], e[i]
}
func (e _Enums) Less(i, j int) bool { // 重写 Less() 方法， 从小到大排序
	return e[i].ConstKey < e[j].ConstKey
}

// UniqueItems 去重
func (e _Enums) UniqueItems() (uniq _Enums) {
	emap := make(map[string]*_enum)
	for _, enum := range e {
		emap[enum.ConstKey] = enum
	}
	uniq = _Enums{}
	for _, enum := range emap {
		uniq = append(uniq, enum)
	}
	return
}

// ColumnNameCamels 获取所有分组
func (e _Enums) ColumnNameCamels() (output []string) {
	columnNameCamelMap := make(map[string]string)
	for _, enum := range e {
		columnNameCamelMap[enum.ColumnNameCamel] = enum.ColumnNameCamel
	}
	output = make([]string, 0)
	for _, columnNameCamel := range columnNameCamelMap {
		output = append(output, columnNameCamel)
	}
	return
}

// GetByGroup 通过分组名称获取enum
func (e _Enums) GetByColumnNameCamel(ColumnNameCamel string) (enums _Enums) {
	enums = _Enums{}
	for _, enum := range e {
		if enum.ColumnNameCamel == ColumnNameCamel {
			enums = append(enums, enum)
		}
	}
	return
}

type Table struct {
	DatabaseConfig DatabaseConfig
	TableName      string
	PrimaryKey     string
	DeleteColumn   string
	Columns        []*Column
	EnumsConst     _Enums
	Comment        string
	TableDef       *executor.TableDef
}

// CamelName 删除表前缀，转换成 camel 格式
func (t *Table) TableNameCamel() (camelName string) {
	name := t.TableNameTrimPrefix()
	camelName = pkg.ToCamel(name)
	return
}
func (t *Table) SnakeCase() (snakeName string) {
	name := t.TableNameTrimPrefix()
	snakeName = pkg.SnakeCase(name)
	return
}
func (t *Table) TableNameTrimPrefix() (name string) {
	name = t.TableName
	if t.DatabaseConfig.TablePrefix != "" {
		name = strings.TrimLeft(name, t.DatabaseConfig.TablePrefix)
	}
	return
}
func (t *Table) PrimaryKeyCamel() (camelName string) {
	primaryKey := t.PrimaryKey
	if t.DatabaseConfig.ColumnPrefix != "" {
		primaryKey = strings.TrimLeft(primaryKey, t.DatabaseConfig.TablePrefix)
	}
	camelName = pkg.ToCamel(primaryKey)
	return
}

func (t *Table) CreatedAtColumn() (createdAtColumn *Column) {
	for _, column := range t.Columns {
		if column.OnCreate {
			return column
		}
	}
	return
}

func (t *Table) GetColumnByCamelName(camelName string) (column *Column) {
	for _, column := range t.Columns {
		if column.CamelName == camelName {
			return column
		}
	}
	return
}

// UpdateAtColumn 获取更新列
func (t *Table) UpdatedAtColumn() (updatedAtColumn *Column) {
	for _, column := range t.Columns {
		if column.OnUpdate {
			return column
		}
	}

	return
}

// DeletedAtColumn 获取删除列
func (t *Table) DeletedAtColumn() (deletedAtColumn *Column) {
	for _, column := range t.Columns {
		if column.OnDelete {
			return column
		}
	}
	return
}

func enumsConst(tablePrefix string, columnPt *Column) (enumsConsts _Enums) {
	prefix := fmt.Sprintf("%s_%s", tablePrefix, columnPt.Name)
	enumsConsts = _Enums{}
	comment := strings.ReplaceAll(columnPt.Comment, " ", ",") // 替换中文逗号(兼容空格和逗号另种分割符号)
	reg, err := regexp.Compile(`\W`)
	if err != nil {
		panic(err)
	}
	for _, constValue := range columnPt.Enums {
		constKey := fmt.Sprintf("%s_%s", prefix, constValue)
		valueFormat := fmt.Sprintf("%s-", constValue) // 枚举类型的comment 格式 value1-title1,value2-title2
		index := strings.Index(comment, valueFormat)
		if index < 0 {
			err := errors.Errorf("column %s(enum) comment except contains %s-xxx,got:%s", columnPt.Name, constValue, comment)
			panic(err)
		}
		title := comment[index+len(valueFormat):]
		comIndex := strings.Index(title, ",")
		if comIndex > -1 {
			title = title[:comIndex]
		} else {
			title = strings.TrimRight(title, " )")
		}
		constKey = reg.ReplaceAllString(constKey, "_") //替换非字母字符
		constKey = strings.ToUpper(constKey)
		enumsConst := &_enum{
			ConstKey:        constKey,
			ConstValue:      constValue,
			Title:           title,
			ColumnNameCamel: columnPt.CamelName,
			Type:            "string",
		}
		enumsConsts = append(enumsConsts, enumsConst)
	}
	return
}
