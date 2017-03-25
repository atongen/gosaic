package controller

import (
	"bytes"
	"github.com/atongen/gosaic/environment"
	"strings"
	"testing"
)

func setupControllerTest() (environment.Environment, *bytes.Buffer, error) {
	var out bytes.Buffer
	env, err := environment.GetTestEnv(&out)
	if err != nil {
		return nil, nil, err
	}

	err = env.Init()
	if err != nil {
		return nil, nil, err
	}

	return env, &out, nil
}

func testResultExpect(t *testing.T, result string, expect []string) {
	for _, e := range expect {
		if !strings.Contains(result, e) {
			t.Fatalf("Expected result to contain '%s', but it did not:\n%s\n", e, result)
		}
	}

	for _, ne := range []string{"fail", "error"} {
		if strings.Contains(strings.ToLower(result), ne) {
			t.Fatalf("Did not expect result to contain: %s, but it did:\n%s\n", ne, result)
		}
	}
}
