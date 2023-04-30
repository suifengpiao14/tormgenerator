package converter

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/generaterepository/pkg"
	"github.com/suifengpiao14/generaterepository/pkg/ddlparser"
	"github.com/suifengpiao14/generaterepository/pkg/tpl2entity"
)

const (
	EOF                  = "\n"
	WINDOW_EOF           = "\r\n"
	HTTP_HEAD_BODY_DELIM = EOF + EOF
)

const (
	TPL_DEFINE_TYPE_CURL_REQUEST  = "curl_request"
	TPL_DEFINE_TYPE_CURL_RESPONSE = "curl_response"
	TPL_DEFINE_TYPE_SQL_SELECT    = "sql_select"
	TPL_DEFINE_TYPE_SQL_UPDATE    = "sql_update"
	TPL_DEFINE_TYPE_SQL_INSERT    = "sql_insert"
	TPL_DEFINE_TYPE_TEXT          = "text"
	CHARACTERISTIC_CURL           = "HTTP/1.1"
	CHARACTERISTIC_SQL_SELECT     = "SELECT"
	CHARACTERISTIC_SQL_UPDATE     = "UPDATE"
	CHARACTERISTIC_SQL_INSERT     = "INSERT"
)

var (
	LeftDelim  = "{{"
	RightDelim = "}}"
)

type EntityDTO struct {
	Name string
	TPL  string
}

type EntityDTOs []*EntityDTO

func (a EntityDTOs) Len() int           { return len(a) }
func (a EntityDTOs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a EntityDTOs) Less(i, j int) bool { return a[i].Name < a[j].Name }

type _EntityElement struct {
	StructName string
	Name       string
	Variables  []*tpl2entity.Variable
	FullName   string
	Type       string
	Torm       string
	OutEntity  *_EntityElement // 输出对象
}

const STRUCT_DEFINE_NANE_FORMAT = "%sEntity"

