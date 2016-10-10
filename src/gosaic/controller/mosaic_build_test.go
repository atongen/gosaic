package controller

import (
	"strings"
	"testing"
)

func TestMosaicBuildRandom(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	err = Index(env, []string{"testdata", "../service/testdata"})
	if err != nil {
		t.Fatalf("Error indexing images: %s\n", err.Error())
	}

	cover, macro := MacroAspect(env, "testdata/jumping_bunny.jpg", 1000, 1000, 2, 3, 10, "", "")
	if cover == nil || macro == nil {
		t.Fatal("Failed to create cover or macro")
	}

	err = PartialAspect(env, macro.Id, -1.0)
	if err != nil {
		t.Fatalf("Error building partial aspects: %s\n", err.Error())
	}

	err = Compare(env, macro.Id)
	if err != nil {
		t.Fatalf("Comparing images: %s\n", err.Error())
	}

	mosaic := MosaicBuild(env, "random", macro.Id, 0)
	if mosaic == nil {
		t.Fatal("Failed to build mosaic")
	}

	expect := []string{
		"Building 150 mosaic partials...",
	}

	testResultExpect(t, out.String(), expect)
}

func TestMosaicBuildBest(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	err = Index(env, []string{"testdata", "../service/testdata"})
	if err != nil {
		t.Fatalf("Error indexing images: %s\n", err.Error())
	}

	cover, macro := MacroAspect(env, "testdata/jumping_bunny.jpg", 1000, 1000, 2, 3, 10, "", "")
	if cover == nil || macro == nil {
		t.Fatal("Failed to create cover or macro")
	}

	err = PartialAspect(env, macro.Id, -1.0)
	if err != nil {
		t.Fatalf("Error building partial aspects: %s\n", err.Error())
	}

	err = Compare(env, macro.Id)
	if err != nil {
		t.Fatalf("Comparing images: %s\n", err.Error())
	}

	mosaic := MosaicBuild(env, "best", macro.Id, 0)
	if mosaic == nil {
		t.Fatal("Failed to build mosaic")
	}

	result := out.String()
	expect := []string{
		"Building 150 mosaic partials...",
	}

	for _, e := range expect {
		if !strings.Contains(result, e) {
			t.Fatalf("Expected result to contain '%s', but it did not\n", e)
		}
	}

	for _, ne := range []string{"fail", "error"} {
		if strings.Contains(strings.ToLower(result), ne) {
			t.Fatalf("Did not expect result to contain: %s, but it did\n", ne)
		}
	}
}
