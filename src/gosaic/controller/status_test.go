package controller

import (
	"bytes"
	"strings"
	"testing"
)

func TestStatus(t *testing.T) {
	var out bytes.Buffer
	env := GetTestEnv(&out)
	env.Init()
	defer env.Close()

	Status(env)
	if !strings.Contains(out.String(), "Status: OK") {
		t.Error("The status was not ok.")
	}
}
