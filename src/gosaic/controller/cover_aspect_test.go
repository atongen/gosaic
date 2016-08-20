package controller

import (
	"strings"
	"testing"
)

func TestCoverAspect(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	CoverAspect(env, "test", 1, 1, 1, 1, 1)

	result := out.String()
	if !strings.Contains(result, "Created cover test with 1 partials") ||
		strings.Contains(result, "Error") {
		t.Fatalf("CoverAspect failed: %s\n", result)
	}

	for _, ne := range []string{"fail", "error"} {
		if strings.Contains(strings.ToLower(result), ne) {
			t.Fatalf("Did not expect result to contain: %s, but it did\n", ne)
		}
	}
}
