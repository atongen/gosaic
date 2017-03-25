package controller

import "testing"

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

	expect := []string{
		"Indexing 1 images...",
	}

	testResultExpect(t, out.String(), expect)
}
