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
	_, macro := MacroAspect(env, "testdata/jumping_bunny.jpg", 1000, 1000, 2, 3, 10)
	Compare(env, macro.Id)

	result := out.String()

	expect := []string{
		"Creating 4 index partials for aspect 2x3",
		"100 / 600 partial comparisons created",
		"200 / 600 partial comparisons created",
		"300 / 600 partial comparisons created",
		"400 / 600 partial comparisons created",
		"500 / 600 partial comparisons created",
		"600 / 600 partial comparisons created",
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
