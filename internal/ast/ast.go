package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
)

type AstFunc struct {
	Name    string
	Params  []*ast.Field
	Results []*ast.Field
}

func (f *AstFunc) String() string {
	return f.Name
}

func ParseDir(dir string) (map[string][]*AstFunc, error) {
	pkgs, err := parser.ParseDir(token.NewFileSet(), dir, goFiles, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] %v scanned", dir)

	pyLibs := map[string][]*AstFunc{}
	for name, pkg := range pkgs {
		log.Printf("[TRACE] Parsing pkg name %v", pkg.Name)

		res, err := parsePkg(name, pkg)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] parsing of pkg %s failed: %v", name, err)
		}

		if pyLibs[name] == nil {
			pyLibs[name] = res
		} else {
			pyLibs[name] = append(pyLibs[name], res...)
		}
	}
	return pyLibs, nil
}

// Specify what files to parser
func goFiles(info os.FileInfo) bool {
	if strings.HasSuffix(info.Name(), ".go") {
		return true
	}
	return false
}

func parsePkg(name string, pkg *ast.Package) ([]*AstFunc, error) {
	astFuncs := []*AstFunc{}
	var err error

	for filePath, f := range pkg.Files {
		log.Printf("[DEBUG] Parsing file Name %v at %v", f.Name, filePath)
		source := path.Base(filePath)

		ast.Inspect(f, func(n ast.Node) bool {
			// handle function declarations without documentation
			if n != nil {
				log.Printf("[TRACE] %s:%v %v", source, n.Pos(), reflect.TypeOf(n))
			}

			if fn, ok := n.(*ast.FuncDecl); ok {
				if fn.Doc != nil && len(fn.Doc.List) > 0 {
					log.Printf("[TRACE] func %s in %s is exported", source, fn.Name.Name)
					for _, comm := range fn.Doc.List {
						log.Printf("[TRACE] func %s in %s comment is %s", source, fn.Name.Name, comm.Text)
						isExported, _err := commentFuncExport(comm.Text)
						if _err != nil {
							log.Printf("[ERROR] failed to parse %s/%s/%s : %v", source, fn.Name.Name, comm.Text, _err)
							err = _err
							return false
						}

						if isExported && ast.IsExported(fn.Name.Name) {
							var results []*ast.Field
							if fn.Type.Results != nil {
								results = fn.Type.Results.List
							}

							astFunc := &AstFunc{
								Name:    fn.Name.Name,
								Params:  fn.Type.Params.List,
								Results: results,
							}
							log.Printf("[DEBUG] func %v is exported in %s", astFunc, name)
							astFuncs = append(astFuncs, astFunc)
						}
					}

				}
			}
			return true
		})
	}

	return astFuncs, err
}

func commentFuncExport(text string) (bool, error) {
	return regexp.MatchString(`(^|^//|[[:space:]])@(pygo)\.(export)($|\W)`, text)
}
