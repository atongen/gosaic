package controller

import (
	"strings"
	"testing"
)

func TestStatus(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Error(err.Error())
	}
	defer env.Close()

	Status(env)
	result := out.String()

	if !strings.Contains(result, "Status: OK") {
		t.Error("The status was not ok.")
	}

	for _, ne := range []string{"fail", "error"} {
		if strings.Contains(strings.ToLower(result), ne) {
			t.Fatalf("Did not expect result to contain: %s, but it did\n", ne)
		}
	}
}
