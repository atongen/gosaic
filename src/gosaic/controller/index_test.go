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

	Index(env, []string{"testdata"})

	result := out.String()

	if !strings.Contains(result, "Indexing 1 images...") ||
		strings.Contains(result, "Error indexing images") {
		t.Errorf("Indexing failed: %s", result)
	}

	for _, ne := range []string{"fail", "error"} {
		if strings.Contains(strings.ToLower(result), ne) {
			t.Fatalf("Did not expect result to contain: %s, but it did\n", ne)
		}
	}
}
