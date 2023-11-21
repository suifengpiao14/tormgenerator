package tormgenerator

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/suifengpiao14/funcs"
	"github.com/suifengpiao14/tormgenerator/converter"
	"github.com/suifengpiao14/tormgenerator/parser/ddlparser"
	"github.com/suifengpiao14/tormgenerator/parser/tormparser"
	"github.com/suifengpiao14/tormgenerator/parser/tplparser"
)

type TormMetaMap map[string]string
type Builder struct {
	pacakgeName string
	ddl         string
	dbConfig    ddlparser.DatabaseConfig
}

func NewBuilder(pacakgeName string, ddl string, dbConfig ddlparser.DatabaseConfig) (builder *Builder) {
	builder = &Builder{
		pacakgeName: pacakgeName,
		ddl:         ddl,
		dbConfig:    dbConfig,
	}
	return
}

func (b *Builder) GetTables() (tables []*ddlparser.Table, err error) {
	tables, err = ddlparser.ParseDDL(b.ddl, b.dbConfig)
	if err != nil {
		return nil, err
	}
	return tables, err
}

// MakeTormMetaWithAllTable 所有的数据表，共用相同的torm生成
func (b *Builder) MakeTormMetaWithAllTable(commonTormMeta string) (tormMetaMap *TormMetaMap, err error) {
	tables, err := ddlparser.ParseDDL(b.ddl, b.dbConfig)
	if err != nil {
		return nil, err
	}
	tormMetaMap = &TormMetaMap{}
	for _, table := range tables {
		(*tormMetaMap)[table.TableNameCamel()] = commonTormMeta
	}
	return tormMetaMap, err
}

func (b *Builder) GenerateModel() (buf *bytes.Buffer, err error) {
	talbes, err := ddlparser.ParseDDL(b.ddl, b.dbConfig)
	if err != nil {
		return nil, err
	}
	for _, table := range talbes {
		for _, column := range table.Columns { // 增加json tag
			column.Tag = fmt.Sprintf("`gorm:\"column:%s\" json:\"%s\"`", column.ColumnName, funcs.ToLowerCamel(column.CamelName))
		}
	}

	modelDTOs, err := converter.GenerateModel(talbes)
	if err != nil {
		return nil, err
	}

	var w bytes.Buffer
	r, err := template.New("").Parse(modelTemplate())
	if err != nil {
		return nil, err
	}
	err = r.Execute(&w, modelDTOs)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (b *Builder) GenerateTormFromMeta(tormMetaMap TormMetaMap) (buf *bytes.Buffer, err error) {

	tables, err := ddlparser.ParseDDL(b.ddl, b.dbConfig)
	if err != nil {
		return nil, err
	}
	var w bytes.Buffer
	for tableName, tormMetaTpl := range tormMetaMap {
		for _, table := range tables {
			if strings.EqualFold(table.TableNameCamel(), strings.ToLower(tableName)) {
				subTables := []*ddlparser.Table{
					table,
				}
				tormDTOs, err := converter.GenerateTormFromDDL(tormMetaTpl, subTables)
				if err != nil {
					return nil, err
				}
				for _, torm := range tormDTOs {
					w.WriteString(torm.TPL)
					w.WriteString(tormparser.EOF)
				}
			}
		}
	}
	return &w, nil
}

func (b *Builder) GenerateSQLTorm(tormText string) (tormStructs tormparser.TormStructs, err error) {
	torms, err := tplparser.ParseDefine(tormText)
	if err != nil {
		return nil, err
	}
	talbes, err := ddlparser.ParseDDL(b.ddl, b.dbConfig)
	if err != nil {
		return nil, err
	}
	tormStructs, err = converter.GenerateTormStructs(torms, talbes)
	if err != nil {
		return nil, err
	}
	sort.Sort(tormStructs)
	return tormStructs, nil
}

func modelTemplate() (tpl string) {
	tpl = `
	package repository
	import (
			"github.com/suifengpiao14/torm/tormfunc"
		)

		{{range $model:=. }}
		{{$model.TPL}}
		{{end}}
	`
	return
}
