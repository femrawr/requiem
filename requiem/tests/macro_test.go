package tests_test

import (
	"fmt"
	"testing"

	"requiem/macro"
)

const TEST_MACRO_FILE_PATH string = ""

func TestFileValidation(test *testing.T) {
	err := macro.ValidateFile(TEST_MACRO_FILE_PATH)
	if err != nil {
		test.Errorf("Failed to validate - %v", err)
		return
	}

	fmt.Println("Validation successful")
}

func TestParseMacro(test *testing.T) {
	parsed, err := macro.ParseMacro(TEST_MACRO_FILE_PATH)
	if err != nil {
		test.Errorf("Failed to parse - %v", err)
		return
	}

	fmt.Printf("Parsed - %s\n", parsed.Encode())
}
