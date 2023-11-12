package converter

import (
	"bytes"
	"sort"
	"text/template"

	"github.com/suifengpiao14/tormgenerator/parser/ddlparser"
	"github.com/suifengpiao14/tormgenerator/parser/tormparser"
	"github.com/suifengpiao14/tormgenerator/parser/tplparser"
)

type EntityDTO struct {
	Name string
	TPL  string
}

type EntityDTOs []*EntityDTO

func (a EntityDTOs) Len() int           { return len(a) }
func (a EntityDTOs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a EntityDTOs) Less(i, j int) bool { return a[i].Name < a[j].Name }

type TormStruct struct {
	StructName  string
	TplIdentity string
	Name        string
	Variables   []*tplparser.Variable
	FullName    string
	Type        string
	Torm        string
	OutEntity   *TormStruct // 输出对象
}

const STRUCT_DEFINE_NANE_FORMAT = "%sEntity"

// GenerateSQLEntity 根据数据表ddl和sql tpl 生成 sql tpl 调用的输入、输出实体
func GenerateSQLEntity(torms tplparser.TPLDefines, tables []*ddlparser.Table) (entityDTOs EntityDTOs, err error) {
	entityDTOs = make(EntityDTOs, 0)
	for _, sqltplDefine := range torms {
		entityElement, err := tormparser.ParserTorm(sqltplDefine, tables)
		if err != nil {
			return nil, err
		}
		tp, err := template.New("").Parse(inputEntityTemplate())
		if err != nil {
			return nil, err
		}
		buf := new(bytes.Buffer)
		err = tp.Execute(buf, entityElement)
		if err != nil {
			return nil, err
		}
		sqlEntity := buf.String()
		entityDTOs = append(entityDTOs, &EntityDTO{
			Name: entityElement.StructName,
			TPL:  sqlEntity,
		})
	}
	sort.Sort(entityDTOs)
	return entityDTOs, nil
}

func inputEntityTemplate() (tpl string) {
	tpl = `
		type {{.StructName}} struct{
			{{range .Variables -}}
				{{.FieldName}} {{.Type}} {{.Tag}} //{{.Comment}}
			{{end -}}
			tormfunc.VolumeMap
		}

		func (t *{{.StructName}}) TplName() string{
			return "{{.FullName}}"
		}

		func (t *{{.StructName}}) TplType() string{
			return "{{.Type}}"
		}

		func (t *{{.StructName}}) Torm() string{
			return {{.Torm}}
		}

		func (t *{{.StructName}}) GetTplIdentity() string{
			return "{{.TplIdentity}}"
		}
		func (t *{{.StructName}}) Exec(ctx context.Context, dst interface{}) (err error) {
			sqls, _, _, err := torm.GetSQL(t.GetTplIdentity(), t.TplName(), t)
			if err != nil {
				return err
			}
			tenantID, err := sqltenantplug.GetTenantID(ctx)
			if err != nil {
				return err
			}
			newsql, err := sqltenantplug.WithTenant(sqls, tenantID)
			if err != nil {
				return err
			}
			err = torm.ExecSQL(ctx, t.GetTplIdentity(), newsql, dst)
			if err != nil {
				return err
			}
			return nil
		}
		
	`
	return
}
