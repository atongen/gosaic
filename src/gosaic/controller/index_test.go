package controller

import (
	"path"
	"runtime"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Error(err.Error())
	}
	defer env.Close()

	_, file, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(file), "testdata")
	Index(env, dir)
	if !strings.Contains(out.String(), "Processing 1 images") ||
		strings.Contains(out.String(), "Error indexing images") {
		t.Errorf("Indexing failed: %s", out.String())
	}
}
