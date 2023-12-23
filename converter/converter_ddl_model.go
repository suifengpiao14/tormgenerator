package converter

import (
	"github.com/suifengpiao14/tormgenerator/parser/ddlparser"
)

type ModelDTO struct {
	Name  string
	Table ddlparser.Table
}

type ModelDTOs []*ModelDTO

func (v ModelDTOs) Len() int { // 重写 Len() 方法
	return len(v)
}
func (v ModelDTOs) Swap(i, j int) { // 重写 Swap() 方法
	v[i], v[j] = v[j], v[i]
}
func (v ModelDTOs) Less(i, j int) bool { // 重写 Less() 方法， 从小到大排序
	return v[i].Name < v[j].Name
}

// GenerateModel 生成 model 文件内容
func GenerateModel(tables []*ddlparser.Table) (modelDTOs ModelDTOs, err error) {
	modelDTOs = make([]*ModelDTO, 0)

	for i := 0; i < len(tables); i++ {
		table := tables[i]
		columns := make([]*ddlparser.Column, 0)
		for _, column := range table.Columns {
			if !table.DatabaseConfig.ExtraConfigs.IsIgnore(ddlparser.ExtraConfig_Domain_model, table.TableName, column.ColumnName) {
				columns = append(columns, column)
			}

		}
		table.Columns = columns
		modelStruct := &ModelDTO{
			Name:  table.TableNameCamel(),
			Table: *table,
		}
		modelDTOs = append(modelDTOs, modelStruct)
	}
	return
}
