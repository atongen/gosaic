package controller

import (
	"bytes"
	"gosaic/environment"
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
