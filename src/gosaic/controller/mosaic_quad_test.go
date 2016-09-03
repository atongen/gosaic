package controller

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestMosaicQuad(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	dir, err := ioutil.TempDir("", "gosaic_test_mosaic_quad")
	if err != nil {
		t.Fatal("Error getting temp dir for mosaic quad test: %s\n", err.Error())
	}
	defer os.RemoveAll(dir)

	Index(env, []string{"testdata", "../service/testdata"})

	mosaic := MosaicQuad(
		env,
		"testdata/jumping_bunny.jpg",
		"Jumping Bunny",
		"random",
		200, 200, 10, 2, 50, -1,
		-1.0,
		filepath.Join(dir, "jumping_bunny_cover.png"),
		filepath.Join(dir, "jumping_bunny_macro.jpg"),
		filepath.Join(dir, "jumping_bunny_mosaic.jpg"),
	)
	if mosaic == nil {
		t.Fatal("Failed to create mosaic")
	}

	expect := []string{
		"Indexing 4 images...",
		"Building 10 macro partial quads...",
		"Building 4 index image partials...",
		"Building 120 partial image comparisons...",
		"Building 30 mosaic partials...",
		"Drawing 30 mosaic partials...",
	}

	testResultExpect(t, out.String(), expect)
}
