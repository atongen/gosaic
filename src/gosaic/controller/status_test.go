package controller

import (
	"bytes"
	"fmt"
	"gosaic/environment"
	"strings"
	"testing"
)

func TestStatus(t *testing.T) {
	fmt.Println("TestStatus")
	var out bytes.Buffer
	env, err := environment.GetTestEnv(&out)
	if err != nil {
		t.Error(err.Error())
	}
	err = env.Init()
	if err != nil {
		t.Error(err.Error())
	}
	defer env.Close()

	Status(env)
	if !strings.Contains(out.String(), "Status: OK") {
		t.Error("The status was not ok.")
	}
}
