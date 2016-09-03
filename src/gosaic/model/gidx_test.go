package model

import "testing"

func gidxWithinTest(width, height int) *Gidx {
	return &Gidx{
		Width:  width,
		Height: height,
	}
}

func TestGidxWithin(t *testing.T) {
	for _, tt := range []struct {
		g *Gidx
		t float64
		a *Aspect
		r bool
	}{
		{gidxWithinTest(1, 1), 0.0, NewAspect(1, 1), true},
		// square aspect
		{gidxWithinTest(100, 500), 1.0, NewAspect(1, 1), true},
		{gidxWithinTest(200, 500), 1.0, NewAspect(1, 1), true},
		{gidxWithinTest(300, 500), 1.0, NewAspect(1, 1), true},
		{gidxWithinTest(400, 500), 1.0, NewAspect(1, 1), true},
		{gidxWithinTest(500, 500), 1.0, NewAspect(1, 1), true},
		{gidxWithinTest(600, 500), 1.0, NewAspect(1, 1), true},
		{gidxWithinTest(700, 500), 1.0, NewAspect(1, 1), true},
		{gidxWithinTest(800, 500), 1.0, NewAspect(1, 1), true},
		{gidxWithinTest(900, 500), 1.0, NewAspect(1, 1), true},
		// portait aspect
		{gidxWithinTest(100, 500), 1.0, NewAspect(1, 2), true},
		{gidxWithinTest(200, 500), 1.0, NewAspect(1, 2), true},
		{gidxWithinTest(300, 500), 1.0, NewAspect(1, 2), true},
		{gidxWithinTest(400, 500), 1.0, NewAspect(1, 2), true},
		{gidxWithinTest(500, 500), 1.0, NewAspect(1, 2), true},
		{gidxWithinTest(600, 500), 1.0, NewAspect(1, 2), true},
		{gidxWithinTest(700, 500), 1.0, NewAspect(1, 2), true},
		{gidxWithinTest(800, 500), 1.0, NewAspect(1, 2), false},
		{gidxWithinTest(900, 500), 1.0, NewAspect(1, 2), false},
		// landscape aspect
		{gidxWithinTest(100, 500), 1.0, NewAspect(2, 1), false},
		{gidxWithinTest(200, 500), 1.0, NewAspect(2, 1), false},
		{gidxWithinTest(300, 500), 1.0, NewAspect(2, 1), false},
		{gidxWithinTest(400, 500), 1.0, NewAspect(2, 1), false},
		{gidxWithinTest(500, 500), 1.0, NewAspect(2, 1), true},
		{gidxWithinTest(600, 500), 1.0, NewAspect(2, 1), true},
		{gidxWithinTest(700, 500), 1.0, NewAspect(2, 1), true},
		{gidxWithinTest(800, 500), 1.0, NewAspect(2, 1), true},
		{gidxWithinTest(900, 500), 1.0, NewAspect(2, 1), true},
	} {
		r := tt.g.Within(tt.t, tt.a)
		if r != tt.r {
			t.Errorf("(%dx%d) within (%dx%d) by %f = %v, want %v",
				tt.g.Width, tt.g.Height, tt.a.Columns, tt.a.Rows, tt.t, r, tt.r)
		}
	}
}
