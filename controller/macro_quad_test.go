package controller

import (
	"fmt"
	"testing"
)

type argTestIn struct {
	width, height, size, minDepth, maxDepth, minArea, maxArea int
}

type argTestOut struct {
	size, minDepth, maxDepth, minArea, maxArea int
	err                                        string
}

func TestMacroQuadFixArgs(t *testing.T) {
	for _, tt := range []struct {
		t argTestIn
		e argTestOut
	}{
		{
			argTestIn{100, 100, -1, -1, -1, -1, -1},
			argTestOut{40, 1, 2, 1225, 0, ""},
		},
		{
			argTestIn{300, 300, -1, -1, -1, -1, -1},
			argTestOut{96, 2, 3, 1225, 0, ""},
		},
		{
			argTestIn{600, 600, -1, -1, -1, -1, -1},
			argTestOut{167, 2, 4, 1225, 0, ""},
		},
		{
			argTestIn{1200, 1200, -1, -1, -1, -1, -1},
			argTestOut{291, 2, 6, 1225, 0, ""},
		},
		{
			argTestIn{2400, 2400, -1, -1, -1, -1, -1},
			argTestOut{506, 3, 8, 1225, 160000, ""},
		},
		{
			argTestIn{3600, 3600, -1, -1, -1, -1, -1},
			argTestOut{700, 3, 9, 1794, 360000, ""},
		},
	} {
		size, minDepth, maxDepth, minArea, maxArea, err := macroQuadFixArgs(
			tt.t.width, tt.t.height, tt.t.size, tt.t.minDepth, tt.t.maxDepth, tt.t.minArea, tt.t.maxArea)
		var errStr string
		if err != nil {
			errStr = err.Error()
		} else {
			errStr = ""
		}
		if tt.e.size != size ||
			tt.e.minDepth != minDepth ||
			tt.e.maxDepth != maxDepth ||
			tt.e.minArea != minArea ||
			tt.e.maxArea != maxArea ||
			tt.e.err != errStr {
			t.Errorf("macroQuadFixArgs(%d, %d, %d, %d, %d, %d, %d) => (%d, %d, %d, %d, %d, %s), expect (%d, %d, %d, %d, %d, %s)",
				tt.t.width, tt.t.height, tt.t.size, tt.t.minDepth, tt.t.maxDepth, tt.t.minArea, tt.t.maxArea,
				size, minDepth, maxDepth, minArea, maxArea, errStr,
				tt.e.size, tt.e.minDepth, tt.e.maxDepth, tt.e.minArea, tt.e.maxArea, tt.e.err)
		}
	}
}

func TestMacroQuad(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	cover, macro := MacroQuad(env, "testdata/jumping_bunny.jpg", 200, 200, 10, -1, 2, 50, -1, "", "")
	if cover == nil || macro == nil {
		fmt.Println(out.String())
		t.Fatal("Failed to create cover or macro")
	}

	expect := []string{
		"Building macro quad with 10 splits, 34 partials, min depth 1, max depth 2, min area 50...",
	}

	testResultExpect(t, out.String(), expect)
}

func TestMacroQuadMinDepthSplits(t *testing.T) {
	for _, tt := range []struct {
		a int
		r int
	}{
		{0, 0},
		{1, 1},
		{2, 5},
		{3, 21},
		{4, 85},
		{5, 341},
	} {
		r := macroQuadMinDepthSplits(tt.a)
		if r != tt.r {
			t.Errorf("macroQuadMinDepthSplits(%d) => %d, want %d", tt.a, r, tt.r)
		}
	}
}
