package utils

import (
	"testing"
)

func TestMain_FindFileInParentFolders(t *testing.T) {
	found, err := FindFileInParentFolders("go.mod", "./tests")
	if err != nil {
		t.Fatalf("%v", err)
	}

	if found == "" {
		t.Fatalf("go.mod should have been found.")
	}
}

func TestMain_AbsModPath(t *testing.T) {
	path, err := AbsModPath("./")
	if err != nil {
		t.Fatalf("%v", err)
	}

	shouldBe := "/internal/utils"

	if path != shouldBe {
		t.Fatalf("match should be %v, was %v", shouldBe, path)
	}
}
