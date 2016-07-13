package controller

import (
	"bytes"
	"fmt"
	"gosaic/environment"
	"path"
	"runtime"
	"strings"
	"testing"
)

func setupIndexTest() (environment.Environment, *bytes.Buffer, error) {
	var out bytes.Buffer
	env, err := environment.GetTestEnv(&out)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting test env: %s\n", err.Error())
	}
	err = env.Init()
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing test env: %s\n", err.Error())
	}
	return env, &out, nil
}

func TestIndex(t *testing.T) {
	env, out, err := setupIndexTest()
	if err != nil {
		t.Error(err.Error())
	}
	defer env.Close()
	_, file, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(file), "testdata")
	Index(env, dir)
	fmt.Println(out.String())
	if !strings.Contains(out.String(), "Processing 1 images") ||
		strings.Contains(out.String(), "Error indexing images") {
		t.Errorf("Indexing failed: %s", out.String())
	}
}
