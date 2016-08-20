package controller

import (
	"strings"
	"testing"
)

func TestMosaicBuild(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	Index(env, []string{"testdata", "../service/testdata"})
	_, macro := MacroAspect(env, "testdata/jumping_bunny.jpg", 1000, 1000, 2, 3, 10)
	PartialAspect(env, macro.Id)
	Compare(env, macro.Id)
	MosaicBuild(env, "Jumping Bunny", macro.Id, 0)

	result := out.String()
	expect := []string{
		"Indexing 4 images...",
		"Building 150 cover partials...",
		"Created cover testdata/jumping_bunny.jpg-",
		"Building 150 macro partials",
		"Created macro for path testdata/jumping_bunny.jpg with cover testdata/jumping_bunny.jpg-",
		"Creating 4 aspect partials for indexed images",
		"Creating 600 partial image comparisons...",
		"Creating mosaic with 150 total partials",
		"Building 150 mosaic partials",
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
