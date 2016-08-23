package controller

import (
	"strings"
	"testing"
)

func TestMacro(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	macroPartialService, err := env.MacroPartialService()
	if err != nil {
		t.Fatalf("Unable to get macro partial service: %s\n", err.Error())
	}

	// build a test cover
	cover := CoverAspect(env, 594, 554, 2, 3, 10)
	if cover == nil {
		t.Fatal("Failed to create cover")
	}
	macro := Macro(env, "testdata/jumping_bunny.jpg", cover.Id, "")
	if macro == nil {
		t.Fatal("Failed to create macro")
	}

	result := out.String()
	expect := []string{
		"Building 160 cover partials...",
		"Created cover macroTest",
		"Building 160 macro partials",
		"Created macro for path testdata/jumping_bunny.jpg with cover macroTest",
	}

	for _, e := range expect {
		if !strings.Contains(result, e) {
			t.Fatalf("Expected result to contain '%s', but it did not", e)
		}
	}

	macroPartials, err := macroPartialService.FindAll("id DESC", 1000, 0, "macro_id = ?", macro.Id)
	if err != nil {
		t.Fatalf("Error FindAll macro partials for created macro: %s\n", err.Error())
	}

	if len(macroPartials) != 160 {
		t.Fatalf("Expected 160 macro partials, but got %d\n", len(macroPartials))
	}

	for _, mp := range macroPartials {
		if len(mp.Pixels) == 0 {
			t.Fatal("macro partial pixels are empty")
		}

		for _, lab := range mp.Pixels {
			if lab.L == 0.0 &&
				lab.A == 0.0 &&
				lab.B == 0.0 &&
				lab.Alpha == 0.0 {
				t.Fatal("macro partial pixel lab is empty")
			}
		}
	}

	for _, ne := range []string{"fail", "error"} {
		if strings.Contains(strings.ToLower(result), ne) {
			t.Fatalf("Did not expect result to contain: %s, but it did\n", ne)
		}
	}
}
