package tplparser

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func getTorm() (tormMetaTpl string) {
	filename := "../../example/torm.tpl"
	content := readFile(filename)
	return content
}
func readFile(filename string) (contnent string) {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	contnent = string(b)
	return contnent
}

func TestGetTable(t *testing.T) {
	t.Run("form", func(t *testing.T) {
		torm := `
		{{define "ServiceGetByServiceID"}}
select * from t_service  where service_id=:ServiceID  and deleted_at is null;
{{end}}

		`
		tplDefine, err := newTPLDefine(torm)
		require.NoError(t, err)
		table := tplDefine.GetTable()
		fmt.Println(table)
	})
	t.Run("form2", func(t *testing.T) {
		torm := "{{define \"ServiceGetByServiceID\"}} select * from `t_service`  where `service_id`=:ServiceID  and `deleted_at` is null; {{end}}"
		tplDefine, err := newTPLDefine(torm)
		require.NoError(t, err)
		table := tplDefine.GetTable()
		fmt.Println(table)
	})

	t.Run("update", func(t *testing.T) {
		torm := `
		{{define "ServiceDel"}}
update t_service set deleted_at={{currentTime .}} where service_id=:ServiceID;
{{end}}

		`
		tplDefine, err := newTPLDefine(torm)
		require.NoError(t, err)
		table := tplDefine.GetTable()
		fmt.Println(table)
	})
	t.Run("update2", func(t *testing.T) {
		torm := "{{define \"ServiceDel\"}}update `t_service` set `deleted_at`={{currentTime .}} where `service_id`=:ServiceID;{{end}}"
		tplDefine, err := newTPLDefine(torm)
		require.NoError(t, err)
		table := tplDefine.GetTable()
		fmt.Println(table)
	})

	t.Run("insert into ", func(t *testing.T) {
		torm := `
		{{define "APIInsert"}}
insert into t_api (api_id,service_id,name,title,tags,uri,summary,description)values
 (:APIID,:ServiceID,:Name,:Title,:Tags,:URI,:Summary,:Description);
{{end}}
		`
		tplDefine, err := newTPLDefine(torm)
		require.NoError(t, err)
		table := tplDefine.GetTable()
		fmt.Println(table)
	})
	t.Run("insert into2 ", func(t *testing.T) {
		torm := "{{define \"ParameterInsert\"}}insert into `t_parameter` (`parameter_id`,`service_id`,`api_id`,`validate_schema_id`,`full_name`,`name`,`title`,`type`,`tag`,`method`,`http_status`,`position`,`example`,`deprecated`,`required`,`serialize`,`explode`,`allow_empty_value`,`allow_reserved`,`description`)values(:ParameterID,:ServiceID,:APIID,:ValidateSchemaID,:FullName,:Name,:Title,:Type,:Tag,:Method,:HTTPStatus,:Position,:Example,:Deprecated,:Required,:Serialize,:Explode,:AllowEmptyValue,:AllowReserved,:Description);{{end}}"
		tplDefine, err := newTPLDefine(torm)
		require.NoError(t, err)
		table := tplDefine.GetTable()
		fmt.Println(table)
	})

	t.Run("empty", func(t *testing.T) {
		torm := `
		{{define "ParameterPaginateWhere"}}
  
{{end}}
`
		tplDefine, err := newTPLDefine(torm)
		require.NoError(t, err)
		table := tplDefine.GetTable()
		fmt.Println(table)
	})

}
