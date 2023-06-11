package converter

import (
	"bytes"
	"text/template"

	"github.com/suifengpiao14/torm/parser/ddlparser"
)

type ModelDTO struct {
	Name string
	TPL  string
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
	tableTpl := modelTemplate()
	tl, err := template.New("").Parse(tableTpl)
	if err != nil {
		return
	}

	for i := 0; i < len(tables); i++ {
		buf := new(bytes.Buffer)
		table := tables[i]
		err = tl.Execute(buf, table)
		if err != nil {
			return
		}
		modelStruct := &ModelDTO{
			Name: table.TableNameCamel(),
			TPL:  buf.String(),
		}
		modelDTOs = append(modelDTOs, modelStruct)
	}
	return
}

func modelTemplate() string {
	return `
	{{- $enumsConst :=.EnumsConst }}
	{{if $enumsConst }}
	const (
		{{range $enumsConst -}}
			{{.ConstKey}}="{{.ConstValue}}"
		{{end}}
		)
	{{end}}
	{{$modelName:= print .TableNameCamel "Model"}}
	type {{$modelName}} struct{
		{{range .Columns -}} 
		{{.CamelName}} {{.Type}} {{if .Tag}} {{.Tag}} {{end}} // {{.Comment}}
		{{end}}
	}
	func (t *{{$modelName}}) TableName()string{
		return "{{.TableName}}"
	}
	func (t *{{$modelName}}) PrimaryKey()string{
		return "{{.PrimaryKey}}"
	}
	func (t *{{$modelName}}) PrimaryKeyCamel()string{
		return "{{.PrimaryKeyCamel}}"
	}
	{{- if $enumsConst}}
		{{- range  $camelName :=$enumsConst.ColumnNameCamels }}
				{{- /* 生成 enumMap 函数 */ -}}
				{{- $enumMapFuncName:=print $camelName "TitleMap" }}
				func (t *{{$modelName}}) {{$enumMapFuncName}}()map[string]string{
					enumMap:=make(map[string]string)
					{{- range $enum:= $enumsConst.GetByColumnNameCamel $camelName }}
						enumMap[{{$enum.ConstKey}}]="{{$enum.Title}}"
					{{- end }}
					return enumMap
				}

				{{- /* 生成 枚举值获取标题 函数*/ -}}
				{{$titleFuncName:=print  $camelName "Title"}}
				func (t *{{$modelName}}) {{$titleFuncName}}()string{
					enumMap:=t.{{$enumMapFuncName}}()
					title,ok:=enumMap[t.{{$camelName}}]
					if !ok{
						msg:="func {{$titleFuncName}} not found title by key "+t.{{$camelName}}
						panic(msg)
					}
					return title
				}

		{{- end }}
	{{- end}}
	`
}
