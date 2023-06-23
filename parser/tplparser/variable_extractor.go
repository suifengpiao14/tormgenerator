package tplparser

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/funcs"
	"github.com/suifengpiao14/helpers"
	"github.com/suifengpiao14/jsonschemaline"
)

const STRUCT_DEFINE_NANE_FORMAT = "%sEntity"

type Variable struct {
	Name       string // 模板中的原始名称
	FieldName  string // 当变量作为结构体的字段时，字段名称
	Comment    string
	Type       string
	Tag        string
	AllowEmpty bool
}

func (v *Variable) NameCamel() (nameCamel string) {
	nameCamel = helpers.ToCamel(v.Name)
	return
}

type Variables []*Variable

func (v Variables) Len() int { // 重写 Len() 方法
	return len(v)
}
func (v Variables) Swap(i, j int) { // 重写 Swap() 方法
	v[i], v[j] = v[j], v[i]
}
func (v Variables) Less(i, j int) bool { // 重写 Less() 方法， 从小到大排序
	return v[i].Name < v[j].Name
}

// UniqueItems 去重
func (v Variables) UniqueItems() (uniq []*Variable) {
	vmap := make(map[string]*Variable)
	for _, variable := range v {
		vmap[variable.Name] = variable
	}
	uniq = make([]*Variable, 0)
	for _, variable := range vmap {
		uniq = append(uniq, variable)
	}
	return
}

func (variables Variables) Lineschema(id string, direction string) (lineschema string, err error) {
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

		kvArr = append(kvArr, fmt.Sprintf("fullname=%s", funcs.ToLowerCamel(v.FieldName)))
		dst := ""
		src := ""
		format := v.Type
		if direction == jsonschemaline.LINE_SCHEMA_DIRECTION_IN {
			dst = v.Name //此处使用驼峰,v.FieldName 被改成蛇型了
		} else if direction == jsonschemaline.LINE_SCHEMA_DIRECTION_OUT {
			src = v.Name
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
		if v.AllowEmpty {
			kvArr = append(kvArr, "allowEmptyValue")
		}
		if v.Comment != "" {
			kvArr = append(kvArr, fmt.Sprintf("description=%s", v.Comment))
		}

		line := strings.Join(kvArr, ",")
		arr = append(arr, line)
	}
	lineschema = strings.Join(arr, "\n")
	return lineschema, err
}

func parseTplVariable(tplContext []byte) (variableList Variables) {
	variableList = make([]*Variable, 0)
	byteArr := tplContext
	// template 模板变量提取
	leftDelim := []byte(LeftDelim)
	rightDelim := []byte(RightDelim)
	itemBegin := false
	itemArr := make([][]byte, 0)
	item := make([]byte, 0)
	byteLen := len(byteArr)
	for i := 0; i < byteLen; i++ {
		c := byteArr[i]
		if c == leftDelim[0] && i+1 < byteLen && byteArr[i+1] == leftDelim[1] && !itemBegin {
			itemBegin = true
			item = make([]byte, 0)
			i++
			continue
		}
		if c == rightDelim[0] && i+1 < byteLen && byteArr[i+1] == rightDelim[1] && itemBegin {
			itemBegin = false
			itemArr = append(itemArr, item)
			i++
			continue
		}
		if itemBegin {
			item = append(item, c)
		}
	}

	// parse define variable
	for _, item := range itemArr {
		variable, _ := parsePrefixVariable(item, byte('.'))
		if variable.Name != "" {
			variable.FieldName = variable.Name

			variableList = append(variableList, &variable)

		}
	}

	// parse sub define variable
	templateNameList := getTemplateNames(string(tplContext))
	for _, templateName := range templateNameList {
		templateName = helpers.ToCamel(templateName)
		variable := Variable{
			Name:       templateName,
			AllowEmpty: false,
		}
		variable.Type = fmt.Sprintf(STRUCT_DEFINE_NANE_FORMAT, variable.NameCamel())
		variableList = append(variableList, &variable)
	}

	variableList = variableList.UniqueItems()
	return
}

func parseCurlTplVariable(tplContext []byte, typ string) (variableList Variables) {
	content := ToEOF(string(tplContext)) // 转换换行符)

	tplVariableList := parseTplVariable([]byte(content))
	variableList = append(variableList, tplVariableList...)

	if typ == TPL_DEFINE_TYPE_CURL_RESPONSE { // parse curl response variable ,curl response 直接采用复制，所以确保 response body 本身符合go语法
		index := strings.Index(content, HTTP_HEAD_BODY_DELIM)
		if index < 0 {
			return
		}
		body := content[index+len(HTTP_HEAD_BODY_DELIM):]
		body = strings.ReplaceAll(body, LeftDelim, "")
		body = strings.ReplaceAll(body, RightDelim, "")
		variable := &Variable{
			Name:      "",
			FieldName: "",
			Type:      body,
		}
		variableList = append(variableList, variable)
		return
	}
	return
}

func parsSqlTplVariable(tplContext []byte) (variableList Variables) {
	subVariableList := parseTplVariable(tplContext)
	variableList = append(variableList, subVariableList...)
	byteArr := tplContext

	// parse sql variable
	for {
		variable, pos := parsePrefixVariable(byteArr, SQL_VARIABLE_DELIM)
		if variable.Name == "" {
			break
		}
		variable.FieldName = variable.Name
		variableList = append(variableList, &variable)
		pos += len(variable.Name)
		byteArr = byteArr[pos:]
	}
	limitVariabeList := getLimitVariable(string(tplContext))
	variableList = append(variableList, limitVariabeList...)
	variableList = variableList.UniqueItems()
	return
}

