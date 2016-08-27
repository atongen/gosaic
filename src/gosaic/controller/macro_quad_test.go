package controller

import "testing"

func TestMacroQuad(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	cover, macro := MacroQuad(env, "testdata/jumping_bunny.jpg", 200, 200, 10, 2, 50, "", "")
	if cover == nil || macro == nil {
		t.Fatal("Failed to create cover or macro")
	}

	expect := []string{
		"Building 10 macro partial quads...",
	}

	testResultExpect(t, out.String(), expect)
}
