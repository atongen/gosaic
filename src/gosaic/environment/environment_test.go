package environment

import (
	"bytes"
	"testing"
)

func setupEnvTest() (Environment, *bytes.Buffer, error) {
	var out bytes.Buffer
	env, err := GetTestEnv(&out)
	if err != nil {
		return nil, nil, err
	}

	err = env.Init()
	if err != nil {
		return nil, nil, err
	}

	return env, &out, nil
}

func TestServices(t *testing.T) {
	env, _, err := setupEnvTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	_, err = env.GidxService()
	if err != nil {
		t.Fatalf("Error getting gidxService: %s\n", err.Error())
	}

	_, err = env.AspectService()
	if err != nil {
		t.Fatalf("Error getting aspectService: %s\n", err.Error())
	}

	_, err = env.GidxPartialService()
	if err != nil {
		t.Fatalf("Error getting gidxPartialService: %s\n", err.Error())
	}

	_, err = env.CoverService()
	if err != nil {
		t.Fatalf("Error getting coverService: %s\n", err.Error())
	}

	_, err = env.CoverPartialService()
	if err != nil {
		t.Fatalf("Error getting coverService: %s\n", err.Error())
	}

	_, err = env.MacroService()
	if err != nil {
		t.Fatalf("Error getting macroService: %s\n", err.Error())
	}

	_, err = env.MacroPartialService()
	if err != nil {
		t.Fatalf("Error getting macroPartialService: %s\n", err.Error())
	}

	_, err = env.PartialComparisonService()
	if err != nil {
		t.Fatalf("Error getting partialComparisonService: %s\n", err.Error())
	}
}
