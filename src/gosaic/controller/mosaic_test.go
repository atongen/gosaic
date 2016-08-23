package controller

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMosaic(t *testing.T) {
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

	mosaic := Mosaic(
		env,
		"testdata/jumping_bunny.jpg",
		"Jumping Bunny",
		"best",
		1000, 1000, 3, 2, 10, -1,
		filepath.Join(dir, "jumping_bunny_mosaic.jpg"),
		filepath.Join(dir, "jumping_bunny_macro.jpg"),
	)
	if mosaic == nil {
		t.Fatal("Failed to create mosaic")
	}

	result := out.String()
	expect := []string{
		"Indexing 4 images...",
		"Building 70 cover partials...",
		"Building 70 macro partials...",
		"Building 4 indexed image partials...",
		"Building 280 partial image comparisons...",
		"Building 70 mosaic partials...",
		"Drawing 70 mosaic partials...",
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
