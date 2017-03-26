package controller

import "testing"

func TestMacro(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	macroPartialService := env.ServiceFactory().MustMacroPartialService()

	// build a test cover
	cover := CoverAspect(env, 594, 554, 2, 3, 10)
	if cover == nil {
		t.Fatal("Failed to create cover")
	}
	macro := Macro(env, "testdata/jumping_bunny.jpg", cover.Id, "")
	if macro == nil {
		t.Fatal("Failed to create macro")
	}

	expect := []string{
		"Building 160 macro partials...",
	}

	testResultExpect(t, out.String(), expect)

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
