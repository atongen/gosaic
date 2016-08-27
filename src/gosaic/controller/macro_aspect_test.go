package controller

import "testing"

func TestMacroAspect(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	cover, macro := MacroAspect(env, "testdata/jumping_bunny.jpg", 1000, 1000, 2, 3, 10, "", "")
	if cover == nil || macro == nil {
		t.Fatal("Failed to create cover or macro")
	}

	expect := []string{
		"Building 150 cover partials...",
		"Building 150 macro partials...",
	}

	testResultExpect(t, out.String(), expect)
}
