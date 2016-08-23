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

	cover := CoverAspect(env, 1, 1, 1, 1, 1)
	if cover == nil {
		t.Fatal("Failed to create cover")
	}

	result := out.String()

	expect := []string{
		"Building 1 cover partials...",
	}

	for _, e := range expect {
		if !strings.Contains(result, e) {
			t.Fatalf("Expected result to contain '%s', but it did not", e)
		}
	}

	for _, ne := range []string{"fail", "error"} {
		if strings.Contains(strings.ToLower(result), ne) {
			t.Fatalf("Did not expect result to contain: %s, but it did\n", ne)
		}
	}
}
