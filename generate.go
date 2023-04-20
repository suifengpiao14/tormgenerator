package generaterepository

import (
	"bytes"
	"io"
	"text/template"

	"github.com/suifengpiao14/generaterepository/converter"
	"github.com/suifengpiao14/generaterepository/pkg/ddlparser"
	"github.com/suifengpiao14/generaterepository/pkg/tpl2entity"
)

func GenerateModel(ddl string, dbConfig ddlparser.DatabaseConfig) (buf *bytes.Buffer, err error) {
	talbes, err := ddlparser.ParseDDL(ddl)
	if err != nil {
		return nil, err
	}
	modelDTOs, err := converter.GenerateModel(talbes)
	if err != nil {
		return nil, err
	}
	var w bytes.Buffer
	for _, model := range modelDTOs {
		w.WriteString(model.TPL)
		w.WriteString(converter.EOF)
	}
	return &w, nil
}

func GenerateTorm(tormTpl *template.Template, ddl string, dbConfig ddlparser.DatabaseConfig) (reader io.Reader, err error) {
	talbes, err := ddlparser.ParseDDLWithConfig(ddl, dbConfig)
	if err != nil {
		return nil, err
	}
	tormDTOs, err := converter.GenerateTorm(tormTpl, talbes)
	if err != nil {
		return nil, err
	}
	var w bytes.Buffer
	for _, torm := range tormDTOs {
		w.WriteString(torm.TPL)
		w.WriteString(converter.EOF)
	}
	return &w, nil
}

func GenerateSQLEntity(tormText string, ddl string, dbConfig ddlparser.DatabaseConfig) (reader io.Reader, err error) {
	torms, err := tpl2entity.ParseDefine(tormText)
	if err != nil {
		return nil, err
	}
	talbes, err := ddlparser.ParseDDLWithConfig(ddl, dbConfig)
	if err != nil {
		return nil, err
	}
	entityDTO, err := converter.GenerateSQLEntity(torms, talbes)
	if err != nil {
		return nil, err
	}
	var w bytes.Buffer
	for _, entity := range entityDTO {
		w.WriteString(entity.TPL)
		w.WriteString(converter.EOF)
	}
	return &w, nil

}
