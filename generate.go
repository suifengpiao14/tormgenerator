package generaterepository

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/suifengpiao14/generaterepository/converter"
	"github.com/suifengpiao14/generaterepository/pkg/ddlparser"
	"github.com/suifengpiao14/generaterepository/pkg/tpl2entity"
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

func (b *Builder) GenerateTorm(tormMetaMap TormMetaMap) (buf *bytes.Buffer, err error) {

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
				tormDTOs, err := converter.GenerateTorm(tormMetaTpl, subTables)
				if err != nil {
					return nil, err
				}
				for _, torm := range tormDTOs {
					w.WriteString(torm.TPL)
					w.WriteString(converter.EOF)
				}
			}
		}
	}
	return &w, nil
}

func (b *Builder) GenerateSQLEntity(tormText string) (buf *bytes.Buffer, err error) {
	torms, err := tpl2entity.ParseDefine(tormText)
	if err != nil {
		return nil, err
	}
	talbes, err := ddlparser.ParseDDL(b.ddl, b.dbConfig)
	if err != nil {
		return nil, err
	}
	entityDTO, err := converter.GenerateSQLEntity(torms, talbes)
	if err != nil {
		return nil, err
	}
	tplText := sqlEntityTemplate()
	r, err := template.New("").Parse(tplText)
	if err != nil {
		return nil, err
	}
	var w bytes.Buffer
	err = r.Execute(&w, entityDTO)
	if err != nil {
		return nil, err
	}
	return &w, nil

}

func modelTemplate() (tpl string) {
	tpl = `
	package repository
	import (
			"github.com/suifengpiao14/gotemplatefunc/templatefunc"
		)

		{{range $model:=. }}
		{{$model.TPL}}
		{{end}}
	`
	return
}

func sqlEntityTemplate() (tpl string) {
	tpl = `
	package repository
	import (
			"github.com/suifengpiao14/gotemplatefunc/templatefunc"
			"text/template"
			"bytes"
		)
		//GetTormTemplate 获取torm 模板 
		func GetTormTemplate()(tormTemplate *template.Template,err error){
			torm:=GetTorm()
			tormTemplate,err= template.New("").Funcs(templatefunc.TemplatefuncMapSQL).Parse(torm)
			if err != nil {
				return nil,err 
			}
			return tormTemplate,nil
		}

		//获取所有torm
		func GetTorm()(torm string){
			var w bytes.Buffer
			{{- range $entity:=. }}
			w.WriteString(new({{$entity.Name}}).Torm())
			{{- end}}
			torm=w.String()
			return torm
		}

		{{range $entity:=. }}
		{{$entity.TPL}}
		{{end}}
	`
	return
}
