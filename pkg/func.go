package pkg

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/scanner"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/iancoleman/strcase"
	"goa.design/goa/codegen"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/imports"
)

// TrimSpaces  去除开头结尾的非有效字符
func TrimSpaces(s string) string {
	return strings.Trim(s, "\r\n\t\v\f ")
}

// 封装 goa.design/goa/v3/codegen 方便后续可定制
func ToCamel(name string) string {
	return codegen.CamelCase(name, true, true)
}

func ToLowerCamel(name string) string {
	return codegen.CamelCase(name, false, true)
}

func SnakeCase(name string) string {
	return codegen.SnakeCase(name)
}

func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func GeneratePackageName(dstDir string) (packageName string, err error) {
	if runtime.GOOS == "windows" { // drop driver name
		index := strings.Index(dstDir, ":")
		if index >= 0 {
			dstDir = dstDir[index+1:]
		}
		dstDir = strings.ReplaceAll(dstDir, "\\", "/")
	}
	absoluteDir := dstDir
	if dstDir[0:1] != "/" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		absoluteDir, err = filepath.Abs(fmt.Sprintf("%s/%s", cwd, dstDir))
		if err != nil {
			return "", err
		}
	}
	if runtime.GOOS == "windows" {
		absoluteDir = strings.ReplaceAll(absoluteDir, "\\", "/")
	}

	basename := path.Base(absoluteDir)
	packageName = strings.ToLower(strcase.ToLowerCamel(basename))
	return

}

// FinalizeGoSource removes unneeded imports from the given Go source file and
// runs go fmt on it.
func FinalizeGoSource(path string) error {
	// Make sure file parses and print content if it does not.
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		content, _ := os.ReadFile(path)
		var buf bytes.Buffer
		scanner.PrintError(&buf, err)
		return fmt.Errorf("%s\n========\nContent:\n%s", buf.String(), content)
	}

	// Clean unused imports
	imps := astutil.Imports(fset, file)
	for _, group := range imps {
		for _, imp := range group {
			path := strings.Trim(imp.Path.Value, `"`)
			if !astutil.UsesImport(file, path) {
				if imp.Name != nil {
					astutil.DeleteNamedImport(fset, file, imp.Name.Name, path)
				} else {
					astutil.DeleteImport(fset, file, path)
				}
			}
		}
	}
	ast.SortImports(fset, file)
	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	if err := format.Node(w, fset, file); err != nil {
		return err
	}
	w.Close()

	// Format code using goimport standard
	bs, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	opt := imports.Options{
		Comments:   true,
		FormatOnly: true,
	}
	bs, err = imports.Process(path, bs, &opt)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bs, os.ModePerm)
}
