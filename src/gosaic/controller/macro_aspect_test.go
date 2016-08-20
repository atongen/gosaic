package controller

import (
	"fmt"
	"strings"
	"testing"
)

func TestMacroAspect(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	MacroAspect(env, "testdata/jumping_bunny.jpg", 1000, 1000, 2, 3, 10)

	result := out.String()
	fmt.Println(result)
	expect := []string{
		"Created cover testdata/jumping_bunny.jpg-",
		"Building 150 cover partials...",
		"Processing 150 macro partials",
		"Built macro for testdata/jumping_bunny.jpg with cover testdata/jumping_bunny.jpg-",
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
