package controller

import "testing"

func TestPartialAspect(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	gidxPartialService := env.ServiceFactory().MustGidxPartialService()

	err = Index(env, []string{"testdata", "../service/testdata"})
	if err != nil {
		t.Fatalf("Error indexing images: %s\n", err.Error())
	}

	cover, macro := MacroAspect(env, "testdata/jumping_bunny.jpg", 594, 554, 2, 3, 10, "", "")
	if cover == nil || macro == nil {
		t.Fatal("Failed to create cover or macro")
	}

	err = PartialAspect(env, macro.Id, -1.0)
	if err != nil {
		t.Fatalf("Error building partial aspects: %s\n", err.Error())
	}

	expect := []string{
		"Building 4 index image partials...",
	}

	testResultExpect(t, out.String(), expect)

	count, err := gidxPartialService.Count()
	if err != nil {
		t.Fatalf("Error counting partial aspects: %s\n", err.Error())
	}

	if count != 4 {
		t.Fatalf("Expected 4 gidx partials, but got %d\n", count)
	}
}

func TestPartialAspectWithin(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	gidxPartialService := env.ServiceFactory().MustGidxPartialService()

	err = Index(env, []string{"testdata", "../service/testdata"})
	if err != nil {
		t.Fatalf("Error indexing images: %s\n", err.Error())
	}

	cover, macro := MacroAspect(env, "testdata/jumping_bunny.jpg", 594, 554, 2, 3, 10, "", "")
	if cover == nil || macro == nil {
		t.Fatal("Failed to create cover or macro")
	}

	err = PartialAspect(env, macro.Id, 0.5)
	if err != nil {
		t.Fatalf("Error building partial aspects: %s\n", err.Error())
	}

	expect := []string{
		"Building 4 index image partials...",
	}

	testResultExpect(t, out.String(), expect)

	count, err := gidxPartialService.Count()
	if err != nil {
		t.Fatalf("Error counting partial aspects: %s\n", err.Error())
	}

	if count != 2 {
		t.Fatalf("Expected 2 gidx partials, but got %d\n", count)
	}
}
