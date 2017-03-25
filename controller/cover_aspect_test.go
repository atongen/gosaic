package controller

import "testing"

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

	expect := []string{
		"Building 1 cover partials...",
	}

	testResultExpect(t, out.String(), expect)
}
