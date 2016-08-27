package controller

import (
	"strings"
	"testing"
)

func TestMacroQuad(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	cover, macro := MacroQuad(env, "testdata/jumping_bunny.jpg", 1000, 1000, 10, 2, 250, "", "")
	if cover == nil || macro == nil {
		t.Fatal("Failed to create cover or macro")
	}

	result := out.String()

	expect := []string{
		"Building 150 cover partials...",
		"Building 150 macro partials...",
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
