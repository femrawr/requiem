package tests_test

import (
	"fmt"
	"testing"

	"requiem/utils"
)

const ROOT_DIR_PATH string = ""
const MAX_TREE_DEPTH int = 2

func TestGenerateFileTree(test *testing.T) {
	tree, err := utils.GenFileTree(ROOT_DIR_PATH, MAX_TREE_DEPTH)
	if err != nil {
		test.Errorf("Failed to generate tree - %s", err)
		return
	}

	fmt.Printf("Generated tree -\n%s", tree)
}
