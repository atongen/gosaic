package controller

import (
	"strings"
	"testing"
)

func TestPartialAspect(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	Index(env, []string{"testdata", "../service/testdata"})
	cover := CoverAspect(env, "macroTest", 594, 554, 2, 3, 10)
	macro := Macro(env, "testdata/jumping_bunny.jpg", cover.Id)
	PartialAspect(env, macro.Id)

	result := out.String()
	expect := []string{
		"Indexing 4 images...",
		"Building 160 cover partials...",
		"Created cover macroTest",
		"Building 160 macro partials",
		"Created macro for path testdata/jumping_bunny.jpg with cover macroTest",
		"Creating 4 aspect partials for indexed images",
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
