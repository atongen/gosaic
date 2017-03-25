package model

type Aspect struct {
	Id      int64 `db:"id"`
	Columns int   `db:"columns"`
	Rows    int   `db:"rows"`
}

func NewAspect(columns int, rows int) *Aspect {
	c, r := calculateAspect(columns, rows)
	return &Aspect{
		Columns: c,
		Rows:    r,
	}
}

func (a *Aspect) RoundWidth(height int) int {
	return round(float64(a.Columns) * float64(height) / float64(a.Rows))
}

func (a *Aspect) RoundHeight(width int) int {
	return round(float64(a.Rows) * float64(width) / float64(a.Columns))
}

func (a *Aspect) Ratio() float64 {
	return float64(a.Columns) / float64(a.Rows)
}

// Scale returns a width and height that is closest
// to the provided width and height, but maintains this
// exact integer aspect ratio, without rounding
func (a *Aspect) Scale(width, height int) (int, int) {
	b := NewAspect(width, height)

	var w, h int

	if a.Ratio() < b.Ratio() {
		// lock width
		if width < a.Columns {
			w = a.Columns
			h = a.Rows
		} else {
			n := round(float64(width) / float64(a.Columns))
			w = n * a.Columns
			h = n * a.Rows
		}
	} else {
		// lock height
		if height < a.Rows {
			h = a.Rows
			w = a.Columns
		} else {
			n := round(float64(height) / float64(a.Rows))
			h = n * a.Rows
			w = n * a.Columns
		}
	}

	return w, h
}

func (a *Aspect) ScaleRound(width, height int) (int, int) {
	b := NewAspect(width, height)

	var w, h int

	if a.Ratio() < b.Ratio() {
		w = width
		h = a.RoundHeight(width)
	} else {
		h = height
		w = a.RoundWidth(height)
	}

	return w, h
}
