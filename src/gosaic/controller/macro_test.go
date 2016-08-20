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

	macroService, err := env.MacroService()
	if err != nil {
		t.Fatalf("Unable to get macro service: %s\n", err.Error())
	}

	macroPartialService, err := env.MacroPartialService()
	if err != nil {
		t.Fatalf("Unable to get macro partial service: %s\n", err.Error())
	}

	// build a test cover
	CoverAspect(env, "macroTest", 594, 554, 2, 3, 10)
	Macro(env, "testdata/jumping_bunny.jpg", "macroTest")

	result := out.String()
	if !strings.Contains(result, "Created cover macroTest") {
		t.Fatal("Cover controller output not found")
	}

	if !strings.Contains(result, "Built macro for testdata/jumping_bunny.jpg with cover macroTest") {
		t.Fatal("Macro controller output not found")
	}

	macro, err := macroService.GetOneBy("md5sum = ?", "40aadb49394b753db0b7ee8bb555623a")
	if err != nil {
		t.Fatal("Error getting created macro: %s\n", err.Error())
	}

	if macro == nil {
		t.Fatal("Unable to find created macro")
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
}
