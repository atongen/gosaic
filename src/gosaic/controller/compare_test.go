package controller

import (
	"strings"
	"testing"
)

func TestCompare(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	Index(env, []string{"testdata", "../service/testdata"})
	cover, macro := MacroAspect(env, "testdata/jumping_bunny.jpg", 1000, 1000, 2, 3, 10, "")
	if cover == nil || macro == nil {
		t.Fatal("Failed to create cover or macro")
	}
	PartialAspect(env, macro.Id)
	Compare(env, macro.Id)

	result := out.String()

	expect := []string{
		"Creating 4 aspect partials for indexed images",
		"Creating 600 partial image comparisons...",
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
