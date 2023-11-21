package converter

import (
	"github.com/suifengpiao14/tormgenerator/parser/ddlparser"
	"github.com/suifengpiao14/tormgenerator/parser/tormparser"
	"github.com/suifengpiao14/tormgenerator/parser/tplparser"
)

const STRUCT_DEFINE_NANE_FORMAT = "%sEntity"

// GenerateTormStructs 根据数据表ddl和sql tpl 生成 sql tpl 调用的输入、输出实体
func GenerateTormStructs(torms tplparser.TPLDefines, tables []*ddlparser.Table) (tormStructs tormparser.TormStructs, err error) {
	tormStructs = make([]tormparser.TormStruct, 0)
	for _, sqltplDefine := range torms {
		entityElement, err := tormparser.ParserTorm(sqltplDefine, tables)
		if err != nil {
			return nil, err
		}
		tormStructs = append(tormStructs, *entityElement)

	}
	return tormStructs, nil
}
