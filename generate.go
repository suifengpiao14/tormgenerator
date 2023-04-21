package generaterepository

import (
	"bytes"
	"fmt"
	"strings"

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

//MakeTormMetaWithAllTable 所有的数据表，共用相同的torm生成
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
	b.writePackageLine(&w)
	for _, model := range modelDTOs {
		w.WriteString(model.TPL)
		w.WriteString(converter.EOF)
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
			if strings.ToLower(table.TableNameCamel()) == strings.ToLower(tableName) {
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
	var w bytes.Buffer
	b.writePackageLine(&w)
	w.WriteString(`import "github.com/suifengpiao14/gotemplatefunc/templatefunc"`)
	w.WriteString(converter.EOF)
	for _, entity := range entityDTO {
		w.WriteString(entity.TPL)
		w.WriteString(converter.EOF)
	}
	return &w, nil

}

func (b *Builder) writePackageLine(w *bytes.Buffer) {
	w.WriteString(fmt.Sprintf("package %s", b.pacakgeName))
	w.WriteString(converter.EOF)
	w.WriteString(converter.EOF)
}
