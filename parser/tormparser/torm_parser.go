package tormparser

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/funcs"
	"github.com/suifengpiao14/jsonschemaline"
	"github.com/suifengpiao14/tormgenerator/parser/ddlparser"
	"github.com/suifengpiao14/tormgenerator/parser/tplparser"
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

type TormStruct struct {
	StructName  string
	TplIdentity string
	Name        string
	Variables   tplparser.Variables
	FullName    string
	Type        string
	Torm        string
	OutEntity   *TormStruct // 输出对象
}

// GetLineschema 获取输入输出schema
func (t TormStruct) GetLineschema() (inschema string, outschema string, err error) {
	inVariables := t.Variables
	inschema, err = inVariables.Lineschema(fmt.Sprintf("%sIn", t.Name), jsonschemaline.LINE_SCHEMA_DIRECTION_IN)
	if err != nil {
		return "", "", err
	}
	var outVariables tplparser.Variables
	if t.OutEntity != nil {
		outVariables = t.OutEntity.Variables
	}
	outschema, err = outVariables.Lineschema(fmt.Sprintf("%sOut", t.Name), jsonschemaline.LINE_SCHEMA_DIRECTION_OUT)
	if err != nil {
		return "", "", err
	}
	return inschema, outschema, nil
}

const STRUCT_DEFINE_NANE_FORMAT = "%sEntity"

// ParserTorm 根据tpl define 和table 提取其中变量,供后续生成go entity 或者 doa 的 template 使用(属于底层函数)
func ParserTorm(sqltplDefineText *tplparser.TPLDefine, tableList []*ddlparser.Table) (tormStuct *TormStruct, err error) {
	variableList := sqltplDefineText.GetVariables()
	variableList, err = formatVariableTypeByTableColumn(variableList, tableList)
	if err != nil {
		return nil, err
	}
	columnArr := parseSQLSelectColumn(sqltplDefineText.Text)
	allColumnVariable := columnsToVariables(tableList)
	columnVariables := make(tplparser.Variables, 0)
	for _, columnVariable := range allColumnVariable {
		for _, columnName := range columnArr {
			if columnName == columnVariable.Name {
				columnVariables = append(columnVariables, columnVariable)
			}
		}
	}
	tplIdentity := ""
	for _, table := range tableList {
		if table.DatabaseConfig.DatabaseName != "" {
			tplIdentity = strings.TrimSpace(table.DatabaseConfig.DatabaseName)
			break
		}
	}

	camelName := sqltplDefineText.NameCamel()
	outName := fmt.Sprintf("%sOut", camelName)
	tormReplacer := strings.NewReplacer(`\r`, "")
	tormStuct = &TormStruct{
		Name:        camelName,
		TplIdentity: tplIdentity,
		Type:        sqltplDefineText.Type(),
		Torm:        tormReplacer.Replace(strconv.Quote(sqltplDefineText.Text)),
		StructName:  fmt.Sprintf(STRUCT_DEFINE_NANE_FORMAT, camelName),
		Variables:   variableList,
		FullName:    sqltplDefineText.Name,
		OutEntity: &TormStruct{
			Name:       outName,
			Type:       sqltplDefineText.Type(),
			StructName: fmt.Sprintf(STRUCT_DEFINE_NANE_FORMAT, camelName),
			Variables:  columnVariables,
			FullName:   fmt.Sprintf("%sOut", sqltplDefineText.Name),
		},
	}
	return tormStuct, nil
}

func parseSQLSelectColumn(sql string) []string {
	grep := regexp.MustCompile(`(?i)select(.+)from`)
	match := grep.FindAllStringSubmatch(sql, -1)
	if len(match) < 1 {
		return make([]string, 0)
	}
	fieldStr := match[0][1]
	out := strings.Split(funcs.StandardizeSpaces(fieldStr), ",")
	return out
}

func columnsToVariables(tableList []*ddlparser.Table) (variables tplparser.Variables) {
	allVariables := make(tplparser.Variables, 0)
	tableColumnMap := make(map[string]*ddlparser.Column)
	columnNameMap := make(map[string]string)
	for _, table := range tableList {
		for _, column := range table.Columns { //todo 多表字段相同，类型不同时，只会取列表中最后一个
			tableColumnMap[column.CamelName] = column
			lname := strings.ToLower(column.CamelName)
			columnNameMap[lname] = column.CamelName // 后续用于检测模板变量拼写错误
			variable := &tplparser.Variable{
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

func parseSQLTPLTableName(sqlTpl string) (tableList []string, err error) {

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

func formatVariableTypeByTableColumn(variableList tplparser.Variables, tableList []*ddlparser.Table) (varaibles tplparser.Variables, err error) {
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
