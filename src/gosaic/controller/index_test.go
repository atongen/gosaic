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

	err = Index(env, []string{"testdata"})
	if err != nil {
		t.Fatalf("Error indexing images: %s\n", err.Error())
	}

	result := out.String()

	expect := []string{
		"Indexing 1 images...",
	}

	for _, e := range expect {
		if !strings.Contains(result, e) {
			t.Fatalf("Expected result to contain '%s', but it did not", e)
		}
	}

	for _, ne := range []string{"fail", "error"} {
		if strings.Contains(strings.ToLower(result), ne) {
			t.Fatalf("Did not expect result to contain: %s, but it did\n", ne)
		}
	}
}
