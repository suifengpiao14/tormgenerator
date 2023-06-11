package tormgenerator

import (
	"fmt"
)

type TengoAPIModel struct {
	APIID        string
	Version      string
	Description  string
	Method       string
	Route        string
	TemplateIDs  string
	MainScript   string
	PreScript    string
	PostScript   string
	InputSchema  string
	OutputSchema string
}

const TengoApiInsertTpl = "insert ignore into `t_tengo_api` (`api_id`,`batch`,`description`,`method`,`route`,`template_ids`,`pre_script`,`main_script`,`post_script`,`input_schema`,`output_schema`) values('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s');"

func GenerateDoaTengoapi(tengoApiModel TengoAPIModel) (tengoApiInsertSql string, err error) {
	tengoApiInsertSql = fmt.Sprintf(TengoApiInsertTpl, tengoApiModel.APIID, tengoApiModel.Version, tengoApiModel.Description, tengoApiModel.Method, tengoApiModel.Route, tengoApiModel.TemplateIDs, tengoApiModel.PreScript, tengoApiModel.MainScript, tengoApiModel.PostScript, tengoApiModel.InputSchema, tengoApiModel.OutputSchema)
	return tengoApiInsertSql, nil
}
