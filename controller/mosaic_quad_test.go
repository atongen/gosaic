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
		200, 200, 10, -1, 2, 50, -1, -1,
		-1.0,
		filepath.Join(dir, "jumping_bunny_cover.png"),
		filepath.Join(dir, "jumping_bunny_macro.jpg"),
		filepath.Join(dir, "jumping_bunny_mosaic.jpg"),
		true,
		false,
	)
	if mosaic == nil {
		t.Fatal("Failed to create mosaic")
	}

	expect := []string{
		"Indexing 4 images...",
		"Building macro quad with 10 splits, 34 partials, min depth 1, max depth 2, min area 50...",
		"Building 4 index image partials...",
		"Building 136 partial image comparisons...",
		"Building 34 mosaic partials...",
		"Drawing 34 mosaic partials...",
	}

	testResultExpect(t, out.String(), expect)

	project, err := envProject(env)
	if err != nil {
		t.Fatalf("Error getting project from environment: %s\n", err.Error())
	} else if project == nil {
		t.Fatalf("Project not found in environment.")
	}

	if !project.IsComplete {
		t.Fatalf("Project not marked complete.")
	}
}
