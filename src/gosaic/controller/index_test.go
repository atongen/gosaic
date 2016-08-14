package controller

import (
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Error(err.Error())
	}
	defer env.Close()

	Index(env, "testdata")
	if !strings.Contains(out.String(), "Processing 1 images") ||
		strings.Contains(out.String(), "Error indexing images") {
		t.Errorf("Indexing failed: %s", out.String())
	}
}
