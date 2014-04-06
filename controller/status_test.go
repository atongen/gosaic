package controller

import (
	"bytes"
	"strings"
	"testing"
)

func setupStatusTest() (*Environment, *bytes.Buffer) {
	var out bytes.Buffer
	env := NewEnvironment("/tmp", &out, ":memory:", 2, false, false)
	env.Init()
	return env, &out
}

func TestStatus(t *testing.T) {
	env, out := setupStatusTest()
	defer env.DB.Close()

	Status(env)
	if !strings.Contains(out.String(), "Status: OK") {
		t.Error("The status was not ok.")
	}
}
