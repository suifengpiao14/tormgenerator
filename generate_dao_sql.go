package tormgenerator

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/suifengpiao14/funcs"
	"github.com/suifengpiao14/tormgenerator/parser/ddlparser"
	"github.com/suifengpiao14/tormgenerator/parser/tormparser"
	"github.com/suifengpiao14/tormgenerator/parser/tplparser"
)

const SourceInsertTpl = "insert ignore into `source` (`source_id`,`source_type`,`config`) values('%s','%s','%s');"
const TemplateInsertTpl = "insert ignore into `template` (`template_id`,`type`,`batch`,`description`,`source_id`,`tpl`) values('%s','SQL','%s','%s','%s','%s');"

func (b *Builder) GenerateDoaSQL(tplDefines tplparser.TPLDefines) (buf *bytes.Buffer, err error) {
	batch := b.dbConfig.Version
	moduleCamel := b.dbConfig.DatabaseName
	tables, err := b.GetTables()
	if err != nil {
		return nil, err
	}
	sourceId := b.dbConfig.DatabaseName
	var w bytes.Buffer
	sourceType := "SQL"
	config := fmt.Sprintf(`{"logLevel":"%s","dsn":"root:123456@tcp(mysql_address:3306)/%s?charset=utf8&timeout=1s&readTimeout=5s&writeTimeout=5s&parseTime=False&loc=Local&multiStatements=true","timeout":30}`, b.dbConfig.LogLevel, b.dbConfig.DatabaseName)
	sourceInsertSql := fmt.Sprintf(SourceInsertTpl, sourceId, sourceType, config)
	w.WriteString(sourceInsertSql)
	w.WriteString(tormparser.EOF)
	for _, tplDefine := range tplDefines {
		defineName := tplDefine.NameCamel()
		tableName := funcs.ToCamel(tplDefine.GetTable())
		templatId := generateTemplateId(moduleCamel, tableName, defineName)
		description := fmt.Sprintf("%s%s%s", moduleCamel, tableName, defineName)
		description = translateDescription(description, tables)
		templateInsertSql := fmt.Sprintf(TemplateInsertTpl, templatId, batch, description, sourceId, tplDefine.Content)
		w.WriteString(templateInsertSql)
		w.WriteString(tormparser.EOF)
	}
	return &w, nil
}

func translateDescription(description string, tables []*ddlparser.Table) (descriptionZH string) {
	extraTranslatemap := make(map[string]string)
	for _, table := range tables {
		extraTranslatemap[table.TableName] = table.Comment
	}
	descriptionZH = Translate(description, extraTranslatemap)
	return descriptionZH
}

func generateTemplateId(dbName string, tableName string, defineName string) (templateId string) {
	templateId = fmt.Sprintf("%s%s%s", dbName, tableName, defineName)
	return templateId
}

var TranslateMapDefault = map[string]string{
	"insert":   "新增",
	"update":   "修改",
	"del":      "删除",
	"list":     "列表",
	"paginate": "分页列表",
	"total":    "总数",
	"get":      "获取",
	"by":       "通过",
	"id":       "ID",
	"all":      "所有",
}

func Translate(name string, extraMap map[string]string) string {
	snakeName := funcs.SnakeCase(name)
	wordArr := strings.Split(snakeName, "_")
	titleArr := make([]string, 0)
	existsWordMap := make(map[string]string)
	formatExtraMap := make(map[string]string)
	for k, v := range extraMap {
		formatExtraMap[funcs.SnakeCase(k)] = v
	}
	for _, word := range wordArr {
		if _, ok := existsWordMap[word]; ok {
			continue
		}

		title, ok := formatExtraMap[word]
		if ok {
			existsWordMap[word] = title
			titleArr = append(titleArr, title)
			continue
		}

		title, ok = TranslateMapDefault[word]
		if ok {
			existsWordMap[word] = title
			titleArr = append(titleArr, title)
			continue
		}
		title = word
		existsWordMap[word] = title
		titleArr = append(titleArr, title)
	}
	out := strings.Join(titleArr, " ")
	return out
}