func parseSQLSelectColumn(sql string) []string {
	grep := regexp.MustCompile(`(?i)select(.+)from`)
	match := grep.FindAllStringSubmatch(sql, -1)
	if len(match) < 1 {
		return make([]string, 0)
	}
	fieldStr := match[0][1]
	out := strings.Split(helpers.StandardizeSpaces(fieldStr), ",")
	return out
}

func getLimitVariable(sqlTplContext string) (variableList []*Variable) {
	variableList = make([]*Variable, 0)
	index := strings.Index(strings.ToLower(sqlTplContext), "limit")
	if index < 0 {
		return
	}
	byteArr := []byte(sqlTplContext[index:])
	// parse sql variable
	sqlVariableDelim := byte(':')

	for {
		variable, pos := parsePrefixVariable(byteArr, sqlVariableDelim)
		if variable.Name == "" {
			break
		}
		variable.FieldName = variable.Name
		variable.Type = "int"
		variable.AllowEmpty = false
		variableList = append(variableList, &variable)
		pos += len(variable.Name)
		byteArr = byteArr[pos:]
	}

	return
}

// 找到第一个变量
func parsePrefixVariable(item []byte, variableStart byte) (variable Variable, pos int) {
	variableBegin := false
	pos = 0
	variableNameByte := make([]byte, 0)
	itemLen := len(item)
	for j := 0; j < itemLen; j++ {
		c := item[j]
		if c == variableStart && j+1 < itemLen && IsNameChar(item[j+1]) { // c=变量标示字符，并且后面跟字符，是变量名的必要条件
			if j == 0 || !IsNameChar(item[j-1]) { // 变量符号开始或者前面不为字符，为变量名的充要条件
				variableBegin = true
				pos = j
				continue
			}
		}
		if variableBegin {
			if IsNameChar(c) {
				variableNameByte = append(variableNameByte, c)
			} else {
				break
			}
		}
	}
	variableName := string(variableNameByte)
	typ := variableSuffix2Type(variableName)
	if typ == "" {
		typ = "string"
	}
	variable = Variable{
		Name:       variableName,
		Type:       typ,
		AllowEmpty: true,
	}
	return
}

func getDefineName(tplDefineText string) (defineName string, err error) {
	delim := []byte("{{define \"")
	tplDefineByte := []byte(tplDefineText)
	index := bytes.Index(tplDefineByte, delim)
	nameByte := make([]byte, 0)
	if index >= 0 {
		index += len(delim)
		for i := index; i < len(tplDefineByte); i++ {
			c := tplDefineByte[i]
			if c != '"' {
				nameByte = append(nameByte, tplDefineByte[i])
			} else {
				break
			}

		}
	}
	defineName = string(nameByte)
	if defineName == "" {
		err = errors.Errorf("define name is empty")
	}
	return
}

func getTemplateNames(sqlTpl string) (templateNameList []string) {
	templateNameList = make([]string, 0)
	reg := regexp.MustCompile(`\{\{template\W+"(\w+)"`)
	if reg == nil {
		panic("regexp err")
	}
	matches := reg.FindAllStringSubmatch(sqlTpl, -1)
	for _, match := range matches {
		templateNameList = append(templateNameList, match[1])
	}
	return
}

// 判断是否可以作为名称的字符
func IsNameChar(c byte) (yes bool) {
	yes = false
	a := byte('a')
	z := byte('z')
	A := byte('A')
	Z := byte('Z')
	underline := byte('_')
	if (a <= c && c <= z) || (A <= c && c <= Z) || c == underline {
		yes = true
	}
	return
}

func ToEOF(s string) string {
	out := strings.ReplaceAll(s, WINDOW_EOF, EOF) // 统一换行符
	return out
}

type VariableSuffixType struct {
	Suffix string
	Type   string
}

type VariableSuffixTypes []*VariableSuffixType

func (v VariableSuffixTypes) Len() int { // 重写 Len() 方法
	return len(v)
}
func (v VariableSuffixTypes) Swap(i, j int) { // 重写 Swap() 方法
	v[i], v[j] = v[j], v[i]
}
func (v VariableSuffixTypes) Less(i, j int) bool { // 重写 Less() 方法， 从长到短排序
	return len(v[i].Suffix) > len(v[j].Suffix)
}

var VariableSuffixTypeList = VariableSuffixTypes{
	&VariableSuffixType{Suffix: "ListInt", Type: "[]int"},
	&VariableSuffixType{Suffix: "ListStr", Type: "[]string"},
	&VariableSuffixType{Suffix: "List", Type: "[]string"},
	&VariableSuffixType{Suffix: "Str", Type: "string"},
	&VariableSuffixType{Suffix: "Int", Type: "int"},
}

func variableSuffix2Type(variableName string) (typ string) {
	typ = ""
	sort.Sort(VariableSuffixTypeList)
	for _, vs2t := range VariableSuffixTypeList {
		if strings.HasSuffix(variableName, vs2t.Suffix) {
			typ = vs2t.Type
			return // 匹配第一个即返回
		}
	}
	return ""
}
