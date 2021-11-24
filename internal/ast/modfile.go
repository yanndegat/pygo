package ast

import (
	"fmt"
	"golang.org/x/mod/modfile"
	"io/ioutil"

	"github.com/yanndegat/pygo/internal/utils"
)

type Mod struct {
	Main string
	Path string
}

func (m Mod) String() string {
	return fmt.Sprintf("main: %s, path: %s", m.Main, m.Path)
}

func (m Mod) Import() string {
	if m.Path == "" {
		return m.Main
	}
	return fmt.Sprintf("%s%s", m.Main, m.Path)
}

func ParseMod(dir string) (*Mod, error) {
	goMod, err := utils.FindFileInParentFolders("go.mod", dir)
	if err != nil {
		return nil, fmt.Errorf("Couldn't find go.mod: %v", err)
	}

	goModText, err := ioutil.ReadFile(goMod)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read modfile %s: %v", goMod, err)
	}

	mainMod := modfile.ModulePath(goModText)
	if mainMod == "" {
		return nil, fmt.Errorf("Oups! no Module defined in modfile %s.", goMod)
	}

	modPath, err := utils.AbsModPath(dir)
	if err != nil {
		return nil, err
	}

	return &Mod{
		Main: mainMod,
		Path: modPath,
	}, nil
}
