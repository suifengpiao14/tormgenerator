package converter

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"
	"github.com/pkg/errors"
	"github.com/suifengpiao14/jsonschemaline"
	"github.com/suifengpiao14/tormrepository/pkg"
	"github.com/suifengpiao14/tormrepository/pkg/ddlparser"
	"github.com/suifengpiao14/tormrepository/pkg/tpl2entity"
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

type EntityElement struct {
	StructName string
	Name       string
	Variables  []*Variable
	FullName   string
	Type       string
	OutEntity  *EntityElement // 输出对象
}

func GetSamePrefixEntityElements(prefix string, entityElementList []*EntityElement) (samePrefixEntityElementList []*EntityElement) {
	samePrefixEntityElementList = make([]*EntityElement, 0)
	for _, entityElement := range entityElementList {
		if strings.HasPrefix(entityElement.Name, prefix) {
			samePrefixEntityElementList = append(samePrefixEntityElementList, entityElement)
		}
	}
	return samePrefixEntityElementList
}

const STRUCT_DEFINE_NANE_FORMAT = "%sEntity"

// SQLEntity 根据数据表ddl和sql tpl 生成 sql tpl 调用的输入、输出实体
func SQLEntity(sqltplDefineText *TPLDefineText, tableList []*Table) (entityStruct string, err error) {
	entityElement, err := SQLEntityElement(sqltplDefineText, tableList)
	if err != nil {
		return "", err
	}
	tp, err := template.New("").Parse(InputEntityTemplate())
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	err = tp.Execute(buf, entityElement)
	if err != nil {
		return
	}
	entityStruct = buf.String()
	return
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

func SQLEntityElement(sqltplDefineText *tpl2entity.TPLDefineText, tableList []*ddlparser.Table) (entityElement *EntityElement, err error) {
	variableList := sqltplDefineText.GetVairables()
	variableList, err = FormatVariableTypeByTableColumn(variableList, tableList)
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
	camelName := sqltplDefineText.FullnameCamel()
	outName := fmt.Sprintf("%sOut", camelName)
	entityElement = &EntityElement{
		Name:       camelName,
		Type:       sqltplDefineText.Type(),
		StructName: fmt.Sprintf(STRUCT_DEFINE_NANE_FORMAT, camelName),
		Variables:  variableList,
		FullName:   sqltplDefineText.Fullname(),
		OutEntity: &EntityElement{
			Name:       outName,
			Type:       sqltplDefineText.Type(),
			StructName: fmt.Sprintf(STRUCT_DEFINE_NANE_FORMAT, camelName),
			Variables:  columnVariables,
			FullName:   fmt.Sprintf("%s%sOut", sqltplDefineText.Namespace, sqltplDefineText.Name),
		},
	}
	return entityElement, nil
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
	sort.Sort(variables)
	return
}

func SqlTplDefineVariable2lineschema(id string, variables []*Variable, direction string) (lineschema string, err error) {
	arr := make([]string, 0)
	if direction == jsonschemaline.LINE_SCHEMA_DIRECTION_IN {
		arr = append(arr, fmt.Sprintf("version=http://json-schema.org/draft-07/schema,id=input,direction=%s", direction))
	} else if direction == jsonschemaline.LINE_SCHEMA_DIRECTION_OUT {
		arr = append(arr, fmt.Sprintf("version=http://json-schema.org/draft-07/schema,id=output,direction=%s", direction))
	}
	for _, v := range variables {
		if v.FieldName == "" { // 过滤匿名字段
			continue
		}
		kvArr := make([]string, 0)

		kvArr = append(kvArr, fmt.Sprintf("fullname=%s", v.FieldName))
		dst := ""
		src := ""
		format := v.Validate.Format
		if direction == jsonschemaline.LINE_SCHEMA_DIRECTION_IN {
			dst = v.Name //此处使用驼峰,v.FieldName 被改成蛇型了
		} else if direction == jsonschemaline.LINE_SCHEMA_DIRECTION_OUT {
			src = v.Validate.DataPathSrc
		}

		if dst != "" {
			kvArr = append(kvArr, fmt.Sprintf("dst=%s", dst))
		}
		if src != "" {
			kvArr = append(kvArr, fmt.Sprintf("src=%s", src))
		}
		if format != "" {
			kvArr = append(kvArr, fmt.Sprintf("format=%s", format))
		}
		kvArr = append(kvArr, "required")

		line := strings.Join(kvArr, ",")
		arr = append(arr, line)
	}
	lineschema = strings.Join(arr, "\n")
	return lineschema, err
}

func SqlTplDefineVariable2Jsonschema(id string, variables []*Variable) (jsonschemaOut string, err error) {
	properties := orderedmap.New()
	//{"$schema":"http://json-schema.org/draft-07/schema#","type":"object","properties":{},"required":[]}
	schema := jsonschema.Schema{
		Version:    "http://json-schema.org/draft-07/schema#",
		Type:       "object",
		ID:         jsonschema.ID(id),
		Properties: properties,
	}
	names := make([]string, 0)
	for _, v := range variables {
		if v.FieldName == "" { // 过滤匿名字段
			continue
		}

		name := v.FieldName
		subSchema := v.Validate
		subSchema.TypeValue = v.Type
		properties.Set(name, subSchema)
		names = append(names, name)
	}
	schema.Required = names
	b, err := schema.MarshalJSON()
	jsonschemaOut = string(b)
	return jsonschemaOut, err
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

type SQLTplNamespace struct {
	Namespace string
	Table     *ddlparser.Table
	Defines   TPLDefineTextList
}

func (s *SQLTplNamespace) String() string { // 这个将第一次模板解析输出的内容，合并成字符串，然后解析出{{define "xxx"}}{{end}}模板
	tplArr := make([]string, 0)
	for _, define := range s.Defines {
		tplArr = append(tplArr, define.Text)
	}
	str := strings.Join(tplArr, EOF)
	tplDefineList := ManualParseDefine(str, "", LeftDelim, RightDelim)
	tplDefineList = tplDefineList.UniqueItems() // 去重
	newTplArr := make([]string, 0)
	for _, tplDefineText := range tplDefineList {
		newTplArr = append(newTplArr, tplDefineText.Text)
	}
	out := strings.Join(newTplArr, EOF)
	return out
}

func (s *SQLTplNamespace) Filename() (out string) {
	out = pkg.SnakeCase(s.Namespace)
	return
}

func ManualParseDefine(content string, namespace string, leftDelim string, rightDelim string) (tplDefineList TPLDefineTextList) {
	// 解析文本
	delim := leftDelim + "define "
	delimLen := len(delim)
	content = pkg.TrimSpaces(content) // 去除开头结尾的非有效字符
	defineList := make([]string, 0)
	for {
		index := strings.Index(content, delim)
		if index >= 0 {
			pos := delimLen + index
			nextIndex := strings.Index(content[pos:], delim)
			if nextIndex >= 0 {
				sepPos := pos + nextIndex
				oneDefine := content[:sepPos]
				defineList = append(defineList, oneDefine)
				content = content[sepPos:]
			} else {
				defineList = append(defineList, content)
				break
			}
		} else {
			break
		}
	}

	tplDefineList = TPLDefineTextList{}

	// 格式化
	for _, tpl := range defineList {
		name, err := GetDefineName(tpl)
		if err != nil {
			panic(err)
		}

		tplDefineText := &TPLDefineText{
			Name:      name,
			Namespace: namespace,
			Text:      tpl,
		}
		tplDefineList = append(tplDefineList, tplDefineText)
	}

	return
}

type EntityTplData struct {
	StructName                  string
	FullName                    string
	ImplementTplEntityInterface bool
	Attributes                  tpl2entity.Variables
}

func FormatVariableTypeByTableColumn(variableList tpl2entity.Variables, tableList []*ddlparser.Table) (varaibles tpl2entity.Variables, err error) {
	varaibles = make(tpl2entity.Variables, 0)
	tableColumnMap := make(map[string]*ddlparser.Column)
	columnTypMap := make(map[string]string)
	columnNameMap := make(map[string]string)
	for _, table := range tableList {
		for _, column := range table.Columns { //todo 多表字段相同，类型不同时，只会取列表中最后一个
			tableColumnMap[column.CamelName] = column
			columnTypMap[column.CamelName] = column.Type
			lname := strings.ToLower(column.CamelName)
			columnNameMap[lname] = column.CamelName // 后续用于检测模板变量拼写错误
		}
	}
	variableList.FormatVariableType(columnTypMap)
	varaibles = variableList
	for _, variable := range variableList { // 使用数据库字段定义校正变量类型
		name := variable.Name
		lname := strings.ToLower(name)
		if columnName, ok := columnNameMap[lname]; ok { // 检测模板中大小写拼写错误
			err = errors.Errorf("spelling mistake: have %s, want %s", name, columnName)
			return
		}
	}
	sort.Sort(varaibles)
	return
}

func InputEntityTemplate() (tpl string) {
	tpl = `
		type {{.StructName}} struct{
			{{range .Variables }}
				{{.FieldName}} {{.Type}} {{.Tag}} 
			{{end}}
			gqt.TplEmptyEntity
		}

		func (t *{{.StructName}}) TplName() string{
			return "{{.FullName}}"
		}

		func (t *{{.StructName}}) TplType() string{
			return "{{.Type}}"
		}
	`
	return
}
