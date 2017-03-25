package util

import "testing"

func TestRound(t *testing.T) {
	for _, tt := range []struct {
		a float64
		r int
	}{
		{-1.9, -2},
		{-1.5, -2},
		{-1.1, -1},
		{-0.5, -1},
		{0.0, 0},
		{0.5, 1},
		{1.1, 1},
		{1.6, 2},
	} {
		r := Round(tt.a)
		if r != tt.r {
			t.Errorf("Round(%f) => %d, want %d", tt.a, r, tt.r)
		}
	}
}

func TestCleanStr(t *testing.T) {
	for _, tt := range []struct {
		a string
		r string
	}{
		{" This is a regular sentence. ", "this_is_a_regular_sentence"},
		{" % dk ## dkdkd EDDED ^% dkfdsf# ", "dk_dkdkd_edded_dkfdsf"},
		{"Keep on cleaning!", "keep_on_cleaning"},
		{`¯\_(ツ)_/¯`, ""},
		{"this_is_unchanged", "this_is_unchanged"},
		{" FFFFUUUUUUUUU !!!! ", "ffffuuuuuuuuu"},
		{`ain't that some "shizznit"?`, "ain_t_that_some_shizznit"},
		{"2016-08-31: That was the day!", "2016-08-31_that_was_the_day"},
	} {
		r := CleanStr(tt.a)
		if r != tt.r {
			t.Errorf("CleanStr(\"%s\") => \"%s\", want \"%s\"", tt.a, r, tt.r)
		}
	}
}
