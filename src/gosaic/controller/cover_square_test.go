package controller

import (
	"strings"
	"testing"
)

func TestCoverSquare(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	CoverSquare(env, "test", 1, 1, 1)

	result := out.String()
	if !strings.Contains(result, "Created cover test with 1 partials") ||
		strings.Contains(result, "Error") {
		t.Fatalf("CoverSquare failed: %s\n", result)
	}
}
