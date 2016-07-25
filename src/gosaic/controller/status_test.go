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
	if !strings.Contains(out.String(), "Status: OK") {
		t.Error("The status was not ok.")
	}
}
