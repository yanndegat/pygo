package ast

import (
	"testing"
)

func TestMain_ParseMod(t *testing.T) {
	mod, err := ParseMod(".")
	if err != nil {
		t.Fatalf("%v", err)
	}

	shouldPath := "/internal/ast"
	shouldMain := "github.com/yanndegat/pygo"

	if mod.Path != shouldPath {
		t.Fatalf("match should be %v, was %v", shouldPath, mod.Path)
	}
	if mod.Main != shouldMain {
		t.Fatalf("match should be %v, was %v", shouldMain, mod.Main)
	}
}
