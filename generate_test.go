package generaterepository

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/generaterepository/pkg/ddlparser"
)

func getDBConfig() (dbcfg *ddlparser.DatabaseConfig) {
	cfgFile := "./example/dbconfig.json"
	b, err := os.ReadFile(cfgFile)
	if err != nil {
		panic(err)
	}
	dbcfg = &ddlparser.DatabaseConfig{}
	err = json.Unmarshal(b, dbcfg)
	if err != nil {
		panic(err)
	}
	return dbcfg
}

func getDDL() (ddl string) {
	ddlsql := "./example/ddl.sql"
	b, err := os.ReadFile(ddlsql)
	if err != nil {
		panic(err)
	}
	ddl = string(b)
	return ddl
}

func TestGenerateModel(t *testing.T) {
	ddl := getDDL()
	dbConfig := getDBConfig()
	buf, err := GenerateModel(ddl, *dbConfig)
	require.NoError(t, err)
	fmt.Println(buf.String())
}
