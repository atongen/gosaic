package controller

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMosaicDraw(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	dir, err := ioutil.TempDir("", "gosaic_test_mosaic_draw")
	if err != nil {
		t.Fatal("Error getting temp dir for mosaic draw test: %s\n", err.Error())
	}
	defer os.RemoveAll(dir)

	Index(env, []string{"testdata", "../service/testdata"})
	cover, macro := MacroAspect(env, "testdata/jumping_bunny.jpg", 1000, 1000, 2, 3, 10, "")
	if cover == nil || macro == nil {
		t.Fatal("Failed to create cover or macro")
	}
	PartialAspect(env, macro.Id)
	Compare(env, macro.Id)
	mosaic := MosaicBuild(env, "Jumping Bunny", "best", macro.Id, 0)
	if mosaic == nil {
		t.Fatal("Failed to build mosaic")
	}
	MosaicDraw(env, mosaic.Id, filepath.Join(dir, "jumping_bunny_mosaic.jpg"))

	result := out.String()
	expect := []string{
		"Indexing 4 images...",
		"Building 150 cover partials...",
		"Created cover testdata/jumping_bunny.jpg-",
		"Building 150 macro partials",
		"Created macro for path testdata/jumping_bunny.jpg with cover testdata/jumping_bunny.jpg-",
		"Creating 4 aspect partials for indexed images",
		"Creating mosaic with 150 total partials",
		"Creating 600 partial image comparisons...",
		"Building 150 mosaic partials",
		"Drawing 150 mosaic partials",
		"Wrote mosaic Jumping Bunny to",
		"jumping_bunny_mosaic.jpg",
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
