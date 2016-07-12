package controller

import (
	"bytes"
	"path"
	"runtime"
	"strings"
	"testing"
)

func setupIndexTest() (Environment, *bytes.Buffer) {
	var out bytes.Buffer
	env := GetTestEnv(&out)
	env.Init()
	return env, &out
}

func TestIndex(t *testing.T) {
	env, out := setupIndexTest()
	defer env.Close()
	_, file, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(file), "..")
	Index(env, dir)
	if !strings.Contains(out.String(), "1 of 1") {
		t.Error("Indexing did not occur")
	}
}
