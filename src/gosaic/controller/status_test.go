package controller

import "testing"

func TestStatus(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Error(err.Error())
	}
	defer env.Close()

	Status(env)

	expect := []string{
		"Status: OK",
	}

	testResultExpect(t, out.String(), expect)
}
