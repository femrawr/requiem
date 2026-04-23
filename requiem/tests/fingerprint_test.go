package tests_test

import (
	"fmt"
	"testing"

	"requiem/funcs"
)

func TestGenerateFingerprint(test *testing.T) {
	finerprint, err := funcs.GenFingerprint()
	if err != nil {
		test.Errorf("Failed to generate fingerprint - %v", err)
		return
	}

	fmt.Printf("Generated fingerprint - %s\n", finerprint)
}
