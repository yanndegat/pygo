package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	maxFolders = 100
)

func AbsModPath(dir string) (string, error) {
	goMod, err := FindFileInParentFolders("go.mod", dir)
	if err != nil {
		return "", fmt.Errorf("Couldn't find go.mod: %v", err)
	}

	if goMod == "" {
		return "", fmt.Errorf("Couldn't find go.mod.")
	}

	fullDir, err := filepath.Abs(filepath.ToSlash(dir))
	if err != nil {
		return "", err
	}
	goModDir, err := filepath.Abs(filepath.Dir(filepath.ToSlash(goMod)))
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(fullDir, goModDir), nil
}

func FindFileInParentFolders(fileName, dir string) (string, error) {
	curDir, err := filepath.Abs(filepath.ToSlash(dir))
	if err != nil {
		return "", err
	}
	prevDir := curDir

	// avoid infinite loop due to symlinks
	for i := 0; i < maxFolders; i++ {
		fileName := filepath.ToSlash(filepath.Join(curDir, fileName))
		file, err := os.Stat(fileName)
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}

		if file != nil {
			return fileName, nil
		}

		prevDir = curDir
		curDir = filepath.ToSlash(filepath.Dir(curDir))
		if curDir == prevDir {
			// reached root
			return "", nil
		}

	}

	return "", fmt.Errorf("MaxFolders(%d) reached! aborting.", maxFolders)
}
