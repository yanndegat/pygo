package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/yanndegat/pygo/internal/ast"
	"github.com/yanndegat/pygo/internal/libfunc"
	"github.com/yanndegat/pygo/internal/logging"
)

const (
	// The parent process will create a file to collect crash logs
	envTmpLogPath = "PYGO_TEMP_LOG_PATH"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	defer logging.PanicHandler()
	log.Printf("[DEBUG] pygo called with args: %v", os.Args)

	tmpLogPath := os.Getenv(envTmpLogPath)
	if tmpLogPath != "" {
		f, err := os.OpenFile(tmpLogPath, os.O_RDWR|os.O_APPEND, 0666)
		if err == nil {
			defer f.Close()

			log.Printf("[DEBUG] Adding temp file log sink: %s", f.Name())
			logging.RegisterSink(f)
		} else {
			log.Printf("[ERROR] Could not open temp log file: %v", err)
		}
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Printf("[ERROR] Oups!: %v", err)
		return 1
	}

	mod, err := ast.ParseMod(pwd)
	if err != nil {
		log.Printf("[ERROR] Couldn't parse mod info: %v", err)
		return 1
	}
	log.Printf("[INFO] mod is %v", mod)

	astLibs, err := ast.ParseDir(pwd)
	if err != nil {
		log.Printf("[ERROR] Couldn't parse dir %s: %v", pwd, err)
		return 1
	}

	libs := map[string][]*libfunc.Func{}
	for lib, astFs := range astLibs {
		fs := []*libfunc.Func{}
		for _, astF := range astFs {
			f, err := libfunc.ConvertFromAstF(lib, astF)
			if err != nil {
				log.Printf("[ERROR] Couldn't convert astFunc %s for lib %s: %v", astF.Name, lib, err)
				return 1
			}
			log.Printf("[DEBUG] adding func to lib %s: %v", lib, f)

			if !f.IsSupported() {
				log.Printf("[WARN] func %v from lib %s is not supported.", f, lib)
				continue
			}
			fs = append(fs, f)

		}
		libs[lib] = fs
	}

	for lib, funcs := range libs {
		// make pygo folder
		pygoDir := filepath.Join(pwd, "pygo")
		err := os.MkdirAll(pygoDir, os.ModePerm)
		if err != nil {
			log.Printf("[ERROR] Couldn't mkdir pygo %s: %v", pygoDir, err)
			return 1
		}
		// generate lib.pygo
		if err := generatePygo(pygoDir, lib, mod, funcs); err != nil {
			log.Printf("[ERROR] Couldn't generate %s.pygo in %s: %v", lib, pygoDir, err)
			return 1
		}

		// generate lib.py
		if err := generatePy(pygoDir, lib, mod, funcs); err != nil {
			log.Printf("[ERROR] Couldn't generate %s.py in %s: %v", lib, pygoDir, err)
			return 1
		}

		// build lib.py.so
		log.Printf("[INFO] build shared lib %s/_%s.so", pygoDir, lib)
		cmd := exec.Command("go", "build", "-buildmode=c-shared",
			"-o", fmt.Sprintf("_%s.so", lib),
			fmt.Sprintf("%s.go", lib))
		cmd.Dir = pygoDir
		cmd.Stderr = log.Writer()
		cmd.Stdout = log.Writer()
		if err := cmd.Run(); err != nil {
			log.Printf("[ERROR] Couldn't generate _%s.so in %s: %v", lib, pygoDir, err)
			return 1
		}
	}
	return 0
}

func generatePygo(dir, lib string, mod *ast.Mod, funcs []*libfunc.Func) error {
	f, err := os.Create(filepath.Join(dir, fmt.Sprintf("%s.go", lib)))
	if err != nil {
		return err
	}

	defer f.Close()
	return pyGoTemplate.Execute(f, struct {
		Timestamp time.Time
		Funcs     []*libfunc.Func
		Lib       string
		Dir       string
		Mod       *ast.Mod
	}{
		Timestamp: time.Now(),
		Lib:       lib,
		Mod:       mod,
		Dir:       dir,
		Funcs:     funcs,
	})
}

func generatePy(dir, lib string, mod *ast.Mod, funcs []*libfunc.Func) error {
	f, err := os.Create(filepath.Join(dir, fmt.Sprintf("%s.py", lib)))
	if err != nil {
		return err
	}

	defer f.Close()
	return pyTemplate.Execute(f, struct {
		Timestamp time.Time
		Funcs     []*libfunc.Func
		Lib       string
		Dir       string
		Mod       *ast.Mod
	}{
		Timestamp: time.Now(),
		Lib:       lib,
		Mod:       mod,
		Dir:       dir,
		Funcs:     funcs,
	})
}

var pyGoTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
// This file was generated by pygo at
// {{ .Timestamp }}
package main

import (
   "C"
   "{{ .Mod.Import }}"
)

{{- range $f := .Funcs }}

//export {{ $f.Name}}
func {{ $f.Name}}({{$f.GoSigArgs}}) {{$f.GoSigRet}} {
   {{ if not $f.IsVoid }}return {{- end }} {{ $f.GoFuncCall }}
}

{{- end }}

func main(){}
`))

var pyTemplate = template.Must(template.New("").Parse(`# Code generated by go generate; DO NOT EDIT.
# This file was generated by pygo at
# {{ .Timestamp }}
from pygo import gofunc

{{- range $f := .Funcs }}


@gofunc(lib="_{{$.Lib}}.so")
def {{ $f.Name }}({{$f.PySig}}): pass
{{- end }}
`))