// GenerateSQLEntity 根据数据表ddl和sql tpl 生成 sql tpl 调用的输入、输出实体
func GenerateSQLEntity(torms tpl2entity.TPLDefines, tables []*ddlparser.Table) (entityDTOs EntityDTOs, err error) {
	entityDTOs = make(EntityDTOs, 0)
	for _, sqltplDefine := range torms {
		entityElement, err := sqlEntityElement(sqltplDefine, tables)
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

func sqlEntityElement(sqltplDefineText *tpl2entity.TPLDefine, tableList []*ddlparser.Table) (entityElement *_EntityElement, err error) {
	variableList := sqltplDefineText.GetVairables()
	variableList, err = formatVariableTypeByTableColumn(variableList, tableList)
	if err != nil {
		return nil, err
	}
	columnArr := parseSQLSelectColumn(sqltplDefineText.Text)
	allColumnVariable := ColumnsToVariables(tableList)
	columnVariables := make(tpl2entity.Variables, 0)
	for _, columnVariable := range allColumnVariable {
		for _, columnName := range columnArr {
			if columnName == columnVariable.Name {
				columnVariables = append(columnVariables, columnVariable)
			}
		}
	}
	camelName := sqltplDefineText.NameCamel()
	outName := fmt.Sprintf("%sOut", camelName)
	entityElement = &_EntityElement{
		Name:       camelName,
		Type:       sqltplDefineText.Type(),
		Torm:       strconv.Quote(sqltplDefineText.Text),
		StructName: fmt.Sprintf(STRUCT_DEFINE_NANE_FORMAT, camelName),
		Variables:  variableList,
		FullName:   sqltplDefineText.Name,
		OutEntity: &_EntityElement{
			Name:       outName,
			Type:       sqltplDefineText.Type(),
			StructName: fmt.Sprintf(STRUCT_DEFINE_NANE_FORMAT, camelName),
			Variables:  columnVariables,
			FullName:   fmt.Sprintf("%sOut", sqltplDefineText.Name),
		},
	}
	return entityElement, nil
}

func parseSQLSelectColumn(sql string) []string {
	grep := regexp.MustCompile(`(?i)select(.+)from`)
	match := grep.FindAllStringSubmatch(sql, -1)
	if len(match) < 1 {
		return make([]string, 0)
	}
	fieldStr := match[0][1]
	out := strings.Split(pkg.StandardizeSpaces(fieldStr), ",")
	return out
}

func ColumnsToVariables(tableList []*ddlparser.Table) (variables tpl2entity.Variables) {
	allVariables := make(tpl2entity.Variables, 0)
	tableColumnMap := make(map[string]*ddlparser.Column)
	columnNameMap := make(map[string]string)
	for _, table := range tableList {
		for _, column := range table.Columns { //todo 多表字段相同，类型不同时，只会取列表中最后一个
			tableColumnMap[column.CamelName] = column
			lname := strings.ToLower(column.CamelName)
			columnNameMap[lname] = column.CamelName // 后续用于检测模板变量拼写错误
			variable := &tpl2entity.Variable{
				Name:      column.Name,
				FieldName: column.Name,
				Type:      column.Type,
			}
			allVariables = append(allVariables, variable)
		}
	}
	sort.Sort(allVariables)
	return allVariables
}

func ParseSQLTPLTableName(sqlTpl string) (tableList []string, err error) {

	updateDelim := "update `?(\\w+)`?"
	updateMatchArr, err := regexpMatch(sqlTpl, updateDelim)
	if err != nil {
		return
	}
	selectDelim := "from `?(\\w+)`?"
	fromMatchArr, err := regexpMatch(sqlTpl, selectDelim)
	if err != nil {
		return
	}
	insertDelim := "into `?(\\w+)`?"
	insertMatchArr, err := regexpMatch(sqlTpl, insertDelim)
	if err != nil {
		return
	}

	tableList = append(tableList, updateMatchArr...)
	tableList = append(tableList, fromMatchArr...)
	tableList = append(tableList, insertMatchArr...)
	return
}

func regexpMatch(s string, delim string) (matcheList []string, err error) {
	reg := regexp.MustCompile(delim)
	if reg == nil {
		err = errors.Errorf("regexp.MustCompile %s is nil", delim)
		return
	}
	matchArr := reg.FindAllStringSubmatch(s, -1)
	for _, matchs := range matchArr {
		matcheList = append(matcheList, matchs[1:]...) // index 0 为匹配对象
	}
	return
}

func formatVariableTypeByTableColumn(variableList tpl2entity.Variables, tableList []*ddlparser.Table) (varaibles tpl2entity.Variables, err error) {
	tableColumnMap := make(map[string]*ddlparser.Column)
	checkSpellingMistakes := make(map[string]string)
	for _, table := range tableList {
		for _, column := range table.Columns { //todo 多表字段相同，类型不同时，只会取列表中最后一个
			tableColumnMap[column.CamelName] = column
			lname := strings.ToLower(column.CamelName)
			checkSpellingMistakes[lname] = column.CamelName // 后续用于检测模板变量拼写错误
		}
	}
	varaibles = variableList
	for i, varia := range varaibles {
		column, ok := tableColumnMap[varia.Name]
		if !ok {
			continue
		}
		varaibles[i].Type = column.Type
		varaibles[i].Comment = column.Comment
	}

	for _, variable := range variableList { // 使用数据库字段定义校正变量类型
		name := variable.Name
		_, ok := tableColumnMap[name]
		if !ok {
			lname := strings.ToLower(name)
			if columnName, ok := checkSpellingMistakes[lname]; name != lname && ok { // 检测模板中大小写拼写错误
				err = errors.Errorf("spelling mistake: have %s, want %s", name, columnName)
				return
			}
		}
	}

	sort.Sort(varaibles)
	return
}

func inputEntityTemplate() (tpl string) {
	tpl = `
		type {{.StructName}} struct{
			{{range .Variables }}
				{{.FieldName}} {{.Type}} {{.Tag}} //{{.Comment}}
			{{end}}
			templatefunc.VolumeMap
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
	`
	return
}
