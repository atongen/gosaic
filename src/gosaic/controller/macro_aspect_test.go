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
		"Added 150 aspect cover partials",
		"Created macro for testdata/jumping_bunny.jpg",
		"Processing 150 macro partials",
		"Built 150 macro partials",
	}

	for _, e := range expect {
		if !strings.Contains(result, e) {
			t.Fatalf("Expected result to contain '%s', but it did not", e)
		}
	}
}
